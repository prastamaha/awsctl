package ecs

import "github.com/aws/aws-sdk-go-v2/aws"

var ecsClusterAliases = []string{"ecsc", "ecscluster", "ecsclusters"}
var ecsServiceAliases = []string{"ecss", "ecsservice", "ecsservices", "ecssvc"}
var ecsCronliases = []string{"ecscron", "ecsscheduledtask", "scheduledtask", "cron"}

type ECS struct {
	AWSConfig aws.Config
}

type ECSClusterList struct {
	Name           string `json:"name" yaml:"name"`
	RunningTasks   int32  `json:"runningTasks" yaml:"RunningTasks"`
	PendingTasks   int32  `json:"pendingTasks" yaml:"pendingTasks"`
	ActiveServices int32  `json:"activeServices" yaml:"activeServices"`
	Status         string `json:"status" yaml:"status"`
}

type ECSClusterDescribe struct {
	Name           string            `json:"name" yaml:"name"`
	Arn            string            `json:"arn" yaml:"arn"`
	RunningTasks   int32             `json:"runningTasks" yaml:"RunningTasks"`
	PendingTasks   int32             `json:"pendingTasks" yaml:"pendingTasks"`
	ActiveServices int32             `json:"activeServices" yaml:"activeServices"`
	Status         string            `json:"status" yaml:"status"`
	Tags           map[string]string `json:"tags" yaml:"tags"`
}

type ECSServicesList struct {
	Name           string `json:"name" yaml:"name"`
	Replicas       string `json:"replicas" yaml:"replicas"`
	Status         string `json:"status" yaml:"status"`
	Type           string `json:"type" yaml:"type"`
	TaskDefinition string `json:"taskDefinition" yaml:"taskDefinition"`
	CreatedAt      string `json:"createdAt" yaml:"createdAt"`
}

type ECSServicesDeploymentDescribe struct {
	Id             string `json:"id" yaml:"id"`
	CreatedAt      string `json:"createdAt" yaml:"createdAt"`
	UpdatedAt      string `json:"updatedAt" yaml:"updatedAt"`
	DesiredCount   int32  `json:"desiredCount" yaml:"desiredCount"`
	FailedTasks    int32  `json:"failedTasks" yaml:"failedTasks"`
	RunningCount   int32  `json:"runningCount" yaml:"runningCount"`
	PendingCount   int32  `json:"pendingCount" yaml:"pendingCount"`
	RolloutState   string `json:"rolloutState" yaml:"rolloutState"`
	Status         string `json:"status" yaml:"status"`
	TaskDefinition string `json:"taskDefinition" yaml:"taskDefinition"`
}

type ECSServicesLoadbalancerDescribe struct {
	ContainerName    string `json:"containerName" yaml:"containerName"`
	ContainerPort    int32  `json:"containerPort" yaml:"containerPort"`
	TargetGroupArn   string `json:"targetGroupArn" yaml:"targetGroupArn"`
	LoadbalancerName string `json:"loadbalancerName" yaml:"loadbalancerName"`
}

type ECSServicesEventDescribe struct {
	CreatedAt string `json:"createdAt" yaml:"createdAt"`
	Message   string `json:"message" yaml:"message"`
}

type ECSServiceDescribe struct {
	Name                 string                            `json:"name" yaml:"name"`
	Arn                  string                            `json:"arn" yaml:"arn"`
	DesiredCount         int32                             `json:"desiredCount" yaml:"desiredCount"`
	RunningCount         int32                             `json:"runningCount" yaml:"runningCount"`
	PendingCount         int32                             `json:"pendingCount" yaml:"pendingCount"`
	EnableExecuteCommand bool                              `json:"enableExecuteCommand" yaml:"enableExecuteCommand"`
	CreatedAt            string                            `json:"createdAt" yaml:"createdAt"`
	Status               string                            `json:"status" yaml:"status"`
	TaskDefinition       string                            `json:"taskDefinition" yaml:"taskDefinition"`
	Deployments          []ECSServicesDeploymentDescribe   `json:"deployments" yaml:"deployments"`
	LoadBalancers        []ECSServicesLoadbalancerDescribe `json:"loadbalancers" yaml:"loadbalancers"`
	SecurityGroups       []string                          `json:"securityGroups" yaml:"securityGroups"`
	Subnets              []string                          `json:"subnets" yaml:"subnets"`
	Events               []ECSServicesEventDescribe        `json:"events" yaml:"events"`
}

type ECSCronList struct {
	Name     string `json:"name" yaml:"name"`
	Schedule string `json:"schedule" yaml:"schedule"`
	State    string `json:"state" yaml:"state"`
}

type ECSCronDescribe struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Arn         string `json:"arn" yaml:"arn"`
	Schedule    string `json:"schedule" yaml:"schedule"`
	State       string `json:"state" yaml:"state"`
	RoleArn     string `json:"roleArn" yaml:"roleArn"`
}
