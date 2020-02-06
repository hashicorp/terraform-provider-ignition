package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/config/v2_4/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDisk() *schema.Resource {
	return &schema.Resource{
		Exists: resourceDiskExists,
		Read:   resourceDiskRead,
		Schema: map[string]*schema.Schema{
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wipe_table": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"partition": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"start": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"type_guid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDiskRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildDisk(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceDiskExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildDisk(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildDisk(d *schema.ResourceData) (string, error) {
	disk := &types.Disk{
		Device:    d.Get("device").(string),
		WipeTable: d.Get("wipe_table").(bool),
	}

	if err := handleReport(disk.ValidateDevice()); err != nil {
		return "", err
	}

	for _, raw := range d.Get("partition").([]interface{}) {
		v := raw.(map[string]interface{})

		partitionLabel := v["label"].(string)
		partitionSize := v["size"].(int)
		partitionStart := v["start"].(int)
		p := types.Partition{
			Label:    &partitionLabel,
			Number:   v["number"].(int),
			Size:     &partitionSize,
			Start:    &partitionStart,
			TypeGUID: v["type_guid"].(string),
		}

		if err := handleReport(p.ValidateSize()); err != nil {
			return "", err
		}

		if err := handleReport(p.ValidateStart()); err != nil {
			return "", err
		}

		if err := handleReport(p.ValidateLabel()); err != nil {
			return "", err
		}

		if err := handleReport(p.ValidateGUID()); err != nil {
			return "", err
		}

		if err := handleReport(p.ValidateTypeGUID()); err != nil {
			return "", err
		}

		if err := handleReport(p.Validate()); err != nil {
			return "", err
		}

		disk.Partitions = append(disk.Partitions, p)
	}

	if err := handleReport(disk.ValidatePartitions()); err != nil {
		return "", err
	}

	b, err := json.Marshal(disk)
	if err != nil {
		return "", err
	}
	d.Set("rendered", string(b))

	return hash(string(b)), nil
}
