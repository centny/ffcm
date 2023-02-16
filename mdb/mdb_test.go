package mdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/Centny/gwf/netw/dtm"
	"w.gdy.io/dyf/mgo"
)

func TestMdb(t *testing.T) {
	mgo.DialShared("mongodb://cny:123@loc.w:27017/cny")
	dbh, err := DefaultDbc("", "")
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
	_, ts, err := dbh.List("", nil, "", 0, 30)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(ts)
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
	_, ts, err = dbh.List("", nil, "", 0, 30)
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
	mgo.DialShared("mongodb://cny:123@loc.w:27017/cny")
	StartTest("../ffcm_s.properties", "../ffcm_c.properties", dtm.NewDoNoneH())
	//
	time.Sleep(time.Second)
	//
	fmt.Println("done...")
}
