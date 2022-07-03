package engine

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var awsConfig = CreateAwsConfig()

func CreateAwsConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	cobra.CheckErr(err)
	return cfg
}
