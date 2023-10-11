package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/subosito/gotenv"

	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
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

	pythonLayer := awslambda.NewLayerVersion(stack, jsii.String("pythonLayer"), &awslambda.LayerVersionProps{
		Code: awslambda.Code_FromAsset(jsii.String("./lambda/python/deployment_package.zip"), nil),
	})

	// create HTTP API Gateway
	gateway := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("glambda-api"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("glambda-api"),
	})

	// add Square Create Booking Lambda
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/create/booking"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"),
			postCreateBookingHandler(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add hello route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/golang/hello"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"),
			getGoHelloHandler(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add NodeJS Andy route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/node/hello"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"),
			getNodeHelloHanlder(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add Python route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    jsii.String("/ai/vertex"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"),
			postGoogleAIVertexHanlder(stack, pythonLayer), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	return stack
}

// Hello World function
func getGoHelloHandler(stack awscdk.Stack) awscdklambdagoalpha.GoFunction {
	return awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GoHello"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda/golang/hello"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

}

// Create Square Booking POST lambda
func postCreateBookingHandler(stack awscdk.Stack) awscdklambdagoalpha.GoFunction {
	return awscdklambdagoalpha.NewGoFunction(stack, jsii.String("CreateSquareBooking"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda/golang/squareApi/createBooking"),
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
}

// Hello World Nodejs Function
func getNodeHelloHanlder(stack awscdk.Stack) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String("NodeHello"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("index.main"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/node/hello/"), nil),
	})
}

// Python 3 Lambda function
func postGoogleAIVertexHanlder(stack awscdk.Stack, layer awslambda.ILayerVersion) awslambda.Function {

	return awslambda.NewFunction(stack, jsii.String("GVertextAI"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Handler: jsii.String("index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/python/googleAI/vertex/"), nil),
		Layers:  &[]awslambda.ILayerVersion{layer},
		Environment: &map[string]*string{
			"VERTEX_JSON": jsii.String(os.Getenv("VERTEX_JSON")),
		},
		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
	})
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
