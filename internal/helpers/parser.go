package helpers

import (
	"os"
)

func ParseFile(filename string) ([]byte, error) {
	qsj, err := os.ReadFile(filename)
	if err != nil {
		return qsj, err
	}

	return qsj, nil
}
