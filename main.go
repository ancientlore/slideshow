package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gonutz/framebuffer"
)

var (
	fFB      = flag.String("fb", "/dev/fb0", "Frame buffer to use.")
	fFolder  = flag.String("folder", ".", "Folder to scan for image files.")
	fTimeout = flag.Duration("timeout", 15*time.Second, "Duration to view each image.")
)

func main() {
	flag.Parse()

	fb, err := framebuffer.Open(*fFB)
	if err != nil {
		panic(err)
	}
	defer fb.Close()

	fileInfo, err := os.ReadDir(*fFolder)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0, len(fileInfo))
	for _, info := range fileInfo {
		if !info.IsDir() && info.Type().IsRegular() && (strings.HasSuffix(info.Name(), ".jpg") || strings.HasSuffix(info.Name(), ".png")) {
			files = append(files, path.Join(*fFolder, info.Name()))
			info.Type()
		}
	}

	black := image.NewUniform(color.RGBA{0, 0, 0, 255})
	draw.Draw(fb, fb.Bounds(), black, image.Point{}, draw.Src)

	for {
		i := rand.Intn(len(files))
		img, err := loadImage(files[i])
		if err != nil {
			log.Print(err)
			continue
		}
		draw.Draw(fb, fb.Bounds(), black, image.Point{}, draw.Src)
		draw.Draw(fb, fb.Bounds(), img, image.Point{}, draw.Src)
		time.Sleep(*fTimeout)
	}
}

func loadImage(fn string) (image.Image, error) {
	reader, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	return m, err
}
