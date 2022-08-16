package deployments

import (
	"fmt"
	"github.com/DandyDev/ploy/engine"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"sync"
)

var FailOnVersionMismatch bool

// TODO: should this command return a non-zero exit code when the versions of one or more deployments don't match?
func Verify(_ *cobra.Command, args []string) {
	deployments, err := LoadDeploymentsFromFile(args)
	cobra.CheckErr(err)
	var wg sync.WaitGroup
	errorChan := make(chan error, len(deployments))
	for _, deployment := range deployments {
		wg.Add(1)
		deployment := deployment
		go func() {
			defer wg.Done()
			errorChan <- verifyDeployment(deployment, FailOnVersionMismatch)
			cobra.CheckErr(err)
		}()
	}
	go func() {
		wg.Wait()
		close(errorChan)
	}()
	var result error
	for err := range errorChan {
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	cobra.CheckErr(result)
}

func verifyDeployment(deploymentConfig engine.Deployment, failOnVersionMismatch bool) error {
	p := CreateDeploymentPrinter(deploymentConfig.Id())
	deploymentEngine := engine.GetEngine(deploymentConfig.Type())
	version, err := deploymentEngine.CheckVersion(deploymentConfig)
	if err != nil {
		return err
	}
	if version != deploymentConfig.Version() {
		p("❌ version '%s' does not match expected version '%s'", version, deploymentConfig.Version())
		if failOnVersionMismatch {
			return fmt.Errorf("version '%s' does not match expected version '%s'", version, deploymentConfig.Version())
		}
	} else {
		p("✅ version '%s' matches expected version '%s'", version, deploymentConfig.Version())
	}
	return nil
}
