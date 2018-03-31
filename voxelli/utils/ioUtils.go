package utils

// Simplifies general IO operations
import "io/ioutil"

func ReadFileAsBytes(path string) []uint8 {
	fileAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return fileAsBytes
}

func ReadFile(path string) string {
	fileAsBytes := ReadFileAsBytes(path)
	return string(fileAsBytes)
}
