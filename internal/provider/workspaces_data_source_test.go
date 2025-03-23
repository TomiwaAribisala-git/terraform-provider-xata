// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspacesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "xata_workspaces" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of workspaces returned
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.#", "9"),
					// Verify the first workspace to ensure all attributes are set
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.Id", "1"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.Name", "markspace"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.Slug", ""),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.Role", "owner"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.Plan", "free"),
				),
			},
		},
	})
}
