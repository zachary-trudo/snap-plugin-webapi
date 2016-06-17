package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestJSON(t *testing.T) {
	buf, err := ioutil.ReadFile("plugins.json")
	if err != nil {
		fmt.Println(err)

	}
	s := string(buf)
	if !isJSON(s) {
		t.Fail()
	}

}

func isJSON(str string) bool {
	var temp interface{}
	return json.Unmarshal([]byte(str), &temp) == nil
}
