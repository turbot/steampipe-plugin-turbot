package turbot

import (
	"context"
	//"fmt"
	"strconv"
	//"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotAwsEc2Instance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "turbot_aws_ec2_instance",
		Description:      "AWS EC2 Instance, from a Turbot Workspace",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("instance_id"),
			//ShouldIgnoreError: isNotFoundError([]string{"InvalidInstanceID.NotFound", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			//ItemFromKey:       instanceFromKey,
			Hydrate: getEc2InstanceResource,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2InstanceResources,
		},
		/// need the regions stuff too... Columns: awsRegionalColumns([]*plugin.Column{
		Columns: []*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The ID of the instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_type",
				Description: "The instance type",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_state",
				Description: "The state of the instance (pending | running | shutting-down | terminated | stopping | stopped)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("State.Name"),
			},
			{
				Name:        "monitoring_state",
				Description: "Indicates whether detailed monitoring is enabled (disabled | enabled)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Monitoring.State"),
			},
			{
				Name:        "disable_api_termination",
				Default:     false,
				Description: "If the value is true, instance can't be terminated through the Amazon EC2 console, CLI, or API",
				Type:        proto.ColumnType_BOOL,
				//Hydrate:     getInstanceDisableAPITerminationData,
				Transform: transform.FromField("DisableApiTermination.Value"),
			},
			{
				Name:        "cpu_options_core_count",
				Description: "The number of CPU cores for the instance",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuOptions.CoreCount"),
			},
			{
				Name:        "cpu_options_threads_per_core",
				Description: "The number of threads per CPU core",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuOptions.ThreadsPerCore"),
			},
			{
				Name:        "ebs_optimized",
				Description: "Indicates whether the instance is optimized for Amazon EBS I/O. This optimization provides dedicated throughput to Amazon EBS and an optimized configuration stack to provide optimal I/O performance. This optimization isn't available with all instance types",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "hypervisor",
				Description: "The hypervisor type of the instance. The value xen is used for both Xen and Nitro hypervisors",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "iam_instance_profile_arn",
				Description: "The Amazon Resource Name (ARN) of IAM instance profile associated with the instance, if applicable",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamInstanceProfile.Arn"),
			},
			{
				Name:        "iam_instance_profile_id",
				Description: "The ID of the instance profile associated with the instance, if applicable",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IamInstanceProfile.Id"),
			},
			{
				Name:        "image_id",
				Description: "The ID of the AMI used to launch the instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_initiated_shutdown_behavior",
				Description: "Indicates whether an instance stops or terminates when you initiate shutdown from the instance (using the operating system command for system shutdown)",
				Type:        proto.ColumnType_STRING,
				//Hydrate:     getInstanceInitiatedShutdownBehavior,
				Transform: transform.FromField("InstanceInitiatedShutdownBehavior.Value"),
			},
			{
				Name:        "kernel_id",
				Description: "The kernel ID",
				Type:        proto.ColumnType_STRING,
				//Hydrate:     getInstanceKernelID,
				Transform: transform.FromField("KernelId.Value"),
			},
			{
				Name:        "key_name",
				Description: "The name of the key pair, if this instance was launched with an associated key pair",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "outpost_arn",
				Description: "The Amazon Resource Name (ARN) of the Outpost, if applicable",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "placement_availability_zone",
				Description: "The Availability Zone of the instance",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.AvailabilityZone"),
			},
			{
				Name:        "placement_group_name",
				Description: "The name of the placement group the instance is in",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.GroupName"),
			},
			{
				Name:        "placement_tenancy",
				Description: "The tenancy of the instance (if the instance is running in a VPC). An instance with a tenancy of dedicated runs on single-tenant hardware",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Placement.Tenancy"),
			},
			{
				Name:        "private_ip_address",
				Description: "The private IPv4 address assigned to the instance",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "private_dns_name",
				Description: "The private DNS hostname name assigned to the instance. This DNS hostname can only be used inside the Amazon EC2 network. This name is not available until the instance enters the running state",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_dns_name",
				Description: "The public DNS name assigned to the instance. This name is not available until the instance enters the running state",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_address",
				Description: "The public IPv4 address, or the Carrier IP address assigned to the instance, if applicable",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ram_disk_id",
				Description: "The RAM disk ID",
				Type:        proto.ColumnType_STRING,
				//Hydrate:     getInstanceRAMDiskID,
				Transform: transform.FromField("RamdiskId.Value"),
			},
			{
				Name:        "root_device_name",
				Description: "The device name of the root device volume (for example, /dev/sda1)",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "root_device_type",
				Description: "The root device type used by the AMI. The AMI can use an EBS volume or an instance store volume",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_dest_check",
				Description: "Specifies whether to enable an instance launched in a VPC to perform NAT. This controls whether source/destination checking is enabled on the instance",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "sriov_net_support",
				Description: "Indicates whether enhanced networking with the Intel 82599 Virtual Function interface is enabled",
				Type:        proto.ColumnType_STRING,
				//Hydrate:     getInstanceSriovNetSupport,
				Transform: transform.FromField("SriovNetSupport.Value"),
			},
			{
				Name:        "state_code",
				Description: "The reason code for the state change",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("State.Code"),
			},
			{
				Name:        "subnet_id",
				Description: "The ID of the subnet in which the instance is running",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_data",
				Description: "The user data of the instance",
				Type:        proto.ColumnType_STRING,
				//Hydrate:     getInstanceUserData,
				Transform: transform.FromField("UserData.Value"),
			},
			{
				Name:        "virtualization_type",
				Description: "The virtualization type of the instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "The ID of the VPC in which the instance is running",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "elastic_gpu_associations",
				Description: "The Elastic GPU associated with the instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "elastic_inference_accelerator_associations",
				Description: "The elastic inference accelerator associated with the instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "block_device_mappings",
				Description: "Block device mapping entries for the instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_interfaces",
				Description: "The network interfaces for the instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product_codes",
				Description: "The product codes attached to this instance, if applicable",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_groups",
				Description: "The security groups for the instance",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_status",
				Description: "The status of an instance. Instance status includes scheduled events, status checks and instance state information",
				Type:        proto.ColumnType_JSON,
				//Hydrate:     getInstanceStatus,
				Transform: transform.FromField("InstanceStatuses[0]"),
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the instance",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},

			/// Standard columns
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.From(getEc2InstanceTurbotTags),
			// },
			// {
			// 	Name:        "title",
			// 	Description: resourceInterfaceDescription("title"),
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.From(getEc2InstanceTurbotTitle),
			// },
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     getAwsEc2InstanceTurbotData,
			// 	Transform:   transform.FromValue(),
			// },
		},
	}
}

// func tableTurbotResource(ctx context.Context) *plugin.Table {
// 	return &plugin.Table{
// 		Name:        "turbot_resource",
// 		Description: "TODO",
// 		List: &plugin.ListConfig{
// 			KeyColumns: plugin.SingleColumn("filter"),
// 			Hydrate:    listResource,
// 		},
// 		Get: &plugin.GetConfig{
// 			KeyColumns: plugin.SingleColumn("id"),
// 			Hydrate:    getResource,
// 		},
// 		Columns: []*plugin.Column{
// 			// Top columns
// 			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the resource."},
// 			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Title"), Description: "Title of the resource."},
// 			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Tags"), Description: "Tags for the resource."},
// 			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the resource."},
// 			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the resource was last modified (created, updated or deleted)."},
// 			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource was first discovered by Turbot. (It may have been created earlier.)"},
// 			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the resource was last updated in Turbot."},
// 			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the resource."},
// 			{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this resource. For the Turbot root resource this is null."},
// 			{Name: "path", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Path"), Description: "Hierarchy path with all identifiers of ancestors of the resource."},
// 			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this resource."},
// 			//{Name: "delete_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.DeleteTimestamp"), Description: "When the resource was deleted from Turbot."},
// 			{Name: "data", Type: proto.ColumnType_JSON, Description: "Resource data."},
// 			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Resource custom metadata."},
// 			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the resource type for this resource."},
// 			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this resource list."},
// 		},
// 	}
// }

// const (
// 	queryResourceList = `
// query resourceList($filter: [String!], $next_token: String) {
// 	resources(filter: $filter, paging: $next_token) {
// 		items {
// 			data
// 			metadata
// 			type {
// 				uri
// 			}
// 			turbot {
// 				id
// 				title
// 				tags
// 				akas
// 				timestamp
// 				createTimestamp
// 				updateTimestamp
// 				versionId
// 				parentId
// 				path
// 				resourceTypeId
// 			}
// 		}
// 		paging {
// 			next
// 		}
// 	}
// }
// `

// 	queryResourceGet = `
// query resourceGet($id: ID!) {
// 	resource(id: $id) {
// 		data
// 		metadata
// 		type {
// 			uri
// 		}
// 		turbot {
// 			id
// 			title
// 			tags
// 			akas
// 			timestamp
// 			createTimestamp
// 			updateTimestamp
// 			versionId
// 			parentId
// 			path
// 			resourceTypeId
// 		}
// 	}
// }
// `
// )

// type ResourcesResponse struct {
// 	Resources struct {
// 		Items  []Resource
// 		Paging struct {
// 			Next string
// 		}
// 	}
// }

// type ResourceResponse struct {
// 	Resource Resource
// }

// type Resource struct {
// 	Turbot   TurbotResourceMetadata
// 	Data     map[string]interface{}
// 	Metadata map[string]interface{}
// 	Type     struct {
// 		URI string
// 	}
// }

// type TurbotResourceMetadata struct {
// 	ID                string
// 	ParentID          string
// 	Akas              []string
// 	Custom            map[string]interface{}
// 	Metadata          map[string]interface{}
// 	Tags              map[string]interface{}
// 	Title             string
// 	VersionID         string
// 	ActorIdentityID   string
// 	ActorPersonaID    string
// 	ActorRoleID       string
// 	ResourceParentAka string
// 	Timestamp         string
// 	CreateTimestamp   string
// 	DeleteTimestamp   string
// 	UpdateTimestamp   string
// 	Path              string
// 	ResourceGroupIDs  []string
// 	ResourceTypeID    string
// 	State             string
// 	Terraform         map[string]interface{}
// }

func listEc2InstanceResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.listResource", "connection_error", err)
		return nil, err
	}

	/*
		filters := []string{"limit:5000"}
		allQuals := d.QueryContext.Quals
		plugin.Logger(ctx).Warn("listResource", "allQuals", allQuals)
		if allQuals["resource_type_id"] != nil {
			rti := allQuals["resource_type_id"]
			for _, q := range rti.Quals {
				plugin.Logger(ctx).Warn("listResource", "q", q)
				plugin.Logger(ctx).Warn("listResource", "q.GetFieldName()", q.GetFieldName())
				plugin.Logger(ctx).Warn("listResource", "q.GetOperator()", q.GetOperator())
				plugin.Logger(ctx).Warn("listResource", "q.GetValue().GetInt64Value()", q.GetValue().GetInt64Value())
				filters = append(filters, fmt.Sprintf("resourceTypeId:%d", q.GetValue().GetInt64Value()))
			}
		}
		filter := strings.Join(filters, " ")
		plugin.Logger(ctx).Warn("listResource", "filter", filter)
	*/

	filter := "resourceType:tmod:@turbot/aws-ec2#/resource/types/instance"

	//quals := d.KeyColumnQuals
	//filter := quals["filter"].GetStringValue()

	plugin.Logger(ctx).Warn("listResource", "filter", filter, "d", d)

	nextToken := ""

	for {
		result := &ResourcesResponse{}
		err = conn.DoRequest(queryResourceList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		plugin.Logger(ctx).Warn("listResource", "result", result, "next", result.Resources.Paging.Next, "err", err)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource.listResource", "query_error", err)
			return nil, err
		}
		for _, r := range result.Resources.Items {
			d.StreamListItem(ctx, r.Data)
		}
		if result.Resources.Paging.Next == "" {
			break
		}
		nextToken = result.Resources.Paging.Next
	}

	return nil, nil
}

func getEc2InstanceResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.getResource", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	idStr := strconv.FormatInt(id, 10)
	plugin.Logger(ctx).Warn("getResource", "id", id, "id.str", idStr, "d", d)

	result := &ResourceResponse{}

	start := time.Now()
	err = conn.DoRequest(queryResourceGet, map[string]interface{}{"id": id}, result)
	plugin.Logger(ctx).Warn("getResource", "time", time.Since(start).Milliseconds(), "result", result, "err", err)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.getResource", "query_error", err)
		return nil, err
	}
	return result.Resource, nil
}
