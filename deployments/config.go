package deployments

import (
	"bytes"
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
	b, err := os.ReadFile(deploymentsConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, deploymentsConfig)
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
	err := marshalYamlToFile(serializableDeployments, deploymentsConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func marshalYamlToFile(in interface{}, path string) (err error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	defer func() {
		err = encoder.Close()
	}()
	encoder.SetIndent(2)
	err = encoder.Encode(in)
	if err != nil {
		return
	}
	err = os.WriteFile(path, buffer.Bytes(), 0644)
	return
}
