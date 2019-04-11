package weather2png

import (
	"image"
	"image/color"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"github.com/golang/freetype"
	"log"
	"io"
	"image/png"
)

type PngWriter struct {
	cxt    *freetype.Context
	width  int
	height int
	img    *image.Gray
	//fgColor *image.Uniform
	//timeSize int
	//dateSize int
	//prefixSize int
	//infoSize int
}

func NewPngWriter(w, h int, fontPath string) *PngWriter {
	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fontData)
	img := image.NewGray(image.Rectangle{image.Point{0, 0,}, image.Point{w, h}})
	for i,_ := range img.Pix{
		img.Pix[i] = 0xff
	}
	freetypeCxt := freetype.NewContext()
	freetypeCxt.SetDst(img)
	freetypeCxt.SetSrc(image.NewUniform(color.Black))
	freetypeCxt.SetFont(font)
	freetypeCxt.SetClip(img.Bounds())

	return &PngWriter{freetypeCxt, w, h, img} //, 80, 40, 24, 40}
}

func (thiz *PngWriter) Reset(writer io.Writer) {
	img := image.NewGray(image.Rectangle{image.Point{0, 0,}, image.Point{thiz.width, thiz.height}})
	for i,_ := range img.Pix{
		img.Pix[i] = 0xff
	}
	thiz.cxt.SetDst(img)
	png.Encode(writer, thiz.img)
	thiz.img = img
}

func (thiz *PngWriter) VerticalLine(x1, y, x2 int) {
	for x := x1; x < x2; x++ {
		thiz.img.Set(x, y, color.Black)
		thiz.img.Set(x, y+1, color.Black)
	}
}

func (thiz *PngWriter) Text(txt string, x, y int, size float64) {
	thiz.cxt.SetFontSize(size)
	_, err := thiz.cxt.DrawString(txt, freetype.Pt(x, y))
	if err != nil {
		log.Println("ERRO: PngWriter draw string failed", err)
	}
}
