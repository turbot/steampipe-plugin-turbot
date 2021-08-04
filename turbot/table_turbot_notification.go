package turbot

import (
	"context"
	"fmt"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotNotification(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_notification",
		Description: "Notifications from the Turbot CMDB.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "notification_type", "control_id", "control_type_id", "control_type_uri", "resource_id", "resource_type_id", "resource_type_uri", "policy_type_id", "policy_type_uri", "filter"}),
			Hydrate:    listNotification,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the resource."},
			{Name: "process_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ProcessID"), Description: "ID of the process that created this notification."},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "Icon for this notification type."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Message for the notification."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Notification data."},
			{Name: "notification_type", Type: proto.ColumnType_STRING, Description: "Type of the notification: resource, action, policySetting, control, grant, activeGrant."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this resource list."},

			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource for this notification."},
			{Name: "resource_new_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceNewVersionID"), Description: "Version ID of the resource after the event."},
			{Name: "resource_old_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceOldVersionID"), Description: "Version ID of the resource before the event."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Resource.Type.Turbot.ID"), Description: "ID of the resource type for this notification."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type for this notification."},
			{Name: "resource_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.Trunk.Title"), Description: ""},
			{Name: "resource_data", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Data"), Description: ""},
			{Name: "resource_akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Turbot.Akas"), Description: ""},
			{Name: "resource_parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Resource.Turbot.ParentID"), Description: ""},
			{Name: "resource_path", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Turbot.Path"), Description: ""},
			{Name: "resource_tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Turbot.Tags"), Description: ""},
			{Name: "resource_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Turbot.Title"), Description: ""},

			{Name: "control_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlID"), Description: "ID of the control for this notification."},
			{Name: "control_new_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlNewVersionID"), Description: "Version ID of the control after the event."},
			{Name: "control_old_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlOldVersionID"), Description: "Version ID of the control before the event."},
			{Name: "control_state", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.State"), Description: ""},
			{Name: "control_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Resource"), Description: ""},
			{Name: "control_details", Type: proto.ColumnType_JSON, Transform: transform.FromField("Control.Details"), Description: ""},
			{Name: "control_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Control.Type.Turbot.ID"), Description: "ID of the control type for this control."},
			{Name: "control_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Type.URI"), Description: "URI of the control type for this control."},
			{Name: "control_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Type.Trunk.Title"), Description: ""},

			//{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the resource was last modified (created, updated or deleted)."},
			//{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the resource was last updated in Turbot."},
			/*
				{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Title"), Description: "Title of the resource."},
				{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the resource."},
				{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Tags"), Description: "Tags for the resource."},
				{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the resource."},
				// Other columns
				{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Notification custom metadata."},
				{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this resource. For the Turbot root resource this is null."},
				{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the resource."},
				{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the resource."},
			*/
		},
	}
}

const (
	queryNotificationList = `
query notificationList($filter: [String!], $next_token: String) {
	notifications(filter: $filter, paging: $next_token) {
		items {

			icon
			message
			notificationType
			data

			control {
				state
				reason
				details
				type {
					uri
					trunk {
						title
					}
					turbot {
						id
					}
				}
			}

			resource {
				data
				metadata
				trunk {
					title
				}
				turbot {
					akas
					parentId
					path
					tags
					title
				}
				type {
					uri
					trunk {
						title
					}
					turbot {
						id
					}
				}
			}

			turbot {
				controlId
				controlNewVersionId
				controlOldVersionId
				createTimestamp
				grantId
				grantNewVersionId
				grantOldVersionId
				id
				policySettingId
				policySettingNewVersionId
				policySettingOldVersionId
				processId
				resourceId
				resourceNewVersionId
				resourceOldVersionId
				type
			}

		}
		paging {
			next
		}
	}
}
`

	queryNotificationGet = `
query resourceGet($id: ID!) {
	resource(id: $id) {
		data
		metadata
		trunk {
			title
		}
		turbot {
			id
			title
			tags
			akas
			timestamp
			createTimestamp
			updateTimestamp
			versionId
			parentId
			path
			resourceTypeId
		}
		type {
			uri
		}
	}
}
`
)

func listNotification(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.listNotification", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals
	filter := ""
	if quals["filter"] != nil {
		filter = quals["filter"].GetStringValue()
		filters = append(filters, filter)
	}
	if quals["id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%d", quals["id"].GetInt64Value()))
	}
	if quals["notification_type"] != nil {
		filters = append(filters, fmt.Sprintf("notificationType:'%s'", escapeQualString(ctx, quals, "notification_type")))
		//filters = append(filters, fmt.Sprintf("notificationType:'"+quals["notification_type"].GetStringValue()+"'"))
	}
	if quals["resource_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceId:%d", quals["resource_id"].GetInt64Value()))
	}
	if quals["resource_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%d resourceTypeLevel:self", quals["resource_type_id"].GetInt64Value()))
	}
	if quals["resource_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:'%s' resourceTypeLevel:self", escapeQualString(ctx, quals, "resource_type_uri")))
	}
	if quals["control_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:%d controlTypeLevel:self", quals["control_type_id"].GetInt64Value()))
	}
	if quals["control_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:'%s' controlTypeLevel:self", escapeQualString(ctx, quals, "control_type_uri")))
	}
	if quals["policy_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%d policyTypeLevel:self", quals["policy_type_id"].GetInt64Value()))
	}
	if quals["policy_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:'%s' policyTypeLevel:self", escapeQualString(ctx, quals, "policy_type_uri")))
	}

	// Default to a very large page size. Page sizes earlier in the filter string
	// win, so this is only used as a fallback.
	pageResults := false
	// Add a limit if they haven't given one in the filter field
	re := regexp.MustCompile(`(^|\s)limit:[0-9]+($|\s)`)
	if !re.MatchString(filter) {
		// The caller did not specify a limit, so set a high limit and page all
		// results.
		pageResults = true
		filters = append(filters, "limit:5000")
	}

	plugin.Logger(ctx).Warn("turbot_resource.listNotification", "filters", filters)

	nextToken := ""
	for {
		result := &NotificationsResponse{}
		err = conn.DoRequest(queryNotificationList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource.listNotification", "query_error", err)
			//return nil, err
		}
		for _, r := range result.Notifications.Items {
			d.StreamListItem(ctx, r)
		}
		if !pageResults || result.Notifications.Paging.Next == "" {
			break
		}
		nextToken = result.Notifications.Paging.Next
	}

	return nil, nil
}
