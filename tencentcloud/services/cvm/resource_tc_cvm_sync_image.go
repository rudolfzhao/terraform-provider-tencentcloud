package cvm

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmSyncImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSyncImageCreate,
		Read:   resourceTencentCloudCvmSyncImageRead,
		Delete: resourceTencentCloudCvmSyncImageDelete,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID. The specified image must meet the following requirement: the images must be in the `NORMAL` state.",
			},

			"destination_regions": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of destination regions for synchronization. Limits: It must be a valid region. For a custom image, the destination region cannot be the source region. For a shared image, the destination region must be the source region, which indicates to create a copy of the image as a custom image in the same region.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Checks whether image synchronization can be initiated.",
			},

			"image_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Destination image name.",
			},

			"image_set_required": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to return the ID of image created in the destination region.",
			},
		},
	}
}

func resourceTencentCloudCvmSyncImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cvm.NewSyncImagesRequest()
	imageId := d.Get("image_id").(string)
	request.ImageIds = []*string{&imageId}

	if v, ok := d.GetOk("destination_regions"); ok {
		destinationRegionsSet := v.(*schema.Set).List()
		for i := range destinationRegionsSet {
			destinationRegions := destinationRegionsSet[i].(string)
			request.DestinationRegions = append(request.DestinationRegions, &destinationRegions)
		}
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("image_set_required"); v != nil {
		request.ImageSetRequired = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().SyncImages(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm syncImages failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(imageId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmSyncImagesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCvmSyncImageRead(d, meta)
}

func resourceTencentCloudCvmSyncImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmSyncImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
