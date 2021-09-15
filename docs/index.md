---
organization: Turbot
category: ["security"]
icon_url: "/images/plugins/turbot/turbot.svg"
brand_color: "#FCC119"
display_name: Turbot
short_name: turbot
description: Steampipe plugin to query resources, controls, policies and more from Turbot.
og_description: Query Turbot with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/turbot-social-graphic.png"
---

# Turbot + Steampipe

[Turbot](https://turbot.com/) is a cloud governance and security platform with a real-time CMDB for cloud resources.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/turbot/tables)**

## Get started

### Install

Download and install the latest Turbot plugin:

```bash
steampipe plugin install turbot
```

### Credentials

Installing the latest turbot plugin will create a config file (`~/.steampipe/config/turbot.spc`) with a single connection named `turbot`. By default, Steampipe will use your [Turbot profiles and credentials](https://turbot.com/v5/docs/reference/cli/installation#setup-your-turbot-credentials) exactly the same as the Turbot CLI and Turbot Terraform provider. In many cases, no extra configuration is required to use Steampipe.

```hcl
connection "turbot" {
  plugin = "turbot"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-turbot
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Advanced configuration options

If you have a `default` profile setup using the Turbot CLI, Steampipe just works with that connection.

For users with multiple workspaces and more complex authentication use cases, here are some examples of advanced configuration options:

### Credentials via key pair

The Turbot plugin allows you set static credentials with the `access_key`, `secret_key`, and `workspace` arguments in any connection profile.

```hcl
connection "turbot" {
  plugin = "turbot"
  workspace  = "https://turbot-acme.cloud.turbot.com/"
  access_key = "c8e2c2ed-1ca8-429b-b369-010e3cf75aac"
  secret_key = "a3d8385d-47f7-40c5-a90c-bfdf5b43c8dd"
}
```

### Credentials via Turbot config profiles

You can use an existing Turbot named profile configured in `/Users/jsmyth/.config/turbot/credentials.yml`. A connect per workspace is a common configuration:

```hcl
connection "turbot_acme" {
  plugin = "turbot"
  profile = "turbot-acme"
}

connection "turbot_dmi" {
  plugin = "turbot"
  profile = "turbot-dmi"
}

```

### Credentials from environment variables

Environment variables provide another way to specify default Turbot CLI credentials:

```sh
export TURBOT_SECRET_KEY=3d397816-575f-4b2a-a470-a96abe29b81a
export TURBOT_ACCESS_KEY=86835f29-1c88-46d9-b6ce-cbe5016842d3
export TURBOT_WORKSPACE=https://turbot-acme.cloud.turbot.com
```

You can also change the default profile to a named profile with the TURBOT_PROFILE environment variable:

```sh
export TURBOT_PROFILE=turbot-acme
```
