package authorizer

import (
	"authorizer/pkg/log"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/jptosso/coraza-libinjection"
	_ "github.com/jptosso/coraza-pcre"
	"github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/seclang"
	"go.uber.org/zap"
)

type APIGatewayV2CustomAuthorizerV2Request struct {
	Version               string                                `json:"version"`
	Type                  string                                `json:"type"`
	RouteArn              string                                `json:"routeArn"`
	IdentitySource        []string                              `json:"identitySource"`
	RouteKey              string                                `json:"routeKey"`
	RawPath               string                                `json:"rawPath"`
	RawQueryString        string                                `json:"rawQueryString"`
	Cookies               []string                              `json:"cookies"`
	Headers               map[string]string                     `json:"headers"`
	QueryStringParameters map[string]string                     `json:"queryStringParameters"`
	RequestContext        events.APIGatewayV2HTTPRequestContext `json:"requestContext"`
	PathParameters        map[string]string                     `json:"pathParameters"`
	StageVariables        map[string]string                     `json:"stageVariables"`
}

// our router
var waf *coraza.Waf

func init() {
	waf = coraza.NewWaf()
	parser, _ := seclang.NewParser(waf)

	files := []string{
		"wafconfig/coraza.conf",
		"wafconfig/coreruleset/crs-setup.conf",
		"wafconfig/coreruleset/rules/*.conf",
	}
	for _, f := range files {
		if err := parser.FromFile(f); err != nil {
			panic(err)
		}
	}

}

func HandleRequest(ctx context.Context, event APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	logger := log.Logger(ctx)

	tx := waf.NewTransaction()

	defer func() {
		tx.ProcessLogging()
		tx.Clean()
	}()

	tx.ProcessConnection(event.RequestContext.HTTP.SourceIP, 55555, "127.0.0.1", 443)

	tx.ProcessURI(fmt.Sprintf("%s?%s", event.RawPath, event.RawQueryString), event.RequestContext.HTTP.Method, event.RequestContext.HTTP.Protocol)

	for k, v := range event.Headers {
		tx.AddRequestHeader(k, v)
	}

	if it := tx.ProcessRequestHeaders(); it != nil {
		logger.Error("Transaction was interrupted", zap.Int("status", it.Status), zap.Int("ruleid", it.RuleID), zap.String("action", it.Action))
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	logger.Error("Transaction passed waf scanning")

	// we've passed the waf
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
	}, nil

}
