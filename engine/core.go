package engine

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var engineRegistry = make(map[string]DeploymentEngine)

type Deployment interface {
	Id() string
	Type() string
	Version() string
}

type BaseDeploymentConfig struct {
	Id_      string `mapstructure:"id"`
	Type_    string `mapstructure:"type"`
	Version_ string `mapstructure:"version"`
}

type DeploymentEngine interface {
	Type() string
	ResolveConfigStruct() Deployment
	Deploy(deploymentConfig Deployment) error
	CheckVersion(deploymentConfig Deployment) (string, error)
}

func RegisterDeploymentEngine(id string, engineConstructor func(config aws.Config) DeploymentEngine) {
	if engineConstructor == nil {
		panic("Engine constructor is nil")
	}
	if _, dup := engineRegistry[id]; dup {
		panic("Register called twice for engine " + id)
	}
	engine := engineConstructor(awsConfig)
	engineRegistry[id] = engine
}

func ListEngines() []string {
	keys := make([]string, 0, len(engineRegistry))
	for key := range engineRegistry {
		keys = append(keys, key)
	}
	return keys
}

func GetEngine(id string) DeploymentEngine {
	engine, ok := engineRegistry[id]
	if !ok {
		panic("Unknown deployment engine " + id)
	}
	return engine
}
