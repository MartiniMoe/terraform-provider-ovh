---
subcategory : "Managed Private Registry (MPR)"
---

# ovh_cloud_project_containerregistry_users (Data Source)

Use this data source to get the list of users of a container registry associated with a public cloud project.

## Example Usage

```terraform
data "ovh_cloud_project_containerregistry" "my_registry" {
  service_name = "XXXXXX"
  registry_id  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"
}

data "ovh_cloud_project_containerregistry_users" "users" {
  service_name = ovh_cloud_project_containerregistry.my_registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.my_registry.id
}
```

## Argument Reference

* `service_name` - The id of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.

* `registry_id` - Registry ID

## Attributes Reference

The following attributes are exported:

* `result` - The list of users of the container registry associated with the project.
  * `id` - User ID
  * `user` - User name
  * `email` - User email
