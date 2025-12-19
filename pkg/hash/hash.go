package hash

import (
	"fmt"

	"go.uber.org/zap"
)

func hash(s string, log *zap.SugaredLogger) (string, error) {
	if len(s) == 0 {
		log.Error("empty input string from user")

		return "", fmt.Errorf("empty input string from user")
	}

	hashValue := 0

	for i := range s {
		hashValue = (hashValue*7 + int(s[i])) % 1000
	}

	log.Infof("hashValue after for:", hashValue)

	alphabit := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	resString := make([]rune, 0, 7)

	for hashValue > 0 && len(resString) != 7 {
		symbolIndex := hashValue % len(alphabit)

		resString = append(resString, rune(alphabit[symbolIndex]))

		hashValue = hashValue / len(alphabit)
	}

	log.Infof("short string after proccessing is:", resString)

	return string(resString), nil
}
