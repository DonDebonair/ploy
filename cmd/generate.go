package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"embed"
)

//go:embed templates/*
var templatesFS embed.FS

type HtmlTemplate struct {
	Title	  string 
	ID        string `yaml:"id"`
	Version   string `yaml:"version"`
	Type      string `yaml:"type"`
	Repository string `yaml:"repository"`
}

type Deployments struct {
	Title string `yaml:"name"`
	Deployments []HtmlTemplate `yaml:"deployments"`
}

func Generate(envFile string) error {
	var deployments Deployments

	yamlData, err := ioutil.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("error loading deployments file: %s", err)
	}

	// unmarshal the environment YAML data into the deployments struct
	err = yaml.Unmarshal(yamlData, &deployments)
	if err != nil {
		return fmt.Errorf("Error parsing YAML file: %s", err)
	}

	// read the env.html template file
	htmlTemplate, err := templatesFS.ReadFile("templates/env.html")
	if err != nil {
		return fmt.Errorf("Error loading env.html template file: %s", err)
	}

	// parse the HTML template
	tmpl, err := template.New("env").Parse(string(htmlTemplate))
	if err != nil {
		return fmt.Errorf("Error parsing HTML template: %s", err)
	}

	// create the output file
	os.MkdirAll("generated", os.ModePerm)
	outFile, err := os.Create(fmt.Sprintf("generated/%s.html", strings.Split(envFile, ".")[0]))
	if err != nil {
		return fmt.Errorf("Error creating output file: %s", err)
	}
	defer outFile.Close()

	// execute the HTML template
	err = tmpl.Execute(outFile, deployments)
	if err != nil {
		return fmt.Errorf("Error executing HTML template: %s", err)
	}

	return nil
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [env-file]",
	Short: "Generate deployment dashboard HTML file",
	Long:  `Generate deployment dashboard HTML file based on the specified environment YAML file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Usage: ploy generate [yaml-file]")
		}

		envFile := args[0]
		
		err := Generate(envFile)
		if err != nil {
			return err
		}

		fmt.Println("HTML file generated successfully.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
