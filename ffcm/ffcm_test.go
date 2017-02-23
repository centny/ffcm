package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
)

func init() {
	ef = func(int) {

	}
}

func TestDim(t *testing.T) {
	ef = func(int) {
	}
	//
	os.Args = []string{"ffcm", "-d", "100", "100", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "300", "100", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "100", "400", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "210", "500", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "21x0", "500", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "21x0", "500"}
	main()
}

func TestSrv(t *testing.T) {
	go func() {
		os.Args = []string{"ffcm", "-s", "./ffcm_s.properties"}
		main()
		panic("done...")
	}()
	go func() {
		os.Args = []string{"ffcm", "-mem", "./ffcm_s_t.properties"}
		main()
		panic("done...")
	}()
	go func() {
		os.Args = []string{"ffcm", "-c", "./ffcm_c.properties"}
		main()
		panic("done...")
	}()
	time.Sleep(2 * time.Second)
	go func() {
		os.Args = []string{"ffcm", "-s", "./ffcm_s.properties"}
		main()
	}()
	go func() {
		os.Args = []string{"ffcm", "-mem", "./ffcm_s_t.properties"}
		main()
	}()
	go func() {
		os.Args = []string{"ffcm", "-c", "./ffcm_c.properties"}
		main()
	}()
	time.Sleep(2 * time.Second)
}

func TestInfo(t *testing.T) {
	os.Args = []string{"ffcm", "-i", "../xx.mp4"}
	main()
	os.Setenv("FFPROBE_C", "/usr/local/bin/ffprobe")
	os.Args = []string{"ffcm", "-i", "../xx.mp4"}
	main()
	os.Args = []string{"ffcm", "-i", "../sdxx.mp4"}
	main()
	os.Args = []string{"ffcm", "-i"}
	main()
}

func TestG(t *testing.T) {
	var ts = httptest.NewMuxServer()
	os.Args = []string{"ffcm", "-g", "http://127.0.0.1:23243"}
	main()
	os.Args = []string{"ffcm", "-g", ts.URL}
	main()
	os.Args = []string{"ffcm", "-g"}
	main()
}

func TestNormal(t *testing.T) {
	os.Args = []string{"ffcm"}
	main()
	os.Args = []string{"ffcm", "-sdsf"}
	main()
}

func TestCov(t *testing.T) {
	os.Setenv("PATH", "/usr/local/bin:"+os.Getenv("PATH"))
	var ts = httptest.NewMuxServer()
	ts.Mux.HFunc("^.*$", func(hs *routing.HTTPSession) routing.HResult {
		ioutil.ReadAll(hs.R.Body)
		return routing.HRES_RETURN
	})
	ecode := 0
	ef = func(c int) {
		ecode = c
	}
	os.Args = []string{"ffmpeg", "-cov_v", ts.URL, "../xx.mp4", "1280", "720", "1024", "768", "tmp/abc.mp4", "out/abc.mp4", "abc.mp4"}
	main()
	if ecode != 0 {
		t.Error("error")
		return
	}
	os.Args = []string{"ffmpeg", "-cov_a", ts.URL, "../xx.amr", "tmp/abc.mp3", "out/abc.mp3", "abc.mp3"}
	main()
	if ecode != 0 {
		t.Error("error")
		return
	}
}
