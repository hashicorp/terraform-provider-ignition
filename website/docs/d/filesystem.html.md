---
layout: "ignition"
page_title: "Ignition: ignition_filesystem"
sidebar_current: "docs-ignition-datasource-filesystem"
description: |-
  Describes the desired state of a system’s filesystem.
---

# ignition\_filesystem

Describes the desired state of a the system’s filesystems to be configured and/or used with the _ignition\_file_ resource.

## Example Usage

```hcl
data "ignition_filesystem" "foo" {
    device = "/dev/disk/by-label/ROOT"
    format = "xfs"
    options = ["-L", "ROOT"]
}
```

## Argument Reference

The following arguments are supported:

* `device` - (Required) The absolute path to the device. Devices are typically referenced by the _/dev/disk/by-*_ symlinks.

* `format` - (Required) The filesystem format (ext4, btrfs, xfs, vfat, or swap).

* `wipe_filesystem` - (Optional)  Whether or not to wipe the device before filesystem creation.

* `label` - (Optional) The label of the filesystem.

* `uuid` - (Optional) The uuid of the filesystem.

* `options` - (Optional) Any additional options to be passed to the format-specific mkfs utility.

* `path` - (Optional) The mount-point of the filesystem while Ignition is running relative to where the root filesystem will be mounted. This is not necessarily the same as where it should be mounted in the real root, but it is encouraged to make it the same.

## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
