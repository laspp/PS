package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"imageCp/cudaImage"
	"os"
	"unsafe"

	"github.com/InternatBlackhole/cudago/cuda"
)

func rgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return gray
}

func main() {
	// read command line arguments
	inputImageStr := flag.String("i", "", "input image")
	outputImageStr := flag.String("o", "", "output image")
	flag.Parse()
	if *inputImageStr == "" || *outputImageStr == "" {
		panic("Missing input or output image arguments\nUsage: go run main.go -i input.png -o output.png")
	}

	//Initialize CUDA API on OS thread
	var err error
	dev, err := cuda.Init(0)
	if err != nil {
		panic(err)
	}
	defer dev.Close()
	//Open image file
	inputFile, err := os.Open(*inputImageStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Read image " + *inputImageStr)
	defer inputFile.Close()
	//Decode image
	inputImage, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}
	//Convert image to grayscale
	inputImageGray := rgbaToGray(inputImage)

	imgSize := inputImageGray.Bounds().Size()
	size := uint64(imgSize.X * imgSize.Y)

	//Allocate memory on the device for input and output image
	imgInDevice, err := cuda.DeviceMemAlloc(size)
	if err != nil {
		panic(err)
	}
	defer imgInDevice.Free()

	imgOutDevice, err := cuda.DeviceMemAlloc(size)
	if err != nil {
		panic(err)
	}
	defer imgOutDevice.Free()

	//Copy image to device
	err = imgInDevice.MemcpyToDevice(unsafe.Pointer(&inputImageGray.Pix[0]), size)
	if err != nil {
		panic(err)
	}
	//Specify grid and block size
	dimBlock := cuda.Dim3{X: 1, Y: 1, Z: 1}
	dimGrid := cuda.Dim3{
		X: 1,
		Y: 1,
		Z: 1,
	}
	//Call the kernel to execute on the device
	err = cudaImage.Process(dimGrid, dimBlock, imgInDevice.Ptr, imgOutDevice.Ptr, int32(imgSize.X), int32(imgSize.Y))
	if err != nil {
		panic(err)
	}
	//Copy image back to host
	imgOutHost := make([]byte, size)
	imgOutDevice.MemcpyFromDevice(unsafe.Pointer(&imgOutHost[0]), size)

	//Save image to file
	outputFile, err := os.Create(*outputImageStr)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	outputImage := image.NewGray(inputImageGray.Bounds().Bounds())
	outputImage.Pix = imgOutHost

	err = png.Encode(outputFile, outputImage)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image saved to " + *outputImageStr)
}
