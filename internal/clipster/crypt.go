// Crypt deals with Fernet encryption and decryption using PBK2DF
package clipster

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"

	"github.com/fernet/fernet-go"
)

// deriveKey from password using a salt via PBK2DF and return urlsafe b64
// cross client compatible by using same parameters and same algos
func deriveKey(user string, pw string, iters int) string {
	salt := "clipster_" + user + "_" + pw
	key, err := pbkdf2.Key(sha256.New, pw, []byte(salt), iters, HASH_LENGTH)
	if err != nil {
		log.Panicln("Error: deriving key", err)
	}
	return base64.URLEncoding.EncodeToString(key)
}

// Encrypt the text using Fernet and the hash_msg key
func Encrypt(text string) (string, error) {
	key, err := fernet.DecodeKey(conf.Hash_msg)
	if err != nil {
		return "", errors.New("no valid encryption key, please edit your credentials")
	}
	tok, err := fernet.EncryptAndSign([]byte(text), key)
	if err != nil {
		return "", err
	}
	return string(tok), nil
}

// Decrypt decrypts a text using hash_msg as a key and Fernet and returns a string
func Decrypt(text string) (string, error) {
	key, err := fernet.DecodeKey(conf.Hash_msg)
	if err != nil {
		return "", errors.New("no valid encryption key, please edit your credentials")
	}
	msg := fernet.VerifyAndDecrypt([]byte(text), 0, []*fernet.Key{key})
	if msg == nil {
		return "", errors.New("could not decrypt clip, was it encrypted with another password?")
	}
	return string(msg), nil
}

// GetLoginHashFromPw returns a hash (string) of the password to be used for authentication
func GetLoginHashFromPw(user string, pw string) string {
	return deriveKey(user, pw, HASH_ITERS_LOGIN)
}

// GetMsgHashFromPw returns a hash (string) of the password to be used as encryption key
func GetMsgHashFromPw(user string, pw string) string {
	return deriveKey(user, pw, HASH_ITERS_MSG)
}
