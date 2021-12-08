package helpers

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveProperties(t *testing.T) {
	type test struct {
		name       string
		properties []interface{}
		excluded   []string
		expected   []interface{}
	}
	tests := []test{
		{
			"No exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{},
			[]interface{}{"a", "b", "c"},
		},
		{
			"String exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"a"},
			[]interface{}{"b", "c"},
		},
		{
			"All excluded",
			[]interface{}{"a", "b", "c"},
			[]string{"a", "b", "c"},
			[]interface{}(nil),
		},
		{
			"Map exclusion",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
			[]string{"c"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}},
		},
		{
			"2 map exclusions",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}, map[string]string{"e": "E", "f": "F"}},
			[]string{"c", "f"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}, map[string]string{"e": "E"}},
		},
		{
			"No matching exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"d"},
			[]interface{}{"a", "b", "c"},
		},
		{
			"No matching exclusions with map",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
			[]string{"e"},
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		result := RemoveProperties(test.properties, test.excluded)
		assert.Equal(t, test.expected, result)
	}
}

func TestGetNullProperties(t *testing.T) {
	type test struct {
		name       string
		properties string
		expected   []interface{}
	}
	tests := []test{
		{
			"Empty object",
			`{
						  "allOf": []
						}`,
			[]interface{}{nil},
		},
		{
			"Single exclusion",
			`{
						  "allOf": [
							{
							  "$ref": "#/definitions/account"
							},
							{
							  "type": "object",
							  "properties": {
								"Id": {
								  "type": "null"
								},
								"foo": {
								  "type": "string"
								}
							  }
							},
							{
							  "type": "object"
							}
						  ]
						}`,
			[]interface{}{"Id"},
		},
		{
			"No exclusion",
			`{
						  "allOf": [
							{
							  "$ref": "#/definitions/account"
							},
							{
							  "type": "object"
							},
							{
							  "type": "object"
							}
						  ]
						}`,
			[]interface{}(nil),
		},
		{
			"Multiple exclusion",
			`{
						  "allOf": [
							{
							  "$ref": "#/definitions/account"
							},
							{
							  "type": "object",
							  "properties": {
								"Id": {
								  "type": "null"
								},
								"foo": {
								  "type": "string"
								},
								"bar": {
								  "type": "null"
								}
							  }
							},
							{
							  "type": "object",
							  "properties": {
								"Id2": {
								  "type": "null"
								},
								"foo2": {
								  "type": "string"
								},
								"bar2": {
								  "type": "null"
								}
							  }
							}
						  ]
						}`,
			[]interface{}([]interface{}{"Id", "bar", "Id2", "bar2"}),
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		var jsonMap map[string]interface{}
		err := json.Unmarshal([]byte(test.properties), &jsonMap)
		if err != nil {
			panic(err)
		}
		var excluded []interface{}
		if value, ok := jsonMap["allOf"]; ok {
			for _, schema := range value.([]interface{}) {
				if res, ok := schema.(map[string]interface{}); ok {
					if res["type"] == "object" {
						for _, element := range GetNullProperties(res) {
							excluded = append(excluded, element)
						}

					}
				}
			}
		}
		assert.ObjectsAreEqual(test.expected, excluded)
	}
}
