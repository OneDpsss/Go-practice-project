package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("image")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}

		pixelArt := pixelate(img, 16)
		outFile, err := os.Create("outImage.jpg")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		jpeg.Encode(outFile, pixelArt, nil)
		http.ServeFile(w, r, "outImage.jpg")
	}
}

func pixelate(img image.Image, size int) image.Image {
	imgBounds := img.Bounds()
	imgWidth := imgBounds.Max.X
	imgHeight := imgBounds.Max.Y

	pixelatedImg := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for x := imgBounds.Min.X; x < imgWidth; x += size {
		for y := imgBounds.Min.Y; y < imgHeight; y += size {
			rect := image.Rect(x, y, x+size, y+size)

			if rect.Max.X > imgWidth {
				rect.Max.X = imgWidth
			}
			if rect.Max.Y > imgHeight {
				rect.Max.Y = imgHeight
			}

			r, g, b, a := calculateMeanAverageColourWithRect(img, rect)

			avgColor := color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), uint8(a / 256)}
			for x2 := rect.Min.X; x2 < rect.Max.X; x2++ {
				for y2 := rect.Min.Y; y2 < rect.Max.Y; y2++ {
					pixelatedImg.Set(x2, y2, avgColor)
				}
			}
		}
	}

	return pixelatedImg
}

func calculateMeanAverageColourWithRect(img image.Image, rect image.Rectangle) (r, g, b, a uint32) {
	var count uint32
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			r += pr
			g += pg
			b += pb
			a += pa
			count++
		}
	}
	if count > 0 {
		r /= count
		g /= count
		b /= count
		a /= count
	}
	return r, g, b, a
}
