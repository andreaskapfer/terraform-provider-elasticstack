---
subcategory: "Security"
layout: ""
page_title: "Elasticstack: elasticstack_elasticsearch_security_role_mapping Resource"
description: |-
  Manage role mappings. See, https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-role-mapping.html
---

# Resource: elasticstack_elasticsearch_security_role_mapping

Manage role mappings. See, https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-put-role-mapping.html

## Example Usage

```terraform
provider "elasticstack" {
  elasticsearch {}
}

resource "elasticstack_elasticsearch_security_role_mapping" "example" {
  name    = "test_role_mapping"
  enabled = true
  roles = [
    "admin"
  ]
  rules = jsonencode({
    any = [
      { field = { username = "esadmin" } },
      { field = { groups = "cn=admins,dc=example,dc=com" } },
    ]
  })
}

output "role" {
  value = elasticstack_elasticsearch_security_role_mapping.example.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The distinct name that identifies the role mapping, used solely as an identifier.
- `rules` (String) The rules that determine which users should be matched by the mapping. A rule is a logical condition that is expressed by using a JSON DSL.

### Optional

- `elasticsearch_connection` (Block List, Max: 1) Used to establish connection to Elasticsearch server. Overrides environment variables if present. (see [below for nested schema](#nestedblock--elasticsearch_connection))
- `enabled` (Boolean) Mappings that have `enabled` set to `false` are ignored when role mapping is performed.
- `metadata` (String) Additional metadata that helps define which roles are assigned to each user. Keys beginning with `_` are reserved for system usage.
- `role_templates` (String) A list of mustache templates that will be evaluated to determine the roles names that should granted to the users that match the role mapping rules.
- `roles` (Set of String) A list of role names that are granted to the users that match the role mapping rules.

### Read-Only

- `id` (String) Internal identifier of the resource

<a id="nestedblock--elasticsearch_connection"></a>
### Nested Schema for `elasticsearch_connection`

Optional:

- `api_key` (String, Sensitive) API Key to use for authentication to Elasticsearch
- `ca_data` (String) PEM-encoded custom Certificate Authority certificate
- `ca_file` (String) Path to a custom Certificate Authority certificate
- `endpoints` (List of String, Sensitive) A list of endpoints the Terraform provider will point to. They must include the http(s) schema and port number.
- `insecure` (Boolean) Disable TLS certificate validation
- `password` (String, Sensitive) A password to use for API authentication to Elasticsearch.
- `username` (String) A username to use for API authentication to Elasticsearch.

## Import

Import is supported using the following syntax:

```shell
terraform import elasticstack_elasticsearch_security_role_mapping.my_role_mapping <cluster_uuid>/<role mapping name>
```
