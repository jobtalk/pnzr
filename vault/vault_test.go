package vault

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/ieee0824/thor/test"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestIsSecret(t *testing.T) {
	plainTextPath := TEST_DIR + "/vaultTestFiles/plain.json"
	chipherTextPath := TEST_DIR + "/vaultTestFiles/chipher.json"

	plain, err := ioutil.ReadFile(plainTextPath)
	if err != nil {
		t.Error(err)
	} else if IsSecret(plain) {
		t.Errorf("The result is illegal. I want %v, but it is actually %v.", false, IsSecret(plain))
	}

	chipher, err := ioutil.ReadFile(chipherTextPath)
	if err != nil {
		t.Error(err)
	} else if !IsSecret(chipher) {
		t.Errorf("The result is illegal. I want %v, but it is actually %v.", true, IsSecret(chipher))
	}
}

func falsification(d []byte) []byte {

	for i, v := range d {
		if i%2 == 0 {
			d[i] <<= uint(rand.Int())
		} else {
			d[i] >>= uint(rand.Int())
		}
	}
	return d
}

func TestEncryption(t *testing.T) {
	var randomStr = randStringRunes(65536)
	var key = "test"
	var invalidKey = "johnDoe"

	encrypter := NewString(randomStr)

	if err := encrypter.Encrypt(key); err != nil {
		t.Error(err)
	} else if string(encrypter.Chipher) == randomStr {
		t.Errorf("cipher text and plain text is match")
	}

	if result, err := encrypter.Decrypt(key); err != nil {
		t.Error(err)
	} else if string(result) != randomStr {
		t.Errorf("cipher text and plain text is not match")
	}

	if _, err := encrypter.Decrypt(invalidKey); err == nil {
		t.Errorf("There is no error with an invalid key.")
	}
	encrypter.Chipher = falsification(encrypter.Chipher)
	if _, err := encrypter.Decrypt(key); err == nil {
		t.Errorf("Tampering is not detected.")
	}

}
