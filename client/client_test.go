package main

import (
	"regexp"
	"testing"
)

func TestGenUuid(t *testing.T) {

	var validUUID = regexp.MustCompile(`[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}`)

	uuid := genUuid()

	if validUUID.MatchString(uuid) == false {
		t.Error("UUID Gen Fail")
	}

}

func TestMd5String(t *testing.T) {

	expectedOutput := `d67c5cbf5b01c9f91932e3b8def5e5f8`

	if md5String("teststring") != expectedOutput {
		t.Error("MD5 String failed")
	}

}
