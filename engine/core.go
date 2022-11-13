package engine

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
)

var engineRegistry = make(map[string]DeploymentEngine)

type Deployment interface {
	Id() string
	Type() string
	Version() string
	PostDeployCommand() []string
	SetVersion(version string)
}

type BaseDeploymentConfig struct {
	Id_                string   `mapstructure:"id"`
	Type_              string   `mapstructure:"type"`
	Version_           string   `mapstructure:"version"`
	PostDeployCommand_ []string `mapstructure:"post-deploy-command,omitempty"`
}

func (d BaseDeploymentConfig) Id() string {
	return d.Id_
}

func (d BaseDeploymentConfig) Type() string {
	return d.Type_
}

func (d BaseDeploymentConfig) Version() string {
	return d.Version_
}

func (d BaseDeploymentConfig) PostDeployCommand() []string {
	return d.PostDeployCommand_
}

func (d *BaseDeploymentConfig) SetVersion(version string) {
	d.Version_ = version
}

type DeploymentEngine interface {
	Type() string
	ResolveConfigStruct() Deployment
	Deploy(deploymentConfig Deployment, printer func(string, ...any)) error
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

func GetEngine(id string) (DeploymentEngine, error) {
	engine, ok := engineRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown deployment engine %s", id)
	}
	return engine, nil
}
