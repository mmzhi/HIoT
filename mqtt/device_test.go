package mqtt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fhmq/hmq/model"
	"testing"
)

func TestDeviceGeneratePassword(t *testing.T) {
	d := Device{
		Device: model.Device{
			ProductId:    "AAAATABB",
			DeviceId:     "GGGG",
			DeviceSecret: "5FBBD6P8",
		},
	}

	var nonce, timestamp, signMethod = "ABCD", 567, "HmacSHA256"

	var p = fmt.Sprintf("clientid=%s\nnonce=%s\ntimestamp=%d", d.ProductId+":"+d.DeviceId, nonce, timestamp)

	h := hmac.New(sha256.New, []byte(d.DeviceSecret))
	h.Write([]byte(p))

	fmt.Printf("%s|%d|%s|%s\n", nonce, timestamp, signMethod, hex.EncodeToString(h.Sum(nil)))
}
