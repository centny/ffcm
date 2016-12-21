package ffcm

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/util"
)

func init() {
	dtm.AddCreator("Video", dtm.FuncCreator(NewAbsV))
}

var FFPROBE_C = "ffprobe"

type Video struct {
	Filename string      `bson:"filename" json:"filename"`
	Duration float64     `bson:"duration" json:"duration"`
	Size     int64       `bson:"size" json:"size"`
	Width    int64       `bson:"width" json:"width"`
	Height   int64       `bson:"height" json:"height"`
	Alias    string      `bson:"alias" json:"alias"`
	Kvs      util.Map    `bson:"-" json:"-"`
	Info     interface{} `bson:"info" json:"info"`
}

func ParseData(data, reg string) ([]*util.Fcfg, error) {
	dataes := regexp.MustCompile(reg).Split(data, -1)
	var cfgs []*util.Fcfg
	for _, data = range dataes {
		data = strings.Trim(data, " \t\n")
		if len(data) < 1 {
			continue
		}
		var stream = util.NewFcfg3()
		err := stream.InitWithData(data)
		if err == nil {
			cfgs = append(cfgs, stream)
		}

	}
	return cfgs, nil
}

func ParseFormat(path string) (*util.Fcfg, error) {
	data, err := util.Exec(FFPROBE_C, "-show_format", path)
	if err != nil {
		return nil, util.Err("exec(%v) error->%v", FFPROBE_C, err)
	}
	cfgs, err := ParseData(data, "\\[[/]*FORMAT\\]")
	var cfg *util.Fcfg
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	return cfg, err
}

func ParseStreams(path string) ([]*util.Fcfg, error) {
	data, err := util.Exec(FFPROBE_C, "-show_streams", path)
	if err != nil {
		return nil, util.Err("exec(%v) error->%v", FFPROBE_C, err)
	}
	return ParseData(data, "\\[[/]*STREAM\\]")
}

func ParseVideo(path string) (*Video, error) {
	var video = &Video{}
	var cfg = util.NewFcfg3()
	//parse format
	format, err := ParseFormat(path)
	if err != nil {
		return nil, err
	}
	video.Filename = format.Val("filename")
	video.Duration = format.FloatValV("duration", 0)
	video.Size = format.Int64ValV("size", 0)
	cfg.Merge2("format", format)
	//parse streams
	streams, err := ParseStreams(path)
	for _, stream := range streams {
		if stream.Val("codec_type") == "video" {
			video.Width = stream.Int64ValV("width", 0)
			video.Height = stream.Int64ValV("height", 0)
		}
		cfg.Merge2("s"+stream.Val("index"), stream)
	}
	video.Kvs = cfg.Map
	return video, err
}

func Dim(whs []string) ([]string, error) {
	if len(whs) < 4 {
		return nil, util.Err("arguments less 4")
	}
	ivs, err := util.ParseInts(whs)
	if err != nil {
		return nil, err
	}
	var tw, th int = ivs[0], ivs[1]
	if tw > ivs[2] {
		tw = ivs[2]
		th = int(float64(ivs[1]) / float64(ivs[0]) * float64(tw))
	}
	if th > ivs[3] {
		th = ivs[3]
		tw = int(float64(ivs[0]) / float64(ivs[1]) * float64(th))
	}
	if (tw % 2) > 0 {
		tw += 1
	}
	if (th % 2) > 0 {
		th += 1
	}
	return []string{fmt.Sprintf("%v", tw), fmt.Sprintf("%v", th)}, nil
}

func Dim2(whs []string) (string, error) {
	vals, err := Dim(whs)
	return strings.Join(vals, "x"), err
}

func VerifyVideo(va, vb string) error {
	videoa, err := ParseVideo(os.Args[2])
	if err != nil {
		return err
	}
	videob, err := ParseVideo(os.Args[2])
	if err != nil {
		return err
	}
	if int(videoa.Duration) != int(videob.Duration) {
		return fmt.Errorf("the duration verify fail to %v(%v),%v(%v)", va, videoa.Duration, vb, videob.Duration)
	}
	fmt.Printf("Verify duration ok by %v(%v),%v(%v)\n", va, videoa.Duration, vb, videob.Duration)
	return nil
}
