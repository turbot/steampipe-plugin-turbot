# Table: turbot_grant

A grant is the assignment of a permission to a Turbot user or group on a resource or resource group. 

## Examples

### List all turbot grants

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
