package sample

import (
	"errors"
	"github.com/louisevanderlith/husk"
	"strings"
	"time"
)

type User struct {
	Name      string `hsk:"size(75)"`
	Verified  bool   `hsk:"default(false)"`
	Email     string `hsk:"size(128)"`
	Password  string `hsk:"min(6)"`
	LoginDate time.Time
	Roles     []Role
}

func (u User) Valid() error {
	err := husk.ValidateStruct(&u)

	if err != nil {
		return err
	}

	if !strings.Contains(u.Email, "@") {
		return errors.New("email is invalid")
	}

	return nil
}
