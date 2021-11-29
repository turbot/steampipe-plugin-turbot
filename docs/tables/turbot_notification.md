# Table: turbot_notification

Notifications represent significant events in the lifecycle of turbot infrastructure, including:

- A history of change for a resource, e.g., my-s3-bucket.
- A log of state changes and actions performed by a control, e.g., the Tags control for my-s3-bucket.
- Changes to policy settings and policy values updated as a result.
- Records of permission grants, activations, deactivations and revocations.

When querying this table, we recommend using at least one of these columns (usually in the `where` clause):

- `id`
- `resource_id`
- `notification_type`
- `control_id`
- `control_type_id`
- `control_type_uri`
- `resource_type_id`
- `resource_type_uri`
- `policy_type_id`
- `policy_type_uri`
- `actor_identity_id`
- `create_timestamp`
- `filter`

For more information on how to construct a `filter`, please see [Notifications examples](https://turbot.com/v5/docs/reference/filter/notifications#examples).

## Examples

### Find all Turbot grants activations in last 1 week using `filter`

```sql
select
  grant_id,
  notification_type,
  grant_permission_type,
  grant_permission_level,
  create_timestamp,
  actor_trunk_title,
  grant_identity_trunk_title,
  grant_end_date,
  grant_identity_profile_id,
  resource_title
from
  turbot_notification
where
  filter = 'notificationType:activeGrant createTimestamp:>T-1w'
  and grant_permission_type = 'Turbot'
order by
  create_timestamp desc,
  notification_type,
  actor_trunk_title,
  resource_title;
```

### Find all AWS grants activations in last 7 days

```sql
select
  grant_id,
  notification_type,
  grant_permission_type,
  grant_permission_level,
  create_timestamp,
  actor_trunk_title,
  grant_identity_trunk_title,
  grant_end_date,
  grant_identity_profile_id,
  resource_title
from
  turbot_notification
where
  notification_type = 'active_grants_created'
  and create_timestamp >= (current_date - interval '7' day)
  and grant_permission_type = 'AWS'
order by
  create_timestamp desc,
  notification_type,
  actor_trunk_title,
  resource_title;
```

### Find all AWS S3 buckets created in last 7 days

```sql
select
  create_timestamp,
  resource_id,
  resource_title,
  resource_trunk_title,
  actor_trunk_title
from
  turbot_notification
where
  notification_type = 'resource_created'
  and create_timestamp >= (current_date - interval '7' day)
  and resource_type_uri = 'tmod:@turbot/aws-s3#/resource/types/bucket'
order by
  create_timestamp desc;
```

### All policy settings set on a given resource or below

```sql
select
  notification_type,
  create_timestamp,
  policy_id,
  policy_trunk_title,
  policy_type_uri,
  resource_trunk_title,
  resource_type_trunk_title,
  policy_read_only,
  policy_secret,
  policy_value
from
  turbot_notification
where
  resource_id = 191382256916538
  and filter = 'notificationType:policySetting level:self,descendant'
order by
  create_timestamp desc;
```

### All policy settings notifications for AWS > Account > Regions policy

```sql
select
  notification_type,
  create_timestamp,
  policy_id,
  resource_id,
  resource_trunk_title,
  policy_value
from
  turbot_notification
where
  policy_type_uri = 'tmod:@turbot/aws#/policy/types/regionsDefault'
  and filter = 'notificationType:policySetting level:self'
order by
  create_timestamp desc;
```

### All notifications for AWS > Account > Budget > Budget control

```sql
select
  notification_type,
  create_timestamp,
  control_id,
  resource_trunk_title,
  control_state,
  control_reason
from
  turbot_notification
where
  control_type_uri = 'tmod:@turbot/aws#/control/types/budget'
  and filter = 'notificationType:control level:self'
order by
  resource_id,
  create_timestamp desc;
```

### Get count of notifications matching filter

```sql
select
  notifications_count
from
  turbot_notification
where
  filter = 'notificationType:policySetting level:self,descendant'
limit 1;
```
