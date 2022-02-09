package turbot

import "time"

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

type ModVersionResponse struct {
	ModVersionSearches struct {
		Items  []ModVersion
		Paging struct {
			Next string
		}
	}
}

type ModVersion struct {
	IdentityName string
	Name         string
	Versions     []ModVersionDetail
}

type ModVersionDetail struct {
	Version string
	Status  string
	Head    ModVersionHead
}

type ModVersionHead struct {
	PeerDependencies []PeerDependency
}

type PeerDependency struct {
	FullName     string
	VersionRange string
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

type GrantInfo struct {
	Grants struct {
		Items []Grant
		Paging struct {
			Next string
		}
	}
}

type Grant struct {
	Resource struct {
		Akas  []string
		Title string
		Trunk struct {
			Title string
		}
		Type struct {
			URI   string
			Trunk struct {
				Title string
			}
		}
		Turbot TurbotControlMetadata
	}
	Identity struct {
		Akas               []string
		Email              string
		Status             string
		GivenName          string
		ProfileID          string
		FamilyName         string
		DisplayName        string
		LastLoginTimestamp *time.Time
		Trunk              struct {
			Title string
		}
	}
	Type struct {
		CategoriUri string
		Category    string
		ModUri      string
		Trunk       struct {
			Title string
		}
		URI string
	}
	Level struct {
		Title string
		URI   string
		Trunk struct {
			Title string
		}
	}
	Turbot TurbotControlMetadata
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
	Category struct {
		Turbot struct {
			ID string
		}
		URI string
	}
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

type PolicyValuesResponse struct {
	PolicyValues struct {
		Items  []PolicyValue
		Paging struct {
			Next string
		}
	}
}

type PolicyValue struct {
	Default               bool
	Value                 interface{}
	State                 string
	Reason                string
	Details               interface{}
	SecretValue           interface{}
	IsCalculated          bool
	Precedence            string
	Type                  PolicyValueType
	Resource              PolicyValueResourceDetails
	DependentControls     interface{}
	DependentPolicyValues interface{}
	Turbot                PolicyValueTurbotProperty
}

type PolicyValueResourceDetails struct {
	Trunk struct {
		Title string
	}
}

type PolicyValueType struct {
	ModURI          string
	DefaultTemplate string
	Title           string
	Trunk           struct {
		Title string
	}
}

type PolicyValueTurbotProperty struct {
	TurbotResourceMetadata
	PolicyTypeId string
	ResourceId   string
	SettingId    string
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
	Status          string
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

type NotificationsGetResponse struct {
	Notification Notification
}

type Notification struct {
	Icon             string
	Message          string
	NotificationType string
	Data             interface{}

	Actor struct {
		Identity struct {
			Trunk struct {
				Title *string
			}
			Turbot struct {
				Title           *string
				ID              *string
				ActorIdentityID *string
			}
		}
	}

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

	ActiveGrant struct {
		Grant GrantNotification
	}

	Grant GrantNotification

	Turbot TurbotNotificationMetadata
}

type TurbotNotificationMetadata struct {
	ControlID                 *string
	ControlNewVersionID       *string
	ControlOldVersionID       *string
	CreateTimestamp           string
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
	GrantID                   *string
	GrantNewVersionID         *string
	GrantOldVersionID         *string
	ActiveGrantsID            *string
	ActiveGrantsNewVersionID  *string
	ActiveGrantsOldVersionID  *string
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

type ActiveGrantInfo struct {
	ActiveGrants struct {
		Items  []ActiveGrant
		Paging struct {
			Next string
		}
	}
}

type ActiveGrant struct {
	Resource struct {
		Akas  []string
		Title string
		Trunk struct {
			Title string
		}
		Type struct {
			URI   string
			Trunk struct {
				Title string
			}
		}
		Turbot TurbotControlMetadata
	}
	Grant struct {
		Identity struct {
			Akas               []string
			Email              string
			Status             string
			GivenName          string
			ProfileID          string
			FamilyName         string
			DisplayName        string
			LastLoginTimestamp *time.Time
			Trunk              struct {
				Title string
			}
		}
		Level struct {
			Title string
			URI   string
			Trunk struct {
				Title string
			}
		}
		Turbot TurbotControlMetadata
	}
	Turbot TurbotResourceMetadata
}

type GrantNotification struct {
	RoleName           *string
	PermissionTypeID   *string
	PermissionLevelId  *string
	ValidToTimestamp   *string
	ValidFromTimestamp *string
	Level              struct {
		Title *string
	}
	Type struct {
		Title *string
	}
	Identity struct {
		Trunk struct {
			Title *string
		}
		ProfileID *string
	}
}
