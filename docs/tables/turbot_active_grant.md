# Table: turbot_active_grant

An active grant is the assignment of a permission to a Turbot user or group on a resource or resource group which is active.  

The `turbot_active_grant` table will only return active grants.  Use the `turbot_grant` table to get a list of all grants.

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

### List active grants for an identity

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

### List active grants for inactive identities

```sql
select
  grant_id,
  identity_status,
  identity_email,
  level_title,
  resource_trunk_title
from
  turbot_active_grant
where
  identity_status = 'Inactive'
```

### List inactive grants

```sql
select 
  grant_id, 
  identity_email,
  level_title,
  level_trunk_title,
  level_uri,
  resource_trunk_title,
  resource_type_trunk_title
from 
  turbot_grant 
where grant_id not in (select grant_id from turbot_active_grant);
```
