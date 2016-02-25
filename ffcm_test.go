package ffcm

import (
	"bytes"
	"fmt"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"regexp"
	"runtime"
	"testing"
	"time"
)

func TestParseVideo(t *testing.T) {
	FFPROBE_C = "/usr/local/bin/ffprobe"
	video, err := ParseVideo("xx.mp4")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(util.S2Json(video))
	video, err = ParseVideo("ffcm.go")
	if err == nil {
		t.Error("error")
		return
	}
	ParseStreams("xx")
	ParseData("@lx:xds", "\\[[/]*STREAM\\]")
}

// func TestParseVideo2(t *testing.T) {
// 	FFPROBE_C = "/usr/local/bin/ffprobe"
// 	video, err := ParseVideo("/Users/vty/Downloads/the.x-files.s02e04.720p.bluray.x264-geckos.mkv")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	fmt.Println(video.Width, video.Height)
// }

type dtcm_s_h struct {
	cw chan int
}

func (d *dtcm_s_h) OnStart(dtcm *dtm.DTCM_S, task *dtm.Task) {
	fmt.Println("OnStart...")
	d.cw <- 1
}
func (d *dtcm_s_h) OnDone(dtcm *dtm.DTCM_S, task *dtm.Task) error {
	fmt.Println("OnDone...")
	d.cw <- 1
	return nil
}

func TestFFCM(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	//
	//
	//
	util.Exec("rm", "-f", "xx_*")
	var fcfg_s = util.NewFcfg3()
	fcfg_s.InitWithFilePath2("ffcm_s.properties", true)

	var sh = &dtcm_s_h{cw: make(chan int, 100)}
	var err error
	go func() {
		err := InitDtcmS(fcfg_s, dtm.MemDbc, sh)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = RunFFCM_S_V(fcfg_s)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}()
	fmt.Println("xxx->")
	time.Sleep(1 * time.Second)
	go RunFFCM_C("ffcm_c.properties")
	time.Sleep(1 * time.Second)
	fmt.Println("xxxx--->a")
	err = SRV.AddTask("xx.mp4")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("xxxx--->b")
	fmt.Println("xxxx--->c")
	<-sh.cw
	<-sh.cw
	util.Exec("rm", "-f", "xx_*")
	at_ts := httptest.NewServer(SRV.AddTaskH)
	res, err := at_ts.G2("/addTask?src=%v", "xx.mp4")
	if res.IntVal("code") != 0 {
		t.Error("error")
		return
	}
	<-sh.cw
	<-sh.cw
	res, err = at_ts.G2("/addTask?src=%v", "")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}
	res, err = at_ts.G2("/addTask?src=%v", "sfsd")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}
	ts := httptest.NewServer(CLIENT.NofityProc)
	ts.PostN("?tid=%v", "text/plain", bytes.NewBufferString(`
		`), "xxx")
	ts.PostN("?tid=%v&duration=1111", "text/plain", bytes.NewBufferString(`
		xx=
		progress=end
		`), "xxx")
	fmt.Println("xxxx--->d")
	//
	util.Exec("cp", "xx.mp4", "xx_mm")
	SRV.AddTask("xx_mm")
	SRV.AddTask("ffcm.sh")
	//
	// RunFFCM_S("ffcm_s.properties", dtm.MemDbc, sh)
	// RunFFCM_S_V(nil, dtm.MemDbc, sh)
	RunFFCM_C("ffcm_c.properties")
	//
	// StopFFCM_C()
	// StopFFCM_S()
	time.Sleep(2 * time.Second)
	//
	// StopFFCM_C()
	// StopFFCM_S()
	//
	res, err = at_ts.G2("/addTask?src=%v", "sfsd")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}
	SRV = nil
	RunFFCM_S_V(nil)
	//
	//

	util.Exec("rm", "-f", "xx_*")
}

func TestDim(t *testing.T) {
	var args []string
	//
	args = []string{"100", "100", "200", "300"}
	fmt.Println(Dim2(args))
	//
	args = []string{"300", "100", "200", "300"}
	fmt.Println(Dim2(args))
	//
	args = []string{"100", "400", "200", "300"}
	fmt.Println(Dim2(args))
	//
	args = []string{"210", "500", "200", "300"}
	fmt.Println(Dim2(args))
	//
	args = []string{"21x0", "500", "200", "300"}
	fmt.Println(Dim2(args))
	//
	args = []string{"21x0", "500"}
	fmt.Println(Dim2(args))
}

func TestMatch(t *testing.T) {
	fmt.Println(regexp.MustCompile("^").MatchString("abc.MKV"))
	fmt.Println(regexp.MustCompile("(?i)^.*\\.[(mkv)(avi)]").MatchString("abc.MKV"))
	fmt.Println(regexp.MustCompile("(?i)^.*\\.(mkv|avi)$").MatchString("abc.avi"))
	fmt.Println(regexp.MustCompile("(?i)^.*\\.(mkv|avi)$").MatchString("abc.aVi"))
	fmt.Println(regexp.MustCompile("(?i)^.*\\.(mkv|avi)$").MatchString("abc.avix"))
	fmt.Println(regexp.MustCompile("(?i)^.*\\.(mkv|avi)$").MatchString("abc.mkv"))
	fmt.Println(regexp.MustCompile("^.*\\.[(mp4)]$").MatchString("xxx.mp4"))
}

// func TestExec(t *testing.T) {
// 	fmt.Println(util.Exec("/usr/local/bin/ffmpeg -i ./xx.mp4 -s `/Users/vty/vgo/bin/ffcm -d 320 240 960 480` ./xx_phone.mp4"))
// }
