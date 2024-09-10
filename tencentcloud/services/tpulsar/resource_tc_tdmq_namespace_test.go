package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqNamespaceResource_basic -v
func TestAccTencentCloudTdmqNamespaceResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_namespace.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqNamespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "environ_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "msg_ttl", "60"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark."),
				),
			},
			{
				Config: testAccTdmqNamespaceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformId, "environ_name", "tf_example"),
					resource.TestCheckResourceAttr(terraformId, "msg_ttl", "300"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
					resource.TestCheckResourceAttr(terraformId, "remark", "remark update."),
				),
			},
			{
				ResourceName:      terraformId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqNamespace = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 60
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}
`

const testAccTdmqNamespaceUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 30
    size_in_mb      = 20
  }
  remark = "remark update."
}
`
