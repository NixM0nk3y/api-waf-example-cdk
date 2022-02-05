package main

import (
	"api-sam-example-cdk/pkg/stacks"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"

	"github.com/kelseyhightower/envconfig"
)

func main() {

	log.Print("Starting Application Build")

	app := awscdk.NewApp(&awscdk.AppProps{
		AnalyticsReporting: jsii.Bool(false),
	})

	applicationProps := stacks.ApplicationProps{}

	err := envconfig.Process("cdk", &applicationProps)

	if err != nil {
		log.Fatal(err.Error())
	}

	applicationProps.StackProps = awscdk.StackProps{
		Env: env(),
	}

	id := fmt.Sprintf("%s%sWafStack", strings.Title(applicationProps.Tenant), strings.Title(applicationProps.Environment))

	stacks.ApplicationStack(app, id, &applicationProps)

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
