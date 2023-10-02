# Welcome to your CDK Go project

This is a blank project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Local Testing Commands

- `cdk synth` emits the synthesized CloudFormation template
- `sam local start-api -t ./cdk.out/GlambdaKitStack.template.json --env-vars .env.json`

## Useful commands

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `go test` run unit tests
- `sam local start-api` Start lambdas locally (need to have template yml file in the root directory)
