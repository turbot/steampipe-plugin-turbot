package turbot

import (
	"context"
	//"fmt"

	//"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsIamRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "aws_iam_role",
		Description:      "AWS IAM Role",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.AnyColumn([]string{"name", "arn"}),
		// 	ShouldIgnoreError: isNotFoundError([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
		// 	// Hydrate:           getIamRole,
		// },
		List: &plugin.ListConfig{
			Hydrate: listAwsIamRoles,
		},

		Columns: awsColumns([]*plugin.Column{

			// "Key" Columns
			{
				Name:        "name",
				Description: "The friendly name that identifies the role.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleName"),
			},
			{
				Name:        "arn",
				Type:        proto.ColumnType_STRING,
				Description: "The Amazon Resource Name (ARN) specifying the role.",
			},
			{
				Name:        "role_id",
				Type:        proto.ColumnType_STRING,
				Description: "The stable and unique string identifying the role.",
			},

			// Other Columns
			{
				Name:        "create_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the role was created.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "A user-provided description of the role.",
			},
			{
				Name:        "instance_profile_arns",
				Description: "A list of instance profiles associated with the role.",
				Type:        proto.ColumnType_JSON,
				// Hydrate:     getAwsIamInstanceProfileData,
				// Transform:   transform.FromValue(),
				// THIS IS MISSING FROM THE TURBOT CMDB!!!!
			},
			{
				Name:        "max_session_duration",
				Description: "The maximum session duration (in seconds) for the specified role. Anyone who uses the AWS CLI, or API to assume the role can specify the duration using the optional DurationSeconds API parameter or duration-seconds CLI parameter.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "path",
				Description: "The path to the role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permissions_boundary_arn",
				Description: "The ARN of the policy used to set the permissions boundary for the role.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PermissionsBoundary.PermissionsBoundaryArn"),
			},
			{
				Name: "permissions_boundary_type",
				Description: "The permissions boundary usage type that indicates what type of IAM resource " +
					"is used as the permissions boundary for an entity. This data type can only have a value of Policy.",
				Type:      proto.ColumnType_STRING,
				Transform: transform.FromField("PermissionsBoundary.PermissionsBoundaryType"),
			},
			{
				Name: "role_last_used_date",
				Type: proto.ColumnType_TIMESTAMP,
				Description: "Contains information about the last time that an IAM role was used. " +
					"Activity is only reported for the trailing 400 days. This period can be " +
					"shorter if your Region began supporting these features within the last year. " +
					"The role might have been used more than 400 days ago.",
				Transform: transform.FromField("RoleLastUsed.LastUsedDate"),
			},
			{
				Name:        "role_last_used_region",
				Description: "Contains the region in which the IAM role was used.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RoleLastUsed.Region"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags that are attached to the role.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "inline_policies",
				Description: "A list of policy documents that are embedded as inline policies for the role..",
				Type:        proto.ColumnType_JSON,
				// Hydrate:     getAwsIamRoleInlinePolicies,
				// We dont store the native format??   Transform: transform.FromField("Policies"),
			},
			{
				Name:        "inline_policies_std",
				Description: "Inline policies in canonical form for the role.",
				Type:        proto.ColumnType_JSON,
				// Hydrate:     getAwsIamRoleInlinePolicies,
				//Transform:   transform.FromValue().Transform(inlinePoliciesToStd),
				Transform: transform.FromField("Policies"),
			},
			{
				Name:        "attached_policy_arns",
				Description: "A list of managed policies attached to the role.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AttachedPolicies").Transform(getAwsIamRoleAttachedPolicyArns),
			},
			{
				Name:        "assume_role_policy",
				Description: "The policy that grants an entity permission to assume the role.",
				Type:        proto.ColumnType_JSON,
				// Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(transform.UnmarshalYAML),
				// We dont store the native format??
			},
			{
				Name:        "assume_role_policy_std",
				Description: "Contains the assume role policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				// Transform:   transform.FromField("AssumeRolePolicyDocument").Transform(unescape).Transform(policyToCanonical),
				Transform: transform.FromField("AssumeRolePolicyDocument"),
			},

			// //** remove this later...
			// {
			// 	Name:        "raw",
			// 	Description: "remove...",
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromValue(),
			// },
		}),
	}
}

func listAwsIamRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listResources(ctx, d, h, "resourceType:tmod:@turbot/aws-iam#/resource/types/role")
}

func getAwsIamRoleAttachedPolicyArns(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var roleArns []string

	items := d.Value.([]interface{})
	for _, item := range items {
		role := item.(map[string]interface{})
		roleArns = append(roleArns, role["PolicyArn"].(string))
	}
	return roleArns, nil
}
