// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

var _ resource.Resource = &ManagedEnvironmentResource{}
var _ resource.ResourceWithImportState = &ManagedEnvironmentResource{}

const SOLUTION_CHECKER_RULES = "meta-remove-dup-reg, meta-avoid-reg-no-attribute, meta-avoid-reg-retrieve, meta-remove-inactive, web-avoid-unpub-api, web-avoid-modals, web-avoid-crm2011-service-odata, web-avoid-crm2011-service-soap, web-avoid-browser-specific-api, web-avoid-2011-api, web-use-relative-uri, web-use-async, web-avoid-window-top, web-use-client-context, web-use-navigation-api, web-use-offline, web-use-grid-api, web-avoid-isactivitytype, meta-avoid-silverlight, meta-avoid-retrievemultiple-annotation, web-remove-debug-script, web-use-strict-mode, web-use-strict-equality-operators, web-avoid-eval, app-formula-issues-high, app-formula-issues-medium, app-formula-issues-low, app-use-delayoutput-text-input, app-reduce-screen-controls, app-include-accessible-label, app-include-alternative-input, app-avoid-autostart, app-include-captions, app-make-focusborder-visible, app-include-helpful-control-setting, app-avoid-interactive-html, app-include-readable-screen-name, app-include-state-indication-text, app-include-tab-order, app-include-tab-index, flow-avoid-recursive-loop, flow-avoid-invalid-reference, flow-outlook-attachment-missing-info, meta-include-missingunmanageddependencies, web-remove-alert, web-remove-console, web-use-global-context, web-use-org-setting, app-testformula-issues-high, app-testformula-issues-medium, app-testformula-issues-low, flow-avoid-connection-mode, web-avoid-with, web-avoid-loadtheme, web-use-getsecurityroleprivilegesinfo, web-sdl-no-cookies, web-sdl-no-document-domain, web-sdl-no-document-write, web-sdl-no-html-method, web-sdl-no-inner-html, web-sdl-no-insecure-url, web-sdl-no-msapp-exec-unsafe, web-sdl-no-postmessage-star-origin, web-sdl-no-winjs-html-unsafe, connector-validate-brandcolor, connector-validate-iconimage, connector-validate-swagger-isproperjson, connector-validate-swagger, connector-validate-swagger-extended, connector-validate-title, connector-validate-connectionparam-isproperjson, connector-validate-connectionparameters, connector-validate-connectionparam-oauth2idp, meta-license-sales-sdkmessages, meta-license-sales-entity-operations, meta-license-sales-customcontrols, web-use-appsidepane-api, meta-license-fieldservice-sdkmessages, meta-license-fieldservice-entity-operations, meta-license-fieldservice-customcontrols, meta-avoid-managed-entity-assets, meta-include-unmanaged-entity-assets, connector-validate-hexadecimalbrandcolor, connector-validate-pngiconimage, connector-validate-iconsize, connector-validate-backgroundwithbrandiconcolor, web-unsupported-syntax"

func NewManagedEnvironmentResource() resource.Resource {
	return &ManagedEnvironmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "managed_environment",
		},
	}
}

func (r *ManagedEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ManagedEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Manages a \"Managed Environment\" and associated settings",
		MarkdownDescription: "Manages a [Managed Environment](https://learn.microsoft.com/power-platform/admin/managed-environment-overview) and associated settings. A Power Platform Managed Environment is a suite of premium capabilities that allows administrators to manage Power Platform at scale with more control, less effort, and more insights. Once an environment is managed, it unlocks additional features across the Power Platform",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique managed environment settings id (guid)",
				Description:         "Unique managed environment settings id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid), of the environment that is managed by these settings",
				Description:         "Unique environment id (guid), of the environment that is managed by these settings",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protection_level": schema.StringAttribute{
				MarkdownDescription: "Protection level",
				Description:         "Protection level",
				Computed:            true,
			},
			"is_usage_insights_disabled": schema.BoolAttribute{
				MarkdownDescription: "[Weekly insights digest for the environment](https://learn.microsoft.com/power-platform/admin/managed-environment-usage-insights)",
				Description:         "Weekly insights digest for the environment",
				Required:            true,
			},
			"is_group_sharing_disabled": schema.BoolAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared. See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details.",
				Description:         "Limits how widely canvas apps can be shared",
				Required:            true,
			},
			"limit_sharing_mode": schema.StringAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared.  See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details",
				Description:         "Limits how widely canvas apps can be shared.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ExcludeSharingToSecurityGroups", "NoLimit"),
				},
			},
			"max_limit_user_sharing": schema.Int64Attribute{
				MarkdownDescription: "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'",
				Description:         "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'. See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details",
				Required:            true,
			},
			"solution_checker_mode": schema.StringAttribute{
				MarkdownDescription: "Automatically verify solution checker results for security and reliability issues before solution import.  See [Solution Checker enforcement](https://learn.microsoft.com/power-platform/admin/managed-environment-solution-checker) for more details.",
				Description:         "Automatically verify solution checker results for security and reliability issues before solution import.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("None", "Warn", "Block"),
				},
			},
			"suppress_validation_emails": schema.BoolAttribute{
				MarkdownDescription: "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Description:         "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Required:            true,
			},
			"solution_checker_rule_overrides": schema.SetAttribute{
				MarkdownDescription: `
				List of rules to exclude from solution checker
				See [Solution Checker enforcement](https://learn.microsoft.com/power-platform/admin/managed-environment-solution-checker) for more details.
				Posible values are:

				| Code                                | Description                                                                                          | Summary                                | Guidance URL                                                                                                                               |
				|-------------------------------------|------------------------------------------------------------------------------------------------------|----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
				| meta-remove-dup-reg                 | Checks for duplicate Dataverse plug-in registrations                                                 | Duplicate plug-in registration         | [https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/do-not-duplicate-plugin-step-registration](https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/do-not-duplicate-plugin-step-registration) |
				| meta-avoid-reg-no-attribute         | Checks for filtering attributes with Dataverse plug-in registrations                                  | Check plug-in filtering attributes     | [https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/include-filtering-attributes-plugin-registration](https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/include-filtering-attributes-plugin-registration) |
				| meta-avoid-reg-retrieve             | Checks for Dataverse plug-ins registered for Retrieve and RetrieveMultiple messages                   | Check plug-ins for Retrieve messages   | [https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/limit-registration-plugins-retrieve-retrievemultiple](https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/business-logic/limit-registration-plugins-retrieve-retrievemultiple) |
				| meta-remove-inactive                | Checks for inactive plug-in configurations in Dataverse                                               | Check inactive plug-ins                | [https://learn.microsoft.com/powerapps/developer/model-driven-apps/best-practices/business-logic/remove-deactivated-disabled-configurations](https://learn.microsoft.com/powerapps/developer/model-driven-apps/best-practices/business-logic/remove-deactivated-disabled-configurations) |
				| web-avoid-unpub-api                 | Checks for usage of unpublished APIs                                                                  | Avoid unpublished APIs                 | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-unpub-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-unpub-api) |
				| web-avoid-modals                    | Checks if modal dialogs are used                                                                      | Check using modal dialogs              | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-modals](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-modals) |
				| web-avoid-crm2011-service-odata     | Checks for usage of the Dynamics CRM 2011 Odata 2.0 endpoint                                          | Avoid CRM 2011 OData endpoint          | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-crm2011-service-odata](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-crm2011-service-odata) |
				| web-avoid-crm2011-service-soap      | Checks for usage of the Dynamics CRM 2011 SOAP endpoint                                               | Avoid CRM 2011 SOAP endpoint           | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-crm2011-service-soap](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-crm2011-service-soap) |
				| web-avoid-browser-specific-api      | Checks for usage of Internet Explorer legacy APIs or browser plug-ins                                  | Avoid browser specific APIs            | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-browser-specific-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-browser-specific-api) |
				| web-avoid-2011-api                  | Checks for usage of the deprecated Dynamics CRM 2011 object model                                      | Avoid CRM 2011 API                     | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-2011-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-2011-api) |
				| web-use-relative-uri                | Checks for usage of absolute Dataverse endpoint URIs                                                  | Use relative URIs                      | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-relative-uri](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-relative-uri) |
				| web-use-async                       | Checks for async pattern in web resources                                                             | Check async pattern                    | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-async](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-async) |
				| web-avoid-window-top                | Checks for usage of window.top API                                                                     | Avoid window.top                       | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-window-top](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-window-top) |
				| web-use-client-context              | Checks if client context is used                                                                       | Use client context                     | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-client-context](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-client-context) |
				| web-use-navigation-api              | Checks if navigation API parameters are used                                                           | Use navigation API                     | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-navigation-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-navigation-api) |
				| web-use-offline                     | Checks if offline mode is used                                                                         | Use offline mode                       | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-offline](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-offline) |
				| web-use-grid-api                    | Checks if the grid APIs are used                                                                       | Use grid API                           | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-grid-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-grid-api) |
				| web-avoid-isactivitytype            | Checks for usage of isActivityType                                                                     | Avoid isActivityType                   | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-isactivitytype](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-isactivitytype) |
				| meta-avoid-silverlight              | Checks for usage of Silverlight                                                                        | Avoid Silverlight                      | [https://learn.microsoft.com/dynamics365/get-started/whats-new/customer-engagement/important-changes-coming#BKMK_Silverlight](https://learn.microsoft.com/dynamics365/get-started/whats-new/customer-engagement/important-changes-coming#BKMK_Silverlight) |
				| meta-avoid-retrievemultiple-annotation | Checks for registering a plugin on RetrieveMultiple of annotation                                      | Check plug-ins for RetrieveMultiple of annotations | [https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/](https://learn.microsoft.com/powerapps/developer/data-platform/best-practices/) |
				| web-remove-debug-script             | Checks for the usage of debug scripts                                                                  | Remove debug scripts                   | [https://learn.microsoft.com/powerapps/developer/model-driven-apps/best-practices/](https://learn.microsoft.com/powerapps/developer/model-driven-apps/best-practices/) |
				| web-use-strict-mode                 | Checks if strict mode is used                                                                          | Use strict mode                        | [https://developer.mozilla.org/docs/Web/JavaScript/Reference/Strict_mode/Transitioning_to_strict_mode](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Strict_mode/Transitioning_to_strict_mode) |
				| web-use-strict-equality-operators   | Checks if strict equality operators are used                                                           | Use strict equality operators          | [https://developer.mozilla.org/docs/Web/JavaScript/Equality_comparisons_and_sameness](https://developer.mozilla.org/docs/Web/JavaScript/Equality_comparisons_and_sameness) |
				| web-avoid-eval                      | Checks for usage of eval function or its functional equivalents                                        | Avoid eval                             | [https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/eval](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Global_Objects/eval) |
				| app-formula-issues-high             | Checks for high severity formula issues in Canvas apps                                                 | Fix high severity formula issues       | [https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference](https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference) |
				| app-formula-issues-medium           | Checks for medium severity formula issues in Canvas apps                                               | Fix medium severity formula issues     | [https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference](https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference) |
				| app-formula-issues-low              | Checks for low severity formula issues in Canvas apps                                                  | Fix low severity formula issues        | [https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference](https://learn.microsoft.com/powerapps/maker/canvas-apps/formula-reference) |
				| app-use-delayoutput-text-input      | Checks if delayed loading is used in Canvas apps                                                       | Use delay load                         | [https://learn.microsoft.com/powerapps/maker/canvas-apps/performance-tips#use-delayed-load](https://learn.microsoft.com/powerapps/maker/canvas-apps/performance-tips#use-delayed-load) |
				| app-reduce-screen-controls          | Checks for excessive controls on a screen in Canvas apps                                               | Reduce screen controls                 | [https://learn.microsoft.com/powerapps/maker/canvas-apps/performance-tips#limit-the-number-of-controls](https://learn.microsoft.com/powerapps/maker/canvas-apps/performance-tips#limit-the-number-of-controls) |
				| app-include-accessible-label        | Checks if accessible labels are included in Canvas apps                                                | Include accessible label               | [https://www.w3.org/WAI/tutorials/forms/labels/](https://www.w3.org/WAI/tutorials/forms/labels/) |
				| app-include-alternative-input       | Checks if all interactive elements are accessible to alternative inputs in Canvas apps                  | Include alternative input              | [https://www.w3.org/WAI/tips/developing/#ensure-that-all-interactive-elements-are-keyboard-accessible](https://www.w3.org/WAI/tips/developing/#ensure-that-all-interactive-elements-are-keyboard-accessible) |
				| app-avoid-autostart                 | Checks for autostart on players within a Canvas app                                                    | Avoid autostart in app                 | [https://digital.gov/2014/06/30/508-accessible-videos-use-a-508-compliant-video-player/](https://digital.gov/2014/06/30/508-accessible-videos-use-a-508-compliant-video-player/) |
				| app-include-captions                | Without captions, people with disabilities may not get any of the information in a video or audio segment | Missing captions                       | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-make-focusborder-visible        | If the focus isn't visible, people who don't use a mouse won't be able to see it when they're interacting with the app | Focus isn't showing                    | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-include-helpful-control-setting | Changing this property setting will give the user better information about the function of the controls in your app | Missing helpful control settings       | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-avoid-interactive-html          | Your app won't work correctly and will not be accessible if you place interactive HTML elements        | If this HTML contains interactive elements, consider using another method, or remove the HTML from this element | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-include-readable-screen-name    | People who are blind, have low vision, or a reading disability rely on screen titles to navigate using the screen reader | Revise screen name                     | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-include-state-indication-text   | Users won't get confirmation of their actions if the state of the control isn't showing                | Add State indication text              | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-include-tab-order               | When a screen reader reads the elements of a slide, it's important that they appear in the order that a user would see them, instead of the order they were added to the slide | Check the order of the screen items    | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| app-include-tab-index               | People who use the keyboard with your app will not be able to access this element without a tab stop   | Missing tab stop                       | [https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues](https://learn.microsoft.com/power-apps/maker/canvas-apps/accessibility-checker#types-of-issues) |
				| connector-validate-pngiconimage         | The icon image is invalid.Icon image should be valid image file and must be submitted in PNG format as icon.png. | The icon image is invalid.Icon image should be valid image file and must be submitted in PNG format as icon.png. | [https://go.microsoft.com/fwlink/?linkid=2244166](https://go.microsoft.com/fwlink/?linkid=2244166)                                         |
				| connector-validate-iconsize             | The icon image size is invalid. Update the icon dimensions to 1:1 within the range of 100x100 to 230x230 pixels. | The icon image size is invalid. Update the icon dimensions to 1:1 within the range of 100x100 to 230x230 pixels. | [https://go.microsoft.com/fwlink/?linkid=2244166](https://go.microsoft.com/fwlink/?linkid=2244166)                                         |
				| connector-validate-backgroundwithbrandiconcolor | The background color of icon image is invalid. Update with consistent background.                    | The background color of icon image is invalid. Update with consistent background. | [https://go.microsoft.com/fwlink/?linkid=2244166](https://go.microsoft.com/fwlink/?linkid=2244166)                                         |
				| web-unsupported-syntax                  |                                                                                                      |                                                                         | [http://go.microsoft.com/fwlink/?LinkID=398563&error=web-unsupported-syntax&client=PAChecker](http://go.microsoft.com/fwlink/?LinkID=398563&error=web-unsupported-syntax&client=PAChecker) |
				| flow-avoid-recursive-loop               | Avoid recursive action as they may result in an infinite trigger loop                                | Avoid recursive action as they may result in an infinite trigger loop   | [https://learn.microsoft.com/flow/error-checker](https://learn.microsoft.com/flow/error-checker)                                           |
				| flow-avoid-invalid-reference            | Include valid references for actions                                                                 | Include valid references for actions                                    | [https://learn.microsoft.com/flow/error-checker](https://learn.microsoft.com/flow/error-checker)                                           |
				| flow-outlook-attachment-missing-info    | Include all required outlook attachment information                                                  | Include all required outlook attachment information                     | [https://learn.microsoft.com/flow/error-checker](https://learn.microsoft.com/flow/error-checker)                                           |
				| meta-include-missingunmanageddependencies | Checks for missing unmanaged dependencies in the solution. Missing unmanaged dependencies will cause a solution to fail to import in a target environment | Checks for missing unmanaged dependencies in the solution. Missing unmanaged dependencies will cause a solution to fail to import in a target environment | [https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import](https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import) |
				| web-remove-alert                        | Checks for usage of alert function or its functional equivalents                                      | Remove alerts                                                           | [https://eslint.org/docs/rules/no-alert](https://eslint.org/docs/rules/no-alert)                                                           |
				| web-remove-console                      | Checks for the usage of methods on console                                                            | Remove console statements                                               | [https://eslint.org/docs/rules/no-console](https://eslint.org/docs/rules/no-console)                                                       |
				| web-use-global-context                  | Checks if global context is used                                                                      | Use global context                                                      | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-global-context](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-global-context) |
				| web-use-org-setting                     | Checks if org settings are used                                                                       | Use organization settings                                               | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-org-setting](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-org-setting) |
				| app-testformula-issues-high             | Checks for high severity formula issues in Test Studio Canvas apps                                    | Fix high severity formula issues for Test Studio                        | [https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions](https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions) |
				| app-testformula-issues-medium           | Checks for medium severity formula issues in Test Studio Canvas apps                                  | Fix medium severity formula issues for Test Studio                      | [https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions](https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions) |
				| app-testformula-issues-low              | Checks for low severity formula issues in Test Studio Canvas apps                                     | Fix low severity formula issues for Test Studio                         | [https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions](https://learn.microsoft.com/powerapps/maker/canvas-apps/working-with-test-studio#test-functions) |
				| flow-avoid-connection-mode              | Use connection references instead of connections.                                                     | Use connection references instead of connections.                       | [https://learn.microsoft.com/powerapps/maker/data-platform/create-connection-reference#updating-a-flow-to-use-connection-references-instead-of-connections](https://learn.microsoft.com/powerapps/maker/data-platform/create-connection-reference#updating-a-flow-to-use-connection-references-instead-of-connections) |
				| web-avoid-with                          | Checks for usage of 'with' operator                                                                   | Avoid 'with' operator                                                   | [https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/with](https://developer.mozilla.org/docs/Web/JavaScript/Reference/Statements/with) |
				| web-avoid-loadtheme                     | Checks for usage of the loadTheme Fluent v8 API                                                       | Avoid LoadTheme API                                                     | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-loadtheme](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/avoid-loadtheme) |
				| web-use-getsecurityroleprivilegesinfo   | Checks for usage of userSettings.securityRolePrivileges                                               | Avoid securityRolePrivileges                                            | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-getsecurityroleprivilegesinfo](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-getsecurityroleprivilegesinfo) |
				| web-sdl-no-cookies                      | HTTP cookies are an old client-side storage mechanism with inherent risks and limitations. Use Web Storage, IndexedDB or other modern methods instead. | Do not use HTTP cookies in modern applications                          | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-cookies.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-cookies.md)               |
				| web-sdl-no-document-domain              | Writes to document.domain property must be reviewed to avoid bypass of same-origin checks. Usage of top level domains such as azurewebsites.net is strictly prohibited. | Do not write to document.domain property                                | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-document-domain.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-document-domain.md) |
				| web-sdl-no-document-write               | Calls to document.write or document.writeln manipulate DOM directly without any sanitization and should be avoided. Use document.createElement() or similar methods instead. | Do not write to DOM directly using document.write or document.writeln methods | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-document-write.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-document-write.md) |
				| web-sdl-no-html-method                  | Direct calls to method html() often (e.g. in jQuery framework) manipulate DOM without any sanitization and should be avoided. Use document.createElement() or similar methods instead. | Do not write to DOM directly using jQuery html() method                 | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-html-method.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-html-method.md)       |
				| web-sdl-no-inner-html                   | Assignments to innerHTML or outerHTML properties manipulate DOM directly without any sanitization and should be avoided. Use document.createElement() or similar methods instead. | Do not write to DOM directly using innerHTML/outerHTML property         | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-inner-html.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-inner-html.md)         |
				| web-sdl-no-insecure-url                 | Insecure protocols such as HTTP or FTP should be replaced by their encrypted counterparts (HTTPS, FTPS) to avoid sending potentially sensitive data over untrusted networks in plaintext. | Do not use insecure URLs                                                | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-insecure-url.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-insecure-url.md)     |
				| web-sdl-no-msapp-exec-unsafe            | Calls to MSApp.execUnsafeLocalFunction() bypass script injection validation and should be avoided.    | Do not bypass script injection validation                               | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-msapp-exec-unsafe.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-msapp-exec-unsafe.md) |
				| web-sdl-no-postmessage-star-origin      | Always provide specific target origin, not * when sending data to other windows using postMessage to avoid data leakage outside of trust boundary. | Do not use * as target origin when sending data to other windows        | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-postmessage-star-origin.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-postmessage-star-origin.md) |
				| web-sdl-no-winjs-html-unsafe            | Calls to WinJS.Utilities.setInnerHTMLUnsafe() and similar methods do not perform any input validation and should be avoided. Use WinJS.Utilities.setInnerHTML() instead. | Do not set HTML using unsafe methods from WinJS.Utilities               | [https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-winjs-html-unsafe.md](https://github.com/microsoft/eslint-plugin-sdl/blob/main/docs/rules/no-winjs-html-unsafe.md) |
				| connector-validate-brandcolor           | Ensure brand color is a valid hexadecimal color and should not be white (#ffffff) or default (#007ee5). | Invalid brand color                                                     | [https://go.microsoft.com/fwlink/?linkid=2244386](https://go.microsoft.com/fwlink/?linkid=2244386)                                         |
				| connector-validate-iconimage            | Ensure a) icon is in PNG format with size below 1MB b) icon's dimensions are 1:1 and fall within the range of 100x100 to 230x230 pixels c) icon matches the brand color with non-transparent, non-white color (#ffffff) and not-default color (#007ee5) background d) logo dimensions are below 70% for image's height & width with consistent background. | Invalid Icon Image Size                                                 | [https://go.microsoft.com/fwlink/?linkid=2244166](https://go.microsoft.com/fwlink/?linkid=2244166)                                         |
				| connector-validate-swagger-isproperjson | Ensure the openapidefinition is a well formatted JSON.                                                | The openapidefinition.json is not a valid JSON.                         | [https://go.microsoft.com/fwlink/?linkid=2244842](https://go.microsoft.com/fwlink/?linkid=2244842)                                         |
				| connector-validate-swagger              | Ensure swagger definition complies with OpenAPI 2.0 standard.                                         | Swagger definition does not confirm to the OpenAPI 2.0 standard.        | [https://go.microsoft.com/fwlink/?linkid=2245509](https://go.microsoft.com/fwlink/?linkid=2245509)                                         |
				| connector-validate-swagger-extended     | Ensure swagger definition complies with OpenAPI 2.0 standard and connectors' extended standard.       | Swagger definition does not confirm to the connector extended standards. | [https://go.microsoft.com/fwlink/?linkid=2245307](https://go.microsoft.com/fwlink/?linkid=2245307)                                         |
				| connector-validate-title                | Ensure connector title is unique and distinguishable from pre-existing connector title.               | Connector title is not unique.                                          | [https://go.microsoft.com/fwlink/?linkid=2247920](https://go.microsoft.com/fwlink/?linkid=2247920)                                         |
				| connector-validate-connectionparam-isproperjson | Ensure the connectionparameters is a well formatted JSON.                                             | The connectionparameters.json is not a valid JSON.                      | [https://go.microsoft.com/fwlink/?linkid=2248011](https://go.microsoft.com/fwlink/?linkid=2248011)                                         |
				| connector-validate-connectionparameters | Ensure the property is updated with appropriate value.                                                | The connectionparameter is not well formed.                             | [https://go.microsoft.com/fwlink/?linkid=2247861](https://go.microsoft.com/fwlink/?linkid=2247861)                                         |
				| connector-validate-connectionparam-oauth2idp | Ensure the identity provider is from the list of supported oauth2 providers.                          | Invalid OAuth2 Identity Provider                                        | [https://go.microsoft.com/fwlink/?linkid=2248012](https://go.microsoft.com/fwlink/?linkid=2248012)                                         |
				| meta-license-sales-sdkmessages          | Dynamics 365 SDK messages require users executing these operations to be licensed with any of the Dynamics 365 licenses entitled to this operation. Check product documentation and the Dynamics 365 Licensing Guide for additional information and license entitlements. | Dynamics 365 SDK messages require users executing these operations to be licensed with any of the appropriate Dynamics 365 license. | [https://go.microsoft.com/fwlink/?linkid=2247983](https://go.microsoft.com/fwlink/?linkid=2247983)                                         |
				| meta-license-sales-entity-operations    | Some operations performed on Dynamics 365 entities require users executing these operations to be licensed with any of the Dynamics 365 licenses entitled to this operation. Check product documentation and the Dynamics 365 Licensing Guide for additional information and license entitlements. | Some operations performed on Dynamics 365 entities require users executing these operations to be licensed with any of the appropriate Dynamics 365 license. | [https://go.microsoft.com/fwlink/?linkid=2248081](https://go.microsoft.com/fwlink/?linkid=2248081)                                         |
				| meta-license-sales-customcontrols       | Some Dynamics 365 Sales controls require users accessing these controls to be licensed with a Dynamics 365 Sales license. Check product documentation for additional information. | Some Dynamics 365 Sales controls require users accessing these controls to be licensed with a Dynamics 365 Sales license. | [https://go.microsoft.com/fwlink/?linkid=2248449](https://go.microsoft.com/fwlink/?linkid=2248449)                                         |
				| web-use-appsidepane-api                 | Checks for usage of legacy Xrm.Panel APIs                                                             | Use the AppSidePane APIs                                                | [https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-appsidepane-api](https://learn.microsoft.com/power-apps/maker/data-platform/powerapps-checker/rules/web/use-appsidepane-api) |
				| meta-license-fieldservice-sdkmessages   | Dynamics 365 SDK messages require users executing these operations to be licensed with any of the Dynamics 365 licenses entitled to this operation. Check product documentation and the Dynamics 365 Licensing Guide for additional information and license entitlements. | Dynamics 365 SDK messages require users executing these operations to be licensed with any of the appropriate Dynamics 365 license. | [https://go.microsoft.com/fwlink/?linkid=2262812](https://go.microsoft.com/fwlink/?linkid=2262812)                                         |
				| meta-license-fieldservice-entity-operations | Some operations performed on Dynamics 365 entities require users executing these operations to be licensed with any of the Dynamics 365 licenses entitled to this operation. Check product documentation and the Dynamics 365 Licensing Guide for additional information and license entitlements. | Some operations performed on Dynamics 365 entities require users executing these operations to be licensed with any of the appropriate Dynamics 365 license. | [https://go.microsoft.com/fwlink/?linkid=2262812](https://go.microsoft.com/fwlink/?linkid=2262812)                                         |
				| meta-license-fieldservice-customcontrols | Some Dynamics 365 Field Service controls require users accessing these controls to be licensed with a Dynamics 365 Field Service license. Check product documentation for additional information. | Some Dynamics 365 Field Service controls require users accessing these controls to be licensed with a Dynamics 365 Field Service license. | [https://go.microsoft.com/fwlink/?linkid=2262812](https://go.microsoft.com/fwlink/?linkid=2262812)                                         |
				| meta-avoid-managed-entity-assets        | Do not add managed entities with all assets to a solution, as this will result in unexpected missing solution dependencies. | Managed entities should not be added to a solution with all assets included. | [https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import](https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import) |
				| meta-include-unmanaged-entity-assets    | Please add unmanaged entities with all assets to a solution, or this will result in sub-component loss due to segmentation. | Unmanaged entities should be added to a solution with all assets included. | [https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import](https://learn.microsoft.com/troubleshoot/power-platform/power-apps/solutions/missing-dependency-on-solution-import) |
				| connector-validate-hexadecimalbrandcolor | The brand color is invalid. Update the brand color with a valid hexadecimal color.                    | The brand color is invalid. Update the brand color with a valid hexadecimal color. | [https://go.microsoft.com/fwlink/?linkid=2244386](https://go.microsoft.com/fwlink/?linkid=2244386)                                         |
				`,
				Description: "List of rules to exclude from solution checker.  See [Solution Checker enforcement](https://learn.microsoft.com/power-platform/admin/managed-environment-solution-checker) for more details.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf(append([]string{""}, strings.Split(SOLUTION_CHECKER_RULES, ", ")...)...)),
				},
			},
			"maker_onboarding_markdown": schema.StringAttribute{
				MarkdownDescription: "First-time Power Apps makers will see this content in the Studio.  See [Maker welcome content](https://learn.microsoft.com/power-platform/admin/welcome-content) for more details.",
				Description:         "First-time Power Apps makers will see this content in the Studio",
				Required:            true,
			},
			"maker_onboarding_url": schema.StringAttribute{
				MarkdownDescription: "Maker onboarding 'Learn more' URL. See [Maker welcome content](https://learn.microsoft.com/power-platform/admin/welcome-content) for more details.",
				Description:         "Maker onboarding 'Learn more' URL",
				Required:            true,
			},
		},
	}
}

func (r *ManagedEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	clientApi := client.Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.ManagedEnvironmentClient = newManagedEnvironmentClient(clientApi)
}

func (r *ManagedEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var solutionCheckerRuleOverrides *string
	if !plan.SolutionCheckerRuleOverrides.IsNull() {
		value := strings.Join(helpers.SetToStringSlice(plan.SolutionCheckerRuleOverrides), ",")
		solutionCheckerRuleOverrides = &value
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard",
		Settings: &environment.SettingsDto{
			ExtendedSettings: environment.ExtendedSettingsDto{
				ExcludeEnvironmentFromAnalysis: strconv.FormatBool(plan.IsUsageInsightsDisabled.ValueBool()),
				IsGroupSharingDisabled:         strconv.FormatBool(plan.IsGroupSharingDisabled.ValueBool()),
				MaxLimitUserSharing:            strconv.FormatInt(plan.MaxLimitUserSharing.ValueInt64(), 10),
				DisableAiGeneratedDescriptions: "false",
				IncludeOnHomepageInsights:      "false",
				LimitSharingMode:               strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:],
				SolutionCheckerMode:            strings.ToLower(plan.SolutionCheckerMode.ValueString()),
				SuppressValidationEmails:       strconv.FormatBool(plan.SuppressValidationEmails.ValueBool()),
				SolutionCheckerRuleOverrides:   "",
				MakerOnboardingUrl:             plan.MakerOnboardingUrl.ValueString(),
				MakerOnboardingMarkdown:        plan.MakerOnboardingMarkdown.ValueString(),
			},
		},
	}

	if solutionCheckerRuleOverrides != nil {
		managedEnvironmentDto.Settings.ExtendedSettings.SolutionCheckerRuleOverrides = *solutionCheckerRuleOverrides
	}

	err := r.ManagedEnvironmentClient.EnableManagedEnvironment(ctx, managedEnvironmentDto, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
	plan.Id = plan.EnvironmentId
	plan.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)
	plan.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
	plan.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
	plan.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
	plan.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
	plan.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
	plan.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)

	if env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides == "" {
		plan.SolutionCheckerRuleOverrides = types.SetNull(types.StringType)
	} else {
		plan.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)

	if env.Properties.GovernanceConfiguration.Settings != nil {
		maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)

		state.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
		state.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
		state.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
		state.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
		state.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
		state.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
		state.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
		state.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)
		if env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides == "" {
			state.SolutionCheckerRuleOverrides = types.SetNull(types.StringType)
		} else {
			state.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))
		}
	} else {
		state.IsGroupSharingDisabled = types.BoolUnknown()
		state.IsUsageInsightsDisabled = types.BoolUnknown()
		state.MaxLimitUserSharing = types.Int64Unknown()
		state.LimitSharingMode = types.StringUnknown()
		state.SolutionCheckerMode = types.StringUnknown()
		state.SuppressValidationEmails = types.BoolUnknown()
		state.MakerOnboardingUrl = types.StringUnknown()
		state.MakerOnboardingMarkdown = types.StringUnknown()
		state.SolutionCheckerRuleOverrides = types.SetUnknown(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ManagedEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var solutionCheckerRuleOverrides *string
	if !plan.SolutionCheckerRuleOverrides.IsNull() {
		value := strings.Join(helpers.SetToStringSlice(plan.SolutionCheckerRuleOverrides), ",")
		solutionCheckerRuleOverrides = &value
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard",
		Settings: &environment.SettingsDto{
			ExtendedSettings: environment.ExtendedSettingsDto{
				ExcludeEnvironmentFromAnalysis: strconv.FormatBool(plan.IsUsageInsightsDisabled.ValueBool()),
				IsGroupSharingDisabled:         strconv.FormatBool(plan.IsGroupSharingDisabled.ValueBool()),
				MaxLimitUserSharing:            strconv.FormatInt(plan.MaxLimitUserSharing.ValueInt64(), 10),
				DisableAiGeneratedDescriptions: "false",
				IncludeOnHomepageInsights:      "false",
				LimitSharingMode:               strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:],
				SolutionCheckerMode:            strings.ToLower(plan.SolutionCheckerMode.ValueString()),
				SuppressValidationEmails:       strconv.FormatBool(plan.SuppressValidationEmails.ValueBool()),
				MakerOnboardingUrl:             plan.MakerOnboardingUrl.ValueString(),
				MakerOnboardingMarkdown:        plan.MakerOnboardingMarkdown.ValueString(),
				SolutionCheckerRuleOverrides:   "",
			},
		},
	}

	if solutionCheckerRuleOverrides != nil {
		managedEnvironmentDto.Settings.ExtendedSettings.SolutionCheckerRuleOverrides = *solutionCheckerRuleOverrides
	}

	err := r.ManagedEnvironmentClient.EnableManagedEnvironment(ctx, managedEnvironmentDto, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
	plan.Id = plan.EnvironmentId
	plan.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)
	plan.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
	plan.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
	plan.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
	plan.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
	plan.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
	plan.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)

	if env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides == "" {
		plan.SolutionCheckerRuleOverrides = types.SetNull(types.StringType)
	} else {
		plan.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ManagedEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ManagedEnvironmentClient.DisableManagedEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}

func (r *ManagedEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
