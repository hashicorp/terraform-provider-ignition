---
layout: "ignition"
page_title: "Ignition: ignition_kernel_arguments"
sidebar_current: "docs-ignition-datasource-kernel-arguments"
description: |-
  Describes the desired kernel arguments.
---

# ignition\_kernel\_arguments

Describes the desired kernel arguments.

## Example Usage

```hcl
data "ignition_kernel_arguments" "foo" {
  shouldexist = ["foo","bar"]
  shouldnotexist = ["baz"]
}
```

## Argument Reference

The following arguments are supported:

* `shouldexist` - (Optional) The list of kernel arguments that should exist.

* `shouldnotexist` - (Optional) The list of kernel arguments that should not exist.

## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
