package internal

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/marius004/phoenix-algo/entities"
	"golang.org/x/crypto/bcrypt"
)

var DefaultCtx = context.Background()

func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

func GenerateJwtToken(jwtSecret string, expires time.Duration, user *entities.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24 * expires).Unix(),
	})

	token, err := claims.SignedString([]byte(jwtSecret))

	return token, err
}

func VerifyToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func IsUserAuthed(user *entities.User) bool {
	return user != nil
}

func IsUserAdmin(user *entities.User) bool {
	return IsUserAuthed(user) && user.IsAdmin
}

func IsUserProposer(user *entities.User) bool {
	return IsUserAdmin(user) || (IsUserAuthed(user) && user.IsProposer)
}

func WriteToFile(path string, data []byte) error {
	// open the file or create a new one in case it does not exist
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	return err
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func MakeDirectory(path string) error {
	err := os.Mkdir(path, 0755)

	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return err
}

func DeleteDirectory(path string) error {
	err := os.RemoveAll(path)

	if errors.Is(err, os.ErrExist) {
		return nil
	}

	return err
}

func RenameDirectory(old, new string) error {
	return os.Rename(old, new)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func ConvertStringToUint(s string) (uint, error) {
	res, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

func CanManageProblem(problem *entities.Problem, user *entities.User) bool {
	if problem == nil {
		return false
	}

	if (IsUserProposer(user) && problem.AuthorId == user.ID) || IsUserAdmin(user) {
		return true
	}

	return false
}
