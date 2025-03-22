// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "xata_workspace" "markspace" {
  name = "markspace"
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify created workspace has Computed attributes filled.
                    resource.TestCheckResourceAttr("xata_workspace.markspace", "Name", "markspace"),
                    resource.TestCheckResourceAttr("xata_workspace.markspace", "Slug", ""),
                    resource.TestCheckResourceAttr("xata_workspace.markspace", "Id", "10"),
                    resource.TestCheckResourceAttr("xata_workspace.markspace", "MemberCount","2"),
                    resource.TestCheckResourceAttr("xata_workspace.markspace", "Plan", "free"),
                ),
            },
            // ImportState testing
            {
                ResourceName:      "xata_workspace.markspace",
                ImportState:       true,
                ImportStateVerify: true,
                // The last_updated attribute does not exist in the Xata
                // API, therefore there is no value for it during import.
                ImportStateVerifyIgnore: []string{"last_updated"},
            },
            // Update and Read testing
            {
                Config: providerConfig + `
resource "xata_workspace" "markspace" {
  name = "narkspace"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify created workspace has Computed attributes filled.
					resource.TestCheckResourceAttr("xata_workspace.markspace", "Name", "narkspace"),
					resource.TestCheckResourceAttr("xata_workspace.markspace", "Slug", ""),
					resource.TestCheckResourceAttr("xata_workspace.markspace", "Id", "10"),
					resource.TestCheckResourceAttr("xata_workspace.markspace", "MemberCount","2"),
					resource.TestCheckResourceAttr("xata_workspace.markspace", "Plan", "free"),
				),
            },
            // Delete testing automatically occurs in TestCase
        },
    })
}