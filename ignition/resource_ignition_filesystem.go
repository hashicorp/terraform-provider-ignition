package ignition

import (
	"fmt"

	"github.com/coreos/ignition/config/v2_1/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFilesystem() *schema.Resource {
	return &schema.Resource{
		Exists: resourceFilesystemExists,
		Read:   resourceFilesystemRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mount": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"create": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"force": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"options": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFilesystemRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildFilesystem(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceFilesystemExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildFilesystem(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildFilesystem(d *schema.ResourceData, c *cache) (string, error) {
	var mount *types.Mount
	if _, ok := d.GetOk("mount"); ok {
		mount = &types.Mount{
			Device: d.Get("mount.0.device").(string),
			Format: d.Get("mount.0.format").(string),
		}

		create, hasCreate := d.GetOk("mount.0.create")
		force, hasForce := d.GetOk("mount.0.force")
		options, hasOptions := d.GetOk("mount.0.options")
		if hasCreate || hasOptions || hasForce {
			mount.Create = &types.Create{
				Force:   force.(bool),
				Options: castSliceInterfaceToCreateOption(options.([]interface{})),
			}
		}

		if !create.(bool) && (hasForce || hasOptions) {
			return "", fmt.Errorf("create should be true when force or options is used")
		}
	}

	var path *string
	if p, ok := d.GetOk("path"); ok {
		str := p.(string)
		path = &str
	}

	if mount != nil && path != nil {
		return "", fmt.Errorf("mount and path are mutually exclusive")
	}

	return c.addFilesystem(&types.Filesystem{
		Name:  d.Get("name").(string),
		Mount: mount,
		Path:  path,
	}), nil
}

func castSliceInterfaceToCreateOption(i []interface{}) []types.CreateOption {
	var o []types.CreateOption
	for _, value := range i {
		if value == nil {
			continue
		}

		o = append(o, types.CreateOption(value.(string)))
	}

	return o
}
