package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	jsonSource, err := ioutil.ReadFile("sample.json")
	if err != nil {
		fmt.Printf("Exception reading json file: %s\n", err)
		return
	}

	var jsonValue interface{}

	err = json.Unmarshal([]byte(jsonSource), &jsonValue)
	if err != nil {
		fmt.Printf("Exception parsing JSON: %s\n", err)
		return
	}

	// TODO - ignore top-level id
	// walk("", jsonValue)

	val, err := getJSONProperty(jsonValue, "properties", "storageProfile", "osDisk", "osType")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else if val == nil {
		fmt.Println("Property not found")
	} else {
		val2, ok := val.(string)
		if ok {
			fmt.Printf("Got value: %v\n", val2)
		} else {
			fmt.Println("Not a string")
		}
	}

}

func walk(currentName string, o interface{}) {
	switch t := o.(type) {
	case map[string]interface{}:
		omap := o.(map[string]interface{})
		name := currentName
		oName := omap["name"]
		if oName != nil {
			if name != "" {
				name = name + "/"
			}
			name = name + oName.(string)
		}

		walkMap(name, omap)

	case float64:
	case bool:
	case string:
		break // ignore strings

	case []interface{}:
		for _, v := range o.([]interface{}) {
			walk(currentName, v)
		}

	default:
		fmt.Printf("unhandled type: %T\n", t)
	}
}

func walkMap(currentName string, o map[string]interface{}) {
	for k, v := range o {
		if k == "id" {
			name := o["name"]
			if name == nil {
				name = "<TODO>"
			}
			fmt.Printf("%s\t [%s]\n", currentName+"/"+name.(string), v)
		} else {
			walk(currentName, v)
		}
	}
}
func getJSONProperty(jsonData interface{}, properties ...string) (interface{}, error) {
	switch jsonData.(type) {
	case map[string]interface{}:
		jsonMap := jsonData.(map[string]interface{})
		name := properties[0]
		jsonSubtree, ok := jsonMap[name]
		if ok {
			if len(properties) == 1 {
				return jsonSubtree, nil
			}
			return getJSONProperty(jsonSubtree, properties[1:]...)
		} else {
			return nil, nil // TODO - error if not found?
		}
	default:
		return nil, nil // TODO - error if not able to walk the tree?
	}

}
