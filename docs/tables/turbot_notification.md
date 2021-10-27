# Table: turbot_notification

Notifications represent significant events in the lifecycle of turbot infrastructure, including:

- History of change for a resource (e.g. my-bucket).
- A log of state changes and actions performed by a control (e.g. my-bucket Tags).
- Changes to policy settings, and the specific policy values they update.
- Records of permission grants, activations, deactivations and revocations.

Queries to this table must specify (usually in the `where` clause) at least one
of these columns: `id`, `resource_id`, `notification_type`, `control_id`, `control_type_id`,
`control_type_uri`, `resource_type_id`, `resource_type_uri`, `policy_type_id`, `policy_type_uri`, `actor_identity_id`, `create_timestamp` or `filter`.

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
