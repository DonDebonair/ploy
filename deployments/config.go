package deployments

import (
	"github.com/DandyDev/ploy/engine"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Deployments struct {
	Deployments []map[string]any `yaml:"deployments"`
}

func LoadDeploymentsFromFile(path string) ([]engine.Deployment, error) {
	deploymentsConfig := &Deployments{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, deploymentsConfig)
	if err != nil {
		return nil, err
	}
	deployments := make([]engine.Deployment, 0, len(deploymentsConfig.Deployments))
	for _, d := range deploymentsConfig.Deployments {
		deploymentEngine := engine.GetEngine(d["type"].(string))
		deploymentConfig := deploymentEngine.ResolveConfigStruct()
		err = mapstructure.Decode(d, deploymentConfig)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deploymentConfig)
	}
	return deployments, nil
}
