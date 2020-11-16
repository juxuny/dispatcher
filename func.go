package main

import "encoding/json"

func toJson(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
