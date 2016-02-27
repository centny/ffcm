package main

import (
	"fmt"
	"github.com/Centny/ffcm"
	"github.com/Centny/ffcm/mdb"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/smartio"
	"github.com/Centny/gwf/util"
	"os"
)

func usage() {
	fmt.Println(`Usage:
	ffcm -d width height max_width max_height	resize the dimensions by less max_width/max_height
	ffcm -c <configure file>					run client mode by configure file
	ffcm -mem <configure file>					run server on memory store mode by configure file
	ffcm -s <configure file>					run server on database store mode by configure file.
	ffcm -i <video file>						print video info
	ffcm -g <http url>							send http get request
		`)
}

var ef = os.Exit

func main() {
	if len(os.Args) < 2 {
		usage()
		ef(1)
		return
	}
	switch os.Args[1] {
	case "-d":
		res, err := ffcm.Dim2(os.Args[2:])
		if err == nil {
			fmt.Println(res)
		} else {
			fmt.Println(err)
			ef(1)
			return
		}
	case "-i":
		if len(os.Args) < 3 {
			usage()
			ef(1)
			return
		}
		var ffprobe = os.Getenv("FFPROBE_C")
		if len(ffprobe) > 0 {
			ffcm.FFPROBE_C = ffprobe
		}
		video, err := ffcm.ParseVideo(os.Args[2])
		if err == nil {
			fmt.Println(video.Width, " ", video.Height, " ", video.Duration)
		} else {
			fmt.Println(err.Error())
			ef(1)
		}
	case "-c":
		var cfg = "conf/ffcm_c.properties"
		if len(os.Args) > 2 {
			cfg = os.Args[2]
		}
		var fcfg_s = util.NewFcfg3()
		fcfg_s.InitWithFilePath2(cfg, true)
		fcfg_s.Print()
		redirect_l(fcfg_s)
		ffcm.RunFFCM_Cv(fcfg_s)
	case "-mem":
		var cfg = "conf/ffcm_s.properties"
		if len(os.Args) > 2 {
			cfg = os.Args[2]
		}
		var fcfg_s = util.NewFcfg3()
		fcfg_s.InitWithFilePath2(cfg, true)
		fcfg_s.Print()
		redirect_l(fcfg_s)
		var err = ffcm.InitDtcmS(fcfg_s, dtm.MemDbc, dtm.NewDoNoneH())
		if err != nil {
			fmt.Println(err)
			return
		}
		ffcm.RunFFCM_S_V(fcfg_s)
	case "-s":
		var cfg = "conf/ffcm_s.properties"
		if len(os.Args) > 2 {
			cfg = os.Args[2]
		}
		var fcfg_s = util.NewFcfg3()
		fcfg_s.InitWithFilePath2(cfg, true)
		fcfg_s.Print()
		redirect_l(fcfg_s)
		var err = ffcm.InitDtcmS(fcfg_s, mdb.MdbH_dc, dtm.NewDoNoneH())
		if err != nil {
			fmt.Println(err)
			return
		}
		ffcm.RunFFCM_S_V(fcfg_s)
	case "-g":
		if len(os.Args) < 3 {
			usage()
			ef(1)
			return
		}
		var res, err = util.HGet("%v", os.Args[2])
		if err == nil {
			fmt.Println(res)
		} else {
			fmt.Printf("request to %v error->%v", os.Args[2], err)
			ef(1)
		}
	default:
		usage()
		ef(1)
	}
}

func redirect_l(fcfg *util.Fcfg) {
	var out_l = fcfg.Val2("out_l", "")
	if len(out_l) > 0 {
		smartio.RedirectStdout3(out_l)
	}
	var err_l = fcfg.Val2("err_l", "")
	if len(err_l) > 0 {
		smartio.RedirectStdout3(err_l)
	}
}
