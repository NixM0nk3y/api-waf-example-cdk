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
		"wafconfig/coraza-additional-rules.conf",
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

	tx.ProcessConnection(event.RequestContext.HTTP.SourceIP, 55555, "1.1.1.1", 443)

	tx.ProcessURI(fmt.Sprintf("%s?%s", event.RawPath, event.RawQueryString), event.RequestContext.HTTP.Method, event.RequestContext.HTTP.Protocol)

	for k, v := range event.Headers {
		tx.AddRequestHeader(k, v)
	}

	// process phase 1
	if it := tx.ProcessRequestHeaders(); it != nil {
		logger.Error("Transaction was interrupted in phase1", zap.Int("status", it.Status), zap.Int("ruleid", it.RuleID), zap.String("action", it.Action))

		al := tx.AuditLog()
		if len(al.Messages) > 0 {
			for _, auditevent := range al.Messages {
				logger.Error("auditevent", zap.Reflect("event", auditevent))
			}
		}

		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	// process phase 2
	if it, _ := tx.ProcessRequestBody(); it != nil {
		logger.Error("Transaction was interrupted in phase2", zap.Int("status", it.Status), zap.Int("ruleid", it.RuleID), zap.String("action", it.Action))

		al := tx.AuditLog()
		if len(al.Messages) > 0 {
			for _, auditevent := range al.Messages {
				logger.Error("auditevent", zap.Reflect("event", auditevent))
			}
		}

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
