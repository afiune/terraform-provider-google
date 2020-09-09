// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceKMSCryptoKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKMSCryptoKeyCreate,
		Read:   resourceKMSCryptoKeyRead,
		Update: resourceKMSCryptoKeyUpdate,
		Delete: resourceKMSCryptoKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceKMSCryptoKeyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceKMSCryptoKeyResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceKMSCryptoKeyUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"key_ring": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: kmsCryptoKeyRingsEquivalent,
				Description: `The KeyRing that this key belongs to.
Format: ''projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}''.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource name for the CryptoKey.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Labels with user-defined metadata to apply to this resource.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"purpose": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENCRYPT_DECRYPT", "ASYMMETRIC_SIGN", "ASYMMETRIC_DECRYPT", ""}, false),
				Description: `The immutable purpose of this CryptoKey. See the
[purpose reference](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyRings.cryptoKeys#CryptoKeyPurpose)
for possible inputs. Default value: "ENCRYPT_DECRYPT" Possible values: ["ENCRYPT_DECRYPT", "ASYMMETRIC_SIGN", "ASYMMETRIC_DECRYPT"]`,
				Default: "ENCRYPT_DECRYPT",
			},
			"rotation_period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: orEmpty(validateKmsCryptoKeyRotationPeriod),
				Description: `Every time this period passes, generate a new CryptoKeyVersion and set it as the primary.
The first rotation will take place after the specified period. The rotation period has
the format of a decimal number with up to 9 fractional digits, followed by the
letter 's' (seconds). It must be greater than a day (ie, 86400).`,
			},
			"version_template": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `A template describing settings for new crypto key versions.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
							Description: `The algorithm to use when creating a version based on this template.
See the [algorithm reference](https://cloud.google.com/kms/docs/reference/rest/v1/CryptoKeyVersionAlgorithm) for possible inputs.`,
						},
						"protection_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM", ""}, false),
							Description:  `The protection level to use when creating a version based on this template. Default value: "SOFTWARE" Possible values: ["SOFTWARE", "HSM"]`,
							Default:      "SOFTWARE",
						},
					},
				},
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKMSCryptoKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	labelsProp, err := expandKMSCryptoKeyLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	purposeProp, err := expandKMSCryptoKeyPurpose(d.Get("purpose"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("purpose"); !isEmptyValue(reflect.ValueOf(purposeProp)) && (ok || !reflect.DeepEqual(v, purposeProp)) {
		obj["purpose"] = purposeProp
	}
	rotationPeriodProp, err := expandKMSCryptoKeyRotationPeriod(d.Get("rotation_period"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rotation_period"); !isEmptyValue(reflect.ValueOf(rotationPeriodProp)) && (ok || !reflect.DeepEqual(v, rotationPeriodProp)) {
		obj["rotationPeriod"] = rotationPeriodProp
	}
	versionTemplateProp, err := expandKMSCryptoKeyVersionTemplate(d.Get("version_template"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_template"); !isEmptyValue(reflect.ValueOf(versionTemplateProp)) && (ok || !reflect.DeepEqual(v, versionTemplateProp)) {
		obj["versionTemplate"] = versionTemplateProp
	}

	obj, err = resourceKMSCryptoKeyEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{KMSBasePath}}{{key_ring}}/cryptoKeys?cryptoKeyId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new CryptoKey: %#v", obj)
	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating CryptoKey: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating CryptoKey %q: %#v", d.Id(), res)

	return resourceKMSCryptoKeyRead(d, meta)
}

func resourceKMSCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{KMSBasePath}}{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("KMSCryptoKey %q", d.Id()))
	}

	res, err = resourceKMSCryptoKeyDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing KMSCryptoKey because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("labels", flattenKMSCryptoKeyLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	if err := d.Set("purpose", flattenKMSCryptoKeyPurpose(res["purpose"], d, config)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	if err := d.Set("rotation_period", flattenKMSCryptoKeyRotationPeriod(res["rotationPeriod"], d, config)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	if err := d.Set("version_template", flattenKMSCryptoKeyVersionTemplate(res["versionTemplate"], d, config)); err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}

	return nil
}

func resourceKMSCryptoKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingProject := ""

	obj := make(map[string]interface{})
	labelsProp, err := expandKMSCryptoKeyLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	rotationPeriodProp, err := expandKMSCryptoKeyRotationPeriod(d.Get("rotation_period"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("rotation_period"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, rotationPeriodProp)) {
		obj["rotationPeriod"] = rotationPeriodProp
	}
	versionTemplateProp, err := expandKMSCryptoKeyVersionTemplate(d.Get("version_template"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("version_template"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, versionTemplateProp)) {
		obj["versionTemplate"] = versionTemplateProp
	}

	obj, err = resourceKMSCryptoKeyUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{KMSBasePath}}{{key_ring}}/cryptoKeys/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating CryptoKey %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("labels") {
		updateMask = append(updateMask, "labels")
	}

	if d.HasChange("rotation_period") {
		updateMask = append(updateMask, "rotationPeriod",
			"nextRotationTime")
	}

	if d.HasChange("version_template") {
		updateMask = append(updateMask, "versionTemplate.algorithm")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating CryptoKey %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating CryptoKey %q: %#v", d.Id(), res)
	}

	return resourceKMSCryptoKeyRead(d, meta)
}

func resourceKMSCryptoKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf(`
[WARNING] KMS CryptoKey resources cannot be deleted from GCP. The CryptoKey %s will be removed from Terraform state,
and all its CryptoKeyVersions will be destroyed, but it will still be present on the server.`, cryptoKeyId.cryptoKeyId())

	// Delete all versions of the key
	if err := clearCryptoKeyVersions(cryptoKeyId, config); err != nil {
		return err
	}

	// Make sure automatic key rotation is disabled if set
	if d.Get("rotation_period") != "" {
		if err := disableCryptoKeyRotation(cryptoKeyId, config); err != nil {
			return fmt.Errorf(
				"While cryptoKeyVersions were cleared, Terraform was unable to disable automatic rotation of key due to an error: %s."+
					"Please retry or manually disable automatic rotation to prevent creation of a new version of this key.", err)
		}
	}

	d.SetId("")
	return nil
}

func resourceKMSCryptoKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.Set("key_ring", cryptoKeyId.KeyRingId.keyRingId())
	d.Set("name", cryptoKeyId.Name)

	return []*schema.ResourceData{d}, nil
}

func flattenKMSCryptoKeyLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenKMSCryptoKeyPurpose(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenKMSCryptoKeyRotationPeriod(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenKMSCryptoKeyVersionTemplate(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["algorithm"] =
		flattenKMSCryptoKeyVersionTemplateAlgorithm(original["algorithm"], d, config)
	transformed["protection_level"] =
		flattenKMSCryptoKeyVersionTemplateProtectionLevel(original["protectionLevel"], d, config)
	return []interface{}{transformed}
}
func flattenKMSCryptoKeyVersionTemplateAlgorithm(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenKMSCryptoKeyVersionTemplateProtectionLevel(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandKMSCryptoKeyLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandKMSCryptoKeyPurpose(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyRotationPeriod(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyVersionTemplate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAlgorithm, err := expandKMSCryptoKeyVersionTemplateAlgorithm(original["algorithm"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAlgorithm); val.IsValid() && !isEmptyValue(val) {
		transformed["algorithm"] = transformedAlgorithm
	}

	transformedProtectionLevel, err := expandKMSCryptoKeyVersionTemplateProtectionLevel(original["protection_level"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProtectionLevel); val.IsValid() && !isEmptyValue(val) {
		transformed["protectionLevel"] = transformedProtectionLevel
	}

	return transformed, nil
}

func expandKMSCryptoKeyVersionTemplateAlgorithm(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyVersionTemplateProtectionLevel(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceKMSCryptoKeyEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// if rotationPeriod is set, nextRotationTime must also be set.
	if d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return nil, fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		obj["nextRotationTime"] = nextRotation
	}

	return obj, nil
}

func resourceKMSCryptoKeyUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// if rotationPeriod is changed, nextRotationTime must also be set.
	if d.HasChange("rotation_period") && d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return nil, fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		obj["nextRotationTime"] = nextRotation
	}

	return obj, nil
}

func resourceKMSCryptoKeyDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	// Take the returned long form of the name and use it as `self_link`.
	// Then modify the name to be the user specified form.
	// We can't just ignore_read on `name` as the linter will
	// complain that the returned `res` is never used afterwards.
	// Some field needs to be actually set, and we chose `name`.
	d.Set("self_link", res["name"].(string))
	res["name"] = d.Get("name").(string)
	return res, nil
}

func resourceKMSCryptoKeyResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_ring": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rotation_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_template": {
				Type:     schema.TypeList,
				Optional: true,
			},
			"self_link": {
				Type: schema.TypeString,
			},
		},
	}
}

func resourceKMSCryptoKeyUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	config := meta.(*Config)
	keyRingId := rawState["key_ring"].(string)
	parsed, err := parseKmsKeyRingId(keyRingId, config)
	if err != nil {
		return nil, err
	}
	rawState["key_ring"] = parsed.keyRingId()

	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
