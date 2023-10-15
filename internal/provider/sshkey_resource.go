package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-clouding/internal/clouding"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SshKeyResource{}
var _ resource.ResourceWithImportState = &SshKeyResource{}

func NewSshKeyResource() resource.Resource {
	return &SshKeyResource{}
}

// SshKeyResource defines the resource implementation.
type SshKeyResource struct {
	client *clouding.API
}

// SshKeyResourceModel describes the resource data model.
type SshKeyResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	PublicKey     types.String `tfsdk:"public_key"`
	PrivateKey    types.String `tfsdk:"private_key"`
	HasPrivateKey types.Bool   `tfsdk:"has_private_key"`
	FingerPrint   types.String `tfsdk:"fingerprint"`
	LastUpdated   types.String `tfsdk:"last_updated"`
}

func (r *SshKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sshkey"
}

func (r *SshKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language Firewall.
		MarkdownDescription: `The SSH keys API provides a set of operations for managing SSH keys used for secure shell authentication when connecting to a clouding virtual machine.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A unique string identifier used to reference a ssh key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the SSH key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The public key of the SSH RSA key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The private key of the SSH RSA key.",
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fingerprint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fingerprint of the SSH RSA key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"has_private_key": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "The SSH key has private key.",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SSH key datetime update",
			},
		},
	}
}

func (r *SshKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clouding.API)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clouding.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SshKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKey := clouding.SshKey{
		Name:          plan.Name.ValueString(),
		PublicKey:     plan.PublicKey.ValueString(),
		PrivateKey:    plan.PrivateKey.ValueString(),
		HasPrivateKey: plan.HasPrivateKey.ValueBool(),
	}
	err := r.client.CreateSshKey(&sshKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ssh key",
			fmt.Sprintf("Error creating ssh key: %s", err),
		)
		return
	}

	// Set into the Terraform state.
	plan.Id = types.StringValue(sshKey.ID)
	plan.Name = types.StringValue(sshKey.Name)
	plan.PublicKey = types.StringValue(sshKey.PublicKey)
	plan.PrivateKey = types.StringValue(sshKey.PrivateKey)
	plan.HasPrivateKey = types.BoolValue(sshKey.HasPrivateKey)
	plan.FingerPrint = types.StringValue(sshKey.Fingerprint)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	//Write logs using tflog
	tflog.Trace(ctx, "Ssh key resource created")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKey, err := r.client.GetSshKeyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Clouding Client Error",
			fmt.Sprintf("Unable to read ssh key id, got error:  %s", err),
		)
		return
	}

	state.Name = types.StringValue(sshKey.Name)
	state.PublicKey = types.StringValue(sshKey.PublicKey)
	state.PrivateKey = types.StringValue(sshKey.PrivateKey)
	state.HasPrivateKey = types.BoolValue(sshKey.HasPrivateKey)
	state.FingerPrint = types.StringValue(sshKey.Fingerprint)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *SshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *SshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSshKey(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Clouding Client Error",
			fmt.Sprintf("Unable to delete ssh key id, got error:  %s", err),
		)
		return
	}

}

func (r *SshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
