# Table: turbot_mod_version

Turbot is designed in such a way that it allows organizations to selectively install policies, controls and guardrails that are associated with particular services. This package of Turbot resources is known as a Mod. Turbot published mods are often focused on a specific service in a specific cloud provider.

## Examples

### Version details for aws mod

```sql
select
  name,
  version,
  status,
  workspace
from 
  turbot_mod_version where name = 'aws';
```

### Get recommended mod version for aws-acm

```sql
select
  name,
  version,
  status
from
  turbot_mod_version where name = 'aws-acm' and status = 'RECOMMENDED';
```

### List available mod versions for aws-acm

```sql
select
  name,
  version,
  status
from
  turbot_mod_version where name = 'aws-acm' and status = 'AVAILABLE';
```
