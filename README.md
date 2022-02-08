# CDK Waf authorizer demo

Demo using coraza in a API gateway authorizer to protect api endpoints

# Stack Setup

As this demo uses a CGO enabled build for a arm64 target, a docker installation configured for arm64 builds will be required.

https://www.docker.com/blog/multi-arch-images/

```bash
go build -v ./cmd/application
🛠️  cmd/application done
✓  stacks/build done
✓  build done
cdk deploy --app ./application
2022/02/08 19:46:26 Starting Application Build
Bundling asset OpenenterpriseProductionWafStack/Hosting/Lambda/Code/Stage...
Bundling asset OpenenterpriseProductionWafStack/Hosting/AuthLambda/Code/Stage...
WARNING: The requested image's platform (linux/arm64/v8) does not match the detected host platform (linux/amd64) and no specific platform was requested
go: downloading github.com/aws/aws-lambda-go v1.28.0
...snip...
go: downloading golang.org/x/text v0.3.6

✨  Synthesis time: 389.27s

OpenenterpriseProductionWafStack: deploying...
[0%] start: Publishing 9389c3589f32bfd1b87004861d601e2975779b9cab5e93b343b6ce714f1be21b:current
[33%] success: Published 9389c3589f32bfd1b87004861d601e2975779b9cab5e93b343b6ce714f1be21b:current
[33%] start: Publishing 6c0316fef24d0df8a9a705c77052001217d864f49af386539d01df54618cd131:current
[66%] success: Published 6c0316fef24d0df8a9a705c77052001217d864f49af386539d01df54618cd131:current
[66%] start: Publishing 03a6956a2874eef697c616680543da701996309b495a79361634bb7569687fd6:current
[100%] success: Published 03a6956a2874eef697c616680543da701996309b495a79361634bb7569687fd6:current
OpenenterpriseProductionWafStack: creating CloudFormation changeset...

 ✅  OpenenterpriseProductionWafStack

✨  Deployment time: 72.56s

Outputs:
OpenenterpriseProductionWafStack.HostingUrlOutput7A35DF00 = https://aaaaaaaaa.execute-api.eu-west-1.amazonaws.com/
Stack ARN:
arn:aws:cloudformation:eu-west-1:074705540277:stack/OpenenterpriseProductionWafStack/d6550f30-867f-11ec-98f3-0a8ba53abf81

✨  Total time: 461.83s

🛠️  deploy/application done
✓  deploy done
```

# License

MIT

## Useful commands

 * `make deploy`             deploy this stack to your default AWS account/region
 * `make waf/test/auth`      call the authorizor with a clean request
 * `make waf/test/authblock` call the authorizor with a blocking request
 * `make waf/test/hello`     call the hello endpoint
 * `make waf/test/version`   call the version endpoint
