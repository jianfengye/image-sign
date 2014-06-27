package signer

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"errors"
	"github.com/golang/glog"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	// "os"
	"strings"
)

const (
	SignTypePng = iota
	SignTypeJpg

	DefaultFontSize = 14
	DefaultDpi      = 72
)

type Signer struct {
	FontSize   float64
	Dpi        float64
	font       *truetype.Font
	startPoint image.Point
	signPoint  image.Point
}

func initFont(fontPath string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}

func NewSigner(fontPath string) *Signer {
	font, err := initFont(fontPath)
	if err != nil {
		return nil
	}
	return &Signer{
		FontSize:   DefaultFontSize,
		font:       font,
		Dpi:        DefaultDpi,
		startPoint: image.ZP,
		signPoint:  image.Point{X: 100, Y: 100},
	}
}

func (this *Signer) SetStartPoint(x int, y int) {
	this.startPoint = image.Pt(x, y)
}

func (this *Signer) SetSignPoint(x int, y int) {
	this.signPoint = image.Pt(x, y)
}

func (this *Signer) Sign(input io.Reader, output io.Writer, text string, format int) error {
	var (
		origin image.Image
		err    error
	)
	switch format {
	case SignTypePng:
		origin, err = png.Decode(input)
	case SignTypeJpg:
		origin, err = jpeg.Decode(input)
	default:
		err = errors.New("not support format, now support png,jpeg")
	}
	if err != nil {
		glog.Errorf("image decode error(%v)", err)
		return err
	}
	dst := image.NewNRGBA(origin.Bounds())
	draw.Draw(dst, dst.Bounds(), origin, image.ZP, draw.Src)
	mask, err := this.drawStringImage(text)
	if err != nil {
		glog.Errorf("drawStringImage error(%v)", err)
		return err
	}
	draw.Draw(dst, mask.Bounds().Add(this.startPoint), mask, image.ZP, draw.Over)
	switch format {
	case SignTypePng:
		err = png.Encode(output, dst)
	case SignTypeJpg:
		err = jpeg.Encode(output, dst, nil)
	default:
		err = errors.New("not support format, now support png,jpeg")
	}
	if err != nil {
		glog.Errorf("image encode error(%v)", err)
		return err
	}
	return nil
}

// 画一个带有text的图片
func (this *Signer) drawStringImage(text string) (image.Image, error) {
	fg, bg := image.Black, image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, this.signPoint.X, this.signPoint.Y))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(this.Dpi)
	c.SetFont(this.font)
	c.SetFontSize(this.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFix32(12)>>8))
	for _, s := range strings.Split(text, "\r\n") {
		_, err := c.DrawString(s, pt)
		if err != nil {
			glog.Errorf("c.DrawString(%s) error(%v)", s, err)
			return nil, err
		}
		pt.Y += c.PointToFix32(12 * 1.5)
	}

	// fff, _ := os.Create("aaa.png")
	// defer fff.Close()
	// png.Encode(fff, rgba)

	return rgba, nil
}
