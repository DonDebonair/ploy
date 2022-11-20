package deployments

import (
	"fmt"
	"strings"
)

func CreateDeploymentPrinter(deploymentId string) func(string, ...any) {
	return func(fmtString string, a ...any) {
		fmtString = "[" + deploymentId + "]" + ": " + fmtString
		if !strings.HasSuffix(fmtString, "\n") {
			fmtString = fmtString + "\n"
		}
		fmt.Printf(fmtString, a...)
	}
}

//func Update[E any](s []E, predicate func(E) bool, updateFunc func(E) E)
