package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudGaapProxyResource_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Thailand"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Jakarta"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_proxy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapProxyResource_update(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-update"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Thailand"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Jakarta"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-new"),
				),
			},
			{
				Config: testAccGaapProxyNewBandwidth,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
				),
			},
			{
				Config: testAccGaapProxyNewConcurrent,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
				),
			},
			{
				Config: testAccGaapProxyUpdateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "tags.test", "test"),
				),
			},
			{
				Config: testAccGaapProxyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "false"),
				),
			},
		},
	})
}

func testAccCheckGaapProxyDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

		proxies, err := service.DescribeProxies(context.TODO(), []string{*id}, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if len(proxies) != 0 {
			return fmt.Errorf("proxy still exists")
		}

		return nil
	}
}

func testAccCheckGaapProxyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no proxy ID is set")
		}

		service := svcgaap.NewGaapService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		proxies, err := service.DescribeProxies(context.TODO(), []string{rs.Primary.ID}, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if len(proxies) == 0 {
			return fmt.Errorf("proxy not found: %s", rs.Primary.ID)
		}

		for _, proxy := range proxies {
			if proxy.ProxyId == nil {
				return errors.New("realserver id is nil")
			}
			if *proxy.ProxyId == rs.Primary.ID {
				*id = rs.Primary.ID
				break
			}
		}

		if *id == "" {
			return fmt.Errorf("proxy not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

const testAccGaapProxyBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Thailand"
  realserver_region = "Jakarta"
  network_type = "normal"
}
`

const testAccGaapProxyUpdateBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-update"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Thailand"
  realserver_region = "Jakarta"
}
`

const testAccGaapProxyNewName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Thailand"
  realserver_region = "Jakarta"
}
`

const testAccGaapProxyNewBandwidth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 2
  access_region     = "Thailand"
  realserver_region = "Jakarta"
}
`

const testAccGaapProxyNewConcurrent = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 10
  access_region     = "Thailand"
  realserver_region = "Jakarta"
}
`

const testAccGaapProxyDisable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 10
  access_region     = "Thailand"
  realserver_region = "Jakarta"
  enable            = false

  tags = {
    "test" = "test"
  }
}
`

const testAccGaapProxyUpdateTags = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 10
  access_region     = "Thailand"
  realserver_region = "Jakarta"
  enable            = false

  tags = {
    "test" = "test"
  }
}
`
