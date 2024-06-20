//go:build !solution

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	height     = 12
	widthColon = 4
	width      = 8
)

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := "http://" + r.Host + r.URL.String()

	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	var timeRequest string
	if len(q["time"]) != 0 && q["time"][0] != "" {
		timeRequest = q["time"][0]
	} else {
		timeRequest = time.Now().Format("15:04:05")
	}
	if len(timeRequest) != 8 {
		http.Error(w, "incorrect time format", 400)
	}
	_, err = time.Parse("15:04:05", timeRequest)
	if err != nil {
		http.Error(w, "incorrect time format", 400)
	}
	var kRequest int
	if len(q["k"]) != 0 && q["k"][0] != "" {
		kRequest, err = strconv.Atoi(q["k"][0])
	} else {
		kRequest = 1
	}
	if err != nil || kRequest < 1 || kRequest > 30 {
		http.Error(w, "k is invalid", 400)
	}
	img := makePicture(timeRequest, kRequest)
	err = png.Encode(w, img)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(200)
}

func getDigitPixels(digit int) string {
	switch digit {
	case 0:
		return Zero
	case 1:
		return One
	case 2:
		return Two
	case 3:
		return Three
	case 4:
		return Four
	case 5:
		return Five
	case 6:
		return Six
	case 7:
		return Seven
	case 8:
		return Eight
	case 9:
		return Nine
	default:
		return ""
	}
}

func drawDigit(img *image.RGBA, digit int32, xStart int, yStart int, k int) (*image.RGBA, int, int) {
	x := xStart
	y := yStart
	digitPixels := getDigitPixels(int(digit - '0'))
	for i, sampleItem := range digitPixels {
		if digitPixels[i] == 10 {
			continue
		}
		img = drawPixel(img, string(sampleItem), x*k, y*k, k)
		if x-xStart == width-1 {
			x = xStart
			y++
		} else {
			x++
		}
	}
	x = xStart + width
	y = 0
	return img, x, y
}

func drawColon(img *image.RGBA, xStart int, yStart int, k int) (*image.RGBA, int, int) {
	x := xStart
	y := yStart
	for i, sampleItem := range Colon {
		if Colon[i] == 10 {
			continue
		}
		img = drawPixel(img, string(sampleItem), x*k, y*k, k)
		if x-xStart == widthColon-1 {
			x = xStart
			y++
		} else {
			x++
		}
	}
	x = xStart + widthColon
	y = 0
	return img, x, y
}

func drawPixel(img *image.RGBA, sign string, xStart int, yStart int, k int) *image.RGBA {
	for y := yStart; y < yStart+k; y++ {
		for x := xStart; x < xStart+k; x++ {
			if sign == "1" {
				img.Set(x, y, Cyan)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}
	return img
}

func makePicture(time string, k int) *image.RGBA {
	resultWidth := (2*widthColon + 6*width) * k
	resultHeight := height * k
	img := image.NewRGBA(image.Rect(0, 0, resultWidth, resultHeight))
	x := 0
	y := 0
	for _, value := range time {
		if value == ':' {
			img, x, y = drawColon(img, x, y, k)
		} else {
			img, x, y = drawDigit(img, value, x, y, k)
		}
	}
	return img
}

func main() {
	portPtr := flag.Int("port", -1, "port")

	flag.Parse()

	if *portPtr == -1 {
		panic("Missing port")
	}

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *portPtr), nil))
}
