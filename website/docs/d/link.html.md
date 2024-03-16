---
layout: "ignition"
page_title: "Ignition: ignition_link"
sidebar_current: "docs-ignition-datasource-link"
description: |-
  Describes a link to be created in a particular filesystem.
---

# ignition\_link

Describes a link to be created in a particular filesystem.

## Example Usage

```hcl
data "ignition_link" "symlink" {
	path = "/symlink"
    target = "/foo"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The absolute path to the link.

* `target` - (Required) The target path of the link.

* `overwrite` - (Optional) Overwrite the link, if it already exists.

* `hard` - (Optional) A symbolic link is created if this is false, a hard one if this is true.

* `uid` - (Optional) The user ID of the owner.

* `gid` - (Optional) The group ID of the owner.

## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
