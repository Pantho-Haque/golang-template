package hash

import (
	"testing"
)

func TestAesEncryptDecrypt(t *testing.T) {
	key := "your-16-or-24-or-32-byte"
	src := "4LHV249"
	encrypted := AesEncrypt(src, key)
	t.Logf("encrypted : %s", encrypted)
	decrypted := AesDecrypt(encrypted, key)
	if decrypted != src {
		t.Error("decryption with wrong key should not produce original plaintext")
	}
}
