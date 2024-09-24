// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccTenantApplicationPackagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_tenant_application_packages" "all_applications" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_visibility", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.catalog_visibility", regexp.MustCompile("^(AdminCenter|All|None|Teams)$")),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.localized_description", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.localized_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.unique_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.version", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.publisher_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.publisher_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.learn_more_url", regexp.MustCompile(helpers.UrlValidStringRegex)),
				),
			},
		},
	})
}

func TestUnitTenantApplicationPackagesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/appmanagement/applicationPackages?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_tenant_application_packages" "all_applications" {
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.#", "245"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_id", "b4b3b295-b5fe-4888-9355-9601c30626b3"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_name", "Microsoft Flow Approvals"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_visibility", "All"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.catalog_visibility", "None"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.localized_description", "Quickly create an approval workflow for your data in SharePoint, Dynamics CRM, OneDrive, and more. As part of the approval process, approvers will get email notifications about pending approvals, and can view and respond to all approval requests in one consolidated approvals center."),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.localized_name", "Microsoft Flow Approvals"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.unique_name", "msdyn_FlowApprovals"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.application_descprition", "An easier way to get and manage approvals."),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.publisher_id", "7ea093b9-252b-4557-a8ae-ad4c1932a412"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.publisher_name", "Microsoft Dynamics 365"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.learn_more_url", "https://go.microsoft.com/fwlink/p/?linkid=847067"),

					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.error_code", "0x80040216"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.error_name", "ApplicationNotVisible"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.message", "The solution is not visible in the catalog."),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.source", "Dataverse"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.status_code", "404"),
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.0.last_error.0.type", "ApplicationNotVisibleException"),
				),
			},
		},
	})
}

func TestUnitTenantApplicationPackagesDataSource_Validate_Filter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/appmanagement/applicationPackages?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/tenant_application_packages/Validate_Read_Filter/get_tenant_applications.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_tenant_application_packages" "all_applications" {
					publisher_name = "Microsoft Dynamics SMB"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.#", "56"),
				),
			},
			{
				Config: `
				data "powerplatform_tenant_application_packages" "all_applications" {
					name = "Healthcare Home Health"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_tenant_application_packages.all_applications", "applications.#", "1"),
				),
			},
		},
	})
}
