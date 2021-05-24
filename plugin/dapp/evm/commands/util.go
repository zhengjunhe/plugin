package commands

import (
	"io/ioutil"
	"os"
)

func readFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return "", err
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}
