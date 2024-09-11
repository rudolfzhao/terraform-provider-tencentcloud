package gaap_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// 测试基本的通道创建
func TestAccTencentCloudGaapProxyResourceMixd_1(t *testing.T) {
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroyMix(proxyId),
		Steps: []resource.TestStep{
			{
				Config: testGaapCrossGZtoSH_1_0,
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
			{
				Config: testGaapCrossGZtoSH_1_1,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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
			{
				Config: testGaapCrossGZtoSH_1_2,
				Check: resource.ComposeTestCheckFunc(

					//check ip rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "ip", "2.5.73.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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

					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
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
				Config: testGaapCrossGZtoSH_1_3,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "1287704"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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
			{
				Config: testGaapCrossGZtoSH_1_4,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "1287704"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "false"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "network_type", "normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "status", "CLOSED"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "scalable", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),

					//check ip rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "ip", "2.5.73.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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
			{
				Config: testGaapCrossGZtoSH_1_5,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "1287704"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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
			{
				Config: testGaapCrossGZtoSH_1_6,
				Check: resource.ComposeTestCheckFunc(
					//check proxy
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", proxyId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "tf-ci-test-gaap-proxy-GZ-SH-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "1287704"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "name", "tf_gaap_test_rs_ip-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_ip-1", "project_id", "1287704"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_ip-1", "id"),

					//check main rs
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "domain", "ah-tencent-zpl.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "name", "tf_gaap_test_rs_domain-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.rs_domain-1", "project_id", "1287704"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_realserver.rs_domain-1", "id"),

					//check tcp l4
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.tcp_l4-1", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "name", "tf-ci-test-gaap-tcp-l4-2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "interval", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.tcp_l4-1", "connect_timeout", "5"),
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
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.udp_l4-1", "name", "tf-ci-test-gaap-udp-l4-2"),
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
