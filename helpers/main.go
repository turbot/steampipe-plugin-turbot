package helpers

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/encryption"
	"reflect"
)

func MergeMaps(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

// given a list of items which may each be either a property or property map, remove the excluded properties
func RemoveProperties(properties []interface{}, excluded []string) []interface{} {
	var result []interface{}
	for _, element := range properties {
		// each element may be either a map, or a single property name
		terraformToTurbotMap, ok := element.(map[string]string)
		if ok {
			// if the element is a map, remove excluded items from map
			result = append(result, RemovePropertiesFromMap(terraformToTurbotMap, excluded))
		} else {
			// otherwise check if this property is excluded and remove if so
			if !SliceContains(excluded, element.(string)) {
				result = append(result, element)
			}
		}
	}
	return result
}

// given a property list, remove the excluded properties
func RemovePropertiesFromMap(propertyMap map[string]string, excluded []string) map[string]string {
	for _, v := range excluded {
		delete(propertyMap, v)
	}
	return propertyMap
}

// no native contains in golang :/
func SliceContains(s []string, searchTerm string) bool {
	for _, v := range s {
		if v == searchTerm {
			return true
		}
	}
	return false

}

func EncryptValue(pgpKey, value string) (string, string, error) {
	encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
	if err != nil {
		return "", "", err
	}
	fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, value, "Secret Key")
	if err != nil {
		return "", "", err
	}
	return fingerprint, encrypted, nil
}

func MapToJsonString(data map[string]interface{}) (string, error) {
	dataBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return "", err
	}
	jsonData := string(dataBytes)
	return jsonData, nil
}

func JsonStringToMap(dataString string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(dataString), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// apply standard formatting to a json string by unmarshalling into a map then marshalling back to JSON
func FormatJson(body string) string {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// ignore error and just return original body
		return body
	}
	body, err := MapToJsonString(data)
	if err != nil {
		// ignore error and just return original body
		return body
	}
	return body

}

// given a json representation of an object, build a map of the property names: property alias -> property path
func PropertyMapFromJson(body string) (map[string]string, error) {
	if body == "" {
		return nil, nil
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}
	var properties = map[string]string{}
	for k := range data {
		properties[k] = k
	}
	return properties, nil
}

// convert a map[string]interface{} to a map[string]string by json encoding any non string fields
func ConvertToStringMap(data map[string]interface{}) (map[string]string, error) {
	var outputMap = map[string]string{}

	for k, v := range data {
		if v != nil {
			if reflect.TypeOf(v).String() != "string" {
				jsonBytes, err := json.MarshalIndent(v, "", " ")
				if err != nil {
					return nil, err
				}
				v = string(jsonBytes)
			}
			outputMap[k] = v.(string)
		}
	}
	return outputMap, nil
}

func GetNullProperties(propertyMap map[string]interface{}) []string {
	var result []string
	if properties, ok := propertyMap["properties"]; ok {
		for id, valueObject := range properties.(map[string]interface{}) {
			if m, ok := valueObject.(map[string]interface{}); ok {
				if m["type"] == "null" {
					result = append(result, id)
				}
			}
		}
	}
	return result
}

// get keys from old map not in new map
func GetOldMapProperties(old, new map[string]interface{}) []interface{} {
	var result []interface{}
	for k := range old {
		if _, ok := new[k]; !ok {
			result = append(result, k)
		}
	}
	return result
}
