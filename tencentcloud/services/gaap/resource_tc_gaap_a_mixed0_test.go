package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// 测试基本的通道创建
func TestAccTencentCloudGaapProxyResourceMixd_0(t *testing.T) {
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: fixTestAccCheckGaapProxyDestroy(proxyId),
		Steps: []resource.TestStep{
			{
				Config: testGaapCrossGZtoSH_0_1,
				Check: resource.ComposeTestCheckFunc(
					//check ip rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "ip", "2.5.73.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),
				),
			},
			{
				Config: testGaapCrossGZtoSH_0_2,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "network_type", "normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "scalable", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testGaapCrossGZtoSH_0_3,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "network_type", "normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "scalable", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "client_ip_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "status", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "create_time"),

					//check udp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.udp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "port", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "status", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.udp_l4-1", "create_time"),
				),
			},
			{
				Config: testGaapCrossGZtoSH_0_4,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "network_type", "normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "scalable", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),

					//check ip rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "ip", "2.5.73.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "client_ip_method", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.0.ip", "2.5.73.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.0.port", "234"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.0.weight", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_bind_set.0.id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "status", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "create_time"),

					//check udp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.udp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "port", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.0.ip", "ah-tencent-zpl.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.0.port", "456"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.0.weight", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.udp_l4-1", "realserver_bind_set.0.id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "status", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.udp_l4-1", "create_time"),
				),
			},
		},
	})
}

func fixTestAccCheckGaapProxyDestroy(id *string) resource.TestCheckFunc {
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

// 创rs
const testGaapCrossGZtoSH_0_1 = `
resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-1"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-1"
}
`

// 创建proxy
const testGaapCrossGZtoSH_0_2 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-1"
bandwidth         = 10
concurrent        = 2
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}
`

// 创建L4+proxy
const testGaapCrossGZtoSH_0_3 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-1"
bandwidth         = 10
concurrent        = 2
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-1"
port            = 9090
realserver_type = "IP"
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-1"
port            = 8080
realserver_type = "DOMAIN"
}
`

// L4+proxy+rs
const testGaapCrossGZtoSH_0_4 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-1"
bandwidth         = 10
concurrent        = 2
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-1"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-1"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-1"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-1"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

// 场景
// Proxy     GZ---SH
// Listener  UDP TCP
// rs        ip domain
// L4+proxy+rs
const testGaapCrossGZtoSH_1_0 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-1"
bandwidth         = 10
concurrent        = 2
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-1"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-1"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-1"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-1"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}
`

////修改命名

const testGaapCrossGZtoSH_1_1 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 10
concurrent        = 2
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

// 修改通道规格
const testGaapCrossGZtoSH_1_2 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

////修改通道项目ID    新增时候，不一样项目ID的RS和proxy不可以绑定，但是可以修改

const testGaapCrossGZtoSH_1_3 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

////禁用通道

const testGaapCrossGZtoSH_1_4 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "false"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

////启用通道

const testGaapCrossGZtoSH_1_5 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

////修改L4健康配置

const testGaapCrossGZtoSH_1_6 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

// //L4健康配置
const testGaapCrossGZtoSH_2_0 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

// 关闭L4监控检查
const testGaapCrossGZtoSH_2_1 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "false"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}`

// //增加L4 数量
const testGaapCrossGZtoSH_2_3 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
}
`

// //增加L4 数量
const testGaapCrossGZtoSH_3_0 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
}
`

// 增加rs 数量
const testGaapCrossGZtoSH_3_1 = `{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-2.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-2.domain
port		= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-3.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-3.domain
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		= 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-2" {
domain			= "ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-2"
project_id		= 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-3" {
domain			= "ah-ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		  = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-2" {
ip				= "2.5.73.12"
name			= "tf_gaap_test_rs_ip-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-3" {
ip				= "2.5.73.13"
name			= "tf_gaap_test_rs_ip-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set.domain
port			= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-2.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-2.ip
port			= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-3.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-3.ip
port			= 234
}
}
`

// 增加rs 数量
const testGaapCrossGZtoSH_4_0 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		= 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-2" {
domain			= "ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-2"
project_id		= 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-3" {
domain			= "ah-ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		  = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-2" {
ip				= "2.5.73.12"
name			= "tf_gaap_test_rs_ip-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-3" {
ip				= "2.5.73.13"
name			= "tf_gaap_test_rs_ip-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-2.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-2.domain
port		= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-3.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-3.domain
port			= 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port			= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-2.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-2.ip
port			= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-3.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-3.ip
port			= 234
}
}
`

//设置通道为wrr 增加权重

const testGaapCrossGZtoSH_4_1 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-2" {
domain			= "ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-3" {
domain			= "ah-ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-2" {
ip				= "2.5.73.12"
name			= "tf_gaap_test_rs_ip-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-3" {
ip				= "2.5.73.13"
name			= "tf_gaap_test_rs_ip-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
scheduler = "wrr"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port		= 234
weight		= 10
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-2.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-2.domain
port			= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-3.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-3.domain
port	    = 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port		= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
scheduler = "wrr"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-2.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-2.ip
port		= 234
weight		= 40
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-3.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-3.ip
port		= 234
weight		= 50
}
}
`

//修改rs Ip/domain  --- 绑定状态的rs无法修改Ip/domain/projectID
//const testGaapCrossGZtoSH_4_x = `
//resource tencentcloud_gaap_proxy "foo" {
//name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
//bandwidth         = 20
//concurrent        = 10
//access_region     = "Guangzhou"
//realserver_region = "Shanghai"
//project_id		   = 1287704
//enable			   = "true"
//}
//
//resource tencentcloud_gaap_realserver "rs_domain-1" {
//domain			= "ah-tencent-zpl5.com"
//name			= "tf_gaap_test_rs_domain-2"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_realserver "rs_domain-1-2" {
//domain			= "ah-ah-tencent-zpl5.com"
//name			= "tf_gaap_test_rs_domain-2-2"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_realserver "rs_domain-1-3" {
//domain			= "ah-ah-ah-tencent-zpl5.com"
//name			= "tf_gaap_test_rs_domain-2-3"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_realserver "rs_ip-1" {
//ip				= "2.5.73.15"
//name			= "tf_gaap_test_rs_ip-2"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_realserver "rs_ip-1-2" {
//ip				= "2.5.73.125"
//name			= "tf_gaap_test_rs_ip-2-2"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_realserver "rs_ip-1-3" {
//ip				= "2.5.73.135"
//name			= "tf_gaap_test_rs_ip-2-3"
//project_id		   = 0
//}
//
//resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
//proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
//protocol        = "TCP"
//name            = "tf-ci-test-gaap-tcp-l4-2"
//port            = 9090
//health_check    = "true"
//interval		 = 10
//connect_timeout = 5
//realserver_type = "IP"
//scheduler = "wrr"
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_ip-1.id
//ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
//port		= 234
//weight		= 15
//}
//}
//
//resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
//proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
//protocol        = "TCP"
//name            = "tf-ci-test-gaap-tcp-l4-2-2"
//port            = 9091
//realserver_type = "DOMAIN"
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_domain-1-2.id
//ip			= tencentcloud_gaap_realserver.rs_domain-1-2.domain
//port			= 234
//}
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_domain-1-3.id
//ip			= tencentcloud_gaap_realserver.rs_domain-1-3.domain
//port	    = 234
//}
//}
//
//resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
//proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
//protocol        = "UDP"
//name            = "tf-ci-test-gaap-udp-l4-2"
//port            = 8080
//realserver_type = "DOMAIN"
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_domain-1.id
//ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
//port		= 456
//}
//}
//
//resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
//proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
//protocol        = "UDP"
//name            = "tf-ci-test-gaap-udp-l4-2-2"
//port            = 8081
//realserver_type = "IP"
//scheduler = "wrr"
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_ip-1-2.id
//ip			= tencentcloud_gaap_realserver.rs_ip-1-2.ip
//port		= 234
//weight		= 95
//}
//realserver_bind_set{
//id			= tencentcloud_gaap_realserver.rs_ip-1-3.id
//ip			= tencentcloud_gaap_realserver.rs_ip-1-3.ip
//port		= 234
//weight		= 55
//}
//}
//`

// 修改rs  名称 权重
const testGaapCrossGZtoSH_4_2 = `
resource tencentcloud_gaap_proxy "foo" {
name              = "tf-ci-test-gaap-proxy-GZ-SH-2"
bandwidth         = 20
concurrent        = 10
access_region     = "Guangzhou"
realserver_region = "Shanghai"
project_id		   = 1287704
enable			   = "true"
}

resource tencentcloud_gaap_realserver "rs_domain-1" {
domain			= "ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-a-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-2" {
domain			= "ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_domain-1-3" {
domain			= "ah-ah-ah-tencent-zpl.com"
name			= "tf_gaap_test_rs_domain-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1" {
ip				= "2.5.73.1"
name			= "tf_gaap_test_rs_ip-a-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-2" {
ip				= "2.5.73.12"
name			= "tf_gaap_test_rs_ip-2-2"
project_id		   = 1287704
}

resource tencentcloud_gaap_realserver "rs_ip-1-3" {
ip				= "2.5.73.13"
name			= "tf_gaap_test_rs_ip-2-3"
project_id		   = 1287704
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2"
port            = 9090
health_check    = "true"
interval		 = 10
connect_timeout = 5
realserver_type = "IP"
scheduler = "wrr"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1.id
ip			= tencentcloud_gaap_realserver.rs_ip-1.ip
port		= 234
weight		= 15
}
}

resource tencentcloud_gaap_layer4_listener "tcp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "TCP"
name            = "tf-ci-test-gaap-tcp-l4-2-2"
port            = 9091
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-2.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-2.domain
port			= 234
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1-3.id
ip			= tencentcloud_gaap_realserver.rs_domain-1-3.domain
port	    = 234
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2"
port            = 8080
realserver_type = "DOMAIN"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_domain-1.id
ip			= tencentcloud_gaap_realserver.rs_domain-1.domain
port		= 456
}
}

resource tencentcloud_gaap_layer4_listener "udp_l4-1-2" {
proxy_id        = tencentcloud_gaap_proxy.foo.proxy_id
protocol        = "UDP"
name            = "tf-ci-test-gaap-udp-l4-2-2"
port            = 8081
realserver_type = "IP"
scheduler = "wrr"
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-2.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-2.ip
port		= 234
weight		= 95
}
realserver_bind_set{
id			= tencentcloud_gaap_realserver.rs_ip-1-3.id
ip			= tencentcloud_gaap_realserver.rs_ip-1-3.ip
port		= 234
weight		= 55
}
}
`
