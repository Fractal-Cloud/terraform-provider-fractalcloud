package provider

import (
	"context"
	"os"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &fractalCloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &fractalCloudProvider{
			version: version,
		}
	}
}

// fractalCloudProvider is the provider implementation.
type fractalCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// fractalCloudProviderModel maps provider schema data to a Go type.
type fractalCloudProviderModel struct {
	Host                 types.String `tfsdk:"host"`
	ServiceAccountId     types.String `tfsdk:"service_account_id"`
	ServiceAccountSecret types.String `tfsdk:"service_account_secret"`
}

// Metadata returns the provider type name.
func (p *fractalCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fractalcloud"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *fractalCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"service_account_id": schema.StringAttribute{
				Optional: true,
			},
			"service_account_secret": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *fractalCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config fractalCloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := "https://api.fractal.cloud"
	serviceAccountId := os.Getenv("FRACTAL_CLOUD_SERVICE_ACCOUNT_ID")
	serviceAccountSecret := os.Getenv("FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ServiceAccountId.IsNull() {
		serviceAccountId = config.ServiceAccountId.ValueString()
	}

	if !config.ServiceAccountSecret.IsNull() {
		serviceAccountSecret = config.ServiceAccountSecret.ValueString()
	}

	if serviceAccountId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_account_id"),
			"Missing Fractal Cloud API Service Account Id",
			"The provider cannot create the Fractal Cloud API client as there is a missing or empty value for the Fractal Cloud API Service Account Id. "+
				"Set the username value in the configuration or use the FRACTAL_CLOUD_SERVICE_ACCOUNT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if serviceAccountSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Fractal Cloud API Password",
			"The provider cannot create the Fractal Cloud API client as there is a missing or empty value for the Fractal Cloud API password. "+
				"Set the password value in the configuration or use the Fractal FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Fractal Cloud client using the configuration values
	client := fractalCloud.NewClient(&fractalCloud.ClientLogger{
		Debug: func(s string) {
			tflog.Debug(ctx, s)
		},
		Information: func(s string) {
			tflog.Info(ctx, s)
		},
		Warning: func(s string) {
			tflog.Warn(ctx, s)
		},
		Error: func(s string) {
			tflog.Error(ctx, s)
		},
	}, &host, &serviceAccountId, &serviceAccountSecret)

	// Make the Fractal Cloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *fractalCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPersonalBoundedContextDataSource,
		NewOrganizationalBoundedContextDataSource,
		NewOrganizationDataSource,
		NewFractalDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *fractalCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPersonalBoundedContext,
		NewOrganizationalBoundedContext,
		NewManagementEnvironment,
		NewOperationalEnvironment,
		NewFractal,
	}
}
