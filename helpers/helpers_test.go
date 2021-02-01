package helpers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRemoveProperties(t *testing.T) {
	type test struct {
		name       string
		properties []interface{}
		excluded   []string
		expected   []interface{}
	}
	tests := []test{
		test{
			"No exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{},
			[]interface{}{"a", "b", "c"},
		},
		test{
			"String exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"a"},
			[]interface{}{"b", "c"},
		},
		test{
			"All excluded",
			[]interface{}{"a", "b", "c"},
			[]string{"a", "b", "c"},
			[]interface{}(nil),
		},
		test{
			"Map exclusion",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
			[]string{"c"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}},
		},
		test{
			"2 map exclusions",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}, map[string]string{"e": "E", "f": "F"}},
			[]string{"c", "f"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}, map[string]string{"e": "E"}},
		},
		test{
			"No matching exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"d"},
			[]interface{}{"a", "b", "c"},
		},
		test{
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
		test{
			"Empty object",
			`{
						  "allOf": []
						}`,
			[]interface{}{nil},
		},
		test{
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
		test{
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
		test{
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
