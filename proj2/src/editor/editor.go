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
	sharpen := image.Sharpen()
	edge := image.Edge()

	//Saves the image to a new file
	gray.Save(filePath + "_gray.png")
	blur.Save(filePath + "_blur.png")
	sharpen.Save(filePath + "_sharpen.png")
	edge.Save(filePath + "_edge.png")
}
