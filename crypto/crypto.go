package crypto

import (
	"os"
	"time"

	"github.com/fernet/fernet-go"
)

var secretByte string = os.Getenv("CRYPTO_SECRET")

func Decrypt(text string) string {
	k := fernet.MustDecodeKeys(secretByte)
	return string(fernet.VerifyAndDecrypt([]byte(text), 24*time.Hour, k))
}
