// Code generated by smithy-go-codegen DO NOT EDIT.

package ecs

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Describes the task sets in the specified cluster and service. This is used when
// a service uses the EXTERNAL deployment controller type. For more information,
// see [Amazon ECS Deployment Types]in the Amazon Elastic Container Service Developer Guide.
//
// [Amazon ECS Deployment Types]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-types.html
func (c *Client) DescribeTaskSets(ctx context.Context, params *DescribeTaskSetsInput, optFns ...func(*Options)) (*DescribeTaskSetsOutput, error) {
	if params == nil {
		params = &DescribeTaskSetsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribeTaskSets", params, optFns, c.addOperationDescribeTaskSetsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribeTaskSetsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribeTaskSetsInput struct {

	// The short name or full Amazon Resource Name (ARN) of the cluster that hosts the
	// service that the task sets exist in.
	//
	// This member is required.
	Cluster *string

	// The short name or full Amazon Resource Name (ARN) of the service that the task
	// sets exist in.
	//
	// This member is required.
	Service *string

	// Specifies whether to see the resource tags for the task set. If TAGS is
	// specified, the tags are included in the response. If this field is omitted, tags
	// aren't included in the response.
	Include []types.TaskSetField

	// The ID or full Amazon Resource Name (ARN) of task sets to describe.
	TaskSets []string

	noSmithyDocumentSerde
}

type DescribeTaskSetsOutput struct {

	// Any failures associated with the call.
	Failures []types.Failure

	// The list of task sets described.
	TaskSets []types.TaskSet

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribeTaskSetsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpDescribeTaskSets{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpDescribeTaskSets{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DescribeTaskSets"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpDescribeTaskSetsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribeTaskSets(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opDescribeTaskSets(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DescribeTaskSets",
	}
}
