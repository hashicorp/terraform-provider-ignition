package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_1/types"
)

func TestIgnitionFilesystem(t *testing.T) {
	testIgnition(t, `
		data "ignition_filesystem" "qux" {
			device = "/qux"
			format = "ext4"
		}
		data "ignition_filesystem" "baz" {
			device = "/baz"
			format = "ext4"
			wipe_filesystem = true
			label = "root"
			uuid = "qux"
			options = ["rw"]
		}
		data "ignition_config" "test" {
			filesystems = [
				data.ignition_filesystem.qux.rendered,
				data.ignition_filesystem.baz.rendered,
			]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Filesystems) != 2 {
			return fmt.Errorf("disks, found %d", len(c.Storage.Filesystems))
		}

		f := c.Storage.Filesystems[0]
		if f.Device != "/qux" {
			return fmt.Errorf("device, found %q", f.Device)
		}

		if string(*f.Format) != "ext4" {
			return fmt.Errorf("format, found %q", *f.Format)
		}

		f = c.Storage.Filesystems[1]

		if f.Device != "/baz" {
			return fmt.Errorf("device, found %q", f.Device)
		}

		if *f.Format != "ext4" {
			return fmt.Errorf("format, found %q", *f.Format)
		}

		if *f.Label != "root" {
			return fmt.Errorf("label, found %q", *f.Label)
		}

		if *f.UUID != "qux" {
			return fmt.Errorf("uuid, found %q", *f.UUID)
		}

		if *f.WipeFilesystem != true {
			return fmt.Errorf("wipe_filesystem, found %t", *f.WipeFilesystem)
		}

		if len(f.Options) != 1 || f.Options[0] != "rw" {
			return fmt.Errorf("options, found %q", f.Options)
		}

		return nil
	})
}

func TestIgnitionFilesystemInvalidPath(t *testing.T) {
	testIgnitionError(t, `
		variable "ignition_filesystem_renders" {
			type = "list"
			default = [""]
		}

		data "ignition_filesystem" "foo" {
			device = "/foo"
			format = "ext4"
			path = "foo"
		}

		data "ignition_config" "test" {
			filesystems = concat(
				[data.ignition_filesystem.foo.rendered],
				var.ignition_filesystem_renders
			)
		}
	`, regexp.MustCompile("absolute"))
}
