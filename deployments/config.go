package deployments

import (
	"github.com/DonDebonair/ploy/engine"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"os"
)

type Deployments struct {
	Deployments []map[string]any `yaml:"deployments"`
}

func LoadDeploymentsFromFile(cliArgs []string) ([]engine.Deployment, error) {
	deploymentsConfigPath := cliArgs[0]
	deploymentsConfig := &Deployments{}
	bytes, err := os.ReadFile(deploymentsConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, deploymentsConfig)
	if err != nil {
		return nil, err
	}
	deployments := make([]engine.Deployment, 0, len(deploymentsConfig.Deployments))
	for _, deployment := range deploymentsConfig.Deployments {
		deploymentEngine, err := engine.GetEngine(deployment["type"].(string))
		if err != nil {
			return nil, err
		}
		deploymentConfig := deploymentEngine.ResolveConfigStruct()
		err = mapstructure.Decode(deployment, deploymentConfig)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deploymentConfig)
	}
	return deployments, nil
}
