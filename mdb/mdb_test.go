package mdb

import (
	"fmt"
	"github.com/Centny/gwf/netw/dtm"
	"testing"
)

func TestMdb(t *testing.T) {
	dbh, err := MdbH_dc("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var task = &dtm.Task{
		Id: "xxx",
	}
	err = dbh.Add(task)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = dbh.Update(task)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts, err := dbh.List()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ts) != 1 {
		t.Error("error")
		return
	}
	if task.Id != ts[0].Id {
		t.Error("error")
		return
	}
	err = dbh.Del(task)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts, err = dbh.List()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ts) > 0 {
		t.Error("error")
		return
	}
	fmt.Println("done...")
}
