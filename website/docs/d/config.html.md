---
layout: "ignition"
page_title: "Ignition: ignition_config"
sidebar_current: "docs-ignition-datasource-config"
description: |-
  Renders an ignition configuration as JSON
---

# ignition\_config

Renders an ignition configuration as JSON. It  contains all the disks, partitions, arrays, filesystems, files, users, groups and units.

## Example Usage

```hcl
data "ignition_config" "example" {
  systemd = [
    data.ignition_systemd_unit.example.rendered,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `disks` - (Optional) The list of disks to be configured and their options.

* `arrays` - (Optional) The list of RAID arrays to be configured.

* `filesystems` - (Optional) The list of filesystems to be configured and/or used in the `ignition_file`, `ignition_directory`, and `ignition_link` resources.

* `files` - (Optional) The list of files to be written.

* `directories` - (Optional) The list of directories to be created.

* `links` - (Optional) The list of links to be created.

* `systemd` - (Optional) The list of systemd units. Describes the desired state of the systemd units.

* `users` - (Optional) The list of accounts to be added.

* `groups` - (Optional) The list of groups to be added.

* `kernel_arguments` - (Optional) A string that describes the desired kernel arguments.

* `tls_ca` - (Optional) The list of additional certificate authorities to be used for TLS verification when fetching over https.

* `merge` - (Optional) A list of the configs to be merged to the current config.

* `replace` - (Optional) A block with config that will replace the current.

The `tls_ca`, `merge` and `replace` blocks support:

* `source` - (Required) The URL of the config. Supported schemes are http, https, tftp, s3, and data. When using http, it is advisable to use the verification option to ensure the contents haven't been modified.

* `compression` - (Optional) The type of compression used on the config (null or gzip). Compression cannot be used with S3.

* `verification` - (Optional) The hash of the config, in the form _\<type\>-\<value\>_ where type is either sha512 or sha256. If compression is specified, the hash describes the decompressed config.

* `http_headers` - (Optional) A list of HTTP headers to be added to the request.

The `http_headers` blocks support:

* `name` - (Required) The header name.

* `value` - (Required) The header contents.


## Attributes Reference

The following attributes are exported:

* `rendered` - The final rendered template.
