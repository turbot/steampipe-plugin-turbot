package turbot

import (
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// column definitions for the common columns
func commonTurbotColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("__turbot.Title"),
			Description: resourceInterfaceDescription("title"),
		},
		{
			Name:        "akas",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("__turbot.Akas"),
			Description: resourceInterfaceDescription("akas"),
		},
		{
			Name:        "tags",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("__turbot.Tags"),
			Description: resourceInterfaceDescription("tags"),
		},
		{
			Name:        "turbot_id",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("__turbot.ID"),
			Description: "Turbot unique identifier.",
		},
	}
}

// column definitions for the common columns
func commonAwsRegionalColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("__metadata.aws.partition"),
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "region",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("__metadata.aws.regionName"),
			Description: "The AWS Region in which the resource is located.",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("__metadata.aws.accountId"),
			Description: "The AWS Account ID in which the resource is located.",
		},
	}
}

// // column definitions for the common columns
// func commonS3Columns() []*plugin.Column {
// 	return []*plugin.Column{
// 		{
// 			Name:        "partition",
// 			Type:        proto.ColumnType_STRING,
// 			Transform:   transform.FromField("__metadata.aws.partition"),
// 			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
// 		},
// 		{
// 			Name:        "account_id",
// 			Type:        proto.ColumnType_STRING,
// 			Transform:   transform.FromField("__metadata.aws.accountId"),
// 			Description: "The AWS Account ID in which the resource is located.",
// 		},
// 	}
// }

func commonAwsColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "partition",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("__metadata.aws.partition"),
			Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
		},
		{
			Name:        "region",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromConstant("global"),
			Description: "The AWS Region in which the resource is located.",
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Description: "The AWS Account ID in which the resource is located.",
			Transform:   transform.FromField("__metadata.aws.accountId"),
		},
	}
}

// append the common aws columns for REGIONAL resources onto the column list
func awsRegionalColumns(columns []*plugin.Column) []*plugin.Column {
	//return append(columns, commonAwsRegionalColumns()...)
	columns = append(columns, commonAwsRegionalColumns()...)
	columns = append(columns, commonTurbotColumns()...)
	return columns

}

// append the common aws columns for GLOBAL resources onto the column list
func awsColumns(columns []*plugin.Column) []*plugin.Column {
	//return append(columns, commonAwsColumns()...)
	columns = append(columns, commonAwsColumns()...)
	columns = append(columns, commonTurbotColumns()...)
	return columns
}

// func awsS3Columns(columns []*plugin.Column) []*plugin.Column {
// 	//return append(columns, commonS3Columns()...)
// 	columns = append(columns, commonS3Columns()...)
// 	columns = append(columns, commonTurbotColumns()...)
// 	return columns
// }

func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."
	}
	return ""
}
