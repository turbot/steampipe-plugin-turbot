# Table: turbot_control_type

List all the cloud control types known to Turbot.

## Examples

### List all control types

```sql
select
  id,
  uri,
  trunk_title
from
  turbot_control_type
order by
  trunk_title
```

### List all control types for AWS S3

```sql
select
  id,
  uri,
  trunk_title
from
  turbot_control_type
where
  mod_uri like 'tmod:@turbot/aws-s3%'
order by
  trunk_title
```

### Count control types by cloud provider

```sql
select
  sum(case when mod_uri like 'tmod:@turbot/aws-%' then 1 else 0 end) as aws,
  sum(case when mod_uri like 'tmod:@turbot/azure-%' then 1 else 0 end) as azure,
  sum(case when mod_uri like 'tmod:@turbot/gcp-%' then 1 else 0 end) as gcp,
  count(*) as total
from
  turbot_control_type
```

### Control types that target AWS > S3 > Bucket

```sql
select
  trunk_title,
  uri,
  targets
from
  turbot_control_type
where
  targets ? 'tmod:@turbot/aws-s3#/resource/types/bucket'
```
