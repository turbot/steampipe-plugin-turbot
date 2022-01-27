# Table: turbot_policy_value

A policy value is the effective policy setting on an instance of a resource type. Every resource that is targeted by a given policy setting will have its own value for that policy, which will be the resultant calculated policy for the "winning" policy in the hierarchy.

Policy settings are inherited through the resource hierarchy, and values for a resource are calculated according to policy settings at or above it in the resource hierarchy. For example, a policy setting at the Turbot level will be inherited by all resources below.

It is recommended that queries to this table should include (usually in the `where` clause) at least one
of these columns: `state`, `policyTypeId`, `resource_type_id`, `resource_type_uri` or `filter`.

## Examples

### List policy values by policy type ID

```sql
select
  id,
  state,
  is_default,
  is_calculated,
  policy_type_id,
  type_mod_uri
from
  turbot_policy_value
where
  policy_type_id = 221505068398189;
```

### List policy values by resource ID

```sql
select
  id,
  state,
  is_default,
  is_calculated,
  resource_id,
  type_mod_uri
from
  turbot_policy_value
where
  resource_id = 161587219904115;
```

### List calculated policy values which are not default

```sql
select
  id,
  state,
  is_default,
  is_calculated,
  resource_type_id,
  type_mod_uri
from
  turbot_policy_value
where
  is_calculated and not default;
```

### Filter for policy values using Turbot filter syntax

```sql
select
  id,
  state,
  is_default,
  is_calculated,
  policy_type_id,
  resource_id,
  resource_type_id
from
  turbot_policy_value
where
  filter = 'state:ok';
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
