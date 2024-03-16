package ignition

import (
	"fmt"
	"testing"

	"github.com/coreos/ignition/v2/config/v3_4/types"
)

func TestIgnitionKargs(t *testing.T) {
	testIgnition(t, `
		data "ignition_kernel_arguments" "foo" {
			shouldexist = ["foo","bar"]
			shouldnotexist = ["baz"]
		}

		data "ignition_config" "test" {
		   kernel_arguments = data.ignition_kernel_arguments.foo.rendered
		}
	`, func(c *types.Config) error {
		if len(c.KernelArguments.ShouldExist) != 2 {
			return fmt.Errorf("shouldexist, found %d", len(c.KernelArguments.ShouldExist))
		}
		if len(c.KernelArguments.ShouldNotExist) != 1 {
			return fmt.Errorf("shouldnotexist, found %d", len(c.KernelArguments.ShouldNotExist))
		}

		if c.KernelArguments.ShouldExist[0] != "foo" {
			return fmt.Errorf("Field ShouldExist didn't match. Expected: %s, Given: %s", "foo", c.KernelArguments.ShouldExist[0])
		}
		if c.KernelArguments.ShouldNotExist[0] != "baz" {
			return fmt.Errorf("Field ShouldNotExist didn't match. Expected: %s, Given: %s", "baz", c.KernelArguments.ShouldNotExist[0])
		}
		return nil
	})
}
