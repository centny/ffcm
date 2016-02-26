package mdb

import (
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/netw/dtm"
	tmgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MdbH struct {
	Name string
	Db   *mgo.MDbs
}

func (m *MdbH) C() *tmgo.Collection {
	return m.Db.C(m.Name)
}

//add task to db
func (m *MdbH) Add(t *dtm.Task) error {
	return m.C().Insert(t)
}

//update task to db
func (m *MdbH) Update(t *dtm.Task) error {
	return m.C().Update(bson.M{"_id": t.Id}, t)
}

//delete task to db
func (m *MdbH) Del(t *dtm.Task) error {
	return m.C().RemoveId(t.Id)
}

//list task from db
func (m *MdbH) List() ([]*dtm.Task, error) {
	var ts []*dtm.Task
	var err = m.C().Find(nil).All(&ts)
	return ts, err
}

//find task by id
func (m *MdbH) Find(id string) (*dtm.Task, error) {
	var ts []*dtm.Task
	var err = m.C().Find(bson.M{"_id": id}).All(&ts)
	if err != nil {
		return nil, err
	}
	if len(ts) < 1 {
		return nil, nil
	} else {
		return ts[0], nil
	}
}

//database creator
func MdbH_dc(uri, name string) (dtm.DbH, error) {
	db, err := mgo.NewMDbs(uri, name)
	return &MdbH{
		Name: "ffcm_task",
		Db:   db,
	}, err
}
