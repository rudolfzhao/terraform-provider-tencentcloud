package apm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmInstanceCreate,
		Read:   resourceTencentCloudApmInstanceRead,
		Update: resourceTencentCloudApmInstanceUpdate,
		Delete: resourceTencentCloudApmInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name Of Instance.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description Of Instance.",
			},

			"trace_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Duration Of Trace Data.",
			},

			"span_daily_counters": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Quota Of Instance Reporting.",
			},

			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Modify the billing mode: `1` means prepaid, `0` means pay-as-you-go, the default value is `0`.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudApmInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = apm.NewCreateApmInstanceRequest()
		response   = apm.NewCreateApmInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trace_duration"); ok {
		request.TraceDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("span_daily_counters"); ok {
		request.SpanDailyCounters = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().CreateApmInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apm instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::apm:%s:uin/:apm-instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApmInstanceRead(d, meta)
}

func resourceTencentCloudApmInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *apm.ApmInstanceDetail
		instanceId = d.Id()
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeApmInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		instance = result
		return nil
	})

	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApmInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.Name != nil {
		_ = d.Set("name", instance.Name)
	}

	if instance.Description != nil {
		_ = d.Set("description", instance.Description)
	}

	if instance.TraceDuration != nil {
		_ = d.Set("trace_duration", instance.TraceDuration)
	}

	if instance.SpanDailyCounters != nil {
		_ = d.Set("span_daily_counters", instance.SpanDailyCounters)
	}

	if instance.PayMode != nil {
		_ = d.Set("pay_mode", instance.PayMode)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "apm", "apm-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApmInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = apm.NewModifyApmInstanceRequest()
		instanceId = d.Id()
	)

	needChange := false
	request.InstanceId = &instanceId

	mutableArgs := []string{"name", "description", "trace_duration", "span_daily_counters", "pay_mode"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("trace_duration"); ok {
			request.TraceDuration = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("span_daily_counters"); ok {
			request.SpanDailyCounters = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("pay_mode"); ok {
			request.PayMode = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update apm instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("apm", "apm-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApmInstanceRead(d, meta)
}

func resourceTencentCloudApmInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.DeleteApmInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
