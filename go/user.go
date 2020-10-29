package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
	Ban      bool
	Level    UserLevel
}

type UserLevel string

const (
	ADMIN  UserLevel = "Admin"
	MEMBER UserLevel = "Member"
	GUEST  UserLevel = "Guest"
)

func (u *User) IsAdmin() bool {
	return u.Level == ADMIN
}

func checkPassword(username, password string) bool {
	if u, ok := GlobalUsers.Load(username); ok {
		if user, ok := u.(User); ok {
			return CheckPasswordHash(password, user.Password)
		}
	}
	return false
}

func getUserByUsername(username string) User {
	if u, ok := GlobalUsers.Load(username); ok {
		if user, ok := u.(User); ok {
			return user
		}
	}
	return User{}
}

func register(username, password string, level UserLevel) error {
	password, err := HashPassword(password)
	if err != nil {
		return err
	}
	if u, ok := GlobalUsers.LoadOrStore(username, User{
		Username: username,
		Password: password,
		Level:    level,
	}); ok {
		return errors.New(fmt.Sprintf("%s is already registed.", username))
	} else if _, ok := u.(User); ok {
		return nil
	}
	return errors.New("register fail")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
