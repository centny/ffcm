package ffcm

import (
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	StartTest2("ffcm_s.properties", "ffcm_c.properties", dtm.NewDoNoneH())
	time.Sleep(2 * time.Second)
	if SRV == nil || CLIENT == nil {
		t.Error("error")
		return
	}
}
