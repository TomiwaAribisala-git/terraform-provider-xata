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
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.#", "4"),
					// Verify the first workspace to ensure all attributes are set
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.id", "Tomiwa-Aribisala-s-workspace-tameub"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.name", "Tomiwa Aribisala's workspace"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.slug", "Tomiwa-Aribisala-s-workspace"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.role", "owner"),
					resource.TestCheckResourceAttr("data.xata_workspaces.test", "workspaces.0.plan", "free"),
				),
			},
		},
	})
}
