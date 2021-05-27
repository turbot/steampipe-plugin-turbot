package turbot

type ResourcesResponse struct {
	Resources struct {
		Items  []Resource
		Paging struct {
			Next string
		}
	}
}

type ResourceResponse struct {
	Resource Resource
}

type Resource struct {
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

type TurbotResourceMetadata struct {
	ID                string
	ParentID          *string
	Akas              []string
	Custom            map[string]interface{}
	Metadata          map[string]interface{}
	Tags              map[string]interface{}
	Title             string
	VersionID         string
	ActorIdentityID   string
	ActorPersonaID    string
	ActorRoleID       string
	ResourceParentAka string
	Timestamp         string
	CreateTimestamp   string
	DeleteTimestamp   *string
	UpdateTimestamp   *string
	Path              string
	ResourceTargetIDs []string
	ResourceGroupIDs  []string
	ResourceTypeID    string
	State             string
	Terraform         map[string]interface{}
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
	}
	Type struct {
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
	Default       bool
	Exception     int
	Input         string
	IsCalculated  bool
	Note          string
	Orphan        int
	Precedence    string
	Template      string
	TemplateInput interface{}
	Type          struct {
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
