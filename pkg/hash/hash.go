package hash

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func hash(id int, log *zap.SugaredLogger) (string, error) {
	if id <= 0 {
		log.Error("hash function error: unexpected nonpositive id")

		return "", fmt.Errorf("unexpected negative id")
	}

	alphabit := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	resString := make([]rune, 0, 7)

	for id > 0 {
		symbolIndex := id % len(alphabit)

		resString = append(resString, rune(alphabit[symbolIndex]))

		id = id / len(alphabit)
	}

	log.Infof("short string after proccessing is:", resString)

	for i := 0; i < len(resString)/2; i++ {
		resString[i], resString[len(resString)-i-1] = resString[len(resString)-i-1], resString[i]
	}

	return string(resString), nil
}

func redirectHash(hashLink string, log *zap.SugaredLogger) (int, error) {
	if len(hashLink) == 0 {
		log.Error("redirect function error: empty input hash string")
		return 0, fmt.Errorf("empty input hash string")
	}
	id := 0

	alphabit := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for _, currentRune := range hashLink {
		runeIndex := strings.IndexRune(alphabit, currentRune)

		if runeIndex == -1 {
			log.Error("redirect hash: symbol not found in 62alphabit")

			return 0, fmt.Errorf("redirect hash: symbol not found in 62alphabit")
		}
		id = id*len(alphabit) + runeIndex
	}
	return id, nil
}
