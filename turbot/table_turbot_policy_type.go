package turbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotPolicyType(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_policy_type",
		Description: "Policy types define the types of controls known to Turbot.",
		List: &plugin.ListConfig{
			Hydrate: listPolicyType,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "uri",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getPolicyType,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the policy type."},
			{Name: "uri", Type: proto.ColumnType_STRING, Description: "URI of the policy type."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the policy type."},
			{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the policy type."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the policy type."},
			{Name: "targets", Type: proto.ColumnType_JSON, Transform: transform.FromField("Targets").Transform(emptyListIfNil), Description: "URIs of the resource types targeted by this policy type."},
			// Other columns
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas").Transform(emptyListIfNil), Description: "AKA (also known as) identifiers for the policy type."},
			{Name: "category_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Category.Turbot.ID"), Description: "ID of the control category for the policy type."},
			{Name: "category_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Category.URI"), Description: "URI of the control category for the policy type."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the policy type was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "default_template", Type: proto.ColumnType_STRING, Description: "Default template used to calculate template-based policy values. Should be a Jinja based YAML string."},
			// TODO - needs raw JSON? {Name: "default_template_input", Type: proto.ColumnType_JSON, Description: "GraphQL query run and passed to the default_template."},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "Icon of the policy type."},
			// TODO - needs raw JSON? {Name: "input", Type: proto.ColumnType_JSON, Description: "GraphQL query run and passed in to calculated policies."},
			{Name: "mod_uri", Type: proto.ColumnType_STRING, Description: "URI of the mod that contains the policy type."},
			{Name: "parent_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this policy type."},
			{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the policy type."},
			{Name: "read_only", Type: proto.ColumnType_BOOL, Description: "If true user-defined policy settings are blocked from being created."},
			// Large and not very useful - {Name: "resolved_schema", Type: proto.ColumnType_JSON, Description: "JSON schema with fully-resolved URI references, defining the allowed schema for policy values for any targeted resources."},
			{Name: "schema", Type: proto.ColumnType_JSON, Description: "JSON schema defining the allowed schema for policy values for any targeted resources."},
			{Name: "secret", Type: proto.ColumnType_BOOL, Description: "JSON schema defining valid values for the policy type."},
			{Name: "secret_level", Type: proto.ColumnType_STRING, Description: "Secret Level: SECRET, CONFIDENTIAL or NONE."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the policy type was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the policy type."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	queryPolicyTypeList = `
query policyTypeList($filter: [String!], $next_token: String) {
	policyTypes(filter: $filter, paging: $next_token) {
		items {
			category {
				turbot {
					id
				}
				uri
			}
			description
			defaultTemplate
			defaultTemplateInput
			icon
			#input
			modUri
			readOnly
			resolvedSchema
			schema
			secret
			secretLevel
			targets
			title
			trunk {
				title
			}
			turbot {
				akas
				categoryId
				createTimestamp
				id
				parentId
				path
				tags
				title
				updateTimestamp
				versionId
			}
			uri
		}
		paging {
			next
		}
	}
}
`

	queryPolicyTypeGet = `
query policyTypeGet($id: ID!) {
	policyType(id: $id) {
		category {
			turbot {
				id
			}
			uri
		}
		description
		defaultTemplate
		defaultTemplateInput
		icon
		#input
		modUri
		readOnly
		resolvedSchema
		schema
		secret
		secretLevel
		targets
		title
		trunk {
			title
		}
		turbot {
			akas
			categoryId
			createTimestamp
			id
			parentId
			path
			tags
			title
			updateTimestamp
			versionId
		}
		uri
	}
}
`
)

func listPolicyType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_type.listPolicyType", "connection_error", err)
		return nil, err
	}

	filter := "limit:5000"
	nextToken := ""

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 5000 {
			filter = fmt.Sprintf("limit:%s", strconv.Itoa(int(*limit)))
		}
	}

	// Additional filters
	if d.KeyColumnQuals["uri"] != nil {
		filter = filter + fmt.Sprintf(" policyTypeId:'%s' policyTypeLevel:self", d.KeyColumnQuals["uri"].GetStringValue())
	}

	for {
		result := &PolicyTypesResponse{}
		err = conn.DoRequest(queryPolicyTypeList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_policy_type.listPolicyType", "query_error", err)
			return nil, err
		}
		for _, r := range result.PolicyTypes.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				return nil, nil
			}
		}
		if result.PolicyTypes.Paging.Next == "" {
			break
		}
		nextToken = result.PolicyTypes.Paging.Next
	}
	return nil, nil
}

func getPolicyType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_type.getPolicyType", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &PolicyTypeResponse{}
	err = conn.DoRequest(queryPolicyTypeGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_type.getPolicyType", "query_error", err)
		return nil, err
	}
	return result.PolicyType, nil
}
