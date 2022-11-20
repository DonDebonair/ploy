package deployments

import (
	"bufio"
	"fmt"
	"github.com/DonDebonair/ploy/engine"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Deploy(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args[0])
	cobra.CheckErr(err)

	errorChan := make(chan error, len(deployments))
	var wg sync.WaitGroup
	for _, deployment := range deployments {
		wg.Add(1)
		go func(deployment engine.Deployment) {
			defer wg.Done()
			errorChan <- doDeployment(deployment)
		}(deployment)
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
		if len(deploymentConfig.PreDeployCommand()) > 0 {
			if err = runDeploymentScript("pre", deploymentConfig.PreDeployCommand(), deploymentConfig.Version(), p); err != nil {
				return err
			}
		}
		if err = deploymentEngine.Deploy(deploymentConfig, p); err != nil {
			return err
		}
		if len(deploymentConfig.PostDeployCommand()) > 0 {
			if err = runDeploymentScript("post", deploymentConfig.PostDeployCommand(), deploymentConfig.Version(), p); err != nil {
				return err
			}
		}
		p("version %s deployed successfully!", deploymentConfig.Version())
		return nil
	} else {
		p("version '%s' matches expected version '%s'. Skipping...", version, deploymentConfig.Version())
		return nil
	}
}

func runDeploymentScript(context string, deploymentCommand []string, version string, p func(string, ...any)) error {
	p("running %s-deployment command %s...", context, strings.Join(deploymentCommand, " "))
	cmd := exec.Command(deploymentCommand[0], deploymentCommand[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("VERSION=%s", version))
	output, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
