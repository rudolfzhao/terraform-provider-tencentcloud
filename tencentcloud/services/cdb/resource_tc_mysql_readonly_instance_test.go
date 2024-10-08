package cdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlReadonlyInstanceResource_basic -v
func TestAccTencentCloudMysqlReadonlyInstanceResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "ro_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.test", "test-tf"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_readonly_instance.mysql_readonly",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(tcacctest.CommonPresetMysql, "read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "read"),
				),
			},
			// update tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(tcacctest.CommonPresetMysql, "readonly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "readonly"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlReadonlyInstance(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role"),
				),
			},
			// update instance_name
			{
				Config: testAccMysqlReadonlyInstance_update(tcacctest.CommonPresetMysql, "mysql-readonly-update", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-update"),
				),
			},
			// update mem_size
			{
				Config: testAccMysqlReadonlyInstance_memSize(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
				),
			},
			{
				Config: testAccMysqlReadonlyInstance_roGroup(tcacctest.CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly_ro_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly_ro_group", "ro_group_id"),
				),
			},
			// // update intranet_port
			// {
			// 	Config: testAccMysqlReadonlyInstance_update(CommonPresetMysql, "mysql-readonly-update", "3361"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
			// 		resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3361"),
			// 	),
			// },
		},
	})
}

func testAccCheckMysqlReadonlyInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_readonly_instance" {
			continue
		}
		instance, err := mysqlService.DescribeRunningDBInstanceById(ctx, rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("mysql instance still exist")
		}
		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
	}
	return nil
}

func testAccCheckMysqlInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mysql instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mysql instance id is not set")
		}

		mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := mysqlService.DescribeDBInstanceById(ctx, rs.Primary.ID)
		if instance == nil {
			return fmt.Errorf("mysql instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccMysqlReadonlyInstance(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_multiTags(mysqlTestCase, value string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
    role = "%s"
  }
}
	`, mysqlTestCase, value)
}

func testAccMysqlReadonlyInstance_update(mysqlTestCase, instance_name, instranet_port string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "%s"
  intranet_port      = %s 
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase, instance_name, instranet_port)
}

func testAccMysqlReadonlyInstance_memSize(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 1000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_roGroup(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly_ro_group" {
  master_instance_id = local.mysql_id
  ro_group_id 		 = tencentcloud_mysql_readonly_instance.mysql_readonly.ro_group_id
  mem_size           = 1000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}
