# Table: turbot_tag

Tags is a unified collection of all tags discovered by Turbot across all
resources in all clouds.

Queries to this table must specify (usually in the `where` clause) at least one
of these columns: `id`, `key`, `value` or `filter`.

## Examples

### Find all resources for the Sales department

```sql
select
  key,
  value,
  resource_ids
from
  turbot_tag
where
  key = 'Department'
  and value = 'Sales';
```

### Find departments with the most tagged resources

```sql
select
  key,
  value,
  jsonb_array_length(resource_ids) as count
from
  turbot_tag
where
  key = 'Department'
order by
  count desc;
```

### List all tags

```sql
select
  *
from
  turbot_tag
where
  -- At least one qualifier must be given, return all with filter = ''
  filter = ''
order by
  key,
  value;
```
