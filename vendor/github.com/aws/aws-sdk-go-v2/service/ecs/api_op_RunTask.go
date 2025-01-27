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

// Starts a new task using the specified task definition.
//
// On March 21, 2024, a change was made to resolve the task definition revision
// before authorization. When a task definition revision is not specified,
// authorization will occur using the latest revision of a task definition.
//
// Amazon Elastic Inference (EI) is no longer available to customers.
//
// You can allow Amazon ECS to place tasks for you, or you can customize how
// Amazon ECS places tasks using placement constraints and placement strategies.
// For more information, see [Scheduling Tasks]in the Amazon Elastic Container Service Developer
// Guide.
//
// Alternatively, you can use StartTask to use your own scheduler or place tasks
// manually on specific container instances.
//
// You can attach Amazon EBS volumes to Amazon ECS tasks by configuring the volume
// when creating or updating a service. For more infomation, see [Amazon EBS volumes]in the Amazon
// Elastic Container Service Developer Guide.
//
// The Amazon ECS API follows an eventual consistency model. This is because of
// the distributed nature of the system supporting the API. This means that the
// result of an API command you run that affects your Amazon ECS resources might
// not be immediately visible to all subsequent commands you run. Keep this in mind
// when you carry out an API command that immediately follows a previous API
// command.
//
// To manage eventual consistency, you can do the following:
//
//   - Confirm the state of the resource before you run a command to modify it.
//     Run the DescribeTasks command using an exponential backoff algorithm to ensure
//     that you allow enough time for the previous command to propagate through the
//     system. To do this, run the DescribeTasks command repeatedly, starting with a
//     couple of seconds of wait time and increasing gradually up to five minutes of
//     wait time.
//
//   - Add wait time between subsequent commands, even if the DescribeTasks
//     command returns an accurate response. Apply an exponential backoff algorithm
//     starting with a couple of seconds of wait time, and increase gradually up to
//     about five minutes of wait time.
//
// [Scheduling Tasks]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduling_tasks.html
// [Amazon EBS volumes]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ebs-volumes.html#ebs-volume-types
func (c *Client) RunTask(ctx context.Context, params *RunTaskInput, optFns ...func(*Options)) (*RunTaskOutput, error) {
	if params == nil {
		params = &RunTaskInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "RunTask", params, optFns, c.addOperationRunTaskMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*RunTaskOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type RunTaskInput struct {

	// The family and revision ( family:revision ) or full ARN of the task definition
	// to run. If a revision isn't specified, the latest ACTIVE revision is used.
	//
	// The full ARN value must match the value that you specified as the Resource of
	// the principal's permissions policy.
	//
	// When you specify a task definition, you must either specify a specific
	// revision, or all revisions in the ARN.
	//
	// To specify a specific revision, include the revision number in the ARN. For
	// example, to specify revision 2, use
	// arn:aws:ecs:us-east-1:111122223333:task-definition/TaskFamilyName:2 .
	//
	// To specify all revisions, use the wildcard (*) in the ARN. For example, to
	// specify all revisions, use
	// arn:aws:ecs:us-east-1:111122223333:task-definition/TaskFamilyName:* .
	//
	// For more information, see [Policy Resources for Amazon ECS] in the Amazon Elastic Container Service Developer
	// Guide.
	//
	// [Policy Resources for Amazon ECS]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/security_iam_service-with-iam.html#security_iam_service-with-iam-id-based-policies-resources
	//
	// This member is required.
	TaskDefinition *string

	// The capacity provider strategy to use for the task.
	//
	// If a capacityProviderStrategy is specified, the launchType parameter must be
	// omitted. If no capacityProviderStrategy or launchType is specified, the
	// defaultCapacityProviderStrategy for the cluster is used.
	//
	// When you use cluster auto scaling, you must specify capacityProviderStrategy
	// and not launchType .
	//
	// A capacity provider strategy can contain a maximum of 20 capacity providers.
	CapacityProviderStrategy []types.CapacityProviderStrategyItem

	// An identifier that you provide to ensure the idempotency of the request. It
	// must be unique and is case sensitive. Up to 64 characters are allowed. The valid
	// characters are characters in the range of 33-126, inclusive. For more
	// information, see [Ensuring idempotency].
	//
	// [Ensuring idempotency]: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/ECS_Idempotency.html
	ClientToken *string

	// The short name or full Amazon Resource Name (ARN) of the cluster to run your
	// task on. If you do not specify a cluster, the default cluster is assumed.
	Cluster *string

	// The number of instantiations of the specified task to place on your cluster.
	// You can specify up to 10 tasks for each call.
	Count *int32

	// Specifies whether to use Amazon ECS managed tags for the task. For more
	// information, see [Tagging Your Amazon ECS Resources]in the Amazon Elastic Container Service Developer Guide.
	//
	// [Tagging Your Amazon ECS Resources]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html
	EnableECSManagedTags bool

	// Determines whether to use the execute command functionality for the containers
	// in this task. If true , this enables execute command functionality on all
	// containers in the task.
	//
	// If true , then the task definition must have a task role, or you must provide
	// one as an override.
	EnableExecuteCommand bool

	// The name of the task group to associate with the task. The default value is the
	// family name of the task definition (for example, family:my-family-name ).
	Group *string

	// The infrastructure to run your standalone task on. For more information, see [Amazon ECS launch types]
	// in the Amazon Elastic Container Service Developer Guide.
	//
	// The FARGATE launch type runs your tasks on Fargate On-Demand infrastructure.
	//
	// Fargate Spot infrastructure is available for use but a capacity provider
	// strategy must be used. For more information, see [Fargate capacity providers]in the Amazon ECS Developer
	// Guide.
	//
	// The EC2 launch type runs your tasks on Amazon EC2 instances registered to your
	// cluster.
	//
	// The EXTERNAL launch type runs your tasks on your on-premises server or virtual
	// machine (VM) capacity registered to your cluster.
	//
	// A task can use either a launch type or a capacity provider strategy. If a
	// launchType is specified, the capacityProviderStrategy parameter must be omitted.
	//
	// When you use cluster auto scaling, you must specify capacityProviderStrategy
	// and not launchType .
	//
	// [Amazon ECS launch types]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/launch_types.html
	// [Fargate capacity providers]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-capacity-providers.html
	LaunchType types.LaunchType

	// The network configuration for the task. This parameter is required for task
	// definitions that use the awsvpc network mode to receive their own elastic
	// network interface, and it isn't supported for other network modes. For more
	// information, see [Task networking]in the Amazon Elastic Container Service Developer Guide.
	//
	// [Task networking]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-networking.html
	NetworkConfiguration *types.NetworkConfiguration

	// A list of container overrides in JSON format that specify the name of a
	// container in the specified task definition and the overrides it should receive.
	// You can override the default command for a container (that's specified in the
	// task definition or Docker image) with a command override. You can also override
	// existing environment variables (that are specified in the task definition or
	// Docker image) on a container or add new environment variables to it with an
	// environment override.
	//
	// A total of 8192 characters are allowed for overrides. This limit includes the
	// JSON formatting characters of the override structure.
	Overrides *types.TaskOverride

	// An array of placement constraint objects to use for the task. You can specify
	// up to 10 constraints for each task (including constraints in the task definition
	// and those specified at runtime).
	PlacementConstraints []types.PlacementConstraint

	// The placement strategy objects to use for the task. You can specify a maximum
	// of 5 strategy rules for each task.
	PlacementStrategy []types.PlacementStrategy

	// The platform version the task uses. A platform version is only specified for
	// tasks hosted on Fargate. If one isn't specified, the LATEST platform version is
	// used. For more information, see [Fargate platform versions]in the Amazon Elastic Container Service
	// Developer Guide.
	//
	// [Fargate platform versions]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html
	PlatformVersion *string

	// Specifies whether to propagate the tags from the task definition to the task.
	// If no value is specified, the tags aren't propagated. Tags can only be
	// propagated to the task during task creation. To add tags to a task after task
	// creation, use the[TagResource] API action.
	//
	// An error will be received if you specify the SERVICE option when running a task.
	//
	// [TagResource]: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_TagResource.html
	PropagateTags types.PropagateTags

	// This parameter is only used by Amazon ECS. It is not intended for use by
	// customers.
	ReferenceId *string

	// An optional tag specified when a task is started. For example, if you
	// automatically trigger a task to run a batch process job, you could apply a
	// unique identifier for that job to your task with the startedBy parameter. You
	// can then identify which tasks belong to that job by filtering the results of a [ListTasks]
	// call with the startedBy value. Up to 128 letters (uppercase and lowercase),
	// numbers, hyphens (-), forward slash (/), and underscores (_) are allowed.
	//
	// If a task is started by an Amazon ECS service, then the startedBy parameter
	// contains the deployment ID of the service that starts it.
	//
	// [ListTasks]: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ListTasks.html
	StartedBy *string

	// The metadata that you apply to the task to help you categorize and organize
	// them. Each tag consists of a key and an optional value, both of which you
	// define.
	//
	// The following basic restrictions apply to tags:
	//
	//   - Maximum number of tags per resource - 50
	//
	//   - For each resource, each tag key must be unique, and each tag key can have
	//   only one value.
	//
	//   - Maximum key length - 128 Unicode characters in UTF-8
	//
	//   - Maximum value length - 256 Unicode characters in UTF-8
	//
	//   - If your tagging schema is used across multiple services and resources,
	//   remember that other services may have restrictions on allowed characters.
	//   Generally allowed characters are: letters, numbers, and spaces representable in
	//   UTF-8, and the following characters: + - = . _ : / @.
	//
	//   - Tag keys and values are case-sensitive.
	//
	//   - Do not use aws: , AWS: , or any upper or lowercase combination of such as a
	//   prefix for either keys or values as it is reserved for Amazon Web Services use.
	//   You cannot edit or delete tag keys or values with this prefix. Tags with this
	//   prefix do not count against your tags per resource limit.
	Tags []types.Tag

	// The details of the volume that was configuredAtLaunch . You can configure the
	// size, volumeType, IOPS, throughput, snapshot and encryption in in [TaskManagedEBSVolumeConfiguration]. The name of
	// the volume must match the name from the task definition.
	//
	// [TaskManagedEBSVolumeConfiguration]: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_TaskManagedEBSVolumeConfiguration.html
	VolumeConfigurations []types.TaskVolumeConfiguration

	noSmithyDocumentSerde
}

type RunTaskOutput struct {

	// Any failures associated with the call.
	//
	// For information about how to address failures, see [Service event messages] and [API failure reasons] in the Amazon Elastic
	// Container Service Developer Guide.
	//
	// [API failure reasons]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/api_failures_messages.html
	// [Service event messages]: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-event-messages.html#service-event-messages-list
	Failures []types.Failure

	// A full description of the tasks that were run. The tasks that were successfully
	// placed on your cluster are described here.
	Tasks []types.Task

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationRunTaskMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpRunTask{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpRunTask{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "RunTask"); err != nil {
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
	if err = addIdempotencyToken_opRunTaskMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpRunTaskValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opRunTask(options.Region), middleware.Before); err != nil {
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

type idempotencyToken_initializeOpRunTask struct {
	tokenProvider IdempotencyTokenProvider
}

func (*idempotencyToken_initializeOpRunTask) ID() string {
	return "OperationIdempotencyTokenAutoFill"
}

func (m *idempotencyToken_initializeOpRunTask) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	if m.tokenProvider == nil {
		return next.HandleInitialize(ctx, in)
	}

	input, ok := in.Parameters.(*RunTaskInput)
	if !ok {
		return out, metadata, fmt.Errorf("expected middleware input to be of type *RunTaskInput ")
	}

	if input.ClientToken == nil {
		t, err := m.tokenProvider.GetIdempotencyToken()
		if err != nil {
			return out, metadata, err
		}
		input.ClientToken = &t
	}
	return next.HandleInitialize(ctx, in)
}
func addIdempotencyToken_opRunTaskMiddleware(stack *middleware.Stack, cfg Options) error {
	return stack.Initialize.Add(&idempotencyToken_initializeOpRunTask{tokenProvider: cfg.IdempotencyTokenProvider}, middleware.Before)
}

func newServiceMetadataMiddleware_opRunTask(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "RunTask",
	}
}
