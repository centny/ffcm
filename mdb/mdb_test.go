package mdb

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/netw/dtm"
	"testing"
	"time"
)

func TestMdb(t *testing.T) {
	dbh, err := MdbH_dc("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	dbh.(*MdbH).C().RemoveAll(nil)
	var task = &dtm.Task{
		Id: "xxx",
	}
	err = dbh.Add(task)
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = dbh.Find(task.Id)
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
	_, err = dbh.Find(task.Id)
	if err != nil {
		t.Error("error")
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
	//
	_, err = MdbH_dc("127.0.0.1:23442", "sdfs")
	if err == nil {
		t.Error("error")
		return
	}
	//
	DefaultDbc("uri", "name")
	//
	mgo.AddDefault("cny:123@loc.w:27017/cny", "cny")
	StartTest("../ffcm_s.properties", "../ffcm_c.properties", dtm.NewDoNoneH())
	//
	time.Sleep(time.Second)
	//
	fmt.Println("done...")
}
