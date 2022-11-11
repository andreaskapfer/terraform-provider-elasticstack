package script

import (
	"context"
	"encoding/json"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceStoredScript() *schema.Resource {
	storedScriptSchema := map[string]*schema.Schema{
		"id": {
			Description: "Internal identifier of the resource",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the stored script.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"script": {
			Description: "Contains the script or search template, its parameters, and its language.",
			Type:        schema.TypeList,
			MaxItems:    1,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"lang": {
						Description: "Script language. For search templates, use mustache.",
						Type:        schema.TypeString,
						Required:    true,
						// TODO Validator for painless, expression, mustache, java
					},
					"source": {
						Description: "The script/search template",
						Type:        schema.TypeList,
						MaxItems:    1,
						Required:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"script": {
									Description:  "The definition of the script.",
									Type:         schema.TypeString,
									ExactlyOneOf: []string{"script", "search_template"},
								},
								"search_template": {
									Description:      "An object containing the search template. The object supports the same parameters as the search API's request body. Also supports Mustache variables. See Search templates.",
									Type:             schema.TypeString,
									ExactlyOneOf:     []string{"script", "search_template"},
									DiffSuppressFunc: utils.DiffJsonSuppress,
									ValidateFunc:     validation.StringIsJSON,
								},
							},
						},
						// TODO Validator for search templates, i.e. lang:mustache
					},
					"params": {
						Description:      "Parameters for the script or search template.",
						Type:             schema.TypeString,
						Optional:         true,
						DiffSuppressFunc: utils.DiffJsonSuppress,
						ValidateFunc:     validation.StringIsJSON,
						Default:          "{}",
					},
				},
			},
		},
	}
	utils.AddConnectionSchema(storedScriptSchema)

	return &schema.Resource{
		Description: "Creates or updates a stored script or search template.",

		Schema: storedScriptSchema,
	}
}

func resourceStoredScriptPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clients.NewApiClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	storedScriptName := d.Get("name").(string)
	id, diags := client.ID(storedScriptName)
	if diags.HasError() {
		return diags
	}
	var storedScript models.StoredScript
	storedScript.Name = storedScriptName

	if v, ok := d.GetOk("script"); ok {
		var script models.Script
		definedScript := v.([]interface{})[0].(map[string]interface{})
		if lang, ok := definedScript["lang"]; ok {
			script.Lang = lang.(string)
		}
		if source, ok := definedScript["source"]; ok {
			definedSource := source.([]interface{})[0].(map[string]interface{})
			if scriptSource, ok := definedSource["script"]; ok {
				script.ScriptSource = scriptSource.(string)
			}
			if searchTemplateSource, ok := definedSource["search_template"]; ok {
				if searchTemplateSource.(string) != "" {
					searchTemplateSourceObject := make(map[string]interface{})
					if err := json.Unmarshal([]byte(searchTemplateSource.(string)), &searchTemplateSourceObject); err != nil {
						return diag.FromErr(err)
					}
					script.SearchTemplateSource = searchTemplateSourceObject
				}
			}
		}
		if params, ok := definedScript["params"]; ok {
			if params.(string) != "" {
				paramsObject := make(map[string]interface{})
				if err := json.Unmarshal([]byte(params.(string)), &paramsObject); err != nil {
					return diag.FromErr(err)
				}
				script.Params = paramsObject
			}
		}
	}

	if diags := client.PutElasticsearchStoredScript(&storedScript); diags.HasError() {
		return diags
	}

	d.SetId(id.String())
	return resourceStoredScriptRead(ctx, d, meta)
}

func resourceStoredScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := clients.NewApiClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	compId, diags := clients.CompositeIdFromStr(d.Id())
	if diags.HasError() {
		return diags
	}
	storedScriptId := compId.ResourceId

	storedScript, diags := client.GetElasticsearchScript(storedScriptId)
	if storedScript == nil && diags == nil {
		d.SetId("")
		return diags
	}
	if diags.HasError() {
		return diags
	}

	// set the fields
	if err := d.Set("name", storedScriptId); err != nil {
		return diag.FromErr(err)
	}
	if storedScript.Script != nil {
		script := make([]map[string]interface{}, 1)
		script[0] = make(map[string]interface{})
		if storedScript.Script.Lang != nil {
			script[0]["lang"] = storedScript.Script.Lang
		}
		if storedScript.Script.ScriptSource != nil {
			script[0]["source"] = []map[string]string{{"script": storedScript.Script.ScriptSource}}
		}
		if storedScript.Script.SearchTemplateSource != nil {
			source, err := json.Marshal(storedScript.Script.SearchTemplateSource)
			if err != nil {
				return diag.FromErr(err)
			}
			script[0]["source"] = []map[string]string{{"search_template": string(source)}}
		}
		if storedScript.Script.Params != nil {
			params, err := json.Marshal(storedScript.Params)
			if err != nil {
				return diag.FromErr(err)
			}
			script[0]["params"] = string(params)
		}
		if err := d.Set("script", script); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func resourceStoredScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := clients.NewApiClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	compId, diags := clients.CompositeIdFromStr(d.Id())
	if diags.HasError() {
		return diags
	}

	if diags := client.DeleteElasticsearchStoredScript(compId.ResourceId); diags.HasError() {
		return diags
	}

	d.SetId("")
	return diags
}
