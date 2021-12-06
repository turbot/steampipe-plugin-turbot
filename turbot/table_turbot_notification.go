package turbot

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotNotification(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_notification",
		Description: "Notifications from the Turbot CMDB.",
		List: &plugin.ListConfig{
			Hydrate: listNotification,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "id", Require: plugin.Optional},
				{Name: "notification_type", Require: plugin.Optional},
				{Name: "control_id", Require: plugin.Optional},
				{Name: "control_type_id", Require: plugin.Optional},
				{Name: "control_type_uri", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "resource_type_id", Require: plugin.Optional},
				{Name: "resource_type_uri", Require: plugin.Optional},
				{Name: "policy_setting_type_id", Require: plugin.Optional},
				{Name: "policy_setting_type_uri", Require: plugin.Optional},
				{Name: "actor_identity_id", Require: plugin.Optional},
				{Name: "create_timestamp", Require: plugin.Optional, Operators: []string{">", ">=", "=", "<", "<="}},
				{Name: "filter", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getNotification,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the notification."},
			{Name: "process_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ProcessID"), Description: "ID of the process that created this notification."},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "Icon for this notification type."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Message for the notification."},
			{Name: "notification_type", Type: proto.ColumnType_STRING, Description: "Type of the notification: resource, action, policySetting, control, grant, activeGrant."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used to search for notifications."},

			// Actor info for the notification
			{Name: "actor_identity_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.Identity.Trunk.Title").NullIfZero(), Description: "Title hierarchy of the actor from the root down to the actor of this event."},
			{Name: "actor_identity_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Actor.Identity.Turbot.ID").NullIfZero(), Description: "Identity ID of the actor that performed this event."},

			// Resource info for notification
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID").NullIfZero(), Description: "ID of the resource for this notification."},
			{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Title of the resource hierarchy from the root down to this resource."},
			{Name: "resource_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Turbot.Title"), Description: "Title of the resource."},
			{Name: "resource_new_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceNewVersionID"), Description: "Version ID of the resource after the event."},
			{Name: "resource_old_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceOldVersionID"), Description: "Version ID of the resource before the event."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Resource.Type.Turbot.ID").NullIfZero(), Description: "ID of the resource type for this notification."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type for this notification."},
			{Name: "resource_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.Trunk.Title"), Description: "Title of the resource type hierarchy from the root down to this resource."},
			{Name: "resource_data", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Data"), Description: "The data for this resource"},
			{Name: "resource_akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Turbot.Akas"), Description: "The globally-unique akas for this resource."},
			{Name: "resource_parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Resource.Turbot.ParentID").NullIfZero(), Description: "The id of the parent resource of this resource."},
			{Name: "resource_path", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Turbot.Path"), Description: "The string of resource ids separated by \".\" from root down to this resource."},
			{Name: "resource_tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resource.Turbot.Tags"), Description: "Tags attached to this resource."},

			// Policy settings notification details
			{Name: "policy_setting_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.PolicySettingID"), Description: "ID of the policy setting for this notification."},
			{Name: "policy_setting_new_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.PolicySettingNewVersionID"), Description: "Version ID of the policy setting after the event."},
			{Name: "policy_setting_old_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.PolicySettingOldVersionID"), Description: "Version ID of the policy setting before the event."},
			{Name: "policy_setting_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("PolicySetting.Type.Turbot.ID").NullIfZero(), Description: "ID of the policy setting type for this notification."},
			{Name: "policy_setting_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("PolicySetting.Type.URI"), Description: "URI of the policy setting type for this notification."},
			{Name: "policy_setting_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("PolicySetting.Type.Trunk.Title"), Description: "This is the title of hierarchy from the root down to this policy type."},
			{Name: "policy_setting_is_calculated", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PolicySetting.isCalculated"), Description: "If true this setting contains calculated inputs e.g. templateInput and template."},
			{Name: "policy_setting_default_template", Type: proto.ColumnType_STRING, Transform: transform.FromField("PolicySetting.DefaultTemplate"), Description: "The Nunjucks template if this setting is for a calculated value."},
			{Name: "policy_setting_default_template_input", Type: proto.ColumnType_STRING, Transform: transform.FromField("PolicySetting.DefaultTemplateInput").Transform(formatPolicyFieldsValue), Description: "The GraphQL Input query if this setting is for a calculated value."},
			{Name: "policy_setting_type_read_only", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PolicySetting.Type.ReadOnly"), Description: "If true user-defined policy settings are blocked from being created."},
			{Name: "policy_setting_type_secret", Type: proto.ColumnType_BOOL, Transform: transform.FromField("PolicySetting.Type.Secret"), Description: "If true policy value will be encrypted."},
			{Name: "policy_setting_value", Type: proto.ColumnType_STRING, Transform: transform.FromField("PolicySetting.Value").Transform(formatPolicyFieldsValue), Description: "The value of the policy setting after this event."},

			// Controls notification details
			{Name: "control_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlID"), Description: "ID of the control for this notification."},
			{Name: "control_new_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlNewVersionID"), Description: "Version ID of the control after the event."},
			{Name: "control_old_version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlOldVersionID"), Description: "Version ID of the control before the event."},
			{Name: "control_state", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.State"), Description: "The current state of the control."},
			{Name: "control_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Reason"), Description: "Optional reason provided at the last state update of this control."},
			{Name: "control_details", Type: proto.ColumnType_JSON, Transform: transform.FromField("Control.Details"), Description: "Optional details provided at the last state update of this control."},
			{Name: "control_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Control.Type.Turbot.ID"), Description: "ID of the control type for this control."},
			{Name: "control_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Type.URI"), Description: "URI of the control type for this control."},
			{Name: "control_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Control.Type.Trunk.Title"), Description: "This is the title of hierarchy from the root down to this control type."},

			// ActiveGrants notification details
			{Name: "active_grant_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.ActiveGrantsID"), Description: "Active grant ID for this notification."},
			{Name: "active_grant_new_version_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.ActiveGrantsNewVersionID"), Description: "Active grant version ID of the grant after the notification."},
			{Name: "active_grant_old_version_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.ActiveGrantsOldVersionID"), Description: "Version ID of the active grant before the event."},
			{Name: "active_grant_valid_to_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: fromField("ActiveGrant.Grant.ValidToTimestamp"), Description: "Optional end date for the active grant to expire."},
			{Name: "active_grant_role_name", Type: proto.ColumnType_STRING, Transform: fromField("ActiveGrant.Grant.RoleName"), Description: "Optional custom roleName for this active grant, when using existing roles rather than Turbot-managed ones."},
			{Name: "active_grant_permission_type_id", Type: proto.ColumnType_INT, Transform: fromField("ActiveGrant.Grant.PermissionTypeID"), Description: "The unique identifier for the active grant permission type."},
			{Name: "active_grant_permission_type_title", Type: proto.ColumnType_STRING, Transform: fromField("ActiveGrant.Grant.Type.Title"), Description: "The name of the active grant permission type."},
			{Name: "active_grant_permission_level_id", Type: proto.ColumnType_INT, Transform: fromField("ActiveGrant.Grant.PermissionLevelId"), Description: "The unique identifier for the active grant permission level."},
			{Name: "active_grant_permission_level_title", Type: proto.ColumnType_STRING, Transform: fromField("ActiveGrant.Grant.Level.Title"), Description: "The name of the active grant permission level."},
			{Name: "active_grant_identity_trunk_title", Type: proto.ColumnType_STRING, Transform: fromField("ActiveGrant.Grant.Identity.Trunk.Title"), Description: "This is the title of hierarchy from the root down to this identity (i.e. Identity whoes access got revoked/permiited) for this active grant."},
			{Name: "active_grant_identity_profile_id", Type: proto.ColumnType_STRING, Transform: fromField("ActiveGrant.Grant.Identity.ProfileID"), Description: "The identity of profile id for this active grant."},

			// Grants notification details
			{Name: "grant_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.GrantID"), Description: "ID of the grant for this notification."},
			{Name: "grant_new_version_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.GrantNewVersionID"), Description: "Version ID of the grant after the event."},
			{Name: "grant_old_version_id", Type: proto.ColumnType_INT, Transform: fromField("Turbot.GrantOldVersionID"), Description: "Version ID of the grant before the event."},
			{Name: "grant_valid_to_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: fromField("Grant.ValidToTimestamp"), Description: "Optional end date for the grant."},
			{Name: "grant_role_name", Type: proto.ColumnType_STRING, Transform: fromField("Grant.RoleName"), Description: "Optional custom roleName for this grant, when using existing roles rather than Turbot-managed ones."},
			{Name: "grant_permission_type_id", Type: proto.ColumnType_INT, Transform: fromField("Grant.PermissionTypeID"), Description: "The unique identifier for the permission type."},
			{Name: "grant_permission_type_title", Type: proto.ColumnType_STRING, Transform: fromField("Grant.Type.Title"), Description: "The name of the permission type."},
			{Name: "grant_permission_level_id", Type: proto.ColumnType_INT, Transform: fromField("Grant.PermissionLevelId"), Description: "The unique identifier for the permission level."},
			{Name: "grant_permission_level", Type: proto.ColumnType_STRING, Transform: fromField("Grant.Level.Title"), Description: "The name of the permission level."},
			{Name: "grant_identity_trunk_title", Type: proto.ColumnType_STRING, Transform: fromField("Grant.Identity.Trunk.Title"), Description: "This is the title of hierarchy from the root down to this identity (i.e. Identity whoes access got revoked/permiited) for this grant."},
			{Name: "grant_identity_profile_id", Type: proto.ColumnType_STRING, Transform: fromField("Grant.Identity.ProfileID"), Description: "The identity profile id for this grant."},
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

					actor {
						identity {
							trunk { title }
							turbot {
								title
								id
								actorIdentityId
							}
						}
					}

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

					policySetting {
						isCalculated
						type {
							uri
							readOnly
							defaultTemplate
							defaultTemplateInput
							secret
							trunk {
								title
							}
							turbot {
								id
							}
						}
						value
					}

					grant {
						roleName
						permissionTypeId
						permissionLevelId
						validToTimestamp
						validFromTimestamp
						level {
							title
						}
						type {
							title
						}
						identity {
							trunk { title }
							profileId: get(path: "profileId")
						}
					}

					activeGrant {
						grant {
							roleName
							permissionTypeId
							permissionLevelId
							validToTimestamp
							validFromTimestamp
							level {
								title
							}
							type {
								title
							}
							identity {
								trunk { title }
								profileId: get(path: "profileId")
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
						grantId
						grantNewVersionId
						grantOldVersionId
						activeGrantsId
						activeGrantsNewVersionId
						activeGrantsOldVersionId
						type
					}

				}
				paging {
					next
				}
			}
		}`

	queryNotificationGet = `
		query notificationGet($id: ID!) {
			notification(id: $id) {
				icon
				message
				notificationType
				data
				actor {
					identity {
						trunk {
							title
						}
						turbot {
							title
							id
							actorIdentityId
						}
					}
				}
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
				policySetting {
					isCalculated
					type {
						uri
						readOnly
						defaultTemplate
						defaultTemplateInput
						secret
						trunk {
							title
						}
						turbot {
							id
						}
					}
					value
				}
				grant {
					roleName
					permissionTypeId
					permissionLevelId
					validToTimestamp
					validFromTimestamp
					level {
						title
					}
					type {
						title
					}
					identity {
						trunk {
							title
						}
						profileId: get(path: "profileId")
					}
				}
				activeGrant {
					grant {
						roleName
						permissionTypeId
						permissionLevelId
						validToTimestamp
						validFromTimestamp
						level {
							title
						}
						type {
							title
						}
						identity {
							trunk {
								title
							}
							profileId: get(path: "profileId")
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
					grantId
					grantNewVersionId
					grantOldVersionId
					activeGrantsId
					activeGrantsNewVersionId
					activeGrantsOldVersionId
					type
				}
			}
		}`
)

func listNotification(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_notification.listNotification", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals
	allQuals := d.Quals
	filter := ""
	if quals["filter"] != nil {
		filter = quals["filter"].GetStringValue()
		filters = append(filters, filter)
	}
	if quals["id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%s", getQualListValues(ctx, quals, "id", "int64")))
	}

	if quals["notification_type"] != nil {
		filters = append(filters, fmt.Sprintf("notificationType:%s", getQualListValues(ctx, quals, "notification_type", "string")))
	}

	if quals["actor_identity_id"] != nil {
		filters = append(filters, fmt.Sprintf("actorIdentityId:%s", getQualListValues(ctx, quals, "actor_identity_id", "int64")))
	}

	if quals["resource_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceId:%s", getQualListValues(ctx, quals, "resource_id", "int64")))
	}

	if quals["resource_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%s resourceTypeLevel:self", getQualListValues(ctx, quals, "resource_type_id", "int64")))
	}

	if quals["resource_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%s resourceTypeLevel:self", getQualListValues(ctx, quals, "resource_type_uri", "string")))
	}

	if quals["control_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:%s controlTypeLevel:self", getQualListValues(ctx, quals, "control_type_id", "int64")))
	}

	if quals["control_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:%s controlTypeLevel:self", getQualListValues(ctx, quals, "control_type_uri", "string")))
	}

	if quals["policy_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%s policyTypeLevel:self", getQualListValues(ctx, quals, "policy_type_id", "int64")))
	}

	if quals["policy_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%s policyTypeLevel:self", getQualListValues(ctx, quals, "policy_type_uri", "string")))
	}

	if allQuals["create_timestamp"] != nil {
		for _, q := range allQuals["create_timestamp"].Quals {
			// Subtracted 1 minute to FilterFrom time and Added 1 minute to FilterTo time to miss any results due to time conersions in steampipe
			switch q.Operator {
			case "=":
				filters = append(filters, fmt.Sprintf("createTimestamp:'%s'", q.Value.GetTimestampValue().AsTime().Format(filterTimeFormat)))
			case ">=", ">":
				filters = append(filters, fmt.Sprintf("createTimestamp:>='%s'", q.Value.GetTimestampValue().AsTime().Add(-1*time.Minute).Format(filterTimeFormat)))
			case "<", "<=":
				filters = append(filters, fmt.Sprintf("createTimestamp:<='%s'", q.Value.GetTimestampValue().AsTime().Add(1*time.Minute).Format(filterTimeFormat)))
			}
		}
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
		var pageLimit int64 = 5000

		// Adjust page limit, if less than default value
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit < pageLimit {
				pageLimit = *limit
			}
		}
		filters = append(filters, fmt.Sprintf("limit:%s", strconv.Itoa(int(pageLimit))))
	}

	plugin.Logger(ctx).Warn("turbot_notification.listNotification", "filters", filters)

	nextToken := ""
	for {
		result := &NotificationsResponse{}
		err = conn.DoRequest(queryNotificationList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_notification.listNotification", "query_error", err)
			// Not returning for function in case of errors because of resources/policies/controls referred might be deleted and
			// graphql queries may fail to retrieve few properties for such items
			// return nil, err
		}
		for _, r := range result.Notifications.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !pageResults || result.Notifications.Paging.Next == "" {
			break
		}
		nextToken = result.Notifications.Paging.Next
	}

	return nil, nil
}

func getNotification(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_notification.getNotification", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetInt64Value()
	result := &NotificationsGetResponse{}
	err = conn.DoRequest(queryNotificationGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_notification.getNotification", "query_error", err)
		return nil, err
	}
	return result.Notification, nil
}

//// TRANFORM FUNCTION

// formatPolicyValue:: Policy value can be a string, hcl or a json.
// It will transform the raw value from api into a string if a hcl or json
func formatPolicyFieldsValue(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var item = d.HydrateItem.(Notification)
	columnName := d.ColumnName
	var value interface{}

	if item.PolicySetting != nil {
		if columnName == "policy_template_input" {
			value = item.PolicySetting.Type.DefaultTemplateInput
		} else {
			value = item.PolicySetting.Value
		}
	}

	if value != nil {
		switch val := value.(type) {
		case string:
			return val, nil
		case []string, map[string]interface{}, interface{}:
			data, err := json.Marshal(val)
			if err != nil {
				return nil, err
			}
			return string(data), nil
		}
	}

	return nil, nil
}

// fromField:: generates a value by retrieving a field or a set of fields from the source item
func fromField(fieldNames ...string) *transform.ColumnTransforms {
	var fieldNameArray []string
	fieldNameArray = append(fieldNameArray, fieldNames...)
	return &transform.ColumnTransforms{Transforms: []*transform.TransformCall{{Transform: fieldValue, Param: fieldNameArray}}}
}

// fieldValue function is intended for the start of a transform chain.
// This returns a field value of either the hydrate call result (if present)  or the root item if not
// the field name is in the 'Param'
func fieldValue(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var item = d.HydrateItem
	var fieldNames []string

	switch p := d.Param.(type) {
	case []string:
		fieldNames = p
	case string:
		fieldNames = []string{p}
	default:
		return nil, fmt.Errorf("'FieldValue' requires one or more string parameters containing property path but received %v", d.Param)
	}

	for _, propertyPath := range fieldNames {
		fieldValue, ok := helpers.GetNestedFieldValueFromInterface(item, propertyPath)
		if ok && !helpers.IsNil(fieldValue) {
			return fieldValue, nil

		}

	}
	return nil, nil
}
