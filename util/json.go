package util

import jsoniter "github.com/json-iterator/go"

func JsonEncode(data interface{}) (string, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func JsonEncodeByte(data interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func JsonDecode(str string, v interface{}) error {
	return JsonDecodeWithByte([]byte(str), v)
}

func JsonDecodeWithByte(data []byte, v interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
