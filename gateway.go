package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

func createHttpApiGateway(gateway awscdkapigatewayv2alpha.HttpApi, stack awscdk.Stack) {
	// add Square Create Booking Lambda
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/create/booking"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), postCreateBookingHandler(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add hello route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/golang/hello"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getGoHelloHandler(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add NodeJS Andy route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/node/hello"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getNodeHelloHanlder(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

	// add Python route to HTTP API Gateway
	gateway.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/ai/vertex"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), postGoogleAIVertexHanlder(stack), &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})

}
