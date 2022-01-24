# Table: turbot_active_grant

A active grant is the assignment of a permission to a Turbot user or group on a resource or resource group which is active. 

## Examples

### Basic info

```sql
select
  grant_id,
  identity_status,
  identity_email,
  identity_profile_id,
  identity_trunk_title,
  level_title,
  resource_trunk_title
from
  turbot_active_grant;
```

### List active grants by identity

```sql
select
  grant_id,
  identity_status,
  identity_email,
  identity_trunk_title,
  level_title,
  resource_trunk_title
from
  turbot_active_grant
where
  identity_email = 'abc@gmail.com'
```