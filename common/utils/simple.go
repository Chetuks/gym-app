package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	lowerCharSet      = "abcdedfghijklmnopqrst"
	upperCharSet      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet    = "!@#$%&*"
	numberSet         = "0123456789"
	allCharSet        = lowerCharSet + upperCharSet + specialCharSet + numberSet
	BearerAccessToken map[string]string
)

func GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	// Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	// Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	// Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func GetAuthxIdentity(c *fiber.Ctx) []string {
	XAuthxIdentity := c.Request().Header.Peek("X-Authx-Identity")
	tokenString := string(XAuthxIdentity)
	tokens := strings.Split(tokenString, ",")
	return tokens
}

/*
Set App Auth Token
*/
func init() {
	// fmt.Println("Inside init of common/utils.go for setting app auth token")
	// _, file, line, _ := runtime.Caller(1)
	// callerFunc := filepath.Base(file) + ":" + strconv.Itoa(line)
	// log.Info().Str("Caller Function", callerFunc).Msg("Inside SetAppAuthToken")
	//
	// //Set Authorization Header - JWT AccessToken
	// clientToken, tknErr := GetClientAccessToken()
	// if tknErr != nil && clientToken.JwtToken.AccessToken == "" {
	// log.Debug().Str("Caller Function", callerFunc).Interface("tknErr ", tknErr).Send()
	// return
	// //return "",tknErr
	// } else if clientToken.JwtToken.AccessToken != "" {
	// log.Debug().Str("Caller Function", callerFunc).Interface("Auth header added ", clientToken.JwtToken.AccessToken).Send()
	// bearerToken := "Bearer " + clientToken.JwtToken.AccessToken
	// BearerAccessToken = bearerToken
	// return
	// }
	fmt.Println("Inside init of common/utils.go for setting app auth token")
	_, file, line, _ := runtime.Caller(1)
	callerFunc := filepath.Base(file) + ":" + strconv.Itoa(line)
	log.Info().Str("Caller Function", callerFunc).Msg("Inside SetAppAuthToken")

}

// StringComparator Custom lexicographic comparison function
func StringComparator(a, b string) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	// If we reach here, they are equal up to minLen
	if len(a) < len(b) {
		return -1
	} else if len(a) > len(b) {
		return 1
	}
	return 0
}

// converts JSON column key to database column name
func jsonToDBColumnName(jsonColumnName string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z0-9])`)

	// Replace uppercase letters with underscores followed by lowercase letters
	dbColumnName := re.ReplaceAllString(jsonColumnName, "${1}_${2}")

	// Convert to lowercase
	dbColumnName = strings.ToLower(dbColumnName)
	return dbColumnName
}

func PanicRecovery(funcName string) {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from panic in %s: %v\n", funcName, r)
	}
}

func FibonacciBackoff(attempt int) time.Duration {
	if attempt <= 1 {
		return 500 * time.Millisecond
	}
	a, b := 500*time.Millisecond, 500*time.Millisecond
	for i := 2; i <= attempt; i++ {
		a, b = b, a+b
	}
	return b
}

func SetDefaultEmptyValues(v interface{}) {
	log.Info().Msg("Setting default empty values for all null fields")
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.Slice:
			if field.IsNil() || field.Len() == 0 {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			}
		}
	}
}
