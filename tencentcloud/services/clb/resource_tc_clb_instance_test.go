package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const BasicClbName = "tf-clb-basic"
const SnatClbName = "tf-clb-snat"
const InternalClbName = "tf-clb-internal"
const InternalClbNameUpdate = "tf-clb-update-internal"
const SingleClbName = "single-open-clb"
const MultiClbName = "multi-open-clb"
const OpenClbName = "tf-clb-open"
const OpenClbNameIpv6 = "tf-clb-open-ipv6"
const OpenClbNameUpdate = "tf-clb-update-open"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_clb_instance
	resource.AddTestSweepers("tencentcloud_clb_instance", &resource.Sweeper{
		Name: "tencentcloud_clb_instance",
		F:    testSweepClbInstance,
	})
}

func testSweepClbInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return err
	}
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := localclb.NewClbService(client)

	res, err := service.DescribeLoadBalancerByFilter(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range res {
		if !tccommon.CheckResourcePersist(*v.LoadBalancerName, *v.CreateTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.LoadBalancerId,
				Name: *v.LoadBalancerName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.LoadBalancerId,
			Name:       *v.LoadBalancerName,
			CreateTime: *v.CreateTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateLoadBalancer")

	if len(res) > 0 {
		for _, v := range res {
			id := *v.LoadBalancerId
			//instanceName := *v.LoadBalancerName
			createTime := tccommon.StringToTime(*v.CreateTime)

			now := time.Now()
			interval := now.Sub(createTime).Minutes()
			// keep not delete
			//if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
			//	continue
			//}
			// less than 30 minute, not delete
			if tccommon.NeedProtect == 1 && int64(interval) < 30 {
				continue
			}
			if err := service.DeleteLoadBalancerById(ctx, id); err != nil {
				log.Printf("Delete %s error: %s", id, err.Error())
				continue
			}
		}
	}
	return nil
}

func TestAccTencentCloudClbInstanceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "clb_name", BasicClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "tags.test", "tf"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "tags.test1", "tf1"),
				),
			},
			{
				ResourceName:            "tencentcloud_clb_instance.clb_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dynamic_vip"},
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_open(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", OpenClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info_vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "tags.test", "tf"),
				),
			},
			{
				Config: testAccClbInstance_update_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", OpenClbNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info_vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "tags.test", "test"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_openIpv6(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_openIpv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open_ipv6"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "clb_name", OpenClbNameIpv6),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "vpc_id", "vpc-mvhjjprd"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "subnet_id", "subnet-2qfyfvv8"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open_ipv6", "address_ip_version", "IPv6FullChain"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_snat(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_snat,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "clb_name", SnatClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "snat_pro", "true"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_internal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_internal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "clb_name", InternalClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "network_type", "INTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "tags.test", "tf1"),
				),
			},
			{
				Config: testAccClbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "clb_name", InternalClbNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "network_type", "INTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "tags.test", "test"),
				),
			},
			{
				Config: testAccClbInstance_updateRecover,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "clb_name", InternalClbNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "network_type", "INTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "delete_protect"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "tags.test", "test"),
				),
			},
			{
				ResourceName:            "tencentcloud_clb_instance.clb_internal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dynamic_vip", "delete_protect"},
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_default_enable(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_default_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.default_enable"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "clb_name", SingleClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "load_balancer_pass_to_target", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "target_region_info_vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "tags.test", "open"),
				),
			},
			{
				Config: testAccClbInstance_default_enable_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.default_enable"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "clb_name", SingleClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "load_balancer_pass_to_target", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.default_enable", "target_region_info_vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.default_enable", "tags.test", "hello"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceResource_multiple_instance(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance__multi_instance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.multiple_instance"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "clb_name", MultiClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "master_zone_id", "100004"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "slave_zone_id", "100003"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "tags.test", "mytest"),
				),
			},
			{
				Config: testAccClbInstance__multi_instance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.multiple_instance"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "clb_name", MultiClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "master_zone_id", "100004"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "slave_zone_id", "100003"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "tags.test", "open"),
				),
			},
		},
	})
}

func testAccCheckClbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_instance" {
			continue
		}

		instance, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if instance != nil && err == nil {
			return fmt.Errorf("[CHECK][CLB instance][Destroy] check: CLB instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB instance][Exists] check: CLB instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB instance][Exists] check: CLB instance id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLB instance][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbInstance_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "` + BasicClbName + `"
  tags = {
    test = "tf"
    test1 = "tf1"
  }
}
`

const testAccClbInstance_snat = `
data "tencentcloud_vpc_instances" "gz3vpc" {
  name = "Default-"
  is_default = true
}

data "tencentcloud_vpc_subnets" "gz3" {
  vpc_id = data.tencentcloud_vpc_instances.gz3vpc.instance_list.0.vpc_id
}

locals {
  keep_clb_subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") if lookup(subnet, "name") == "keep-clb-sub"]
  subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") ]
  subnet_for_clb_snat = concat(local.keep_clb_subnets, local.subnets)
}

resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "` + SnatClbName + `"
  snat_pro     = true
  snat_ips {
	subnet_id = local.subnet_for_clb_snat.0
  }
  snat_ips {
    subnet_id = local.subnet_for_clb_snat.1
  }
}
`

const testAccClbInstance_internal = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-internal-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua-ci-temp-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_clb_instance" "clb_internal" {
  network_type = "INTERNAL"
  clb_name     = "` + InternalClbName + `"
  vpc_id       = tencentcloud_vpc.foo.id
  subnet_id    = tencentcloud_subnet.subnet.id
  project_id   = 0

  tags = {
    test = "tf1"
  }
}
`

const testAccClbInstance_open = `
resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "` + OpenClbName + `"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = ["sg-if748odn"]

  tags = {
    test = "tf"
  }
}
`

const testAccClbInstance_update = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-internal-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "tf-example-subnet-inc"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_clb_instance" "clb_internal" {
  network_type = "INTERNAL"
  clb_name     = "` + InternalClbNameUpdate + `"
  vpc_id       = tencentcloud_vpc.foo.id
  subnet_id    = tencentcloud_subnet.subnet.id
  project_id   = 0
  delete_protect = true
  tags = {
    test = "test"
  }
}
`
const testAccClbInstance_updateRecover = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-internal-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "tf-example-subnet-inc"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_clb_instance" "clb_internal" {
  network_type = "INTERNAL"
  clb_name     = "` + InternalClbNameUpdate + `"
  vpc_id       = tencentcloud_vpc.foo.id
  subnet_id    = tencentcloud_subnet.subnet.id
  project_id   = 0
  delete_protect = false
  tags = {
    test = "test"
  }
}
`
const testAccClbInstance_update_open = `

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "` + OpenClbNameUpdate + `"
  vpc_id                    = tencentcloud_vpc.foo.id
  project_id                = 0
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = ["sg-if748odn"]

  tags = {
    test = "test"
  }
}
`

const testAccClbInstance_default_enable = `
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "keep-sdk-feature-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_vpc" "foo" {
  name         = "clb-instance-default-vpc"
  cidr_block   = "10.0.0.0/16"

  tags = {
    "test" = "mytest"
  }
}

resource "tencentcloud_clb_instance" "default_enable" {
  network_type                 = "OPEN"
  clb_name                     = "` + SingleClbName + `"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.foo.id
  load_balancer_pass_to_target = true

  security_groups              = ["sg-if748odn"]
  target_region_info_region    = "ap-guangzhou"
  target_region_info_vpc_id    = tencentcloud_vpc.foo.id

  tags = {
    test = "open"
  }
}
`

const testAccClbInstance_default_enable_open = `
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "keep-sdk-feature-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_vpc" "foo" {
  name         = "clb-instance-default-vpc"
  cidr_block   = "10.0.0.0/16"

  tags = {
    "test" = "mytest"
  }
}

resource "tencentcloud_clb_instance" "default_enable" {
  network_type                 = "OPEN"
  clb_name                     = "` + SingleClbName + `"
  project_id                   = 0
  vpc_id                       = tencentcloud_vpc.foo.id
  load_balancer_pass_to_target = true

  security_groups              = ["sg-if748odn"]
  target_region_info_region    = "ap-guangzhou"
  target_region_info_vpc_id    = tencentcloud_vpc.foo.id

  tags = {
    test = "hello"
  }
}
`

const testAccClbInstance__multi_instance = `
resource "tencentcloud_clb_instance" "multiple_instance" {
  network_type              = "OPEN"
  clb_name                  = "` + MultiClbName + `"
  master_zone_id = "100004"
  slave_zone_id = "100003"

  tags = {
    test = "mytest"
  }
}
`

const testAccClbInstance__multi_instance_update = `
resource "tencentcloud_clb_instance" "multiple_instance" {
  network_type              = "OPEN"
  clb_name                  = "` + MultiClbName + `"
  master_zone_id = "100004"
  slave_zone_id = "100003"

  tags = {
    test = "open"
  }
}
`

const testAccClbInstance_openIpv6 = `
resource "tencentcloud_clb_instance" "clb_open_ipv6" {
	clb_name           = "` + OpenClbNameIpv6 + `"
	network_type       = "OPEN"
	project_id         = 0
	vpc_id             = "vpc-mvhjjprd"
	subnet_id          = "subnet-2qfyfvv8"
	address_ip_version = "IPv6FullChain"
}
`
