package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_4/types"
)

func TestIgnitionSystemdUnit(t *testing.T) {
	testIgnition(t, `
		data "ignition_systemd_unit" "foo" {
			name = "foo.service"
			content = "[Match]\nName=eth0\n\n[Network]\nAddress=10.0.1.7\n"
			enabled = true
			mask = true

			dropin {
				name = "foo.conf"
				content = "[Match]\nName=eth0\n\n[Network]\nAddress=10.0.1.7\n"
			}
		}

		data "ignition_config" "test" {
			systemd = [data.ignition_systemd_unit.foo.rendered]
		}
	`, func(c *types.Config) error {
		if len(c.Systemd.Units) != 1 {
			return fmt.Errorf("systemd, found %d", len(c.Systemd.Units))
		}

		u := c.Systemd.Units[0]

		if u.Name != "foo.service" {
			return fmt.Errorf("name, found %q", u.Name)
		}

		if *u.Contents != "[Match]\nName=eth0\n\n[Network]\nAddress=10.0.1.7\n" {
			return fmt.Errorf("content, found %q", *u.Contents)
		}

		if *u.Mask != true {
			return fmt.Errorf("mask, found %t", *u.Mask)
		}

		if *u.Enabled == false {
			return fmt.Errorf("enabled, found %t", *u.Enabled)
		}

		if len(u.Dropins) != 1 {
			return fmt.Errorf("dropins, found %v", u.Dropins)
		}

		return nil
	})
}

func TestIgnitionSystemdUnitEmptyContentWithDropIn(t *testing.T) {
	testIgnition(t, `
		data "ignition_systemd_unit" "foo" {
			name = "foo.service"
			dropin {
				name = "foo.conf"
				content = "[Match]\nName=eth0\n\n[Network]\nAddress=10.0.1.7\n"
			}
		}

		data "ignition_config" "test" {
			systemd = [data.ignition_systemd_unit.foo.rendered]
		}
	`, func(c *types.Config) error {
		if len(c.Systemd.Units) != 1 {
			return fmt.Errorf("systemd, found %d", len(c.Systemd.Units))
		}

		u := c.Systemd.Units[0]

		if u.Name != "foo.service" {
			return fmt.Errorf("name, found %q", u.Name)
		}

		if u.Contents != nil {
			return fmt.Errorf("content, found %q", *u.Contents)
		}

		if len(u.Dropins) != 1 {
			return fmt.Errorf("dropins, found %v", u.Dropins)
		}

		return nil
	})
}

// #11325
func TestIgnitionSystemdUnit_emptyContent(t *testing.T) {
	testIgnition(t, `
		data "ignition_systemd_unit" "foo" {
			name = "foo.service"
			enabled = true
		}

		data "ignition_config" "test" {
			systemd = [data.ignition_systemd_unit.foo.rendered]
		}
	`, func(c *types.Config) error {
		if len(c.Systemd.Units) != 1 {
			return fmt.Errorf("systemd, found %d", len(c.Systemd.Units))
		}

		u := c.Systemd.Units[0]
		if u.Name != "foo.service" {
			return fmt.Errorf("name, expected 'foo.service', found %q", u.Name)
		}
		if u.Contents != nil {
			return fmt.Errorf("expected empty content, found %q", *u.Contents)
		}
		if len(u.Dropins) != 0 {
			return fmt.Errorf("expected 0 dropins, found %v", u.Dropins)
		}
		return nil
	})
}

func TestIgnitionSystemUnitInvalidName(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_systemd_unit" "foo" {
			name = "foo"
			enabled = true
		}

		data "ignition_config" "test" {
			systemd = [data.ignition_systemd_unit.foo.rendered]
		}
	`, regexp.MustCompile("invalid"))
}

func TestIgnitionSystemUnitInvalidContent(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_systemd_unit" "foo" {
			name = "foo.service"
			enabled = true
			content = "[foo"
		}

		data "ignition_config" "test" {
			systemd = [data.ignition_systemd_unit.foo.rendered]
		}
	`, regexp.MustCompile("unable to find end of section"))
}
