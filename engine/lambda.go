package engine

import (
	"context"
	"github.com/avast/retry-go/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"strings"
)

type LambdaDeployment struct {
	BaseDeploymentConfig  `mapstructure:",squash"`
	VersionEnvironmentKey string `mapstructure:"version-environment-key"`
}

type LambdaDeploymentEngine struct {
	LambdaClient *lambda.Client
}

func (engine *LambdaDeploymentEngine) Type() string {
	return "lambda"
}

func (engine *LambdaDeploymentEngine) ResolveConfigStruct() Deployment {
	return &LambdaDeployment{}
}

func (engine *LambdaDeploymentEngine) Deploy(config Deployment, _ func(string, ...any)) error {
	lambdaConfig := config.(*LambdaDeployment)
	// TODO: we get the image of the deployed Lambda, but Deploy() is always called after CheckVersion() in deployments.doDeployment(). So this is a bit of a waste...
	functionInfo, err := engine.LambdaClient.GetFunction(
		context.Background(),
		&lambda.GetFunctionInput{
			FunctionName: aws.String(lambdaConfig.Id()),
		},
	)
	if err != nil {
		return err
	}
	if lambdaConfig.VersionEnvironmentKey != "" {
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
	}
	image := strings.Split(*functionInfo.Code.ImageUri, ":")[0]
	imageUri := image + ":" + config.Version()

	// we wrap updating the function code in a retry loop to handle the case where the function is still updating
	// this happens for example when we've updated the configuration
	err = retry.Do(func() error {
		_, err := engine.LambdaClient.UpdateFunctionCode(
			context.Background(),
			&lambda.UpdateFunctionCodeInput{
				FunctionName: aws.String(lambdaConfig.Id()),
				ImageUri:     aws.String(imageUri),
				Publish:      true,
			},
		)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (engine *LambdaDeploymentEngine) CheckVersion(config Deployment) (string, error) {
	functionInfo, err := engine.LambdaClient.GetFunction(
		context.Background(),
		&lambda.GetFunctionInput{
			FunctionName: aws.String(config.Id()),
		},
	)
	if err != nil {
		return "", err
	}
	return strings.Split(*functionInfo.Code.ImageUri, ":")[1], nil
}

func init() {
	RegisterDeploymentEngine("lambda", func(awsConfig aws.Config) DeploymentEngine {
		return &LambdaDeploymentEngine{LambdaClient: lambda.NewFromConfig(awsConfig)}
	})
}
