package mdb

import (
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/ffcm"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MdbH struct {
	Name string
	MC   func(string) *tmgo.Collection
}

func (m *MdbH) C() *tmgo.Collection {
	return m.MC(m.Name)
}

//add task to db
func (m *MdbH) Add(t *dtm.Task) error {
	log.D("MdbH add task by id(%v)", t.Id)
	t.Time = util.Now()
	return m.C().Insert(t)
}

//update task to db
func (m *MdbH) Update(t *dtm.Task) error {
	t.Time = util.Now()
	return m.C().Update(bson.M{"_id": t.Id}, t)
}

//delete task to db
func (m *MdbH) Del(t *dtm.Task) error {
	log.D("MdbH delete task by id(%v)", t.Id)
	return m.C().RemoveId(t.Id)
}

//list task from db
func (m *MdbH) List(running []string, status string, skip, limit int) (int, []*dtm.Task, error) {
	var and = []bson.M{}
	and = append(and, bson.M{
		"$or": []bson.M{
			bson.M{"mid": ""},
			bson.M{"mid": bson.M{"$exists": false}},
		},
	})
	and = append(and, bson.M{
		"$or": []bson.M{
			bson.M{
				"next": bson.M{
					"$lt": util.Now(),
				},
			},
			bson.M{"next": bson.M{"$exists": false}},
		},
	})
	if len(running) > 0 {
		and = append(and, bson.M{
			"_id": bson.M{
				"$nin": running,
			},
		})
	}
	if len(status) > 0 {
		and = append(and, bson.M{
			"status": status,
		})
	}
	var ts []*dtm.Task
	var err = m.C().Find(bson.M{"$and": and}).Skip(skip).Limit(limit).All(&ts)
	if err != nil {
		return 0, nil, err
	}
	var mid = util.MID()
	var rts []*dtm.Task
	for _, t := range ts {
		_, err = m.C().Find(bson.M{
			"$and": []bson.M{
				bson.M{"_id": t.Id},
				bson.M{"$or": []bson.M{
					bson.M{"mid": ""},
					bson.M{"mid": bson.M{"$exists": false}},
				}},
			},
		}).Apply(tmgo.Change{
			Update: bson.M{
				"$set": bson.M{
					"mid": mid,
				},
			},
		}, nil)
		if err == nil {
			rts = append(rts, t)
		} else if err == tmgo.ErrNotFound {
			continue
		} else {
			return 0, nil, err
		}
	}
	var total int = 0
	if err == nil {
		total, err = m.C().Count()
	}
	return total, ts, err
}

//find task by id
func (m *MdbH) Find(id string) (*dtm.Task, error) {
	var ts []*dtm.Task
	var err = m.C().Find(bson.M{"_id": id}).All(&ts)
	var task *dtm.Task
	if err == nil && len(ts) > 0 {
		task = ts[0]
	}
	return task, err
}

//database creator
func MdbH_dc(uri, name string) (dtm.DbH, error) {
	db, err := mgo.NewMDbs(uri, name)
	if err == nil {
		return &MdbH{
			Name: "ffcm_task",
			MC:   db.C,
		}, err
	} else {
		return nil, err
	}
}

func DefaultDbc(uri, name string) (dtm.DbH, error) {
	return &MdbH{
		Name: "ffcm_task",
		MC:   mgo.C,
	}, nil
}

func StartTest(cfgs, cfgc string, h dtm.DTCM_S_H) {
	ffcm.StartTest(cfgs, cfgc, DefaultDbc, h)
}
