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
	Turbot   TurbotResourceMetadata
	Data     map[string]interface{}
	Metadata map[string]interface{}
	Type     struct {
		URI string
	}
}

type TurbotResourceMetadata struct {
	ID                string
	ParentID          string
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
	State   string
	Reason  string
	Details interface{}
	Type    struct {
		URI string
	}
	Turbot TurbotControlMetadata
}

type TurbotControlMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp string
	UpdateTimestamp string
	ControlTypeID   string
	ResourceID      string
}
