package impl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fhmq/hmq/plugins/database"
	"testing"
)

func TestDeviceValid(t *testing.T) {

	d := device{
		database.Device{
			ProductId:    "123456",
			DeviceId:     "ABCDEFG",
			DeviceSecret: "1a2b3c",
		},
	}

	nonce := "c3b2a1"
	timestamp := "0"
	// c3b2a1|0|HmacSHA256|f403df2d3a2e8510d7ee50b80f024e2e15e96cb1b60b90e8d4281b67e0952272
	var p = fmt.Sprintf("clientid=%s\nnonce=%s\ntimestamp=%s",
		d.ProductId+":"+d.DeviceId, nonce, timestamp)
	h := hmac.New(sha256.New, []byte(d.DeviceSecret))
	h.Write([]byte(p))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}
