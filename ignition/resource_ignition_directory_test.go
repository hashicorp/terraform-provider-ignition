package ignition

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_4/types"
)

func TestIgnitionDirectory(t *testing.T) {
	testIgnition(t, `
		data "ignition_directory" "foo" {
			path = "/foo"
			mode = 420
			uid = 42
			gid = 84
		}

		data "ignition_directory" "bar" {
			path = "/bar"
			overwrite = true
		}

		data "ignition_config" "test" {
			directories = [
				data.ignition_directory.foo.rendered,
				data.ignition_directory.bar.rendered,
			]
		}
	`, func(c *types.Config) error {
		if len(c.Storage.Directories) != 2 {
			return fmt.Errorf("arrays, found %d", len(c.Storage.Directories))
		}

		f := c.Storage.Directories[0]
		if f.Path != "/foo" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != false {
			return fmt.Errorf("overwrite, found %t", *f.Overwrite)
		}

		if int(*f.Mode) != 420 {
			return fmt.Errorf("mode, found %q", *f.Mode)
		}

		if *f.User.ID != 42 {
			return fmt.Errorf("uid, found %q", *f.User.ID)
		}

		if *f.Group.ID != 84 {
			return fmt.Errorf("gid, found %q", *f.Group.ID)
		}

		f = c.Storage.Directories[1]
		if f.Path != "/bar" {
			return fmt.Errorf("path, found %q", f.Path)
		}

		if *f.Overwrite != true {
			return fmt.Errorf("overwrite, found %t", *f.Overwrite)
		}

		return nil
	})
}

func TestIgnitionDirectoryInvalidMode(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "/foo"
			mode = 999999
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("illegal file mode"))
}

func TestIgnitionDirectoryInvalidPath(t *testing.T) {
	testIgnitionError(t, `
		data "ignition_directory" "foo" {
			path = "foo"
			mode = 999999
		}

		data "ignition_config" "test" {
			directories = [data.ignition_directory.foo.rendered]
		}
	`, regexp.MustCompile("path not absolute"))
}
