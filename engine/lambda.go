package engine

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"strings"
)

type LambdaDeployment struct {
	BaseDeploymentConfig `mapstructure:",squash"`
}

func (d LambdaDeployment) Id() string {
	return d.Id_
}

func (d LambdaDeployment) Type() string {
	return d.Type_
}

func (d LambdaDeployment) Version() string {
	return d.Version_
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

func (engine *LambdaDeploymentEngine) Deploy(config Deployment) error {
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
