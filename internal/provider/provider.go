package provider

import (
    "os"
    "context"
    "github.com/xataio/xata-go/xata"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &xataProvider{}
)

// xataProviderModel maps provider schema data to a Go type.
type xataProviderModel struct {
    Apikey types.String `tfsdk:"apikey"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &xataProvider{
            version: version,
        }
    }
}

// xataProvider is the provider implementation.
type xataProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// Metadata returns the provider type name.
func (p *xataProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "xata"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *xataProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "apikey": schema.StringAttribute{
                Optional:  true,
                Sensitive: true,
            },
        },
    }
}

func (p *xataProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    tflog.Info(ctx, "Configuring Xata client")

    // Retrieve provider data from configuration
    var config xataProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.Apikey.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("apikey"),
            "Unknown Xata API Key",
            "The provider cannot create the Xata API client as there is an unknown configuration value for the Xata API Key. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the XATA_API_KEY environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    apikey := os.Getenv("XATA_API_KEY")

    if !config.Apikey.IsNull() {
        apikey = config.Apikey.ValueString()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if apikey == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("apikey"),
            "Missing Xata API Key",
            "The provider cannot create the Xata API client as there is a missing or empty value for the Xata API Key. "+
                "Set the password value in the configuration or use the XATA_API_KEY environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    ctx = tflog.SetField(ctx, "xata_apikey", apikey)
    ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "xata_apikey")

    tflog.Debug(ctx, "Creating Xata client")

    // Create a new Xata client using the configuration values
    client, err := xata.NewWorkspacesClient(xata.WithAPIKey("apikey")) 
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Xata API Client",
            "An unexpected error occurred when creating the Xata API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Xata Client Error: "+err.Error(),
        )
        return
    }

    // Make the Xata client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client

    tflog.Info(ctx, "Configured Xata client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *xataProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource {
      NewWorkspacesDataSource,
    }
}  

// Resources defines the resources implemented in the provider.
func (p *xataProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewWorkspaceResource,
    }
}