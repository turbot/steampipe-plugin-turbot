select
  create_timestamp,
  notification_type,
  actor_title,
  resource_title,
  grant_id,
  grant_id_1,
  grant_new_version_id,
  grant_new_version_id_1,
  grant_old_version_id,
  grant_old_version_id_1,
  grant_permission_level,
  grant_permission_level_1,
  grant_permission_level_id,
  grant_permission_level_id_1,
  grant_permission_type,
  grant_permission_type_1,
  grant_permission_type_id,
  grant_permission_type_id_1,
  grant_role_name,
  grant_role_name_1
from
  turbot.turbot_notification
where
  filter = 'notificationType:grant,activeGrant limit:1000'
order by
  create_timestamp,
  notification_type,
  actor_title,
  resource_title

{
  "Icon": "",
  "Message": "",
  "NotificationType": "grant_deleted",
  "Data": null,
  "Actor": {
    "Identity": {
      "Turbot": {
        "Title": "Lalit Bhardwaj",
        "ID": "191382681670329",
        "ActorIdentityID": "191378558966314"
      }
    }
  },
  "Control": {
    "State": "",
    "Reason": "",
    "Details": null,
    "Type": {
      "URI": null,
      "Turbot": {
        "ID": null
      },
      "Trunk": {
        "Title": null
      }
    }
  },
  "Resource": {
    "Data": {
      "title": "Turbot"
    },
    "Metadata": {},
    "Type": {
      "URI": "tmod:@turbot/turbot#/resource/types/turbot",
      "Turbot": {
        "ID": "191378502345748"
      },
      "Trunk": {
        "Title": "Turbot"
      }
    },
    "Trunk": {
      "Title": "Turbot"
    },
    "Turbot": {
      "Akas": [
        "tmod:@turbot/turbot#/"
      ],
      "ParentID": "",
      "Path": "191378502301697",
      "Tags": {},
      "Title": "Turbot"
    }
  },
  "PolicySetting": null,
  "ActiveGrant": {
    "grant": {
      "level": {},
      "type": {}
    }
  },
  "Grant": {
    "roleName": "",
    "permissionTypeId": "191378554894730",
    "permissionLevelId": "191378555110866",
    "level": {
      "title": "Owner"
    },
    "type": {
      "title": "Turbot"
    }
  },
  "Turbot": {
    "ControlID": null,
    "ControlNewVersionID": null,
    "ControlOldVersionID": null,
    "CreateTimestamp": "2021-10-12T15:27:32.642Z",
    "ID": "237590172328106",
    "PolicySettingID": null,
    "PolicySettingNewVersionID": null,
    "PolicySettingOldVersionID": null,
    "ProcessID": null,
    "ResourceID": "191378502301697",
    "ResourceNewVersionID": "191378505865386",
    "ResourceOldVersionID": null,
    "ResourceTypeID": null,
    "Timestamp": "",
    "UpdateTimestamp": null,
    "VersionID": "",
    "GrantID": "213669171873729",
    "GrantNewVersionID": "237590172319913",
    "GrantOldVersionID": "213669171873729",
    "ActiveGrantsID": null,
    "ActiveGrantsNewVersionID": null,
    "ActiveGrantsOldVersionID": null
  }
}


{
  "Icon": "",
  "Message": "",
  "NotificationType": "grant_deleted",
  "Data": null,
  "Actor": {
    "Identity": {
      "Turbot": {
        "Title": "Lalit Bhardwaj",
        "ID": "191382681670329",
        "ActorIdentityID": "191378558966314"
      }
    }
  },
  "Control": {
    "State": "",
    "Reason": "",
    "Details": null,
    "Type": {
      "URI": null,
      "Turbot": {
        "ID": null
      },
      "Trunk": {
        "Title": null
      }
    }
  },
  "Resource": {
    "Data": {
      "title": "Turbot"
    },
    "Metadata": {},
    "Type": {
      "URI": "tmod:@turbot/turbot#/resource/types/turbot",
      "Turbot": {
        "ID": "191378502345748"
      },
      "Trunk": {
        "Title": "Turbot"
      }
    },
    "Trunk": {
      "Title": "Turbot"
    },
    "Turbot": {
      "Akas": [
        "tmod:@turbot/turbot#/"
      ],
      "ParentID": "",
      "Path": "191378502301697",
      "Tags": {},
      "Title": "Turbot"
    }
  },
  "PolicySetting": null,
  "ActiveGrant": {
    "grant": {
      "level": {},
      "type": {}
    }
  },
  "Grant": {
    "roleName": "",
    "permissionTypeId": "191378554894730",
    "permissionLevelId": "191378555110866",
    "level": {
      "title": "Owner"
    },
    "type": {
      "title": "Turbot"
    }
  },
  "Turbot": {
    "ControlID": null,
    "ControlNewVersionID": null,
    "ControlOldVersionID": null,
    "CreateTimestamp": "2021-10-12T15:27:32.642Z",
    "ID": "237590172328106",
    "PolicySettingID": null,
    "PolicySettingNewVersionID": null,
    "PolicySettingOldVersionID": null,
    "ProcessID": null,
    "ResourceID": "191378502301697",
    "ResourceNewVersionID": "191378505865386",
    "ResourceOldVersionID": null,
    "ResourceTypeID": null,
    "Timestamp": "",
    "UpdateTimestamp": null,
    "VersionID": "",
    "grantId": "213669171873729",
    "GrantNewVersionID": "237590172319913",
    "GrantOldVersionID": "213669171873729",
    "ActiveGrantsNewVersionID": null,
    "ActiveGrantsOldVersionID": null
  }
}

{
  "Icon": "",
  "Message": "",
  "NotificationType": "grant_deleted",
  "Data": null,
  "Actor": {
    "Identity": {
      "Turbot": {
        "Title": "Lalit Bhardwaj",
        "ID": "191382681670329",
        "ActorIdentityID": "191378558966314"
      }
    }
  },
  "Control": {
    "State": "",
    "Reason": "",
    "Details": null,
    "Type": {
      "URI": null,
      "Turbot": {
        "ID": null
      },
      "Trunk": {
        "Title": null
      }
    }
  },
  "Resource": {
    "Data": {
      "title": "Turbot"
    },
    "Metadata": {},
    "Type": {
      "URI": "tmod:@turbot/turbot#/resource/types/turbot",
      "Turbot": {
        "ID": "191378502345748"
      },
      "Trunk": {
        "Title": "Turbot"
      }
    },
    "Trunk": {
      "Title": "Turbot"
    },
    "Turbot": {
      "Akas": [
        "tmod:@turbot/turbot#/"
      ],
      "ParentID": "",
      "Path": "191378502301697",
      "Tags": {},
      "Title": "Turbot"
    }
  },
  "PolicySetting": null,
  "ActiveGrant": {
    "grant": {
      "level": {},
      "type": {}
    }
  },
  "Grant": {
    "roleName": "",
    "permissionTypeId": "191378554894730",
    "permissionLevelId": "191378555110866",
    "level": {
      "title": "Owner"
    },
    "type": {
      "title": "Turbot"
    }
  },
  "Turbot": {
    "ControlID": null,
    "ControlNewVersionID": null,
    "ControlOldVersionID": null,
    "CreateTimestamp": "2021-10-12T15:27:32.642Z",
    "ID": "237590172328106",
    "PolicySettingID": null,
    "PolicySettingNewVersionID": null,
    "PolicySettingOldVersionID": null,
    "ProcessID": null,
    "ResourceID": "191378502301697",
    "ResourceNewVersionID": "191378505865386",
    "ResourceOldVersionID": null,
    "ResourceTypeID": null,
    "Timestamp": "",
    "UpdateTimestamp": null,
    "VersionID": "",
    "GrantID": "213669171873729",
    "GrantNewVersionID": "237590172319913",
    "GrantOldVersionID": "213669171873729",
    "ActiveGrantsNewVersionID": null,
    "ActiveGrantsOldVersionID": null
  }
}

{
  "Icon": "",
  "Message": "",
  "NotificationType": "grant_deleted",
  "Data": null,
  "Actor": {
    "Identity": {
      "Turbot": {
        "Title": "Lalit Bhardwaj",
        "ID": "191382681670329",
        "ActorIdentityID": "191378558966314"
      }
    }
  },
  "Control": {
    "State": "",
    "Reason": "",
    "Details": null,
    "Type": {
      "URI": null,
      "Turbot": {
        "ID": null
      },
      "Trunk": {
        "Title": null
      }
    }
  },
  "Resource": {
    "Data": {
      "title": "Turbot"
    },
    "Metadata": {},
    "Type": {
      "URI": "tmod:@turbot/turbot#/resource/types/turbot",
      "Turbot": {
        "ID": "191378502345748"
      },
      "Trunk": {
        "Title": "Turbot"
      }
    },
    "Trunk": {
      "Title": "Turbot"
    },
    "Turbot": {
      "Akas": [
        "tmod:@turbot/turbot#/"
      ],
      "ParentID": "",
      "Path": "191378502301697",
      "Tags": {},
      "Title": "Turbot"
    }
  },
  "PolicySetting": null,
  "ActiveGrant": {
    "grant": {
      "level": {},
      "type": {}
    }
  },
  "Grant": {
    "roleName": "",
    "permissionTypeId": "191378554894730",
    "permissionLevelId": "191378555110866",
    "level": {
      "title": "Owner"
    },
    "type": {
      "title": "Turbot"
    }
  },
  "Turbot": {
    "ControlID": null,
    "ControlNewVersionID": null,
    "ControlOldVersionID": null,
    "CreateTimestamp": "2021-10-12T15:27:32.642Z",
    "ID": "237590172328106",
    "PolicySettingID": null,
    "PolicySettingNewVersionID": null,
    "PolicySettingOldVersionID": null,
    "ProcessID": null,
    "ResourceID": "191378502301697",
    "ResourceNewVersionID": "191378505865386",
    "ResourceOldVersionID": null,
    "ResourceTypeID": null,
    "Timestamp": "",
    "UpdateTimestamp": null,
    "VersionID": "",
    "GrantID": "213669171873729",
    "GrantNewVersionID": "237590172319913",
    "GrantOldVersionID": "213669171873729",
    "ActiveGrantsNewVersionID": null,
    "ActiveGrantsOldVersionID": null
  }
}

select create_timestamp, notification_type, actor_title, resource_title, grant_id, grant_id_1, grant_id_2 from turbot.turbot_notification where filter = 'notificationType:grant,activeGrant limit:1' order by create_timestamp, notification_type, actor_title, resource_title limit 5