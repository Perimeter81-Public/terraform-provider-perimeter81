resource "checkpointsase_object_services" "web_ports" {
  name        = "webPorts"
  description = "Common HTTP/HTTPS ports"

  protocols {
    protocol   = "tcp"
    value      = [80, 443]
    value_type = "list"
  }
}
