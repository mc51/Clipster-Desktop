// Crypt deals with Fernet encryption and decryption using PBK2DF
package clipster

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/fernet/fernet-go"
	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(user string, pw string, iters int) string {
	// deriveKey from password using a salt via PBK2DF and return urlsafe b64
	// cross client compatible by using same parameters and same algos
	salt := "clipster_" + user + "_" + pw
	key := pbkdf2.Key([]byte(pw), []byte(salt), iters, HASH_LENGTH, sha256.New)
	key_b64 := base64.URLEncoding.EncodeToString(key)
	log.Println("Ok: derived key string b64", key_b64)
	return key_b64
}

func Encrypt(text string) string {
	// Encrypt the text using Fernet and the hash_msg key
	key := fernet.MustDecodeKeys(conf.Hash_msg)
	tok, err := fernet.EncryptAndSign([]byte(text), key[0])
	if err != nil {
		log.Panicln("Error:", err)
	}
	log.Println("Ok: encrypted token string", string(tok))
	return string(tok)
}

func Decrypt(text string) string {
	// Decrypt decrypts a text using hash_msg as a key and Fernet and returns a string
	key := fernet.MustDecodeKeys(conf.Hash_msg)
	msg := fernet.VerifyAndDecrypt([]byte(text), 0, key)
	log.Println("Ok: decrypted text", string(msg))
	return string(msg)
}

func GetLoginHashFromPw(user string, pw string) string {
	// GetLoginHashFromPw returns a hash (string) of the password to be used for authentication
	hash := deriveKey(user, pw, HASH_ITERS_LOGIN)
	return hash
}

func GetMsgHashFromPw(user string, pw string) string {
	// GetMsgHashFromPw returns a hash (string) of the password to be used as encryption key
	hash := deriveKey(user, pw, HASH_ITERS_MSG)
	return hash
}
