package turbot

import (
	"context"
	"strings"

	//"fmt"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotAwsEc2Instance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "aws_ec2_instance",
		Description:      "AWS EC2 Instance, from a Turbot Workspace",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("turbot_id"),
		// 	Hydrate:    getEc2InstanceResource,
		// },
		List: &plugin.ListConfig{
			Hydrate: listEc2InstanceResources,
		},
		/// need the regions stuff too... Columns: awsRegionalColumns([]*plugin.Column{
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "instance_id",
				Description: "The ID of the instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) specifying the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromValue().Transform(ec2InstanceArn),
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
				Transform:   transform.FromField("KernelId.Value"),
			},
			{
				Name:        "key_name",
				Description: "The name of the key pair, if this instance was launched with an associated key pair",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_time",
				Description: "The time the instance was launched.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "metadata_options",
				Description: "The metadata options for the instance.",
				Type:        proto.ColumnType_JSON,
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

			//** remove this later...
			{
				Name:        "raw",
				Description: "remove...",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listEc2InstanceResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listResources(ctx, d, h, "resourceType:tmod:@turbot/aws-ec2#/resource/types/instance")

}

func getEc2InstanceResource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	id := quals["turbot_id"].GetInt64Value()
	return getResourceById(ctx, d, h, id)
}

// TRANSFORM
func ec2InstanceArn(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(map[string]interface{})

	metaData := instance["__metadata"].(map[string]interface{})
	awsMetaData := metaData["aws"].(map[string]interface{})

	return strings.Join([]string{"arn",
		awsMetaData["partition"].(string),
		"ec2",
		awsMetaData["regionName"].(string),
		awsMetaData["accountId"].(string),
		"instance/" +
			instance["InstanceId"].(string),
	}, ":"), nil

}
