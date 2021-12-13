# Table: turbot_policy_type

List all the policy types known to Turbot.

## Examples

### List all policy types

```sql
select
  id,
  uri,
  trunk_title
from
  turbot_policy_type
order by
  trunk_title;
```

### List all policy types, descriptions and available settings

```sql
select
  trunk_title as policy_name,
  description,
  schema ->> 'enum' as available_settings,
  schema ->> 'default' as default_setting,  
  schema ->> 'type' as data_type,
  uri as policy_uri
from
  turbot_policy_type
order by
  trunk_title;
```

### List all policy types for AWS S3

```sql
select
  id,
  uri,
  trunk_title
from
  turbot_policy_type
where
  mod_uri like 'tmod:@turbot/aws-s3%'
order by
  trunk_title;
```

### Count policy types by cloud provider

```sql
select
  sum(case when mod_uri like 'tmod:@turbot/aws-%' then 1 else 0 end) as aws,
  sum(case when mod_uri like 'tmod:@turbot/azure-%' then 1 else 0 end) as azure,
  sum(case when mod_uri like 'tmod:@turbot/gcp-%' then 1 else 0 end) as gcp,
  count(*) as total
from
  turbot_policy_type;
```

### Policy types that target AWS > S3 > Bucket

```sql
select
  trunk_title,
  uri,
  targets
from
  turbot_policy_type
where
  targets ? 'tmod:@turbot/aws-s3#/resource/types/bucket';
```
