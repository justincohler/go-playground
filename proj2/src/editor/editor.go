package main

import (
	"imgutil"
	"os"
)

func main() {

	//Assumes the user specifies a file as the first argument
	filePath := os.Args[1]

	//Loads the png image and returns the image or an error
	image, _ := imgutil.Load(filePath)

	//Performs a grayscale filtering effect on the image
	gray := image.Grayscale()
	blur := image.Blur()

	//Saves the image to a new file
	blur.Save(filePath + "_blur.png")
	gray.Save(filePath + "_gray.png")
}
