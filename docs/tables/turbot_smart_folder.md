# Table: turbot_smart_folder

Smart folders in Turbot allow groups of policies to be applied across a
collection of resources.

## Examples

### List all smart folders

```sql
select
  id,
  title
from
  turbot_smart_folder;
```

### List smart folders with their policy settings

```sql
select
  sf.trunk_title as smart_folder,
  pt.trunk_title as policy,
  ps.id,
  ps.precedence,
  ps.is_calculated,
  ps.value
from
  turbot_smart_folder as sf
  left join turbot_policy_setting as ps on ps.resource_id = sf.id
  left join turbot_policy_type as pt on pt.id = ps.policy_type_id
order by
  smart_folder;
```

### List smart folders with their attached resources

Get each smart folder with an array of the resources attached to it:

```sql
select
  title,
  attached_resource_ids
from
  turbot_smart_folder
order by
  title;
```

Create a row per smart folder and resource:

```sql
select
  sf.title as smart_folder,
  sf_resource_id
from
  turbot_smart_folder as sf,
  jsonb_array_elements(sf.attached_resource_ids) as sf_resource_id
order by
  smart_folder,
  sf_resource_id;
```

Unfortunately, this query to join the smart folder with its resources does not
work yet due to issues with qualifier handling in the Steampipe Postgres FDW:

```sql
select
  sf.title as smart_folder,
  r.trunk_title as resource,
  r.id
from
  turbot_smart_folder as sf
  cross join jsonb_array_elements(sf.attached_resource_ids) as sf_resource_id
  left join turbot_resource as r on r.id = sf_resource_id::bigint;
```
