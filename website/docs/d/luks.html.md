---
layout: "ignition"
page_title: "Ignition: ignition_luks"
sidebar_current: "docs-ignition-datasource-luks"
description: |-
  Describes the desired state of the system’s luks.
---

# ignition\_luks

Describes the desired state of the system’s luks.

## Example Usage

```hcl
data "ignition_luks" "luks" {
	name = "data"
	device = "/dev/sda2"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for the resulting luks device.

* `device` - (Required) The absolute path to the device. Devices are typically referenced by the /dev/disk/by-* symlinks.

* `discard` - (Optional) Whether to issue discard commands to the underlying block device when blocks are freed. Enabling this improves performance and device longevity on SSDs and space utilization on thinly provisioned SAN devices, but leaks information about which disk blocks contain data. If omitted, it defaults to false.

* `label` - (Optional) The label of the luks device.

* `open_options` - (Optional) Any additional options to be passed to cryptsetup luksOpen. Supported options will be persistently written to the luks volume.

* `options` - (Optional) Any additional options to be passed to cryptsetup luksFormat.

* `uuid` - (Optional) The uuid of the luks device.

* `wipe_volume` - (Optional) Whether or not to wipe the device before luks creation.

* `key_file` - (Optional) Options related to the contents of the key file.

* `clevis` - (Optional) Describes the clevis configuration for the luks device.

The `key_file` block supports:

* `source` - (Required) The URL of the key file. Supported schemes are http, https, tftp, s3, arn, gs, and data. When using http, it is advisable to use the verification option to ensure the contents haven’t been modified.

* `compression` - (Optional) The type of compression used on the key file (null or gzip). Compression cannot be used with S3.

* `http_headers` - (Optional) A list of HTTP headers to be added to the request. Available for http and https source schemes only.

* `verification` - (Optional) The hash of the config, in the form _\<type\>-\<value\>_ where type is either sha512 or sha256. If compression is specified, the hash describes the decompressed file.

The `http_headers` block supports:

* `name` - (Required) The header name.

* `value` - (Required) The header contents.

The `clevis` block supports:

* `tang` - (Optional) describes a tang server. Every server must have a unique url.

* `tpm2` - (Optional) Whether or not to use a tpm2 device.

* `threshold` - (Optional) Sets the minimum number of pieces required to decrypt the device. Default is 1.

* `custom` - (Optional) Overrides the clevis configuration. The pin & config will be passed directly to clevis luks bind. If specified, all other clevis options must be omitted.

The `tang` block supports:

* `url` - (Required) Url of the tang server.

* `thumbprint` - (Optional) Thumbprint of a trusted signing key.

* `advertisement` - (Optional) The advertisement JSON. If not specified, the advertisement is fetched from the tang server during provisioning.

The `custom` block supports:

* `pin` - (Optional) The clevis pin.

* `config` - (Optional) The clevis configuration JSON.

* `needs_network` - (Optional) Whether or not the device requires networking.


## Attributes Reference

The following attributes are exported:

* `rendered` - The rendered template to reference this resource in _ignition_config_.
