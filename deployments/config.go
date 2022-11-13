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

func LoadDeploymentsFromFile(deploymentsConfigPath string) ([]engine.Deployment, error) {
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

func WriteDeploymentsToFile(deploymentsConfigPath string, deployments []engine.Deployment) error {
	deploymentMaps := make([]map[string]any, 0)
	for _, deployment := range deployments {
		var deploymentMap map[string]any
		err := mapstructure.Decode(deployment, &deploymentMap)
		if err != nil {
			return err
		}
		deploymentMaps = append(deploymentMaps, deploymentMap)
	}
	serializableDeployments := Deployments{Deployments: deploymentMaps}
	result, err := yaml.Marshal(serializableDeployments)
	if err != nil {
		return err
	}
	err = os.WriteFile(deploymentsConfigPath, result, 0644)
	if err != nil {
		return err
	}
	return nil
}
