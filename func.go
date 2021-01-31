package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

func toJson(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func checkToken(file string, token string) (ok bool, err error) {
	data, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return false, errors.Wrap(err, "read file error")
	}
	list := strings.Split(string(data), "\n")
	for _, item := range list {
		if strings.Index(item, "#") == 0 {
			continue
		}
		if strings.TrimSpace(item) == "" {
			continue
		}
		if strings.TrimSpace(item) == token {
			return true, nil
		}
	}
	return false, nil
}
