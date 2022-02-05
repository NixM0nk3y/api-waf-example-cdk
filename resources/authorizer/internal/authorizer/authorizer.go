package authorizer

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func HandleRequest(ctx context.Context, event APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {

	var token = "allow"

	switch strings.ToLower(token) {
	case "allow":
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: true,
		}, nil
	case "deny":
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	case "unauthorized":
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, errors.New("unauthorized") // Return a 401 Unauthorized response
	default:
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, errors.New("error: Invalid token")
	}
}
