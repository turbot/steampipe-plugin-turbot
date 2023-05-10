# Table: turbot_resource

Resources in Turbot represent cloud configuration items such as users,
networks, servers, etc.

It is recommended that queries to this table should include (usually in the `where` clause) at least one
of these columns: `id`, `resource_type_id`, `resource_type_uri` or `filter`.

## Examples

### List all AWS IAM Roles

```sql
select
  id,
  title,
  create_timestamp,
  metadata,
  data
from
  turbot_resource
where
  resource_type_uri = 'tmod:@turbot/aws-iam#/resource/types/role';
```

### List all S3 buckets with a given Owner tag

```sql
select
  id,
  title,
  tags
from
  turbot_resource
where
  resource_type_uri = 'tmod:@turbot/aws-s3#/resource/types/bucket'
  and tags ->> 'Owner' = 'Jane';
```

### Get a specific resource by ID

```sql
select
  id,
  title,
  create_timestamp,
  metadata,
  data
from
  turbot_resource
where
  id = 216005088871602;
```

### Filter for resources using Turbot filter syntax

```sql
select
  resource_type_uri,
  count(*)
from
  turbot_resource
where
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/iam"'
group by
  resource_type_uri
order by
  count desc;
```

### Search for AWS IAM Roles by name (Turbot side)

This query will ask Turbot to filter the resources down to the given `filter`,
limiting the results by name.

```sql
select
  id,
  title,
  create_timestamp,
  metadata,
  data
from
  turbot_resource
where
  resource_type_uri = 'tmod:@turbot/aws-iam#/resource/types/role'
  and filter = 'admin';
```

### Search for AWS IAM Roles by name (Steampipe side)

This query gathers all the AWS IAM roles from Turbot and then uses Postgres
level filters to limit the results.

```sql
select
  id,
  title,
  create_timestamp,
  metadata,
  data
from
  turbot_resource
where
  resource_type_uri = 'tmod:@turbot/aws-iam#/resource/types/role'
  and title ilike '%admin%';
```

### Search for console logins within 7 days

```sql
select
  id,
  title,
  data ->> 'email' as email,
  array_to_string(regexp_matches(trunk_title, '^Turbot > (.*) >'), ' ' ) as "directory",
  trunk_title,
  to_char((data ->> 'lastLoginTimestamp') :: timestamp, 'YYYY-MM-DD HH24:MI') as "last_login"
from
  turbot_resource
where
  filter = 'resourceTypeId:"tmod:@turbot/turbot-iam#/resource/types/profile" $.lastLoginTimestamp:>=T-7d';
```

### Search for resources created within 7 days, join with count of controls in alarm state

This query gathers all the recently created resources within 1 day, and joins the resource's current count of controls in alarm state.

```sql 
select
  r.id,
  r.title,
  r.trunk_title,
  r.resource_type_uri,
  to_char(r.create_timestamp, 'YYYY-MM-DD HH24:MI') as create_timestamp,
  count(c.*) as alarm_count 
from
  turbot_resource as r 
  left join
    turbot_control as c 
    on r.id = c.resource_id 
    and c.state = 'alarm' 
where
  r.filter = 'notificationType:resource timestamp:>=T-7d' 
group by
  r.id,
  r.title,
  r.trunk_title,
  r.resource_type_uri,
  r.create_timestamp 
order by
  r.create_timestamp desc;
```

### Extract all resources from Turbot

WARNING - This is a large query and may take minutes to run. It is not recommended and may timeout.
It's included here as a reference for those who need to extract all data.

```sql
select
  *
from
  turbot_resource;
```
