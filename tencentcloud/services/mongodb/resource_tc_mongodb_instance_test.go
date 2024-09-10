package mongodb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmongodb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mongodb_instance
	resource.AddTestSweepers("tencentcloud_mongodb_instance", &resource.Sweeper{
		Name: "tencentcloud_mongodb_instance",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svcmongodb.NewMongodbService(client)

			instances, err := service.DescribeInstancesByFilter(ctx, "", -1)
			if err != nil {
				return err
			}

			var isolated []string

			for i := range instances {
				ins := instances[i]
				id := *ins.InstanceId
				name := *ins.InstanceName

				if strings.HasPrefix(name, tcacctest.KeepResource) || strings.HasPrefix(name, tcacctest.DefaultResource) {
					continue
				}

				created, err := time.Parse("2006-01-02 15:04:05", *ins.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
					continue
				}
				log.Printf("%s (%s) will Isolated", id, name)
				err = service.IsolateInstance(ctx, id)
				if err != nil {
					continue
				}
				isolated = append(isolated, id)
			}

			log.Printf("Offline isolated instance %v", isolated)
			for _, id := range isolated {
				err = service.OfflineIsolatedDBInstance(ctx, id, true)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudMongodbInstanceResource_PostPaid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "memory"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "volume"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "prepaid_period"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "security_groups.0", "sg-if748odn"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag"},
			},
			{
				Config: testAccMongodbInstance_updateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "512"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccMongodbInstance_updateNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "512"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "node_num", "5"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccMongodbInstance_updateSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "security_groups.0", "sg-05f7wnhn"),
				),
			},
			{
				Config: testAccMongodbInstance_updateMaintenance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "maintenance_start", "05:00"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "maintenance_end", "06:00"),
				),
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceResource_MultiZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstance_multiZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb_mutil_zone"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "node_num", "5"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "availability_zone_list.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "hidden_zone", "ap-guangzhou-6"),
				),
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceResource_Prepaid(t *testing.T) {
	// Avoid to set Parallel to make sure EnvVar secure
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstancePrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "instance_name", "tf-mongodb-test-prepaid"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "memory"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "volume"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.test", "test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "auto_renew_flag", "1"),
				),
			},
			{
				Config: testAccMongodbInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "instance_name", "tf-mongodb-test-prepaid-update"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
			{
				Config: testAccMongodbInstancePrepaid_updateMaintenance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "maintenance_start", "05:00"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "maintenance_end", "06:00"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb_prepaid",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag", "prepaid_period"},
			},
		},
	})
}

func testAccCheckMongodbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_instance" {
			continue
		}

		_, has, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("mongodb instance still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckMongodbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mongodb instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance id is not set")
		}
		mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("mongodb instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccMongodbInstance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    test = "test"
  }
}
`

const testAccMongodbInstance_updateConfig = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateNode = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateSecurityGroup = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = ["sg-05f7wnhn"]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateMaintenance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = ["sg-05f7wnhn"]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
  maintenance_start = "05:00"
  maintenance_end = "06:00"
}
`

const testAccMongodbInstancePrepaid = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    test = "test-prepaid"
  }
}
`

const testAccMongodbInstancePrepaid_update = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid-update"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    prepaid = "prepaid"
  }
}
`

const testAccMongodbInstancePrepaid_updateMaintenance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid-update"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    prepaid = "prepaid"
  }
  maintenance_start = "05:00"
  maintenance_end = "06:00"
}
`

const testAccMongodbInstance_multiZone = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-multi-zone-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-multi-zone-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_mutil_zone" {
  instance_name   = "mongodb-mutil-zone-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  node_num = 5
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
  tags = {
    test = "test"
  }
}
`
