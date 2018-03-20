package main

// Simplifies general IO operations
import "io/ioutil"

func readFile(path string) string {
	fileAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(fileAsBytes)
}
