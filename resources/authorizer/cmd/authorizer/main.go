package main

import (
	"context"

	"authorizer/internal/config"

	"authorizer/internal/authorizer"
	"authorizer/pkg/awsemf"
	"authorizer/pkg/log"
	"authorizer/pkg/version"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/prozz/aws-embedded-metrics-golang/emf"
	"go.uber.org/zap"
)

// Handler is
func Handler(ctx context.Context, req authorizer.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	logger := log.LoggerWithLambdaRqID(ctx)

	xray.Configure(xray.Config{
		LogLevel:       "warn",
		ServiceVersion: version.Version,
	})

	logger.Info("application handler")

	logger.Debug("recieved event", zap.Reflect("req", req))

	ctx = awsemf.NewEMF(ctx)

	vctx := config.ReadEnvConfig(ctx, "APPLICATION")

	metricslogger := awsemf.GetEMF(ctx)
	defer metricslogger.Log()

	metricslogger.Dimension("openenterprise", "lambda").MetricAs("Invokes", 1, emf.None)

	return authorizer.HandleRequest(vctx, req)
}

func main() {
	lambda.Start(Handler)
}
