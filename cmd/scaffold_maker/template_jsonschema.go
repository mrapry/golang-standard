package main

//SchemaCreateTemplate template schema create
const SchemaCreateTemplate = `{
    "$schema": "http://json-schema.org/draft-07/schema",
    "title": "JSON Schema for create {{$.module}}",
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "minLength": 3
      }
    },
    "required": [
      "name"
    ],
    "additionalProperties": true
}`

//SchemaGetAllTemplate template schema get all
const SchemaGetAllTemplate = `{
    "$schema": "http://json-schema.org/draft-07/schema",
    "title": "JSON Schema for get all {{$.module}} parameter",
    "type": "object",
    "properties": {
        "page": {
            "type": "number",
            "default": 1,
            "minimum": 0
        },
        "limit": {
            "type": "number",
            "default": 10,
            "minimum": 1
        },
        "orderBy": {
            "type": "string",
            "enum": ["name"]
        },
        "sort": {
            "type": "string",
            "enum": ["asc", "desc"]
        },
        "search": {
            "type": "string"
        },
        "showAll": {
            "type": "boolean"
        }
    },
    "dependencies": {
        "sortBy": ["orderBy"]
    },
    "additionalProperties": true
}
`

//SchemaUpdateTemplate template schema update
const SchemaUpdateTemplate = `{
    "$schema": "http://json-schema.org/draft-07/schema",
    "title": "JSON Schema for {{$.module}} book",
    "type": "object",
    "properties": {
      "name": {
        "type": "string",
        "minLength": 3
      }
    },
    "required": [
      "name"
    ],
    "additionalProperties": true
  }`
