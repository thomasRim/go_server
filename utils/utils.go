package utils

import (
	"time"
	"fmt"
	"encoding/hex"
	"crypto/rand"
	"log"
	"golang.org/x/crypto/bcrypt"
	"os"
)

/*
	***	Time operations
 */

func TimeInMillis() uint32 {
	millis := time.Now().UnixNano() / 1000000
	return uint32(millis)
}

func TimeInUTC() string {
	t := time.Now()
	return t.Format("2006-01-02T15:04:05.999Z")
}

/*
	***	UUID operations
 */

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	// "crypto/rand" provides the function Read
	n, err := rand.Read(uuid)

	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	//
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	//
	uuid[6] = uuid[6]&^0xf0 | 0x40
	//fmt.Printf("check uuid:\t%s\n", hex.EncodeToString(uuid))
	return fmt.Sprintf("%s", hex.EncodeToString(uuid)), nil
}

/*
	*** Password operation
 */

func HashPassword(pass string) ([]byte,error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	return hash,nil
}

func CheckPasswordWithHash(pass string, hash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(pass)); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

/*
	*** File operations
 */

func CheckFileExist(name string) (bool, error) {
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return false,err
	}
	return true,nil
}