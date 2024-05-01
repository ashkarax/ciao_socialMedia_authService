package hashpassword_authSvc

import (
	"errors"
	"fmt"

	interface_hash_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/hash_password/interface"
	"golang.org/x/crypto/bcrypt"
)

type HashUtil struct{}

func NewHashUtil() interface_hash_authSvc.IhashPassword {
	return &HashUtil{}
}

func (hashUtil *HashUtil) HashPassword(password string) string {

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err, "problem at hashing ")
	}
	fmt.Println(HashedPassword)
	return string(HashedPassword)
}

func (hashUtil *HashUtil) CompairPassword(hashedPassword string, plainPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err != nil {
		return errors.New("passwords does not match")
	}

	return nil
}
