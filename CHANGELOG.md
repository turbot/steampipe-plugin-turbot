## v0.3.0 [2022-01-27]

_What's new?_

- New tables added
  - [turbot_active_grant](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_active_grant) ([#24](https://github.com/turbot/steampipe-plugin-turbot/pull/24))
  - [turbot_grant](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_grant) ([#15](https://github.com/turbot/steampipe-plugin-turbot/pull/15))
  - [turbot_policy_value](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_policy_value) ([#31](https://github.com/turbot/steampipe-plugin-turbot/pull/31))

_Enhancements_

- Added an example to `turbot_tag` document to find tags with empty values ([#21](https://github.com/turbot/steampipe-plugin-turbot/pull/21))

## v0.2.0 [2021-12-13]

_What's new?_

- New tables added
  - [turbot_notification](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_notification) ([#9](https://github.com/turbot/steampipe-plugin-turbot/pull/9))

## v0.1.0 [2021-11-26]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) and Go version 1.17 ([#17](https://github.com/turbot/steampipe-plugin-turbot/pull/17))
- Added additional optional key quals, filter support and context cancellation handling across all the tables ([#5](https://github.com/turbot/steampipe-plugin-turbot/pull/5))
- Added `workspace` column across all the tables to identify Turbot workspace ([#5](https://github.com/turbot/steampipe-plugin-turbot/pull/5))

## v0.0.3 [2021-09-22]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v161--2021-09-21) ([#7](https://github.com/turbot/steampipe-plugin-turbot/pull/7))
- `resource_type_*` columns of `turbot_resource` table should now limit on the exact resource type

## v0.0.2 [2021-05-27]

_Bug fixes_

- Tidy up example on Overview page

## v0.0.1 [2021-05-27]

_What's new?_

- New tables added
  - [turbot_control](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_control)
  - [turbot_control_type](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_control_type)
  - [turbot_policy_setting](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_policy_setting)
  - [turbot_policy_type](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_policy_type)
  - [turbot_resource](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_resource)
  - [turbot_resource_type](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_resource_type)
  - [turbot_smart_folder](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_smart_folder)
  - [turbot_tag](https://hub.steampipe.io/plugins/turbot/turbot/tables/turbot_tag)
