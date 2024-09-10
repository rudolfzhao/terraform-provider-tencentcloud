package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqNamespaceRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqNamespaceRoleAttachmentCreate,
		Read:   resourceTencentCloudTdmqNamespaceRoleAttachmentRead,
		Update: resourceTencentCloudTdmqNamespaceRoleAttachmentUpdate,
		Delete: resourceTencentCloudTdmqNamespaceRoleAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"environ_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq namespace.",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of tdmq role.",
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The permissions of tdmq role.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of tdmq cluster.",
			},
			//compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.create")()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tdmqService = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		environId   string
		roleName    string
		permissions []*string
		clusterId   string
	)

	if temp, ok := d.GetOk("environ_id"); ok {
		environId = temp.(string)
		if len(environId) < 1 {
			return fmt.Errorf("environ_id should be not empty string")
		}
	}

	if temp, ok := d.GetOk("role_name"); ok {
		roleName = temp.(string)
		if len(roleName) < 1 {
			return fmt.Errorf("role_name should be not empty string")
		}
	}

	if v, ok := d.GetOk("permissions"); ok {
		for _, id := range v.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	}

	if temp, ok := d.GetOk("cluster_id"); ok {
		clusterId = temp.(string)
		if len(clusterId) < 1 {
			return fmt.Errorf("cluster_id should be not empty string")
		}
	}

	err := tdmqService.CreateTdmqNamespaceRoleAttachment(ctx, environId, roleName, permissions, clusterId)
	if err != nil {
		return err
	}

	d.SetId(environId + tccommon.FILED_SP + roleName)

	return resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tdmqService = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}

	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := tdmqService.DescribeTdmqNamespaceRoleAttachment(ctx, environId, roleName, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			return nil
		}

		_ = d.Set("environ_id", info.EnvironmentId)
		_ = d.Set("role_name", info.RoleName)
		_ = d.Set("permissions", info.Permissions)
		_ = d.Set("create_time", info.CreateTime)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.update")()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		permissions []*string
	)

	immutableArgs := []string{"environ_id", "role_name", "cluster_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}

	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	old, now := d.GetChange("permissions")
	if d.HasChange("permissions") {
		for _, id := range now.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	} else {
		for _, id := range old.([]interface{}) {
			permissions = append(permissions, helper.String(id.(string)))
		}
	}

	d.Partial(true)

	if err := service.ModifyTdmqNamespaceRoleAttachment(ctx, environId, roleName, permissions, clusterId); err != nil {
		return err
	}

	d.Partial(false)
	return resourceTencentCloudTdmqNamespaceRoleAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqNamespaceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_namespace_role_attachment.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("environment role id is borken, id is %s", d.Id())
	}

	environId := idSplit[0]
	roleName := idSplit[1]
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.DeleteTdmqNamespaceRoleAttachment(ctx, environId, roleName, clusterId); err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == svcvpc.VPCNotFound {
					return nil
				}
			}

			return resource.RetryableError(err)
		}

		return nil
	})

	return err
}
