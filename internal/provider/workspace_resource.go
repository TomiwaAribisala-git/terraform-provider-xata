package provider

import (
	"fmt"
    "context"
	"github.com/xataio/xata-go/xata"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &workspaceResource{}
    _ resource.ResourceWithConfigure = &workspaceResource{}
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
    Name        types.String    `tfsdk:"name"`
    Slug        types.String    `tfsdk:"slug"`
	ID          types.String    `tfsdk:"id"`
    MemberCount types.Int64     `tfsdk:"memberCount"`
    Plan        types.String    `tfsdk:"plan"`
}
 
// Metadata returns the resource type name.
func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Schema defines the schema for the resource.
func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                Required: true,
            },
            "slug": schema.StringAttribute{
                Optional: true,
            },
            "id": schema.StringAttribute{
                Computed: true,
            },
            "memberCount": schema.Int64Attribute{
                Computed: true,
            },
            "plan": schema.StringAttribute{
                Computed: true,
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

    // Create new order
    workspace, err := r.client.Create(ctx, &workspaceRequest)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating workspace",
            "Could not create workspace, unexpected error: "+err.Error(),
        )
        return
    }

    // Map response body to schema and populate Computed attribute values
	plan.Name = types.StringValue(workspace.Name)
	plan.Slug = types.StringValue(*workspace.Slug)
	plan.ID = types.StringValue(workspace.Id)
	plan.MemberCount = types.Int64Value(int64(workspace.MemberCount))
	plan.Plan = types.StringValue(workspace.Plan.String())

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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