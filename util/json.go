package util

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

//JSONSerialize desc
//@method JSONSerialize desc: golang object Serialized to json string
//@param  (interface{}) json object
//@return (string) json string
func JSONSerialize(obj interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(&obj)

	if err != nil {
		return fmt.Sprintf("{ \"code\" : -1, \"message\" : \"json marshl error:%s\"}", err.Error())
	}

	return string(data)
}

//JSONUnSerialize desc
//@method JSONUnSerialize desc: Reverse the json string [byte] into a golang object
//@param  ([]byte) json []byte
//@param  (interface{}) out json object
//@return (error)
func JSONUnSerialize(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, v)
	return err
}

//JSONUnFormSerialize desc
//@method JSONUnFormSerialize desc: Reverse the json string into a golang object
//@param  (string) json string
//@param  (interface{}) out json object
func JSONUnFormSerialize(data string, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(data, v)
	return err
}
