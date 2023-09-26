package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/subosito/gotenv"

	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func init() {
	if err := gotenv.Load(); err != nil {
		log.Println("no .env file")
	}
}

type GlambdaKitStackProps struct {
	awscdk.StackProps
}

func NewGlambdaKitStack(scope constructs.Construct, id string, props *GlambdaKitStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	layer := awslambda.NewLayerVersion(stack, jsii.String("MyLayer"), &awslambda.LayerVersionProps{
		Code: awslambda.Code_FromAsset(jsii.String("./lambda/python/deployment_package.zip"), nil),
	})

	// hello Lambda function
	getHelloHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("hello"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda/hello"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

	// Create Square Booking POST lambda
	createSquareBooking := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("createSquareBooking"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda/squareApi/createBooking"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
		Environment: &map[string]*string{
			"SQUARE_API_KEY":            jsii.String(os.Getenv("SQUARE_API_KEY")),
			"SERVICE_VARIATION_VERSION": jsii.String(os.Getenv("SERVICE_VARIATION_VERSION")),
			"SERVICE_VARIATION_ID":      jsii.String(os.Getenv("SERVICE_VARIATION_ID")),
			"TEAM_MEMBER_ID":            jsii.String(os.Getenv("TEAM_MEMBER_ID")),
			"CUSTOMER_ID":               jsii.String(os.Getenv("CUSTOMER_ID")),
			"LOCATION_ID":               jsii.String(os.Getenv("LOCATION_ID")),
		},
	})

	// NodeJS Andy Lambda function
	getAndyHandler := awslambda.NewFunction(stack, jsii.String("node"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("index.main"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/node/"), nil),
	})

	// Python 3 Lambda function
	getPythonHanlder := awslambda.NewFunction(stack, jsii.String("PythonLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Handler: jsii.String("index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/python"), nil),
		Layers:  &[]awslambda.ILayerVersion{layer},
		Environment: &map[string]*string{
			"VERTEX_JSON": jsii.String(os.Getenv("VERTEX_JSON")),
		},
		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
	})

	// create HTTP API Gateway
	httpApiGateway := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("glambda-api"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("glambda-api"),
	})

	// add Square Create Booking Lambda
	httpApiGateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/create/booking"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), createSquareBooking, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add hello route to HTTP API Gateway
	httpApiGateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/hello"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getHelloHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add NodeJS Andy route to HTTP API Gateway
	httpApiGateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/node"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getAndyHandler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add Python route to HTTP API Gateway
	httpApiGateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/python"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getPythonHanlder, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGlambdaKitStack(app, "GlambdaKitStack", &GlambdaKitStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
