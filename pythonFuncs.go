package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

// Python 3 Lambda function
func postGoogleAIVertexHanlder(stack awscdk.Stack) awslambda.Function {
	layer := awslambda.NewLayerVersion(stack, jsii.String("MyLayer"), &awslambda.LayerVersionProps{
		Code: awslambda.Code_FromAsset(jsii.String("./lambda/python/googleAI/vertex/deployment_package.zip"), nil),
	})
	return awslambda.NewFunction(stack, jsii.String("PythonLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_9(),
		Handler: jsii.String("index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("./lambda/python/googleAI/vertex"), nil),
		Layers:  &[]awslambda.ILayerVersion{layer},
		Environment: &map[string]*string{
			"VERTEX_JSON": jsii.String(os.Getenv("VERTEX_JSON")),
		},
		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
	})
}
