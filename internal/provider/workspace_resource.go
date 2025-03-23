// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xataio/xata-go/xata"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &workspaceResource{}
	_ resource.ResourceWithConfigure   = &workspaceResource{}
	_ resource.ResourceWithImportState = &workspaceResource{}
)

// NewWorkspaceResource is a helper function to simplify the provider implementation.
func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

// workspaceResource is the resource implementation.
type workspaceResource struct {
	client xata.WorkspacesClient
}

// workspaceResourceModel maps the resource schema data.
type workspaceResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Id          types.String `tfsdk:"id"`
	MemberCount types.Int64  `tfsdk:"membercount"`
	Plan        types.String `tfsdk:"plan"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Configure adds the provider configured client to the resource.
func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(xata.WorkspacesClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected xata.WorkspacesClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a workspace.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the worskpace.",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "Slug Identifier of the worskpace.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Numeric Identifier of the worskpace.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"membercount": schema.Int64Attribute{
				Description: "Member Count of the workspace.",
				Computed:    true,
			},
			"plan": schema.StringAttribute{
				Description: "Tier of the worskpace.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the workspace.",
				Computed:    true,
			},
		},
	}
}

// Create a new resource.
func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan workspaceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	workspaceRequest := xata.WorkspaceMeta{
		Name: plan.Name.ValueString(),
		Slug: xata.String(plan.Slug.ValueString()),
	}

	// Create new workspace
	workspace, err := r.client.Create(ctx, &workspaceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Xata workspace",
			fmt.Sprintf("Could not create workspace, unexpected error: %s", err.Error()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(workspace.Name)
	plan.Slug = types.StringValue(*workspace.Slug)
	plan.Id = types.StringValue(workspace.Id)
	plan.MemberCount = types.Int64Value(int64(workspace.MemberCount))
	plan.Plan = types.StringValue(workspace.Plan.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get workspace Id
	var id types.String
	diags := req.State.GetAttribute(ctx, path.Root("id"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get existing workspace for a given Id
	workspaceInfo, err := r.client.GetWithWorkspaceID(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Xata workspace",
			fmt.Sprintf("Could not read workspace, unexpected error: %s", err.Error()),
		)
		return
	}

	// Map API response to resource model
	workspace := workspaceResourceModel{
		Name:        types.StringValue(workspaceInfo.Name),
		Slug:        types.StringValue(*workspaceInfo.Slug),
		Id:          types.StringValue(workspaceInfo.Id),
		MemberCount: types.Int64Value(int64(workspaceInfo.MemberCount)),
		Plan:        types.StringValue((workspaceInfo.Plan.String())),
	}

	// Return workspace info
	diags = resp.State.Set(ctx, &workspace)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource information.
func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get workspace Id
	var id types.String
	diags := req.State.GetAttribute(ctx, path.Root("id"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan workspaceResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	workspaceRequest := xata.UpdateWorkspaceRequest{
		WorkspaceID: xata.String(id.ValueString()),
		Payload: &xata.WorkspaceMeta{
			Name: plan.Name.ValueString(),
			Slug: xata.String(plan.Slug.ValueString()),
		},
	}

	// Update existing workspace
	updatedWorkspace, err := r.client.Update(ctx, workspaceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Xata Workspace",
			fmt.Sprintf("Could not update workspace, unexpected error: %s", err.Error()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(updatedWorkspace.Name)
	plan.Slug = types.StringValue(*updatedWorkspace.Slug)
	plan.Id = types.StringValue(updatedWorkspace.Id)
	plan.MemberCount = types.Int64Value(int64(updatedWorkspace.MemberCount))
	plan.Plan = types.StringValue(updatedWorkspace.Plan.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get workspace Id
	var id types.String
	diags := req.State.GetAttribute(ctx, path.Root("id"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete workspace
	err := r.client.Delete(ctx, id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Xata Workspace",
			fmt.Sprintf("Could not delete workspace, unexpected error: %s", err.Error()),
		)
		return
	}
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
