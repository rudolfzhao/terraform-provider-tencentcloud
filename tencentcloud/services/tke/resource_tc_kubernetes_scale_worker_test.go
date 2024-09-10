package tke_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var testTkeScaleWorkerResourceName = "tencentcloud_kubernetes_scale_worker"
var testTkeScaleWorkerResourceKey = testTkeScaleWorkerResourceName + ".test_scale"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_kubernetes_scale_worker
	resource.AddTestSweepers("tencentcloud_kubernetes_scale_worker", &resource.Sweeper{
		Name: "tencentcloud_kubernetes_scale_worker",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctke.NewTkeService(client)

			clusters, err := service.DescribeClusters(ctx, "", tcacctest.DefaultTkeClusterName)

			if err != nil {
				return err
			}

			if len(clusters) == 0 {
				log.Println("no cluster names " + tcacctest.DefaultTkeClusterName)
				return nil
			}

			clusterId := clusters[0].ClusterId

			_, workers, err := service.DescribeClusterInstances(ctx, clusterId)

			if err != nil {
				return err
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance

			cvmService := svccvm.NewCvmService(client)
			instanceIds := make([]string, 0)
			for i := range workers {
				worker := workers[i]
				if worker.NodePoolId != "" {
					continue
				}
				instance, err := cvmService.DescribeInstanceById(ctx, worker.InstanceId)
				if err != nil {
					continue
				}

				if !tccommon.CheckResourcePersist(*instance.InstanceName, worker.CreatedTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   worker.InstanceId,
						Name: *instance.InstanceName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         worker.InstanceId,
					Name:       *instance.InstanceName,
					CreateTime: worker.CreatedTime,
				})

				created, err := time.Parse(tccommon.TENCENTCLOUD_COMMON_TIME_LAYOUT, worker.CreatedTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(*instance.InstanceName, &created) {
					continue
				}
				instanceIds = append(instanceIds, worker.InstanceId)
			}

			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateClusterInstances")

			if len(instanceIds) == 0 {
				return nil
			}
			err = service.DeleteClusterInstances(ctx, clusterId, instanceIds)
			if err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudKubernetesScaleWorkerResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTkeScaleWorkerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeScaleWorkerInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeScaleWorkerExists(testTkeScaleWorkerResourceKey),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "data_disk.0.file_system", "xfs"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "data_disk.0.mount_target", "/data1"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "data_disk.0.auto_format_and_mount", "true"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_config.#"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "worker_config.0.data_disk.0.encrypt", "true"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "worker_instances_list.#", "1"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "worker_instances_list.0.instance_role"),
					resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "unschedulable"),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "pre_start_user_script", "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="),
					resource.TestCheckResourceAttr(testTkeScaleWorkerResourceKey, "user_script", "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="),
				),
			},
			// gpu_args依赖于指定镜像ID，但账号没有镜像ID权限，暂时注释
			//{
			//	Config: testAccTkeScaleWorkerInstanceGpuInsTypeUpdate,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckTkeScaleWorkerExists(testTkeScaleWorkerResourceKey),
			//		resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "cluster_id"),
			//		resource.TestCheckResourceAttrSet(testTkeScaleWorkerResourceKey, "gpu_args.#"),
			//	),
			//},
		},
	})
}

func testAccCheckTkeScaleWorkerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTkeScaleWorkerResourceName {
			continue
		}
		instanceId := rs.Primary.Attributes["worker_instances_list.0.instance_id"]
		clusterId := rs.Primary.Attributes["cluster_id"]

		if clusterId == "" || instanceId == "" {
			return fmt.Errorf("miss worker_instances_list.0.instance_id[%s] or cluster_id[%s]", instanceId, clusterId)
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctke.NewTkeService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InvalidParameter.ClusterNotFound" {
						return nil
					}
				}
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}

		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				return fmt.Errorf("cvm %s found in DescribeClusterInstances", instanceId)
			}
		}
		log.Printf("[DEBUG]instance %s delelte ok", instanceId)

	}
	return nil
}

func testAccCheckTkeScaleWorkerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("tke worker scale instance %s is not found", n)
		}
		instanceId := rs.Primary.Attributes["worker_instances_list.0.instance_id"]
		clusterId := rs.Primary.Attributes["cluster_id"]

		if clusterId == "" || instanceId == "" {
			return fmt.Errorf("miss worker_instances_list.0.instance_id[%s] or cluster_id[%s]", instanceId, clusterId)
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svctke.NewTkeService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() == "InvalidParameter.ClusterNotFound" {
						return nil
					}
				}

				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}

		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				log.Printf("[DEBUG]instance %s create ok", instanceId)
				return nil
			}
		}
		return fmt.Errorf("cvm %s not found in DescribeClusterInstances", instanceId)
	}
}

const testAccTkeScaleWorkerInstance = testAccTkeCluster + `
resource "tencentcloud_kubernetes_scale_worker" "test_scale" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  extra_args = [
    "root-dir=/var/lib/kubelet"
  ]

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  unschedulable = 0
  pre_start_user_script   = "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="
  user_script   = "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="

  data_disk {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
    file_system           = "xfs"
    mount_target          = "/data1"
    auto_format_and_mount = true
  }

  worker_config {
    count                      				= 1
    availability_zone          				= "ap-guangzhou-3"
    instance_type              				= "S2.LARGE16"
    subnet_id                  				= local.subnet_id1
    system_disk_type           				= "CLOUD_SSD"
    system_disk_size           				= 50
    internet_charge_type       				= "TRAFFIC_POSTPAID_BY_HOUR"
    security_group_ids                      = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
	  encrypt = true
    }

    enhanced_security_service 				= false
    enhanced_monitor_service  				= false
    user_data                 				= "dGVzdA=="
    password                  				= "AABBccdd1122"
  }
}
`

const testAccTkeScaleWorkerInstanceGpuInsTypeUpdate = testAccTkeCluster + `
resource "tencentcloud_kubernetes_scale_worker" "test_scale" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id

  extra_args = [
    "root-dir=/var/lib/kubelet"
  ]

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  unschedulable = 0
  pre_start_user_script   = "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="
  user_script   = "IyEvYmluL3NoIGVjaG8gImhlbGxvIHdvcmxkIg=="

  worker_config {
    count                = 1
    availability_zone    = "ap-guangzhou-3"
    instance_type        = "GN6S.LARGE20"
    subnet_id            = local.subnet_id1
    system_disk_type     = "CLOUD_SSD"
    system_disk_size     = 50
    internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
    security_group_ids   = [local.sg_id]
    img_id               = "img-oyd1zdra"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "AABBccdd1122"
  }

  gpu_args {
    mig_enable = false
    driver     = {
      name    = "NVIDIA-Linux-x86_64-470.182.03.run"
      version = "470.182.03"
    }
    cuda = {
      name    = "cuda_11.4.3_470.82.01_linux.run"
      version = "11.4.3"
    }
    cudnn = {
      name    = "cudnn-11.4-linux-x64-v8.2.4.15.tgz"
      version = "8.2.4"
    }
  }
}
`
