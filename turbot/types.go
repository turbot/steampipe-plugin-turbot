package turbot

type ResourcesResponse struct {
	Resources struct {
		Items  []Resource
		Paging struct {
			Next string
		}
	}
}

type TurbotIDObject struct {
	Turbot struct {
		ID string
	}
}

type ResourceResponse struct {
	Resource Resource
}

type Resource struct {
	AttachedResources struct {
		Items []TurbotIDObject
	}
	Data     map[string]interface{}
	Metadata map[string]interface{}
	Trunk    struct {
		Title string
	}
	Turbot TurbotResourceMetadata
	Type   struct {
		URI string
	}
}

type ResourceTypesResponse struct {
	ResourceTypes struct {
		Items  []ResourceType
		Paging struct {
			Next string
		}
	}
}

type ResourceTypeResponse struct {
	ResourceType ResourceType
}

type ResourceType struct {
	Category struct {
		Turbot struct {
			ID string
		}
	}
	CategoryURI string
	Description string
	Icon        string
	ModURI      string
	Title       string
	Trunk       struct {
		Title string
	}
	Turbot TurbotResourceMetadata
	URI    string
}

type ControlTypesResponse struct {
	ControlTypes struct {
		Items  []ControlType
		Paging struct {
			Next string
		}
	}
}

type ControlTypeResponse struct {
	ControlType ControlType
}

type ControlType struct {
	Category struct {
		Turbot struct {
			ID string
		}
		URI string
	}
	Description string
	Icon        string
	ModURI      string
	Targets     []string
	Title       string
	Trunk       struct {
		Title string
	}
	Turbot TurbotResourceMetadata
	URI    string
}

type PolicyTypesResponse struct {
	PolicyTypes struct {
		Items  []PolicyType
		Paging struct {
			Next string
		}
	}
}

type PolicyTypeResponse struct {
	PolicyType PolicyType
}

type PolicyType struct {
	CategoryURI          string
	Description          string
	DefaultTemplate      string
	DefaultTemplateInput interface{}
	Icon                 string
	Input                interface{}
	ModURI               string
	ReadOnly             bool
	ResolvedSchema       interface{}
	Schema               interface{}
	Secret               bool
	SecretLevel          string
	Targets              []string
	Title                string
	Trunk                struct {
		Title string
	}
	Turbot TurbotResourceMetadata
	URI    string
}

type TurbotResourceMetadata struct {
	ActorIdentityID   string
	ActorPersonaID    string
	ActorRoleID       string
	Akas              []string
	CategoryID        string
	CreateTimestamp   string
	Custom            map[string]interface{}
	DeleteTimestamp   *string
	ID                string
	Metadata          map[string]interface{}
	ParentID          *string
	Path              string
	ResourceGroupIDs  []string
	ResourceParentAka string
	ResourceTargetIDs []string
	ResourceTypeID    string
	State             string
	Tags              map[string]interface{}
	Terraform         map[string]interface{}
	Timestamp         *string
	Title             string
	UpdateTimestamp   *string
	VersionID         string
}

type ControlsResponse struct {
	Controls struct {
		Items  []Control
		Paging struct {
			Next string
		}
	}
}

type ControlResponse struct {
	Control Control
}

type Control struct {
	State    string
	Reason   string
	Details  interface{}
	Resource struct {
		Type struct {
			URI string
		}
		Trunk struct {
			Title string
		}
	}
	Type struct {
		Trunk struct {
			Title string
		}
		URI string
	}
	Turbot TurbotControlMetadata
}

type TurbotControlMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp *string
	UpdateTimestamp *string
	ControlTypeID   string
	ResourceID      string
	ResourceTypeID  string
}

type PolicySettingsResponse struct {
	PolicySettings struct {
		Items  []PolicySetting
		Paging struct {
			Next string
		}
	}
}

type PolicySettingResponse struct {
	PolicySetting PolicySetting
}

type PolicySetting struct {
	Default      bool
	Exception    int
	Input        string
	IsCalculated bool
	Note         string
	Orphan       int
	Precedence   string
	Resource     struct {
		Trunk struct {
			Title string
		}
	}
	Template      string
	TemplateInput interface{}
	Type          struct {
		Trunk struct {
			Title string
		}
		URI string
	}
	Turbot             TurbotPolicySettingMetadata
	ValidFromTimestamp *string
	ValidToTimestamp   *string
	Value              interface{}
	ValueSource        interface{}
}

type TurbotPolicySettingMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp *string
	UpdateTimestamp *string
	PolicyTypeID    string
	ResourceID      string
}

type NotificationsResponse struct {
	Notifications struct {
		Items  []Notification
		Paging struct {
			Next string
		}
	}
}

type NotificationResponse struct {
	Notification Notification
}

type Notification struct {
	Icon             string
	Message          string
	NotificationType string
	Data             interface{}

	Control struct {
		State   string
		Reason  string
		Details interface{}
		Type    struct {
			URI    *string
			Turbot struct {
				ID *string
			}
			Trunk struct {
				Title *string
			}
		}
	}

	Resource struct {
		Data     interface{}
		Metadata interface{}
		Type     struct {
			URI    string
			Turbot struct {
				ID string
			}
			Trunk struct {
				Title string
			}
		}
		Trunk struct {
			Title string
		}
		Turbot struct {
			Akas     []string
			ParentID string
			Path     string
			Tags     interface{}
			Title    string
		}
	}

	PolicySetting *struct {
		isCalculated *bool
		Type         struct {
			URI                  *string
			ReadOnly             *bool
			DefaultTemplate      *string
			DefaultTemplateInput interface{}
			Secret               *bool
			Trunk                struct {
				Title *string
			}
			Turbot struct {
				ID string
			}
		}
		Value interface{}
	}

	Turbot TurbotNotificationMetadata
}

type TurbotNotificationMetadata struct {
	ControlID                 *string
	ControlNewVersionID       *string
	ControlOldVersionID       *string
	CreateTimestamp           string
	GrantID                   *string
	GrantNewVersionID         *string
	GrantOldVersionID         *string
	ID                        string
	PolicySettingID           *string
	PolicySettingNewVersionID *string
	PolicySettingOldVersionID *string
	ProcessID                 *string
	ResourceID                *string
	ResourceNewVersionID      *string
	ResourceOldVersionID      *string
	ResourceTypeID            *string
	Timestamp                 string
	UpdateTimestamp           *string
	VersionID                 string
}

type TagsResponse struct {
	Tags struct {
		Items  []Tag
		Paging struct {
			Next string
		}
	}
}

type Tag struct {
	Key       string
	Value     string
	Resources TagResources
	Turbot    TurbotTagMetadata
}

type TagResources struct {
	Items []TagResource
}

type TagResource struct {
	Turbot struct {
		ID string
	}
}

type TurbotTagMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp *string
	UpdateTimestamp *string
}
