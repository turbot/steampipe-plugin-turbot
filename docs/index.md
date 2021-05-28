---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/turbot.svg"
brand_color: "#FF9900"
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
  key,
  value,
  resource_ids
from
  turbot_tag
```

```
+------------+-----------+---------------+
| title      | value     | resource_ids  |
+------------+-----------+---------------+
| Department | Sales     | [111,222]     |
| Department | Warehouse | [333,444,555] |
| Owner      | Jim       | [111]         |
| Owner      | Daryl     | [333,555]     |
+------------+-----------+---------------+
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

| Item              | Description                                                                                                                                                                                                                                                                                                                                                    |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials       | Specify a named profile from an Turbot credential file with the `profile` argument.                                                                                                                                                                                                                                                                            |
| Permissions       | Grant the `Turbot/ReadOnly` permission to your user or role.                                                                                                                                                                                                                                                                                                   |
| Radius            | Each connection represents a single Turbot workspace.                                                                                                                                                                                                                                                                                                          |
| Resolution        | 1. Credentials specified in environment variables e.g. `AWS_ACCESS_KEY_ID`.<br />2. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable.<br />3. Credentials for the Default profile from the credential file.<br />4. EC2 Instance Role Credentials (if running on an ec2 instance) |
| Region Resolution | 1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable<br />2. The region specified in the active profile (`AWS_PROFILE` or `default`).                                                                                                                                                                                                              |

### Configuration

Installing the latest turbot plugin will create a config file (`~/.steampipe/config/turbot.spc`) with a single connection named `turbot`:

```hcl
connection "turbot" {
  plugin  = "turbot"
  profile = "default"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-turbot
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Advanced configuration options

If you have a `default` profile setup using the Turbot CLI Steampipe just works with that connection.

For users with multiple accounts and more complex authentication use cases, here are some examples of advanced configuration options:

The Turbot plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in any connection profile. You may also specify one or more regions with the `regions` argument. An AWS connection may connect to multiple regions, however be aware that performance may be negatively affected by both the number of regions and the latency to them.

### Credentials via key pair

```hcl
connection "aws_account_x" {
  plugin      = "aws"
  secret_key  = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key  = "ASIA3ODZSWFYSN2PFHPJ"
  regions     = ["us-east-1" , "us-west-2"]
}
```

### Credentials via AWS config profiles

Named profile from an AWS credential file with the `profile` argument. A connect per profile is a common configuration:

```hcl
# credentials via profile
connection "aws_account_y" {
  plugin      = "aws"
  profile     = "profile_y"
  regions     = ["us-east-1", "us-west-2"]
}

# credentials via profile
connection "aws_account_z" {
  plugin      = "aws"
  profile     = "profile_z"
  regions     = ["ap-southeast-1", "ap-southeast-2"]
}
```

### Credentials from environment variables

```sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=eu-west-1
export AWS_SESSION_TOKEN=AQoDYXdzEJr...
export AWS_ROLE_SESSION_NAME=steampipe@myaccount
```

### Credentials from an EC2 instance role

If you are running Steampipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Steampipe will automatically use the associated IAM role without other credentials:

```hcl
connection "aws" {
  plugin      = "aws"
  regions     = ["eu-west-1", "eu-west-2"]
}
```
