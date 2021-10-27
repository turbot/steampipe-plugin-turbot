### Turbot grants updates in last 1 week

```sql
select
  notification_type,
  create_timestamp,
  actor_trunk_title as granted_by,
  grant_identity_trunk_title as granted_to,
  grant_end_date,
  grant_identity_profile_id,
  resource_title,
  grant_id,
  grant_permission_type,
  grant_permission_level
from
  turbot.turbot_notification
where
  filter = 'notificationType:activeGrant createTimestamp:>T-1w'
  and grant_permission_type = 'Turbot'
order by
  create_timestamp desc,
  notification_type,
  actor_trunk_title,
  resource_title;

```

```sql
select
  notification_type,
  create_timestamp,
  actor_trunk_title as granted_by,
  grant_identity_trunk_title as granted_to,
  grant_end_date,
  grant_identity_profile_id,
  resource_title,
  grant_id,
  grant_permission_type,
  grant_permission_level
from
  turbot.turbot_notification
where
  filter = 'notificationType:active_grants_created createTimestamp:>T-1w'
  and grant_permission_type = 'Turbot'
order by
  create_timestamp desc,
  notification_type,
  actor_trunk_title,
  resource_title;
```

```sql
 select
  notification_type,
  create_timestamp,
  actor_trunk_title as granted_by,
  grant_identity_trunk_title as granted_to,
  grant_end_date,
  grant_identity_profile_id,
  resource_title,
  grant_id,
  grant_permission_type,
  grant_permission_level
from
  turbot.turbot_notification
where
  notification_type = 'active_grants_created' and
  create_timestamp >= (current_date - interval '7' day) and
  grant_permission_type = 'Turbot'
order by
  create_timestamp desc,
  notification_type,
  actor_trunk_title,
  resource_title
```

```sql
 select
  notification_type,
  create_timestamp,
  actor_trunk_title as granted_by,
  grant_identity_trunk_title as granted_to,
  grant_end_date,
  grant_identity_profile_id,
  resource_title,
  grant_id,
  grant_permission_type,
  grant_permission_level
from
  turbot.turbot_notification
where
 id = 238901643750068
```
