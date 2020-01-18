---
layout: "ignition"
page_title: "Provider: Ignition"
sidebar_current: "docs-ignition-index"
description: |-
  The Ignition provider is used to generate Ignition configuration files used by CoreOS Linux.
---

# Ignition Provider

The Ignition provider is used to generate [Ignition](https://coreos.com/ignition/docs/latest/) configuration files. _Ignition_ is the provisioning utility used by [CoreOS](https://coreos.com/) Linux.

The ignition provider is what we call a _logical provider_ and doesn't manage any _physical_ resources. It generates configurations files to be used by other resources.

Use the navigation to the left to read about the available resources.

## Ignition versions

The current Ignition Config Spec version supported by this provider is `3.0.0`. For older versions you should use previous [releases](https://github.com/terraform-providers/terraform-provider-ignition/releases) of this provider:

* terraform-provider-ignition `[1.0.0,2.0.0)` - Ignition `0.34.0` / Config Spec Version `2.1.0`
* terraform-provider-ignition `>=2.0.0` - Ignition `2.1.1` / Config Spec Version `3.0.0`

## Example Usage

This config will write a single service unit (shown below) with the contents of an example service. This unit will be enabled as a dependency of multi-user.target and therefore start on boot

```hcl
# Systemd unit data resource containing the unit definition
data "ignition_systemd_unit" "example" {
  name = "example.service"
  content = "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
}

# Ignition config include the previous defined systemd unit data resource
data "ignition_config" "example" {
  systemd = [
    data.ignition_systemd_unit.example.rendered,
  ]
}

# Create a CoreOS server using the Ignition config.
resource "aws_instance" "web" {
  # ...

  user_data = data.ignition_config.example.rendered
}
```
