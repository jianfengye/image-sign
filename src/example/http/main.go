package main

import (
	"bytes"
	"errors"
	"flag"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"os"
	"signer"
)

var (
	sign        *signer.Signer
	imgPath     string
	fontPath    string
	contentType string
	srcBytes    []byte
	signFormat  int
	startX      int
	startY      int
	signX       int
	signY       int
)

func init() {
	flag.StringVar(&imgPath, "img", "img/src.jpg", "source image path")
	flag.StringVar(&fontPath, "font", "font/FZFSK.TTF", "source font path")
	flag.IntVar(&signFormat, "format", signer.SignTypeJpg, "image format: 0 png, 1 jpg")
	flag.IntVar(&startX, "startx", 88, "start point: x")
	flag.IntVar(&startY, "starty", 222, "start point: y")
	flag.IntVar(&signX, "signx", 100, "sign point: x")
	flag.IntVar(&signY, "signy", 100, "sign point: y")
}

func initImg() error {
	srcFile, err := os.Open(imgPath)
	if err != nil {
		glog.Errorf("os.Open(\"%s\") error(%v)", imgPath, err)
		return err
	}
	defer srcFile.Close()
	srcBytes, err = ioutil.ReadAll(srcFile)
	if err != nil {
		glog.Errorf("ioutil.ReadAll(srcFile) error(%v)", err)
		return err
	}
	return nil
}

func initContentType() error {
	switch signFormat {
	case signer.SignTypePng:
		contentType = "image/png"
	case signer.SignTypeJpg:
		contentType = "image/jpeg"
	default:
		return errors.New("not support format, please 0 png, 1 jpeg")
	}
	return nil
}

func initSign() error {
	sign = signer.NewSigner(fontPath)
	if sign == nil {
		return errors.New("sign init error")
	}
	sign.SetStartPoint(startX, startY)
	sign.SetSignPoint(signX, signY)
	return nil
}

func main() {
	if err := initImg(); err != nil {
		glog.Errorf("initImg error(%v)", err)
		return
	}
	if err := initSign(); err != nil {
		glog.Errorf("initSign error(%v)", err)
		return
	}
	if err := initContentType(); err != nil {
		glog.Errorf("initContentType error(%v)", err)
		return
	}
	http.HandleFunc("/", index)
	http.HandleFunc("/p", handler)
	if err := http.ListenAndServe(":8099", nil); err != nil {
		glog.Errorf("http.ListenAndServe(\":8099\", nil) error(%v)", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	str := "<meta charset=\"utf-8\"><h3>golang 图片</h3><img border=\"1\" src=\"/p?t=aaa\" onclick=\"this.src='/p?t=bbb'\" />"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(str))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/p" {
		return
	}
	r.ParseForm()
	t := r.FormValue("t")
	if t == "" {
		glog.Errorf("t == \"\"")
		return
	}
	w.Header().Set("Content-Type", contentType)
	bs := make([]byte, len(srcBytes))
	copy(bs, srcBytes)
	err := sign.Sign(bytes.NewReader(bs), w, t, signFormat)
	if err != nil {
		glog.Errorf("sign.Sign(\"%s\") error(%v)", t, err)
		return
	}
}
