package ffcm

import (
	"github.com/Centny/gwf/netw/dtm"
	"time"
)

func StartTest(cfgs, cfgc string, dbc dtm.DB_C, h dtm.DTCM_S_H) {
	go RunFFCM_S(cfgs, dbc, h)
	time.Sleep(time.Second)
	go RunFFCM_C(cfgc)
}

func StartTest2(cfgs, cfgc string, h dtm.DTCM_S_H) {
	StartTest(cfgs, cfgc, dtm.MemDbc, h)
}
