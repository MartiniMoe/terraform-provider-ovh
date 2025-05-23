// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func CloudProjectGatewayInterfaceDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Id",
		},
		"interface_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Interface ID",
		},
		"ip": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "IP of the interface",
		},
		"network_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Network ID of the interface",
		},
		"region": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Region name",
		},
		"service_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Required:    true,
			Description: "Service name",
		},
		"subnet_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Subnet ID of the interface",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}
