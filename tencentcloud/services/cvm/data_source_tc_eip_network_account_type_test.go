package cvm_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEipNetworkAccountTypeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipNetworkAccountTypeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eip_network_account_type.network_account_type"),
					resource.TestCheckResourceAttrWith("data.tencentcloud_eip_network_account_type.network_account_type", "network_account_type", func(value string) error {
						switch value {
						case "STANDARD", "LEGACY":
						default:
							return fmt.Errorf("invalid network type: %q", value)
						}
						return nil
					}),
				),
			},
		},
	})
}

const testAccEipNetworkAccountTypeDataSource = `

data "tencentcloud_eip_network_account_type" "network_account_type" {}

`
