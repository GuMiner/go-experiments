package commonIo

// Simplifies general IO operations
import (
	"image"
	"io/ioutil"
	"os"

	_ "image/png"
)

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

func ReadImageFromFile(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	img, _, imageError := image.Decode(file)
	if imageError != nil {
		panic(imageError)
	}

	return img
}
