package turbot

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsS3Bucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "aws_s3_bucket",
		Description:      "AWS S3 Bucket",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("turbot_id"),
		// 	Hydrate:    getS3Bucket,
		// },
		List: &plugin.ListConfig{
			Hydrate: listS3Buckets,
		},
		/// need the regions stuff too... Columns: awsRegionalColumns([]*plugin.Column{
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The user friendly name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The ARN of the AWS S3 Bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(s3BucketArn),
			},
			{
				Name:        "creation_date",
				Description: "The date and tiem when bucket was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "bucket_policy_is_public",
				Description: "The policy status for an Amazon S3 bucket, indicating whether the bucket is public.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("PolicyStatus.IsPublic"),
			},
			{
				Name:        "versioning_enabled",
				Description: "The versioning state of a bucket.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Versioning.Status").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "versioning_mfa_delete",
				Description: "The MFA Delete status of the versioning state.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Versioning.MFADelete").Transform(handleNilString).Transform(transform.ToBool),
			},
			{
				Name:        "block_public_acls",
				Description: "Specifies whether Amazon S3 should block public access control lists (ACLs) for this bucket and objects in this bucket.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PublicAccessBlock.BlockPublicAcls"),
			},
			{
				Name:        "block_public_policy",
				Description: "Specifies whether Amazon S3 should block public bucket policies for this bucket. If TRUE it causes Amazon S3 to reject calls to PUT Bucket policy if the specified bucket policy allows public access.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PublicAccessBlock.BlockPublicPolicy"),
			},
			{
				Name:        "ignore_public_acls",
				Description: "Specifies whether Amazon S3 should ignore public ACLs for this bucket and objects in this bucket. Setting this element to TRUE causes Amazon S3 to ignore all public ACLs on this bucket and objects in this bucket.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PublicAccessBlock.IgnorePublicAcls"),
			},
			{
				Name:        "restrict_public_buckets",
				Description: "Specifies whether Amazon S3 should restrict public bucket policies for this bucket. Setting this element to TRUE restricts access to this bucket to only AWS service principals and authorized users within this account if the bucket has a public policy.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PublicAccessBlock.RestrictPublicBuckets"),
			},
			{
				Name:        "server_side_encryption_configuration",
				Description: "The default encryption configuration for an Amazon S3 bucket.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Encryption.ServerSideEncryptionConfiguration"),
			},
			{
				Name:        "acl",
				Description: "The access control list (ACL) of a bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lifecycle_rules",
				Description: "The lifecycle configuration information of the bucket.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Lifecycle"),
			},
			{
				Name:        "logging",
				Description: "The logging status of a bucket and the permissions users have to view and modify that status.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy",
				Description: "The resource IAM access document for the bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				//Transform:   transform.FromField("Policy").Transform(policyToCanonical),
			},
			{
				Name:        "replication",
				Description: "The replication configuration of a bucket.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to bucket.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			//** remove this later...
			{
				Name:        "raw",
				Description: "remove...",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION
func listS3Buckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listResources(ctx, d, h, "resourceType:tmod:@turbot/aws-s3#/resource/types/bucket")
}

// TRANSFORM
func s3BucketArn(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(map[string]interface{})

	metaData := instance["__metadata"].(map[string]interface{})
	awsMetaData := metaData["aws"].(map[string]interface{})

	return strings.Join([]string{"arn",
		awsMetaData["partition"].(string),
		"s3",
		"",
		"",
		instance["Name"].(string),
	}, ":"), nil

}
