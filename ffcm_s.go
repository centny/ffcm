package ffcm

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"path/filepath"
	"strings"
)

var SRV *FFCM_S = nil

func InitDtcmS(fcfg *util.Fcfg, dbc dtm.DB_C, h dtm.DTCM_S_H) error {
	var err error
	SRV, err = NewFFCM_S(fcfg, dbc, h)
	return err
}

// func RunFFCM_S(cfg string, init bool) error {
// 	if DTCM_S == nil {
// 		return util.Err("server is not running")
// 	}
// 	var fcfg = util.NewFcfg3()
// 	fcfg.InitWithFilePath2(cfg, true)
// 	fcfg.Print()
// 	return RunFFCM_S_V(fcfg)
// }

func RunFFCM_S_V(fcfg *util.Fcfg) error {
	if SRV == nil {
		return util.Err("server is not running")
	}
	var ffprobe_c = fcfg.Val("ffprobe_c")
	if len(ffprobe_c) > 0 {
		FFPROBE_C = ffprobe_c
	}
	var w_dir = fcfg.Val("w_dir")
	if len(w_dir) > 0 {
		SRV.W_DIR = w_dir
	}
	var listen = fcfg.Val("listen")
	SRV.Hand("", routing.Shared)
	log.D("listen web server on %v", listen)
	return routing.ListenAndServe(listen)
}

type FFCM_S struct {
	S     *dtm.DTCM_S
	W_DIR string
}

func NewFFCM_S(fcfg *util.Fcfg, dbc dtm.DB_C, h dtm.DTCM_S_H) (*FFCM_S, error) {
	dtcm, err := dtm.StartDTCM_S(fcfg, dbc, h)
	return NewFFCM_Sv(dtcm), err
}

func NewFFCM_Sv(dtcm *dtm.DTCM_S) *FFCM_S {
	return &FFCM_S{S: dtcm, W_DIR: "."}
}
func (f *FFCM_S) AddTask(src string) error {
	return f.AddTaskV(src, nil)
}
func (f *FFCM_S) AddTaskV(src string, info interface{}) error {
	var ext = filepath.Ext(src)
	if len(ext) < 1 {
		return util.Err("invalid file(%v)", src)
	}
	video, err := ParseVideo(filepath.Join(f.W_DIR, src))
	if err != nil {
		log.D("FFCM parse video by src(%v) error->%v", src, err)
		return err
	}
	video.Info = info
	var tw, th, dur = video.Width, video.Height, int64(video.Duration * 1000000)
	log.D("FFCM add task by src(%v),width(%v),height(%v),duration(%v)", src, tw, th, dur)
	return f.S.AddTask(video, src, strings.TrimSuffix(src, ext), tw, th, dur)
}

func (f *FFCM_S) AddTaskH(hs *routing.HTTPSession) routing.HResult {
	var src string = hs.RVal("src")
	if len(src) < 1 {
		return hs.MsgResE3(2, "arg-err", "src argument is empty")
	}
	var err = f.AddTask(src)
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		err = util.Err("AddTask error->%v", err)
		log.E("%v", err)
		return hs.MsgResErr2(3, "srv-err", err)
	}
}

func (f *FFCM_S) Hand(pre string, mux *routing.SessionMux) {
	mux.H("^"+pre+"/status(\\?.*)?", f.S)
	mux.HFunc("^"+pre+"/v/addTask(\\?.*)?", f.AddTaskH)
	mux.HFunc("^"+pre+"/n/addTask(\\?.*)?", f.S.AddTaskH)
}
