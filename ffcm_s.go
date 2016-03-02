package ffcm

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"path/filepath"
	"strings"
)

var SRV *dtm.DTCM_S = nil

func InitDtcmS(fcfg *util.Fcfg, dbc dtm.DB_C, h dtm.DTCM_S_H) error {
	var err error
	SRV, err = dtm.StartDTCM_S(fcfg, dbc, h)
	return err
}

func RunFFCM_S(cfg string, dbc dtm.DB_C, h dtm.DTCM_S_H) error {
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath2(cfg, true)
	fcfg.Print()
	var err = InitDtcmS(fcfg, dbc, h)
	if err == nil {
		err = RunFFCM_S_V(fcfg)
	}
	return err
}

func RunFFCM_S_V(fcfg *util.Fcfg) error {
	if SRV == nil {
		return util.Err("server is not running")
	}
	var ffprobe_c = fcfg.Val("ffprobe_c")
	if len(ffprobe_c) > 0 {
		FFPROBE_C = ffprobe_c
	}
	var listen = fcfg.Val("listen")
	routing.H("^/status(\\?.*)?", SRV)
	routing.HFunc("^/addTask(\\?.*)?", SRV.AddTaskH)
	routing.Shared.Print()
	log.D("listen web server on %v", listen)
	return routing.ListenAndServe(listen)
}

type AbsV struct {
	*dtm.AbsN
}

func NewAbsV(sec string, cfg *util.Fcfg) (dtm.Abs, error) {
	var n, err = dtm.NewAbsN(sec, cfg)
	return &AbsV{AbsN: n.(*dtm.AbsN)}, err
}

func (a *AbsV) Build(dtcm *dtm.DTCM_S, id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error) {
	var src = fmt.Sprintf("%v", args[0])
	video, err := ParseVideo(filepath.Join(a.WDir, src))
	if err != nil {
		err = util.Err("AbsV parse video by src(%v) error->%v", src, err)
		log.E("%v", err)
		return nil, nil, nil, err
	}
	video.Info = info
	video.Alias = a.Alias
	var mv, _ = util.Json2Map(util.S2Json(video))
	var tw, th, dur = video.Width, video.Height, int64(video.Duration * 1000000)
	var dst interface{}
	if len(args) > 1 {
		dst = args[1]
	} else {
		dst = strings.TrimSuffix(src, filepath.Ext(src))
	}
	log.D("AbsV adding task by src(%v),width(%v),height(%v),duration(%v)", src, tw, th, dur)
	return id, mv, []interface{}{
		src, dst, tw, th, dur,
	}, nil
}
