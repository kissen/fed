package util

import "encoding/base64"

func DecodeBase64ToString(b64s string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(b64s)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
