package Model

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
	"sync"
)

const (
	WATERMARK = `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`
)

var watermark image.Image
var ilock sync.Locker

func init() {
	fb, err := base64.StdEncoding.DecodeString(WATERMARK)
	if err != nil {
		panic(err)
		return
	}
	watermark, err = png.Decode(bytes.NewBuffer(fb))
	if err != nil {
		panic(err)
	}
	ilock = new(sync.Mutex)
}
func AddWaterMark(b64 string) (result []byte, err error) {
	ilock.Lock()
	defer ilock.Unlock()
	spl := strings.Split(b64, ",")
	if len(spl) > 1 {
		b64 = spl[1]
	}
	tmp, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return
	}
	fb := bytes.NewBuffer(tmp)
	img, err := jpeg.Decode(fb)
	if err != nil {
		return
	}
	//把水印写到右下角，并向0坐标各偏移10个像素
	offset := image.Pt(img.Bounds().Dx()-watermark.Bounds().Dx()-20, img.Bounds().Dy()-watermark.Bounds().Dy()-40)
	b := img.Bounds()
	m := image.NewNRGBA(b)

	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	jpeg.Encode(fb, m, &jpeg.Options{100})
	result = fb.Bytes()
	return
}
