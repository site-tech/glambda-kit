package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

// Hello World Nodejs Function
func getNodeHelloHanlder(stack awscdk.Stack) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String("node"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("index.main"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/node/hello/"), nil),
	})
}
