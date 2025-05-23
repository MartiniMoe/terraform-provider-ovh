data "ovh_iploadbalancing" "lb" {
  service_name = "ip-1.2.3.4"
  state        = "ok"
}

resource "ovh_iploadbalancing_tcp_farm" "farm_name" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  port         = 8080
  zone         = "all"
}

resource "ovh_iploadbalancing_tcp_farm_server" "backend" {
  service_name           = "${data.ovh_iploadbalancing.lb.service_name}"
  farm_id                = "${ovh_iploadbalancing_tcp_farm.farm_name.id}"
  display_name           = "mybackend"
  address                = "4.5.6.7"
  status                 = "active"
  port                   = 80
  proxy_protocol_version = "v2"
  weight                 = 2
  probe                  = true
  ssl                    = false
  backup                 = true
}

resource "ovh_iploadbalancing_refresh" "mylb" {
  service_name = "${data.ovh_iploadbalancing.lb.service_name}"
  keepers = [
    "${ovh_iploadbalancing_tcp_farm_server.backend.*.address}",
  ]
}
