package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	//不明用处
	dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	//字体，需要支持中文
	fontfile = flag.String("fontfile", "/usr/share/fonts/wqy-zenhei/wqy-zenhei.ttc", "filename of the ttf font")
	//好像没用
	hinting = flag.String("hinting", "none", "none | full")
	size    = flag.Float64("size", 12, "font size in points")
	//间距
	spacing = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	//字体颜色
	wonb = flag.Bool("whiteonblack", false, "white text on a black background")
)

var text = []string{
	"某某",
	// "Test",
}

func main() {
	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	a, err := os.Open("test.png")
	if err != nil {
		fmt.Println(err)
	}
	afile, err := png.Decode(a)
	if err != nil {
		fmt.Println(err)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, afile.Bounds().Max.X, afile.Bounds().Max.Y))
	draw.Draw(rgba, afile.Bounds(), afile, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(afile.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.第二个参数应该是字的下边缘坐标
	pt := freetype.Pt(220, 340+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//否则图片不全
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
