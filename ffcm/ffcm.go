package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Centny/ffcm"
	"github.com/Centny/ffcm/mdb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/smartio"
	"github.com/Centny/gwf/util"
)

func usage() {
	fmt.Println(`Usage:
	ffcm -d width height max_width max_height	resize the dimensions by less max_width/max_height
	ffcm -c <configure file>					run client mode by configure file
	ffcm -mem <configure file>					run server on memory store mode by configure file
	ffcm -s <configure file>					run server on database store mode by configure file.
	ffcm -i <video file>						print video info
	ffcm -g <http url>							send http get request
	ffcm -c <process url> <input> <width> <height> <max_width> <max_height> <tmp> <out> <result> 
	ffcm -verify <video first> <video second>	verify the two vidoe is having same duration.
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
	case "-cov_v":
		if len(os.Args) < 11 {
			usage()
			ef(1)
			return
		}
		err := do_cov_video(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			ef(1)
			return
		}
	case "-cov_a":
		if len(os.Args) < 7 {
			usage()
			ef(1)
			return
		}
		err := do_cov_audio(os.Args[2:])
		if err != nil {
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
		err := ffcm.RunFFCM_Cv(fcfg_s)
		fmt.Println(err)
		smartio.ResetStd()
		time.Sleep(time.Second)
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
		if err == nil {
			err = ffcm.RunFFCM_S_V(fcfg_s)
		}
		fmt.Println(err)
		smartio.ResetStd()
		time.Sleep(time.Second)
	case "-verify":
		if len(os.Args) < 4 {
			usage()
			ef(1)
			return
		}
		_, err := ffcm.VerifyVideo(os.Args[2], os.Args[3])
		if err == nil {
			fmt.Println("Verify Success")
			ef(0)
		} else {
			fmt.Println("Verify Fail")
			fmt.Println(err)
			ef(1)
		}
		return
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
		if err == nil {
			err = ffcm.RunFFCM_S_V(fcfg_s)
		}
		fmt.Println(err)
		smartio.ResetStd()
		time.Sleep(time.Second)
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
			fmt.Printf("request to %v error->%v\n", os.Args[2], err)
			ef(1)
		}
	default:
		usage()
		ef(1)
	}
}

func redirect_l(fcfg *util.Fcfg) {
	var out_l = fcfg.Val2("out_l", "")
	var err_l = fcfg.Val2("err_l", "")
	fmt.Printf("redirect stdout to file(%v) and stderr to file(%v)\n", out_l, err_l)
	if len(out_l) > 0 {
		smartio.RedirectStdout3(out_l)
	}
	if len(err_l) > 0 {
		smartio.RedirectStderr3(err_l)
	}
	log.SetWriter(os.Stdout)
}

func do_cov_video(args []string) error {
	fmt.Printf("run_ff arguments list:\n"+
		"	%v\n	%v\n	%v\n	%v\n	%v\n	%v\n	%v\n	%v\n	%v\n",
		args[0], args[1], args[2],
		args[3], args[4], args[5],
		args[6], args[7], args[8])
	err := os.MkdirAll(filepath.Dir(args[6]), os.ModePerm)
	if os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(filepath.Dir(args[7]), os.ModePerm)
	if os.IsNotExist(err) {
		return err
	}
	res, err := ffcm.Dim2(args[2:6])
	if err != nil {
		return err
	}
	exe := exec.Command("ffmpeg", "-progress", args[0], "-i", args[1], "-s", res, "-y", args[6])
	exe.Stderr = os.Stderr
	exe.Stdout = os.Stdout
	err = exe.Run()
	if err != nil {
		return err
	}
	_, err = ffcm.VerifyVideo(args[1], args[6])
	if err != nil {
		return err
	}
	src, err := os.OpenFile(args[6], os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Printf("do copy %v to %v...\n", args[6], args[7])
	_, err = util.Copyp(args[7], src)
	if err != nil {
		src.Close()
		return err
	}
	src.Close()
	fmt.Printf("do clear tmp file %v...\n", args[6])
	err = os.Remove(args[6])
	if err != nil {
		fmt.Printf("[Warning]remove tmp file %v fail", args[6])
	}
	fmt.Printf(`
----------------result----------------
[text]
%v
[/text]
	
	`, args[8])
	return nil
}

func do_cov_audio(args []string) error {
	fmt.Printf("run_ff arguments list:\n"+
		"	%v\n	%v\n	%v\n	%v\n	%v\n",
		args[0], args[1], args[2], args[3], args[4])
	err := os.MkdirAll(filepath.Dir(args[2]), os.ModePerm)
	if os.IsNotExist(err) {
		return err
	}
	err = os.MkdirAll(filepath.Dir(args[3]), os.ModePerm)
	if os.IsNotExist(err) {
		return err
	}
	// exe := exec.Command("ffmpeg", "-progress", args[0], "-i", args[1], "-y", args[2])
	exe := exec.Command("ffmpeg", "-i", args[1], "-y", args[2])
	exe.Stderr = os.Stderr
	exe.Stdout = os.Stdout
	err = exe.Run()
	if err != nil {
		return err
	}
	src, err := os.OpenFile(args[2], os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Printf("do copy %v to %v...\n", args[2], args[3])
	_, err = util.Copyp(args[3], src)
	if err != nil {
		src.Close()
		return err
	}
	src.Close()
	fmt.Printf("do clear tmp file %v...\n", args[2])
	err = os.Remove(args[2])
	if err != nil {
		fmt.Printf("[Warning]remove tmp file %v fail", args[2])
	}
	fmt.Printf(`
----------------result----------------
[text]
%v
[/text]
	
	`, args[4])
	return nil
}
