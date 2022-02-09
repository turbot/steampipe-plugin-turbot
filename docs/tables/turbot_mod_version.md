# Table: turbot_mod_version

Turbot mod version table provides essential information pertaining to the different versions of mods(packaged collection of policies, controls, and guardrails that are associated with a particular cloud service) that are available in the registry.

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

### List mod versions using the filter syntax

```sql
select
  name,
  version,
  status
from
  turbot_mod_version where filter = 'aws-x';
```