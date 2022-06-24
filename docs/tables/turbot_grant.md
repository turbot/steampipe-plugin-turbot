# Table: turbot_grant

A grant is the assignment of a permission to a Turbot user or group on a resource or resource group. 

## Examples

### Basic info

```sql
select
  id,
  identity_status,
  identity_email,
  identity_profile_id,
  identity_trunk_title,
  level_title,
  resource_trunk_title
from
  turbot_grant;
```

### List grants for an identity

```sql
select
  id,
  identity_email,
  identity_family_name,
  level_title,
  level_trunk_title,
from
  turbot_grant
where
  identity_email = 'xyz@gmail.com';
```

### List SuperUser grants

```sql
select
  id,
  identity_email,
  identity_family_name,
  level_title,
  resource_trunk_title
from
  turbot_grant
where
  level_uri  = 'tmod:@turbot/turbot-iam#/permission/levels/superuser';
```

### List grants for inactive identities

```sql
select
  id,
  identity_email,
  identity_status,
  resource_trunk_title
from
  turbot_grant
where
  identity_status = 'Inactive';
```
