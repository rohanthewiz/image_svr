package main

import (
	"image"
	"image/color"
	"log"

	"github.com/savsgio/atreugo/v9"
	"github.com/disintegration/imaging"
)

func main() {
	cfg := &atreugo.Config{
		Addr: "0.0.0.0:8000",
	}
	svr := atreugo.New(cfg)

	// Register a route
	svr.Path("GET", "/", processPic)

	err := svr.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func processPic(ctx *atreugo.RequestCtx) (err error) {
	src, err := imaging.Open("pics/us.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Crop the original image to 300x300px size using the center anchor.
	src = imaging.CropAnchor(src, 600, 600, imaging.Center)

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, 400, 0, imaging.Lanczos)

	// Create a blurred version of the image.
	img1 := imaging.Blur(src, 5)

	// Create a grayscale version of the image with higher contrast and sharpness.
	img2 := imaging.Grayscale(src)
	img2 = imaging.AdjustContrast(img2, 20)
	img2 = imaging.Sharpen(img2, 2)

	// Create an inverted version of the image.
	img3 := imaging.Invert(src)

	// Create an embossed version of the image using a convolution filter.
	img4 := imaging.Convolve3x3(
		src,
		[9]float64{
			-1, -1, 0,
			-1, 1, 1,
			0, 1, 1,
		},
		nil,
	)

	// Create a new image and paste the four produced images into it.
	dst := imaging.New(800, 800, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, img1, image.Pt(0, 0))
	dst = imaging.Paste(dst, img2, image.Pt(0, 400))
	dst = imaging.Paste(dst, img3, image.Pt(400, 0))
	dst = imaging.Paste(dst, img4, image.Pt(400, 400))

	// Save the resulting image as JPEG.
	err = imaging.Save(dst, "pics/us_processed.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	return ctx.TextResponse("Picture has been processed")
}
