package engine

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"strings"
)

type EcsDeployment struct {
	BaseDeploymentConfig `mapstructure:",squash"`
	Cluster              string `mapstructure:"cluster"`
}

func (d EcsDeployment) Id() string {
	return d.Id_
}

func (d EcsDeployment) Type() string {
	return d.Type_
}

func (d EcsDeployment) Version() string {
	return d.Version_
}

type ECSDeploymentEngine struct {
	ECSClient *ecs.Client
}

func (engine *ECSDeploymentEngine) Type() string {
	return "ecs"
}

func (engine *ECSDeploymentEngine) ResolveConfigStruct() Deployment {
	return &EcsDeployment{}
}

func (engine *ECSDeploymentEngine) Deploy(config Deployment) error {
	//TODO implement me
	panic("implement me")
}

// TODO: error handling if service can't be found
// TODO: deal with task definitions without a service (i.e. one-off tasks)
func (engine *ECSDeploymentEngine) CheckVersion(config Deployment) (string, error) {
	ecsConfig := config.(*EcsDeployment)
	services, err := engine.ECSClient.DescribeServices(
		context.Background(),
		&ecs.DescribeServicesInput{
			Services: []string{ecsConfig.Id()},
			Cluster:  aws.String(ecsConfig.Cluster),
		},
	)
	if err != nil {
		return "", err
	}
	taskDefinitionArn := services.Services[0].Deployments[0].TaskDefinition
	taskDefinition, err := engine.ECSClient.DescribeTaskDefinition(context.Background(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: taskDefinitionArn,
	})
	if err != nil {
		return "", err
	}
	return strings.Split(*taskDefinition.TaskDefinition.ContainerDefinitions[0].Image, ":")[1], nil
}

func init() {
	RegisterDeploymentEngine("ecs", func(awsConfig aws.Config) DeploymentEngine {
		return &ECSDeploymentEngine{ECSClient: ecs.NewFromConfig(awsConfig)}
	})
}
