
package main

import  (
        "testing"
        "regexp"
        )

func TestGenUuid(t *testing.T) {

        var validUUID = regexp.MustCompile(`[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}`)

        uuid := genUuid()

        if validUUID.MatchString(uuid) == false {
             t.Error("UUID Gen Fail")
        }

}




