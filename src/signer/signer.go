package signer

import(
    "image/png"
    "io"
    "image"
    "image/draw"
    "image/color"
    "code.google.com/p/freetype-go/freetype"
    "io/ioutil"
    "strings"
)

const (
    DefaultFontSize = 12
    DefaultDpi = 72
)

type Signer struct {
    // 字体大小
    FontSize float64
    // DPI
    Dpi float64

    input io.Reader
    output io.Writer
    // 字体路径
    fontPath string
    startPoint image.Point
}

func NewSigner(input io.Reader, output io.Writer, fontPath string) *Signer{
    return &Signer{
        FontSize: DefaultFontSize, 
        input: input, 
        output: output,
        fontPath: fontPath,
        Dpi: DefaultDpi,
        startPoint: image.ZP }
}

func (this *Signer) SetFontSize(size float64) {
    this.FontSize = size;
}

func (this *Signer) SetDpi(dpi float64) {
    this.Dpi = dpi;
}

func (this *Signer) SetStartPoint(x int, y int) {
    this.startPoint = image.Pt(x, y)
}

func (this *Signer) Sign(text string) error {
    // TODO: 目前只支持png
    origin, err := png.Decode(this.input)
    if err != nil {
        return err
    }

    dst := image.NewRGBA(image.Rect(0,0, origin.Bounds().Dx(), origin.Bounds().Dy()))
    draw.Draw(dst, dst.Bounds(), origin, image.ZP, draw.Src)

    // 画一个新的带有字体的图片
    mask, err := this.drawStringImage(text)
    if err != nil {
        return err
    }

    src := image.NewUniform(color.RGBA{255,255,255,50})
    draw.DrawMask(dst, dst.Bounds(), src, this.startPoint, mask, image.Point{50,50}, draw.Over)

    err = png.Encode(this.output, dst)
    if err != nil {
        return err
    }
    return nil
}

// 画一个带有text的图片
func (this *Signer) drawStringImage(text string) (image.Image, error) {
    fontBytes, err := ioutil.ReadFile(this.fontPath)
    if err != nil {
        return nil, err
    }

    font, err := freetype.ParseFont(fontBytes)
    if err != nil {
        return nil, err
    }

    fg, bg :=  image.White, image.Black 
    rgba := image.NewRGBA(image.Rect(0, 0, 900, 900))
    draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

    c := freetype.NewContext()
    c.SetDPI(this.Dpi)
    c.SetFont(font)
    c.SetFontSize(this.FontSize)
    c.SetClip(rgba.Bounds())
    c.SetDst(rgba)
    c.SetSrc(fg)

    // Draw the text.
    pt := freetype.Pt(10, 10+int(c.PointToFix32(12)>>8))
    for _, s := range strings.Split(text, "\r\n") {
        _, err = c.DrawString(s, pt)
        pt.Y += c.PointToFix32(12 * 1.5)
    }

    return rgba, nil
}