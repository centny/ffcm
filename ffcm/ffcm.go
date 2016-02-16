package main

import (
	"fmt"
	"github.com/Centny/ffcm"
	"github.com/Centny/gwf/netw/dtm"
	"os"
)

func usage() {
	fmt.Println(`Usage:
	ffcm -d width height max_width max_height	resize the dimensions by less max_width/max_height
	ffcm -c <configure file>					run client mode by configure file
	ffcm -mem <configure file>					run server on memory store mode by configure file
	ffcm -s <configure file>					run server on database store mode by configure file.
	ffcm -i <video file>						print video info
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
		fmt.Println(os.Args[2])
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
		ffcm.RunFFCM_C(cfg)
	case "-mem":
		var cfg = "conf/ffcm_s.properties"
		if len(os.Args) > 2 {
			cfg = os.Args[2]
		}
		ffcm.RunFFCM_S(cfg, dtm.MemDbc, dtm.NewDoNoneH())
	case "-s":
		var cfg = "conf/ffcm_s.properties"
		if len(os.Args) > 2 {
			cfg = os.Args[2]
		}
		ffcm.RunFFCM_S(cfg, dtm.MemDbc, dtm.NewDoNoneH())
	default:
		usage()
		ef(1)
	}
}
