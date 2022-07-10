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

func LoadDeploymentsFromFile(cliArgs []string) ([]engine.Deployment, error) {
	deploymentsConfigPath := cliArgs[0]
	deploymentsConfig := &Deployments{}
	bytes, err := ioutil.ReadFile(deploymentsConfigPath)
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
