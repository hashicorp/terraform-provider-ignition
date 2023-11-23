package ignition

import (
	"fmt"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_4/types"
)

func TestIgnitionLuks(t *testing.T) {
	testIgnition(t, `
		data "ignition_luks" "foo2" {
			name = "foo"
			device = "/foo"
			clevis {
				custom {
					pin = "4242"
					config = "cfg"
					needs_network = true
				}
			}
		}

		data "ignition_luks" "foo" {
			name = "foo"
			device = "/foo"
			label = "FOO"
			discard = true
			open_options = ["there"]
			options = ["aes"]
			uuid = "uuid"
			wipe_volume = true
			clevis {
				tpm2 = true
				threshold = 42
				tang {
					url = "url"
					thumbprint = "thumbprint"
					advertisement = "advertisement"
				}
			}
			key_file {
				source = "https://bar"
			}
		}

		data "ignition_config" "test" {
			luks = [
					data.ignition_luks.foo.rendered,
					data.ignition_luks.foo2.rendered,
			]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Luks) != 2 {
			return fmt.Errorf("luks, found %d", len(c.Storage.Luks))
		}

		a := c.Storage.Luks[0]
		if a.Name != "foo" {
			return fmt.Errorf("name, found %q", a.Name)
		}

		if *a.Device != "/foo" {
			return fmt.Errorf("device, found %q", *a.Device)
		}

		if !*a.Discard {
			return fmt.Errorf("discard, found %t", *a.Discard)
		}

		if *a.Label != "FOO" {
			return fmt.Errorf("label, found %q", *a.Label)
		}

		if *a.UUID != "uuid" {
			return fmt.Errorf("uuid, found %q", *a.UUID)
		}

		if !*a.WipeVolume {
			return fmt.Errorf("wipeVolume, found %t", *a.WipeVolume)
		}

		if a.OpenOptions[0] != "there" {
			return fmt.Errorf("open_options, found %q", a.OpenOptions)
		}

		if a.Options[0] != "aes" {
			return fmt.Errorf("options, found %q", a.Options)
		}

		if *a.KeyFile.Source != "https://bar" {
			return fmt.Errorf("key_file.source, found %q", *a.KeyFile.Source)
		}

		if !*a.Clevis.Tpm2 {
			return fmt.Errorf("clevis.tmp2, found %t", *a.Clevis.Tpm2)
		}

		if *a.Clevis.Threshold != 42 {
			return fmt.Errorf("clevis.threshold, found %d", *a.Clevis.Threshold)
		}

		b := c.Storage.Luks[1]
		if *b.Clevis.Custom.Pin != "4242" {
			return fmt.Errorf("clevis.custom.pin, found %q", *b.Clevis.Custom.Pin)
		}
		if *b.Clevis.Custom.Config != "cfg" {
			return fmt.Errorf("clevis.custom.cfg, found %q", *b.Clevis.Custom.Config)
		}
		if !*b.Clevis.Custom.NeedsNetwork {
			return fmt.Errorf("clevis.custom.needs_network, found %t", *b.Clevis.Custom.NeedsNetwork)
		}

		if a.Clevis.Tang[0].URL != "url" {
			return fmt.Errorf("clevis.tang.url, found %q", b.Clevis.Tang[0].URL)
		}
		if *a.Clevis.Tang[0].Thumbprint != "thumbprint" {
			return fmt.Errorf("clevis.tang.thumbprint, found %q", *b.Clevis.Tang[0].Thumbprint)
		}
		if *a.Clevis.Tang[0].Advertisement != "advertisement" {
			return fmt.Errorf("clevis.tang.advertisement, found %q", *b.Clevis.Tang[0].Advertisement)
		}

		return nil
	})
}
