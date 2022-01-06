# Table: turbot_resource

Resources in Turbot represent cloud configuration items such as users,
networks, servers, etc.

It is recommended that queries to this table should include (usually in the `where` clause) at least one
of these columns: `id`, `resource_type_id`, `resource_type_uri` or `filter`.

A Policy Value is the effective policy setting on an instance of a resource type. Every resource that is targeted by a given policy setting will have its own value for that policy, which will be the resultant calculated policy for the "winning" policy in the hierarchy.

Policy settings are inherited through the resource hierarchy, and values for a resource are calculated according to policy settings at or above it in the resource hierarchy. For example, a policy setting at the Turbot level will be inherited by all resources below.

It is recommended that queries to this table should include (usually in the `where` clause) at least one
of these columns: `state`,`policyTypeId`, `resource_type_id`, `resource_type_uri` or `filter`.

## Examples

### List policy values by policy type id

```sql
select
  id,
  state,
  default,
  is_calculated,
  policy_type_id,
  policy_value_type_mod_uri
from
  turbot_policy_value
where
  policy_type_id = 123456789;
```

### List policy values by resource id

```sql
select
  id,
  state,
  default,
  is_calculated,
  resource_id,
  policy_value_type_mod_uri
from
  turbot_policy_value
where
  resource_id = 123456789;
```

### List policy values by resource type id

```sql
select
  id,
  state,
  default,
  is_calculated,
  resource_type_id,
  policy_value_type_mod_uri
from
  turbot_policy_value
where
  resource_type_id = 123456789;
```

### Filter for policy values using Turbot filter syntax

```sql
select
  id,
  policy_type_id
  count(*)
from
  turbot_policy_value
where
  filter = 'policyTypeId:123456789'
group by
  policy_type_id
order by
  count desc;
```

### Extract all policy values from Turbot

WARNING - This is a large query and may take minutes to run. It is not recommended and may timeout.
It's included here as a reference for those who need to extract all data.

```sql
select
  *
from
  turbot_policy_value;
```
