// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xataio/xata-go/xata"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &workspacesDataSource{}
	_ datasource.DataSourceWithConfigure = &workspacesDataSource{}
)

// workspacesDataSourceModel maps the data source schema data.
type workspacesDataSourceModel struct {
	Workspaces []workspacesModel `tfsdk:"workspaces"`
}

// workspacesModel maps workspaces schema data.
type workspacesModel struct {
	ID types.String `tfsdk:"id"`
	//Unique_id   types.String              `tfsdk:"unique_id"`
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
	Role types.String `tfsdk:"role"`
	Plan types.String `tfsdk:"plan"`
}

// workspacesDataSource is the data source implementation.
type workspacesDataSource struct {
	client xata.WorkspacesClient
}

// NewWorkspacesDataSource is a helper function to simplify the provider implementation.
func NewWorkspacesDataSource() datasource.DataSource {
	return &workspacesDataSource{}
}

// Metadata returns the data source type name.
func (d *workspacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

// Schema defines the schema for the data source.
func (d *workspacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of workspaces.",
		Attributes: map[string]schema.Attribute{
			"workspaces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Numeric Identifier of each worskpace.",
							Computed:    true,
						},
						//"unique_id": schema.StringAttribute{
						//Computed: true,
						//},
						"name": schema.StringAttribute{
							Description: "Name of each worskpace.",
							Computed:    true,
						},
						"slug": schema.StringAttribute{
							Description: "Slug identifier of each worskpace.",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "User role	status of each worskpace.",
							Computed:    true,
						},
						"plan": schema.StringAttribute{
							Description: "Tier of each worskpace.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *workspacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(xata.WorkspacesClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *xata.workspaceCli, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *workspacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workspacesDataSourceModel

	workspaceresponse, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read User Workspaces",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, workspace := range workspaceresponse.Workspaces {
		workspaceState := workspacesModel{
			ID: types.StringValue(workspace.Id),
			// Unique_id:  types.StringValue(*workspace.Unique_id),
			Name: types.StringValue(workspace.Name),
			Slug: types.StringValue(workspace.Slug),
			Role: types.StringValue(workspace.Role.String()),
			Plan: types.StringValue(workspace.Plan.String()),
		}

		state.Workspaces = append(state.Workspaces, workspaceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
