package ffcm

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"path/filepath"
	"strings"
)

var DTCM_S *dtm.DTCM_S = nil
var W_DIR string = "."

func InitDtcmS(fcfg *util.Fcfg, dbc dtm.DB_C, h dtm.DTCM_S_H) error {
	var err error
	DTCM_S, err = dtm.StartDTCM_S(fcfg, dbc, h)
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
	if DTCM_S == nil {
		return util.Err("server is not running")
	}
	var ffprobe_c = fcfg.Val("ffprobe_c")
	if len(ffprobe_c) > 0 {
		FFPROBE_C = ffprobe_c
	}
	var w_dir = fcfg.Val("w_dir")
	if len(w_dir) > 0 {
		W_DIR = w_dir
	}
	var listen = fcfg.Val("listen")
	routing.H("^/status(\\?.*)?", DTCM_S)
	routing.HFunc("^/addTask(\\?.*)?", AddTaskH)
	log.D("listen web server on %v", listen)
	return routing.ListenAndServe(listen)
}

func AddTask(src string) error {
	if DTCM_S == nil {
		panic("server is not running")
	}
	var ext = filepath.Ext(src)
	if len(ext) < 1 {
		return util.Err("invalid file(%v)", src)
	}
	video, err := ParseVideo(filepath.Join(W_DIR, src))
	if err != nil {
		log.D("FFCM parse video by src %v error->", src, err)
		return err
	}
	var tw, th, dur = video.Width, video.Height, int64(video.Duration * 1000000)
	log.D("FFCM add task by src(%v),width(%v),height(%v),duration(%v)", src, tw, th, dur)
	return DTCM_S.AddTask(video, src, strings.TrimSuffix(src, ext), tw, th, dur)
}

func AddTaskH(hs *routing.HTTPSession) routing.HResult {
	if DTCM_S == nil {
		return hs.MsgResE3(1, "srv-err", "servicer is not running")
	}
	var src string = hs.RVal("src")
	if len(src) < 1 {
		return hs.MsgResE3(2, "arg-err", "src argument is empty")
	}
	var err = AddTask(src)
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		err = util.Err("AddTask error->%v", err)
		log.E("%v", err)
		return hs.MsgResErr2(3, "srv-err", err)
	}
}
