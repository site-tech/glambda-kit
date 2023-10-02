package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Hello World function
func getGoHelloHandler(stack awscdk.Stack) awscdklambdagoalpha.GoFunction {
	return awscdklambdagoalpha.NewGoFunction(stack, jsii.String("hello"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("./lambda/golang/hello"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

}

// Create Square Booking POST lambda
func postCreateBookingHandler(stack awscdk.Stack) awscdklambdagoalpha.GoFunction {
	return awscdklambdagoalpha.NewGoFunction(stack, jsii.String("createSquareBooking"), &awscdklambdagoalpha.GoFunctionProps{
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
