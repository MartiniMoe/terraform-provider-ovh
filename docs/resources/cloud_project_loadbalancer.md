---
subcategory : "Load Balancer (Public Cloud / Octavia)"
---

# ovh_cloud_project_loadbalancer

## Example Usage

```terraform
resource "ovh_cloud_project_loadbalancer" "lb" {
  service_name = "<public cloud project ID>"
  region_name = "GRA9"
  flavor_id = "<loadbalancer flavor ID>"
  network = {
    private = {
      network = {
        id = element([for region in ovh_cloud_project_network_private.mypriv.regions_attributes: region if "${region.region}" == "GRA9"], 0).openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.myprivsub.id
      }
    }
  }
  description = "My new LB"
  listeners = [
    {
      port = "34568"
      protocol = "tcp"
    },
    {
      port = "34569"
      protocol = "udp"
    }
  ]
}
```

### Example usage with network and subnet creation

```terraform
resource "ovh_cloud_project_network_private" "priv" {
  service_name  = "<public cloud project ID>"
  vlan_id       = "10"
  name          = "my_priv"
  regions       = ["GRA9"]
}

resource "ovh_cloud_project_network_private_subnet" "privsub" {
  service_name  = ovh_cloud_project_network_private.priv.service_name
  network_id    = ovh_cloud_project_network_private.priv.id
  region        = "GRA9"
  start         = "10.0.0.2"
  end           = "10.0.255.254"
  network       = "10.0.0.0/16"
  dhcp          = true
}

resource "ovh_cloud_project_loadbalancer" "lb" {
  service_name = ovh_cloud_project_network_private_subnet.privsub.service_name
  region_name = ovh_cloud_project_network_private_subnet.privsub.region
  flavor_id = "<loadbalancer flavor ID>"
  network = {
    private = {
      network = {
        id = element([for region in ovh_cloud_project_network_private.priv.regions_attributes: region if "${region.region}" == "GRA9"], 0).openstackid
        subnet_id = ovh_cloud_project_network_private_subnet.privsub.id
      }
    }
  }
  description = "My new LB"
  listeners = [
    {
      port = "34568"
      protocol = "tcp"
    },
    {
      port = "34569"
      protocol = "udp"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `flavor_id` (String) Loadbalancer flavor id
- `network` (Attributes) Network information to create the loadbalancer (see [below for nested schema](#nestedatt--network))
- `region_name` (String) Region name
- `service_name` (String) Service name

### Optional

- `description` (String) Description of the loadbalancer
- `listeners` (Attributes List) Listeners to create with the loadbalancer (see [below for nested schema](#nestedatt--listeners))
- `name` (String) Name of the resource

### Read-Only

- `created_at` (String) The UTC date and timestamp when the resource was created
- `floating_ip` (Attributes) Information about floating IP (see [below for nested schema](#nestedatt--floating_ip))
- `id` (String) ID of the resource
- `operating_status` (String) Operating status of the resource
- `provisioning_status` (String) Provisioning status of the resource
- `region` (String) Region of the resource
- `updated_at` (String) UTC date and timestamp when the resource was created
- `vip_address` (String) IP address of the Virtual IP
- `vip_network_id` (String) Openstack ID of the network for the Virtual IP
- `vip_subnet_id` (String) ID of the subnet for the Virtual IP

<a id="nestedatt--network"></a>

### Nested Schema for `network`

Required:

- `private` (Attributes) Information to private network (see [below for nested schema](#nestedatt--network--private))

<a id="nestedatt--network--private"></a>

### Nested Schema for `network.private`

Required:

- `network` (Attributes) Network to associate (see [below for nested schema](#nestedatt--network--private--network))

Optional:

- `floating_ip` (Attributes) Floating IP to associate (see [below for nested schema](#nestedatt--network--private--floating_ip))
- `floating_ip_create` (Attributes) Floating IP to create (see [below for nested schema](#nestedatt--network--private--floating_ip_create))
- `gateway` (Attributes) Gateway to associate (see [below for nested schema](#nestedatt--network--private--gateway))
- `gateway_create` (Attributes) Gateway to create (see [below for nested schema](#nestedatt--network--private--gateway_create))

<a id="nestedatt--network--private--network"></a>

### Nested Schema for `network.private.network`

Required:

- `id` (String) Private network ID
- `subnet_id` (String) Subnet ID

<a id="nestedatt--network--private--floating_ip"></a>

### Nested Schema for `network.private.floating_ip`

Optional:

- `id` (String) ID of the floatingIp

<a id="nestedatt--network--private--floating_ip_create"></a>

### Nested Schema for `network.private.floating_ip_create`

Optional:

- `description` (String) Description for the floatingIp

<a id="nestedatt--network--private--gateway"></a>

### Nested Schema for `network.private.gateway`

Optional:

- `id` (String) ID of the gateway

<a id="nestedatt--network--private--gateway_create"></a>

### Nested Schema for `network.private.gateway_create`

Optional:

- `model` (String) Model of the gateway
- `name` (String) Name of the gateway

<a id="nestedatt--listeners"></a>

### Nested Schema for `listeners`

Required:

- `port` (Number) Listener port
- `protocol` (String) Protocol for the listener

Optional:

- `allowed_cidrs` (List of String) The allowed CIDRs
- `description` (String) The description of the listener
- `name` (String) Name of the listener
- `pool` (Attributes) Listener pool (see [below for nested schema](#nestedatt--listeners--pool))
- `secret_id` (String) Secret ID to get certificate for SSL listener creation
- `timeout_client_data` (Number) Timeout client data of the listener
- `timeout_member_data` (Number) Timeout member data of the listener
- `tls_versions` (List of String) TLS versions of the listener

<a id="nestedatt--listeners--pool"></a>

### Nested Schema for `listeners.pool`

Optional:

- `algorithm` (String) Pool algorithm to split traffic between members
- `health_monitor` (Attributes) Pool health monitor (see [below for nested schema](#nestedatt--listeners--pool--health_monitor))
- `members` (Attributes List) Pool members (see [below for nested schema](#nestedatt--listeners--pool--members))
- `name` (String) Name of the pool
- `protocol` (String) Protocol for the pool
- `session_persistence` (Attributes) Pool session persistence (see [below for nested schema](#nestedatt--listeners--pool--session_persistence))

<a id="nestedatt--listeners--pool--health_monitor"></a>

### Nested Schema for `listeners.pool.health_monitor`

Optional:

- `delay` (Number) Duration between sending probes to members, in seconds
- `http_configuration` (Attributes) Monitor HTTP configuration (see [below for nested schema](#nestedatt--listeners--pool--health_monitor--http_configuration))
- `max_retries` (Number) Number of successful checks before changing the operating status of the member to ONLINE
- `max_retries_down` (Number) Number of allowed check failures before changing the operating status of the member to ERROR
- `monitor_type` (String) Type of the monitor
- `name` (String) The name of the resource
- `operating_status` (String) The operating status of the resource
- `provisioning_status` (String) The provisioning status of the resource
- `timeout` (Number) Maximum time, in seconds, that a monitor waits to connect before it times out. This value must be less than the delay value

<a id="nestedatt--listeners--pool--health_monitor--http_configuration"></a>

### Nested Schema for `listeners.pool.health_monitor.http_configuration`

Optional:

- `domain_name` (String) Domain name, which be injected into the HTTP Host Header to the backend server for HTTP health check
- `expected_codes` (String) Status codes expected in response from the member to declare it healthy; The list of HTTP status codes expected in response from the member to declare it healthy. Specify one of the following values: * A single value, such as 200; * A list, such as 200, 202; * A range, such as 200-204
- `http_method` (String) HTTP method that the health monitor uses for requests
- `http_version` (String) HTTP version that the health monitor uses for requests
- `url_path` (String) HTTP URL path of the request sent by the monitor to test the health of a backend member

<a id="nestedatt--listeners--pool--members"></a>

### Nested Schema for `listeners.pool.members`

Optional:

- `address` (String) IP address of the resource
- `name` (String) Name of the member
- `protocol_port` (Number) Protocol port number for the resource
- `weight` (Number) Weight of a member determines the portion of requests or connections it services compared to the other members of the pool. Between 1 and 256.

<a id="nestedatt--listeners--pool--session_persistence"></a>

### Nested Schema for `listeners.pool.session_persistence`

Optional:

- `cookie_name` (String) Cookie name, only applicable to session persistence through cookie
- `type` (String) Type of session persistence

<a id="nestedatt--floating_ip"></a>

### Nested Schema for `floating_ip`

Read-Only:

- `id` (String) ID of the resource
- `ip` (String) IP Address of the resource

## Import

A load balancer in a public cloud project can be imported using the `service_name`, `region_name` and `id` attributes. Using the following configuration:

```terraform
import {
  id = "<service_name>/<region_name>/<id>"
  to = ovh_cloud_project_loadbalancer.lb
}
```

You can then run:

```bash
$ terraform plan -generate-config-out=lb.tf
$ terraform apply
```

The file `lb.tf` will then contain the imported resource's configuration, that can be copied next to the `import` block above. See https://developer.hashicorp.com/terraform/language/import/generating-configuration for more details.
