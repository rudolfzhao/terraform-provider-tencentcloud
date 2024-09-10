Provide a resource to configure kubernetes cluster authentication info.

~> **NOTE:** Only available for cluster version >= 1.20

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                           = tencentcloud_kubernetes_cluster.managed_cluster.id
  jwks_uri                             = "https://${tencentcloud_kubernetes_cluster.managed_cluster.id}.ccs.tencent-cloud.com/openid/v1/jwks"
  issuer                               = "https://${tencentcloud_kubernetes_cluster.managed_cluster.id}.ccs.tencent-cloud.com"
  auto_create_discovery_anonymous_auth = true
}
```

Use the TKE default issuer and jwks_uri

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

# if you want to use tke default issuer and jwks_uri, please set use_tke_default to true and set issuer to empty string.
resource "tencentcloud_kubernetes_auth_attachment" "test_use_tke_default_auth_attach" {
  cluster_id                           = tencentcloud_kubernetes_cluster.managed_cluster.id
  auto_create_discovery_anonymous_auth = true
  use_tke_default                      = true
}
```

Use OIDC Config
```
resource "tencentcloud_kubernetes_auth_attachment" "test_auth_attach" {
  cluster_id                              = tencentcloud_kubernetes_cluster.managed_cluster.id
  use_tke_default                         = true
  auto_create_discovery_anonymous_auth    = true
  auto_create_oidc_config                 = true
  auto_install_pod_identity_webhook_addon = true
}

data "tencentcloud_cam_oidc_config" "oidc_config" {
  name       = tencentcloud_kubernetes_cluster.managed_cluster.id
  depends_on = [
    tencentcloud_kubernetes_auth_attachment.test_auth_attach
  ]
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}

output "identity_url" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_url
}

```

Import

tke cluster authentication can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_auth_attachment.test cls-xxx
```