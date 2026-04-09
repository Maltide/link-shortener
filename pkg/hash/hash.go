// Package hash provides functions for generating and decoding short URL hash.
package hash

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Maltide/link-shortener/pkg/helpers"
	"go.uber.org/zap"
)

var pair_counter int = 0

// Hash generates a short link based on the row-id from DB. Clears the 'links' table when the limit is reached.
func Hash(id int, log *zap.SugaredLogger, db *sql.DB) (string, error) {
	if id <= 0 {
		log.Error("hash function error: unexpected nonpositive id")

		return "", fmt.Errorf("unexpected nonpositive id")
	}

	if pair_counter >= 1000000 {
		err := helpers.ClearDB(db, log)
		if err != nil {
			return "", err
		}
		pair_counter = 0
	}

	alphabit := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	resString := make([]rune, 0, 7)

	for id > 0 {
		symbolIndex := id % len(alphabit)

		resString = append(resString, rune(alphabit[symbolIndex]))

		id = id / len(alphabit)
	}

	for i := 0; i < len(resString)/2; i++ {
		resString[i], resString[len(resString)-i-1] = resString[len(resString)-i-1], resString[i]
	}

	log.Infof("short string after proccessing is:%v", resString)

	pair_counter++

	return string(resString), nil
}

// RedirectHash converts a short link back to an id and find original-url from DB.
func RedirectHash(hashLink string, log *zap.SugaredLogger) (int, error) {
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
