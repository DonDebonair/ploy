package deployments

import (
	"fmt"
	"github.com/DandyDev/ploy/engine"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Deploy(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args)
	cobra.CheckErr(err)

	errorChan := make(chan error, len(deployments))
	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		deployment := deployment
		go func() {
			defer wg.Done()
			errorChan <- doDeployment(deployment)
		}()
	}
	wg.Wait()
	close(errorChan)
	var result *multierror.Error
	for err := range errorChan {
		result = multierror.Append(result, err)
	}
	cobra.CheckErr(result.ErrorOrNil())
}

func doDeployment(deploymentConfig engine.Deployment) error {
	p := CreateDeploymentPrinter(deploymentConfig.Id())
	deploymentEngine, err := engine.GetEngine(deploymentConfig.Type())
	if err != nil {
		return err
	}
	p("checking deployed version...")
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		p("version '%s' does not match expected version '%s'. Deploying new version...", version, deploymentConfig.Version())
		err = deploymentEngine.Deploy(deploymentConfig, p)
		if err != nil {
			return err
		}
		if len(deploymentConfig.PostDeployCommand()) > 0 {
			commandString := strings.Join(deploymentConfig.PostDeployCommand(), " ")
			p("running post-deployment command %s...", commandString)
			if err != nil {
				return err
			}
			cmd := exec.Command(deploymentConfig.PostDeployCommand()[0], deploymentConfig.PostDeployCommand()[1:]...)
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "VERSION="+deploymentConfig.Version())
			output, err := cmd.Output()
			if err != nil {
				return err
			}
			fmt.Println(string(output))
		}
		p("version %s deployed successfully!", deploymentConfig.Version())
		return nil
	} else {
		p("version '%s' matches expected version '%s'. Skipping...", version, deploymentConfig.Version())
		return nil
	}
}
