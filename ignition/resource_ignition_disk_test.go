package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_1/types"
)

func TestIgnitionDisk(t *testing.T) {
	testIgnition(t, `
		data "ignition_disk" "foo" {
			device = "/foo"
			partition {
				label = "qux"
				sizemib = 42
				startmib = 2048
				type_guid = "01234567-89AB-CDEF-EDCB-A98765432101"
			}
		}

		data "ignition_config" "test" {
			disks = [data.ignition_disk.foo.rendered]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Disks) != 1 {
			return fmt.Errorf("disks, found %d", len(c.Storage.Disks))
		}

		d := c.Storage.Disks[0]
		if d.Device != "/foo" {
			return fmt.Errorf("name, found %q", d.Device)
		}

		if len(d.Partitions) != 1 {
			return fmt.Errorf("partition, found %d", len(d.Partitions))
		}

		p := d.Partitions[0]
		if string(*p.Label) != "qux" {
			return fmt.Errorf("partition.0.label, found %q", *p.Label)
		}

		if int(*p.SizeMiB) != 42 {
			return fmt.Errorf("partition.0.sizemib, found %q", *p.SizeMiB)
		}

		if int(*p.StartMiB) != 2048 {
			return fmt.Errorf("partition.0.start, found %q", *p.StartMiB)
		}

		if string(*p.TypeGUID) != "01234567-89AB-CDEF-EDCB-A98765432101" {
			return fmt.Errorf("partition.0.type_guid, found %q", *p.TypeGUID)
		}

		return nil
	})
}

func TestIgnitionDiskInvalidDevice(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_disk" "foo" {
			device = "a"
		}

		data "ignition_config" "test" {
			disks = [data.ignition_disk.foo.rendered]
		}
	`, regexp.MustCompile("path not absolute"))
}

func TestIgnitionDiskInvalidPartition(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_disk" "foo" {
			device = "/foo"
			partition {
				label = "qux"
				sizemib = 42
				startmib = 2048
				type_guid =  "01234567-89AB-CDEF-EDCB-A98765432101"
			}
			partition {
				label = "bar"
				sizemib = 42
				startmib = 2048
				type_guid =  "01234567-89AB-CDEF-EDCB-A98765432101"
			}
		}

		data "ignition_config" "test" {
			disks = [data.ignition_disk.foo.rendered]
		}
	`, regexp.MustCompile("overlap"))
}
