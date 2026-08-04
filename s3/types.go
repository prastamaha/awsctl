package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var bucketAliases = []string{"bucket", "buckets", "s3"}

type S3 struct {
	AWSConfig aws.Config
}

type BucketList struct {
	Name         string `json:"name" yaml:"name"`
	CreationDate string `json:"creationDate" yaml:"creationDate"`
}

type BucketDescribe struct {
	Name              string            `json:"name" yaml:"name"`
	Tags              map[string]string `json:"tags" yaml:"tags"`
	Versioning        string            `json:"versioning" yaml:"versioning"`
	Encryption        string            `json:"encryption" yaml:"encryption"`
	PublicAccessBlock *PublicAccess     `json:"publicAccessBlock" yaml:"publicAccessBlock"`
}

type PublicAccess struct {
	BlockPublicAcls       bool `json:"blockPublicAcls" yaml:"blockPublicAcls"`
	IgnorePublicAcls      bool `json:"ignorePublicAcls" yaml:"ignorePublicAcls"`
	BlockPublicPolicy     bool `json:"blockPublicPolicy" yaml:"blockPublicPolicy"`
	RestrictPublicBuckets bool `json:"restrictPublicBuckets" yaml:"restrictPublicBuckets"`
}
