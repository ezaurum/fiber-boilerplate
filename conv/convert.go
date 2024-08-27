package conv

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func ToInt64(param string) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 10, 64)
	if nil != err {
		return 0, false
	}
	return id, true
}

func ToInt64FromHex(param string) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 16, 64)
	if nil != err {
		return 0, false
	}
	return id, true
}

func ToInt(param string) (int, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(param), 10, 32)
	if nil != err {
		return 0, false
	}
	return int(id), true
}

func ToFloat64(param string) (float64, bool) {
	id, err := strconv.ParseFloat(strings.TrimSpace(param), 64)
	if nil != err {
		return 0, false
	}
	return id, true
}

func ToBase64Json(v interface{}) (string, error) {
	marshal, err := json.Marshal(v)
	if nil != err {
		return "", errors.Wrapf(err, "marshal to json failed")
	}

	encoded := base64.StdEncoding.EncodeToString(marshal)
	return encoded, nil
}

func FromBase64Json(base64string string, v interface{}) error {
	decodeString, err := base64.StdEncoding.DecodeString(base64string)
	if nil != err {
		return err
	}
	return json.Unmarshal(decodeString, v)
}
