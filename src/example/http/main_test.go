package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	server *httptest.Server
)

func inittest() {
	server = httptest.NewServer(http.HandlerFunc(handler))
	initImg()
	initSign()
	initContentType()
}

func TestPost(t *testing.T) {
	inittest()
	resp, err := http.PostForm(server.URL+"/p", url.Values{"t": []string{"ttttt"}})
	if err != nil {
		t.Errorf("http.PostForm error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll error(%v)", err)
		t.FailNow()
	}
	f, err := os.Create("img/test.jpg")
	if err != nil {
		t.Errorf("os.Create error(%v)", err)
		t.FailNow()
	}
	defer f.Close()
	_, err = f.Write(bs)
	if err != nil {
		t.Errorf("f.Write error(%v)", err)
		t.FailNow()
	}
}

func BenchmarkPost(b *testing.B) {
	b.StopTimer()
	inittest()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := http.PostForm(server.URL+"/p", url.Values{"t": []string{"ttttt"}})
		if err != nil {
			b.Errorf("http.PostForm error(%v)", err)
			b.FailNow()
		}
	}
}
