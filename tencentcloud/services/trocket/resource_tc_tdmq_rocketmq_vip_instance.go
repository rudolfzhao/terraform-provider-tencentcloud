package trocket

import (
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqRocketmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRocketmqVipInstanceCreate,
		Read:   resourceTencentCloudTdmqRocketmqVipInstanceRead,
		Update: resourceTencentCloudTdmqRocketmqVipInstanceUpdate,
		Delete: resourceTencentCloudTdmqRocketmqVipInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},
			"spec": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance specification: Universal type, rocket-vip-basic-0, Basic type: `rocket-vip-basic-1`, Standard type: `rocket-vip-basic-2`, Advanced Type I: `rocket-vip-basic-3`, Advanced Type II: `rocket-vip-basic-4`.",
			},
			"node_count": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 20),
				Description:  "Number of nodes, minimum 2, maximum 20.",
			},
			"storage_size": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerMin(200),
				Description:  "Single node storage space, in GB, minimum 200GB.",
			},
			"zone_ids": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The Zone ID list for node deployment, such as Guangzhou Zone 1, is 100001. For details, please refer to the official website of Tencent Cloud.",
			},
			"vpc_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "VPC information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet ID.",
						},
					},
				},
			},
			"time_span": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Purchase period, in months.",
			},
			"ip_rules": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Public IP access control rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_rule": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address block information.",
						},
						"allow": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to allow or deny.",
						},
						"remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remark.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqVipInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request   = tdmq.NewCreateRocketMQVipInstanceRequest()
		response  = tdmq.NewCreateRocketMQVipInstanceResponse()
		clusterId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec"); ok {
		request.Spec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("node_count"); ok {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		request.StorageSize = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		for i := range zoneIdsSet {
			zoneIds := zoneIdsSet[i].(string)
			request.ZoneIds = append(request.ZoneIds, &zoneIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "vpc_info"); ok {
		vpcInfo := tdmq.VpcInfo{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcInfo.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			vpcInfo.SubnetId = helper.String(v.(string))
		}
		request.VpcInfo = &vpcInfo
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ip_rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ipRule := tdmq.PublicAccessRule{}
			if v, ok := dMap["ip_rule"]; ok {
				ipRule.IpRule = helper.String(v.(string))
			}

			if v, ok := dMap["allow"]; ok {
				ipRule.Allow = helper.Bool(v.(bool))
			}

			if v, ok := dMap["remark"]; ok {
				ipRule.Remark = helper.String(v.(string))
			}

			request.IpRules = append(request.IpRules, &ipRule)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateRocketMQVipInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId

	// wait
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRocketmqVipInstancesByFilter(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *result.Status == svctdmq.RocketMqVipInsSuccess {
			return nil
		} else if *result.Status == svctdmq.RocketMqVipInsRunning {
			return resource.RetryableError(fmt.Errorf("tdmq rocketmqVipInstance status is running"))
		} else {
			e = fmt.Errorf("tdmq rocketmqVipInstance status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)
	return resourceTencentCloudTdmqRocketmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		clusterId = d.Id()
	)

	rocketmqVipInstanceDetail, err := service.DescribeTdmqRocketmqVipInstanceById(ctx, clusterId)
	if err != nil {
		return err
	}

	if rocketmqVipInstanceDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	rocketmqVipInstances, err := service.DescribeTdmqRocketmqVipInstancesByFilter(ctx, clusterId)
	if err != nil {
		return err
	}

	if rocketmqVipInstances == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRocketmqVipInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rocketmqVipInstanceDetail.ClusterInfo.ClusterName != nil {
		_ = d.Set("name", rocketmqVipInstanceDetail.ClusterInfo.ClusterName)
	}

	if rocketmqVipInstanceDetail.InstanceConfig.NodeCount != nil {
		_ = d.Set("node_count", rocketmqVipInstanceDetail.InstanceConfig.NodeCount)
	}

	if rocketmqVipInstanceDetail.InstanceConfig.NodeDistribution != nil {
		tmpList := []interface{}{}
		for _, v := range rocketmqVipInstanceDetail.InstanceConfig.NodeDistribution {
			tmpList = append(tmpList, *v.ZoneId)
		}
		_ = d.Set("zone_ids", tmpList)
	}

	if rocketmqVipInstanceDetail.ClusterInfo.Vpcs != nil {
		vpcInfoMap := map[string]interface{}{}
		if rocketmqVipInstanceDetail.ClusterInfo.Vpcs[0].VpcId != nil {
			vpcInfoMap["vpc_id"] = rocketmqVipInstanceDetail.ClusterInfo.Vpcs[0].VpcId
		}

		if rocketmqVipInstanceDetail.ClusterInfo.Vpcs[0].SubnetId != nil {
			vpcInfoMap["subnet_id"] = rocketmqVipInstanceDetail.ClusterInfo.Vpcs[0].SubnetId
		}

		_ = d.Set("vpc_info", []interface{}{vpcInfoMap})
	}

	if rocketmqVipInstances.SpecName != nil {
		_ = d.Set("spec", rocketmqVipInstances.SpecName)
	}

	if rocketmqVipInstances.MaxStorage != nil && rocketmqVipInstanceDetail.InstanceConfig.NodeCount != nil {
		storageSize := *rocketmqVipInstances.MaxStorage / *rocketmqVipInstanceDetail.InstanceConfig.NodeCount
		_ = d.Set("storage_size", storageSize)
	}

	return nil
}

func resourceTencentCloudTdmqRocketmqVipInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request   = tdmq.NewModifyRocketMQInstanceSpecRequest()
		clusterId = d.Id()
	)

	request.InstanceId = &clusterId

	immutableArgs := []string{"zone_ids", "vpc_info", "time_span", "ip_rules"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("spec") {
		if v, ok := d.GetOk("spec"); ok {
			request.Specification = helper.String(v.(string))
		}
	}

	if d.HasChange("node_count") {
		if v, ok := d.GetOkExists("node_count"); ok {
			request.NodeCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("storage_size") {
		if v, ok := d.GetOkExists("storage_size"); ok {
			request.StorageSize = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRocketMQInstanceSpec(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	// sleep - fix in the future
	time.Sleep(20 * time.Second)

	// wait
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		result, e := service.DescribeTdmqRocketmqVipInstancesByFilter(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *result.Status == svctdmq.RocketMqVipInsSuccess {
			return nil
		} else if *result.Status == svctdmq.RocketMqVipInsUpdate {
			return resource.RetryableError(fmt.Errorf("tdmq rocketmqVipInstance status is updating"))
		} else {
			e = fmt.Errorf("tdmq rocketmqVipInstance status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update tdmq rocketmqVipInstance failed, reason:%+v", logId, err)
		return err
	}

	// update name
	clusterRequest := tdmq.NewModifyRocketMQClusterRequest()
	clusterRequest.ClusterId = &clusterId
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			clusterRequest.ClusterName = helper.String(v.(string))
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRocketMQCluster(clusterRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cluster name failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqVipInstanceRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqVipInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rocketmq_vip_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		clusterId = d.Id()
	)

	// delete
	if err := service.DeleteTdmqRocketmqVipInstanceById(ctx, clusterId); err != nil {
		return err
	}

	// wait status is 2
	deleteFlag := false
	request := tdmq.NewDescribeRocketMQVipInstanceDetailRequest()
	request.ClusterId = &clusterId
	err := resource.Retry(tccommon.WriteRetryTimeout*6, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().DescribeRocketMQVipInstanceDetail(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.Instance" {
					deleteFlag = true
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result.Response.ClusterInfo.InstanceStatus != nil && *result.Response.ClusterInfo.InstanceStatus == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("tdmq rocketmqVipInstance first deleting"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cluster failed, reason:%+v", logId, err)
		return err
	}

	if deleteFlag {
		return nil
	}

	// delete again
	if err = service.DeleteTdmqRocketmqVipInstanceById(ctx, clusterId); err != nil {
		return err
	}

	// wait done
	err = resource.Retry(tccommon.WriteRetryTimeout*6, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().DescribeRocketMQVipInstanceDetail(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "ResourceNotFound.Instance" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return resource.RetryableError(fmt.Errorf("tdmq rocketmqVipInstance second deleting"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cluster failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
