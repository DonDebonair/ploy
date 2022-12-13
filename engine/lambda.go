package engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"strings"
	"text/template"
	"time"
)

type LambdaDeployment struct {
	BaseDeploymentConfig  `mapstructure:",squash"`
	VersionEnvironmentKey string `mapstructure:"version-environment-key,omitempty"`
	Bucket                string `mapstructure:"bucket,omitempty"`
	KeyTemplate           string `mapstructure:"key,omitempty"`
	VersionTag            string `mapstructure:"version-tag,omitempty"`
}

type LambdaDeploymentEngine struct {
	LambdaClient *lambda.Client
}

type keyTemplateVars struct {
	Id      string
	Version string
}

func (engine *LambdaDeploymentEngine) Type() string {
	return "lambda"
}

func (engine *LambdaDeploymentEngine) ResolveConfigStruct() Deployment {
	return &LambdaDeployment{}
}

func (engine *LambdaDeploymentEngine) Deploy(config Deployment, p func(string, ...any)) error {
	lambdaConfig := config.(*LambdaDeployment)
	getFunctionInput := &lambda.GetFunctionInput{
		FunctionName: aws.String(lambdaConfig.Id()),
	}
	waiter := lambda.NewFunctionUpdatedV2Waiter(engine.LambdaClient)
	// TODO: in case of PackageType `image`, we get the image of the deployed Lambda, but Deploy() is always called after CheckVersion() in deployments.doDeployment(). So this is a bit of a waste...
	functionInfo, err := engine.LambdaClient.GetFunction(
		context.Background(),
		getFunctionInput,
	)
	if err != nil {
		return err
	}
	if lambdaConfig.VersionEnvironmentKey != "" {
		p("Updating function configuration")
		environment := functionInfo.Configuration.Environment.Variables
		environment[lambdaConfig.VersionEnvironmentKey] = lambdaConfig.Version()
		_, err = engine.LambdaClient.UpdateFunctionConfiguration(
			context.Background(),
			&lambda.UpdateFunctionConfigurationInput{
				FunctionName: aws.String(lambdaConfig.Id()),
				Environment: &types.Environment{
					Variables: environment,
				},
			},
		)
		// Wait for function configuration to be updated
		err = waiter.Wait(context.Background(), getFunctionInput, 1*time.Minute)
		if err != nil {
			return err
		}
	}
	switch functionInfo.Configuration.PackageType {
	case types.PackageTypeImage:
		p("Image deployment detected")
		err = deployImage(lambdaConfig, functionInfo, engine.LambdaClient, waiter, p)
	case types.PackageTypeZip:
		p("Zip deployment detected")
		err = deployZip(lambdaConfig, functionInfo, engine.LambdaClient, waiter, p)
	default:
		err = fmt.Errorf("package type %s not supported", functionInfo.Configuration.PackageType)
	}
	return err
}

func deployImage(
	deploymentConfig *LambdaDeployment,
	functionInfo *lambda.GetFunctionOutput,
	lambdaClient *lambda.Client,
	waiter *lambda.FunctionUpdatedV2Waiter,
	p func(string, ...any),
) error {
	getFunctionInput := &lambda.GetFunctionInput{
		FunctionName: aws.String(deploymentConfig.Id()),
	}
	image := strings.Split(*functionInfo.Code.ImageUri, ":")[0]
	imageUri := image + ":" + deploymentConfig.Version()
	p("Updating to version '%s' by setting new image uri '%s'", deploymentConfig.Version(), imageUri)
	_, err := lambdaClient.UpdateFunctionCode(
		context.Background(),
		&lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(deploymentConfig.Id()),
			ImageUri:     aws.String(imageUri),
			Publish:      true,
		},
	)
	if err != nil {
		return err
	}
	err = waiter.Wait(context.Background(), getFunctionInput, 1*time.Minute)
	return err
}

func deployZip(
	deploymentConfig *LambdaDeployment,
	functionInfo *lambda.GetFunctionOutput,
	lambdaClient *lambda.Client,
	waiter *lambda.FunctionUpdatedV2Waiter,
	p func(string, ...any),
) error {
	getFunctionInput := &lambda.GetFunctionInput{
		FunctionName: aws.String(deploymentConfig.Id()),
	}
	s3Key, err := resolveS3Key(deploymentConfig)
	if err != nil {
		return err
	}
	p("Updating to version '%s' with new zip file '%s'", deploymentConfig.Version(), s3Key)
	_, err = lambdaClient.UpdateFunctionCode(
		context.Background(),
		&lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(deploymentConfig.Id()),
			S3Bucket:     aws.String(deploymentConfig.Bucket),
			S3Key:        aws.String(s3Key),
			Publish:      true,
		},
	)
	if err != nil {
		return err
	}
	err = waiter.Wait(context.Background(), getFunctionInput, 1*time.Minute)
	if err != nil {
		return err
	}
	p("Setting version tag %s to %s", deploymentConfig.VersionTag, deploymentConfig.Version())
	_, err = lambdaClient.TagResource(
		context.Background(),
		&lambda.TagResourceInput{
			Resource: functionInfo.Configuration.FunctionArn,
			Tags:     map[string]string{deploymentConfig.VersionTag: deploymentConfig.Version()},
		},
	)
	return err
}

// Resolve key template to actual S3 key
func resolveS3Key(lambdaConfig *LambdaDeployment) (string, error) {
	keyTemplate, err := template.New(lambdaConfig.Id()).Parse(lambdaConfig.KeyTemplate)
	if err != nil {
		return "", err
	}
	var resolvedKey bytes.Buffer
	templateVars := keyTemplateVars{
		Id:      lambdaConfig.Id(),
		Version: lambdaConfig.Version(),
	}
	err = keyTemplate.Execute(&resolvedKey, templateVars)
	if err != nil {
		return "", err
	}
	return resolvedKey.String(), nil
}

func (engine *LambdaDeploymentEngine) CheckVersion(config Deployment) (string, error) {
	lambdaConfig := config.(*LambdaDeployment)
	functionInfo, err := engine.LambdaClient.GetFunction(
		context.Background(),
		&lambda.GetFunctionInput{
			FunctionName: aws.String(lambdaConfig.Id()),
		},
	)
	if err != nil {
		return "", err
	}
	var version string
	switch functionInfo.Configuration.PackageType {
	case types.PackageTypeImage:
		version = getImageVersion(functionInfo)
		err = nil
	case types.PackageTypeZip:
		version, err = getZipVersion(functionInfo, lambdaConfig)
	default:
		err = fmt.Errorf("package type %s not supported", functionInfo.Configuration.PackageType)
	}
	return version, err
}

func getImageVersion(functionInfo *lambda.GetFunctionOutput) string {
	return strings.Split(*functionInfo.Code.ImageUri, ":")[1]
}

func getZipVersion(funtionInfo *lambda.GetFunctionOutput, lambdaConfig *LambdaDeployment) (string, error) {
	version, present := funtionInfo.Tags[lambdaConfig.VersionTag]
	if present == false {
		return "", fmt.Errorf("version tag '%s' not found", lambdaConfig.VersionTag)
	}
	return version, nil
}

func init() {
	RegisterDeploymentEngine("lambda", func(awsConfig aws.Config) DeploymentEngine {
		return &LambdaDeploymentEngine{LambdaClient: lambda.NewFromConfig(awsConfig)}
	})
}
