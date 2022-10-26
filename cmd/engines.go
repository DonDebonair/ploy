package cmd

import (
	"fmt"
	"github.com/DandyDev/ploy/engine"

	"github.com/spf13/cobra"
)

// enginesCmd represents the engines command
var enginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "List supported deployment engines",
	Long:  `List the deployment engines that ploy supports`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deployment engines:")
		for _, e := range engine.ListEngines() {
			fmt.Printf("- %s\n", e)
		}
	},
}

func init() {
	rootCmd.AddCommand(enginesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enginesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enginesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
