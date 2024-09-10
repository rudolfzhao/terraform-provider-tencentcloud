package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAlbServerAttachment() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.15.0. Please use 'tencentcloud_clb_attachment' instead.",
		Create:             resourceTencentCloudAlbServerAttachmentCreate,
		Read:               resourceTencentCloudAlbServerAttachmentRead,
		Delete:             resourceTencentCloudAlbServerAttachmentDelete,
		Update:             resourceTencentCloudAlbServerAttachmentUpdate,

		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "loadbalancer ID.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "listener ID.",
			},
			"location_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "location ID, only support for layer 7 loadbalancer.",
			},
			"backends": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    100,
				Description: "list of backend server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A list backend instance ID (CVM instance ID).",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 65535),
							Description:  "The port used by the backend server. Valid value range: [1-65535].",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "Weight of the backend server. Valid value range: [0-100]. Default to 10.",
						},
					},
				},
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol type, http or tcp.",
			},
		},
	}
}

func resourceTencentCloudAlbServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := clb.NewRegisterTargetsRequest()

	loadbalancerId := d.Get("loadbalancer_id").(string)
	listenerId := d.Get("listener_id").(string)
	request.LoadBalancerId = helper.String(loadbalancerId)
	request.ListenerId = helper.String(listenerId)

	location_id := ""
	if v, ok := d.GetOk("location_id"); ok {
		location_id = v.(string)
		if location_id != "" {
			request.LocationId = helper.String(location_id)
		}
	}

	for _, inst_ := range d.Get("backends").(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		requestId := ""
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().RegisterTargets(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId = *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(retryErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create alb attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}

	id := fmt.Sprintf("%v:%v:%v", loadbalancerId, listenerId, location_id)
	d.SetId(id)

	return resourceTencentCloudAlbServerAttachmentRead(d, meta)
}

func resourceTencentCloudAlbServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return fmt.Errorf("id of resource.tencentcloud_alb_server_attachment is wrong")
	}

	clbId := items[0]
	listenerId := items[1]
	locationId := items[2]

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var instance *clb.ListenerBackend
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read alb attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("loadbalancer_id", clbId)
	_ = d.Set("listener_id", listenerId)
	_ = d.Set("protocol_type", instance.Protocol)

	if *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *instance.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
		_ = d.Set("location_id", locationId)
		if len(instance.Rules) > 0 {
			for _, loc := range instance.Rules {
				if locationId == "" || locationId == *loc.LocationId {
					_ = d.Set("backends", flattenBackendList(loc.Targets))
				}
			}
		}
	} else if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_UDP ||
		*instance.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL || *instance.Protocol == CLB_LISTENER_PROTOCOL_QUIC {
		_ = d.Set("backends", flattenBackendList(instance.Targets))
	}

	return nil
}

func resourceTencentCloudAlbServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return fmt.Errorf("id of resource.tencentcloud_alb_server_attachment is wrong")
	}

	clbId := items[0]
	listenerId := items[1]
	locationId := items[2]

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteAttachmentById(ctx, clbId, listenerId, locationId, d.Get("backends").(*schema.Set).List(), "", "")
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s reason[%s]\n", logId, err.Error())
		return err
	}

	return nil
}

func resourceTencentCloudAlbServerAttachementRemove(d *schema.ResourceData, meta interface{}, remove []interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.remove")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return fmt.Errorf("id %s of resource.tencentcloud_alb_server_attachment is wrong", d.Id())
	}
	clbId := items[0]
	listenerId := items[1]
	locationId := items[2]
	request := clb.NewDeregisterTargetsRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)
	if locationId != "" {
		request.LocationId = helper.String(locationId)
	}
	for _, inst_ := range remove {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		requestId := ""
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().DeregisterTargets(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId = *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(retryErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s remove clb attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudAlbServerAttachementAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.add")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	listenerId := d.Get("listener_id").(string)
	clbId := d.Get("loadbalancer_id").(string)
	locationId := ""
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)

	if v, ok := d.GetOk("location_id"); ok {
		locationId = v.(string)
		if locationId != "" {
			request.LocationId = helper.String(locationId)
		}
	}

	for _, inst_ := range add {
		inst := inst_.(map[string]interface{})
		request.Targets = append(request.Targets, clbNewTarget(inst["instance_id"], inst["eni_ip"], inst["port"], inst["weight"]))
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		requestId := ""
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().RegisterTargets(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId = *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(retryErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s add clb attachment failed, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudAlbServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_server_attachment.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	if d.HasChange("backends") {
		o, n := d.GetChange("backends")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(remove) > 0 {
			err := resourceTencentCloudAlbServerAttachementRemove(d, meta, remove)
			if err != nil {
				return err
			}
		}
		if len(add) > 0 {
			err := resourceTencentCloudAlbServerAttachementAdd(d, meta, add)
			if err != nil {
				return err
			}
		}
		return resourceTencentCloudAlbServerAttachmentRead(d, meta)
	}

	return nil
}
