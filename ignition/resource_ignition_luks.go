package ignition

import (
	"encoding/json"
	"github.com/coreos/ignition/v2/config/v3_4/types"
	vcontext_path "github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLuks() *schema.Resource {
	return &schema.Resource{
		Exists: resourceLuksExists,
		Read:   resourceLuksRead,
		Schema: map[string]*schema.Schema{
			"clevis": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tang": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"thumbprint": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"advertisement": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"tpm2": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							ForceNew: true,
							Optional: true,
						},
						"custom": {
							Type:     schema.TypeList,
							ForceNew: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pin": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"config": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"needs_network": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"discard": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"key_file": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     configReferenceResource,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"open_options": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"options": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"wipe_volume": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLuksRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildLuks(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceLuksExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildLuks(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildLuks(d *schema.ResourceData) (string, error) {
	luks := &types.Luks{}

	if _, hasClevis := d.GetOk("clevis"); hasClevis {
		if tang, hasTang := d.GetOk("clevis.0.tang"); hasTang {
			luks.Clevis.Tang = castSliceInterfaceToTangs(tang.([]interface{}))
		}
		if tpm2, hasTpm2 := d.GetOk("clevis.0.tpm2"); hasTpm2 {
			b := tpm2.(bool)
			luks.Clevis.Tpm2 = &b
		}
		if threshold, hasThreshold := d.GetOk("clevis.0.threshold"); hasThreshold {
			i := threshold.(int)
			luks.Clevis.Threshold = &i
		}
		if custom, hasCustom := d.GetOk("clevis.0.custom"); hasCustom {
			c := custom.([]interface{})[0].(map[string]interface{})
			if config, hasConfig := c["config"]; hasConfig {
				str := config.(string)
				luks.Clevis.Custom.Config = &str
			}

			if needsNetwork, hasNeedsNetwork := c["needs_network"]; hasNeedsNetwork {
				b := needsNetwork.(bool)
				luks.Clevis.Custom.NeedsNetwork = &b
			}

			if pin, hasPin := c["pin"]; hasPin {
				str := pin.(string)
				luks.Clevis.Custom.Pin = &str
			}
		}
	}

	if device, hasDevice := d.GetOk("device"); hasDevice {
		str := device.(string)
		luks.Device = &str
	}

	if _, hasKeyFile := d.GetOk("key_file"); hasKeyFile {
		var keyFile types.Resource
		if err := fillResource(d, &keyFile, "key_file"); err != nil {
			return "", err
		}
		luks.KeyFile = keyFile
	}

	if discard, hasDiscard := d.GetOk("discard"); hasDiscard {
		b := discard.(bool)
		luks.Discard = &b
	}

	if label, hasLabel := d.GetOk("label"); hasLabel {
		str := label.(string)
		luks.Label = &str
	}

	luks.Name = d.Get("name").(string)

	for _, value := range d.Get("open_options").([]interface{}) {
		luks.OpenOptions = append(luks.OpenOptions, types.OpenOption(value.(string)))
	}

	for _, value := range d.Get("options").([]interface{}) {
		luks.Options = append(luks.Options, types.LuksOption(value.(string)))
	}

	if uuid, hasUuid := d.GetOk("uuid"); hasUuid {
		str := uuid.(string)
		luks.UUID = &str
	}

	if wipeVolume, hasWipeVolume := d.GetOk("wipe_volume"); hasWipeVolume {
		b := wipeVolume.(bool)
		luks.WipeVolume = &b
	}

	//

	b, err := json.Marshal(luks)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), handleReport(luks.Validate(vcontext_path.ContextPath{}))
}

func castSliceInterfaceToTangs(tangs []interface{}) []types.Tang {
	var res []types.Tang
	for _, t := range tangs {
		if t == nil {
			continue
		}
		tt := t.(map[string]interface{})
		tang := types.Tang{
			URL: tt["url"].(string),
		}
		if thumbprint, hasThumbprint := tt["thumbprint"]; hasThumbprint {
			str := thumbprint.(string)
			tang.Thumbprint = &str
		}
		if advertisement, hasAdvertisement := tt["advertisement"]; hasAdvertisement {
			str := advertisement.(string)
			tang.Advertisement = &str
		}

		res = append(res, tang)
	}
	return res
}
