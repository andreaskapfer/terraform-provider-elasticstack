package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"net/http"
)

func (a *ApiClient) PutElasticsearchStoredScript(script *models.StoredScript) diag.Diagnostics {
	var diags diag.Diagnostics
	scriptBytes, err := json.Marshal(script)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[TRACE] sending request to ES: %s", scriptBytes)
	res, err := a.es.PutScript(script.Name, bytes.NewReader(scriptBytes))
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()
	if diags := utils.CheckError(res, "Unable to create or update a stored script"); diags.HasError() {
		return diags
	}
	return diags
}

func (a *ApiClient) GetElasticsearchStoredScript(scriptId string) (*models.StoredScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	res, err := a.es.GetScript(scriptId)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find a stored script in the cluster.",
			Detail:   fmt.Sprintf("Unable to get stored script: '%s' from the cluster.", scriptId),
		})
		return nil, diags
	}
	if diags := utils.CheckError(res, "Unable to get a stored script."); diags.HasError() {
		return nil, diags
	}

	// unmarshal our response to proper type
	var storedScript models.StoredScript
	if err := json.NewDecoder(res.Body).Decode(&storedScript); err != nil {
		return nil, diag.FromErr(err)
	}
	log.Printf("[TRACE] Fetch stored scripts from ES API: %#+v", storedScript)

	return &storedScript, diags
}
