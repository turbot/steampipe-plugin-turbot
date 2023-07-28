# :warning: DEPRECATED

The Turbot plugin has been deprecated as part of our [renaming](https://turbot.com/blog/2023/07/introducing-turbot-guardrails-and-pipes) of Turbot to Turbot Guardrails. Please use the [Turbot Guardrails plugin](https://hub.steampipe.io/plugins/turbot/guardrails) instead.

---
![image](https://hub.steampipe.io/images/plugins/turbot/turbot-social-graphic.png)

# Turbot Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, identity and more from Turbot.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/turbot)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/turbot/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-turbot/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install turbot
```

Run a query:

```sql
select
  trunk_title,
  uri
from
  turbot_resource_type;
```

```
+---------------------------------+---------------------------------------------------------+
| trunk_title                     | uri                                                     |
+---------------------------------+---------------------------------------------------------+
| Turbot > IAM > Access Key       | tmod:@turbot/turbot-iam#/resource/types/accessKey       |
| GCP > Monitoring > Alert Policy | tmod:@turbot/gcp-monitoring#/resource/types/alertPolicy |
| AWS > IAM > Access Key          | tmod:@turbot/aws-iam#/resource/types/accessKey          |
| AWS > EC2 > AMI                 | tmod:@turbot/aws-ec2#/resource/types/ami                |
| AWS > SSM > Association         | tmod:@turbot/aws-ssm#/resource/types/association        |
| GCP > Network > Address         | tmod:@turbot/gcp-network#/resource/types/address        |
+---------------------------------+---------------------------------------------------------+
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-turbot.git
cd steampipe-plugin-turbot
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/turbot.spc
```

Try it!

```shell
steampipe query
> .inspect turbot
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-turbot/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Turbot Plugin](https://github.com/turbot/steampipe-plugin-turbot/labels/help%20wanted)
