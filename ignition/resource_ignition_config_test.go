package ignition

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_4/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestIgnitionFileReplace(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			replace {
				source = "foo"
				verification = "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
			}
		}
	`, func(c *types.Config) error {
		r := c.Ignition.Config.Replace
		if r.Source == nil {
			return fmt.Errorf("unable to find replace config")
		}

		if *r.Source != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *r.Source)
		}

		if *r.Verification.Hash != "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" {
			return fmt.Errorf("config.replace.verification, found %q", *r.Verification.Hash)
		}
		if r.Compression != nil {
			return fmt.Errorf("config.replace.compression, found %q", *r.Compression)
		}

		return nil
	})
}

func TestIgnitionFileReplaceWithCompression(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			replace {
				source = "foo"
				compression = "gzip"
				verification = "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
			}
		}
	`, func(c *types.Config) error {
		r := c.Ignition.Config.Replace
		if r.Source == nil {
			return fmt.Errorf("unable to find replace config")
		}

		if *r.Source != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *r.Source)
		}

		if *r.Verification.Hash != "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" {
			return fmt.Errorf("config.replace.verification, found %q", *r.Verification.Hash)
		}
		if *r.Compression != "gzip" {
			return fmt.Errorf("config.replace.compression, found %q", *r.Compression)
		}

		return nil
	})
}

func TestIgnitionFileMerge(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			merge {
				source = "foo"
				verification = "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
			}

		    merge {
		    	source = "foo"
		    	verification = "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
			}
		}
	`, func(c *types.Config) error {
		a := c.Ignition.Config.Merge
		if len(a) != 2 {
			return fmt.Errorf("unable to find merge config, expected 2")
		}

		if string(*a[0].Source) != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *a[0].Source)
		}

		if *a[0].Verification.Hash != "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" {
			return fmt.Errorf("config.replace.verification, found %q", *a[0].Verification.Hash)
		}

		return nil
	})
}

func TestIgnitionFileReplaceNoVerification(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			replace {
				source = "foo"
			}
		}
	`, func(c *types.Config) error {
		r := c.Ignition.Config.Replace
		if r.Source == nil {
			return fmt.Errorf("unable to find replace config")
		}

		if string(*r.Source) != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *r.Source)
		}

		if r.Verification.Hash != nil {
			return fmt.Errorf("verification hash should be nil")
		}

		return nil
	})
}

func TestIgnitionFileMergeNoVerification(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			merge {
				source = "foo"
			}

			merge {
				source = "foo"
			}
		}
	`, func(c *types.Config) error {
		a := c.Ignition.Config.Merge
		if len(a) != 2 {
			return fmt.Errorf("unable to find merge config, expected 2")
		}

		if string(*a[0].Source) != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *a[0].Source)
		}

		if a[0].Verification.Hash != nil {
			return fmt.Errorf("verification hash should be nil")
		}

		return nil
	})
}

func TestIgnitionConfigReplaceWithHttpHeaders(t *testing.T) {
	testIgnition(t, `
		data "ignition_config" "test" {
			replace {
				source = "foo"
				http_headers {
					name = "Authorization"
					value = "Basic token"
				}
				http_headers {
					name = "Cache-Control"
					value = "no-cache"
				}
			}
		}
	`, func(c *types.Config) error {
		r := c.Ignition.Config.Replace
		if r.Source == nil {
			return fmt.Errorf("unable to find replace config")
		}

		if *r.Source != "foo" {
			return fmt.Errorf("config.replace.source, found %q", *r.Source)
		}

		if len(r.HTTPHeaders) != 2 {
			return fmt.Errorf("config.replace.http_headers, found %d", len(r.HTTPHeaders))
		}

		if string(r.HTTPHeaders[0].Name) != "Authorization" {
			return fmt.Errorf("config.replace.http_headers[0].name, found %q", r.HTTPHeaders[0].Name)
		}

		if string(*r.HTTPHeaders[1].Value) != "no-cache" {
			return fmt.Errorf("config.replace.http_headers[1].value, found %q", *r.HTTPHeaders[1].Value)
		}

		return nil
	})
}

func TestIgnitionConfigSecurityTLSCAs(t *testing.T) {
	testIgnition(t, `
	data "ignition_config" "test" {
		replace {
			source = "foo"
		}
		tls_ca {
			source = "bar"
		}
		tls_ca {
			source = "baz"
			verification = "sha512-012345678912789abcdef"
		}
	}
	`, func(c *types.Config) error {
		cas := c.Ignition.Security.TLS.CertificateAuthorities
		if len(cas) != 2 {
			return fmt.Errorf("config.tls_ca, found %d", len(cas))
		}
		if string(*cas[0].Source) != "bar" {
			return fmt.Errorf("config.tls_ca.source, found %q", *cas[0].Source)
		}
		if string(*cas[1].Source) != "baz" {
			return fmt.Errorf("config.tls_ca.source, found %q", *cas[1].Source)
		}
		if cas[1].Verification.Hash == nil || *cas[1].Verification.Hash != "sha512-012345678912789abcdef" {
			return fmt.Errorf("config.tls_ca.verification, found %q", *cas[1].Verification.Hash)
		}

		return nil
	})
}

func TestIgnitionConfigDisks(t *testing.T) {
	testIgnition(t, `
	variable "ignition_disk_renders" {
		type = list
		default = [""]
	}

	data "ignition_disk" "test" {
		device = "/dev/sda"
		partition {
			startmib = 2048
			sizemib = 20480
		}
	 }

	data "ignition_config" "test" {
		disks = concat([data.ignition_disk.test.rendered],
			var.ignition_disk_renders)
	}
	`, func(c *types.Config) error {
		f := c.Storage.Disks[0]
		if f.Device != "/dev/sda" {
			return fmt.Errorf("device, found %q", f.Device)
		}
		return nil
	})
}

func TestIgnitionConfigLuks(t *testing.T) {
	testIgnition(t, `
	variable "ignition_luks_renders" {
		type = list
		default = [""]
	}

	data "ignition_luks" "test" {
		device = "/dev/sda"
		name = "data"
	 }

	data "ignition_config" "test" {
		luks = concat([data.ignition_luks.test.rendered],
			var.ignition_luks_renders)
	}
	`, func(c *types.Config) error {
		f := c.Storage.Luks[0]
		if *f.Device != "/dev/sda" {
			return fmt.Errorf("device, found %q", *f.Device)
		}
		return nil
	})
}

func TestIgnitionConfigArrays(t *testing.T) {
	testIgnition(t, `
	variable "ignition_array_renders" {
		type = list
		default = [""]
	}

	data "ignition_raid" "md" {
		name = "data"
		level = "stripe"
		devices = [
			"/dev/disk/by-partlabel/raid.1.1",
			"/dev/disk/by-partlabel/raid.1.2"
		]
	}

	data "ignition_config" "test" {
		arrays = concat([data.ignition_raid.md.rendered],
			var.ignition_array_renders)
	}
	`, func(c *types.Config) error {
		f := c.Storage.Raid[0]
		if f.Name != "data" {
			return fmt.Errorf("device, found %q", f.Name)
		}
		return nil
	})
}

func TestIgnitionConfigFilesystems(t *testing.T) {
	testIgnition(t, `
	variable "ignition_filesystem_renders" {
		type = list
		default = [""]
	}

	data "ignition_filesystem" "test" {
		path = "/test"
		device = "/dev/sda"
		format = "ext4"
	 }

	data "ignition_config" "test" {
		filesystems = concat(
			[data.ignition_filesystem.test.rendered],
			var.ignition_filesystem_renders
		)
	}
	`, func(c *types.Config) error {
		f := c.Storage.Filesystems[0]
		if string(*f.Path) != "/test" {
			return fmt.Errorf("device, found %q", *f.Path)
		}
		return nil
	})
}

func TestIgnitionConfigFiles(t *testing.T) {
	testIgnition(t, `
	variable "ignition_file_renders" {
		type = list
		default = [""]
	}

	data "ignition_file" "test" {
		path = "/hello.txt"
		content {
			content = "Hello World!"
		}
	 }

	data "ignition_config" "test" {
		files = concat(
			[data.ignition_file.test.rendered],
			var.ignition_file_renders,
		)
	}
	`, func(c *types.Config) error {
		f := c.Storage.Files[0]
		if f.Path != "/hello.txt" {
			return fmt.Errorf("device, found %q", f.Path)
		}
		return nil
	})
}

func TestIgnitionConfigSystemd(t *testing.T) {
	testIgnition(t, `
	variable "ignition_systemd_renders" {
		type = list
		default = [""]
	}

	data "ignition_systemd_unit" "test" {
		name = "example.service"
		content = "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
	}

	data "ignition_config" "test" {
		systemd = concat(
			[data.ignition_systemd_unit.test.rendered],
			var.ignition_systemd_renders,
		)
	}
	`, func(c *types.Config) error {
		f := c.Systemd.Units[0]
		if f.Name != "example.service" {
			return fmt.Errorf("device, found %q", f.Name)
		}
		return nil
	})
}

func TestIgnitionConfigUsers(t *testing.T) {
	testIgnition(t, `
	variable "ignition_user_renders" {
		type = list
		default = [""]
	}

	data "ignition_user" "test" {
		name = "foo"
		home_dir = "/home/foo/"
		shell = "/bin/bash"
	}

	data "ignition_config" "test" {
		users = concat(
			[data.ignition_user.test.rendered],
			var.ignition_user_renders
		)
	}
	`, func(c *types.Config) error {
		f := c.Passwd.Users[0]
		if f.Name != "foo" {
			return fmt.Errorf("device, found %q", f.Name)
		}
		return nil
	})
}

func TestIgnitionConfigGroups(t *testing.T) {
	testIgnition(t, `
	variable "ignition_group_renders" {
		type = list
		default = [""]
	}

	data "ignition_group" "test" {
		name = "test"
	}

	data "ignition_config" "test" {
		groups = concat(
			[data.ignition_group.test.rendered],
			var.ignition_group_renders
		)
	}
	`, func(c *types.Config) error {
		f := c.Passwd.Groups[0]
		if f.Name != "test" {
			return fmt.Errorf("device, found %q", f.Name)
		}
		return nil
	})
}

func testIgnitionError(t *testing.T, input string, expectedErr *regexp.Regexp) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testTemplate, input),
				ExpectError: expectedErr,
			},
		},
	})
}

func testIgnition(t *testing.T, input string, assert func(*types.Config) error) {
	check := func(s *terraform.State) error {
		got := s.RootModule().Outputs["rendered"].Value.(string)

		c := &types.Config{}
		err := json.Unmarshal([]byte(got), c)
		if err != nil {
			return err
		}

		return assert(c)
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testTemplate, input),
				Check:  check,
			},
		},
	})
}

var testTemplate = `
%s

output "rendered" {
	value = data.ignition_config.test.rendered
}

`
