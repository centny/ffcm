package mdb

import (
	"github.com/Centny/ffcm"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/util"
	"w.gdy.io/dyf/mgo"
	"w.gdy.io/dyf/mgo/bson"
)

type MdbH struct {
	Name string
	MC   func(string) *mgo.Collection
}

func (m *MdbH) C() *mgo.Collection {
	return m.MC(m.Name)
}

// add task to db
func (m *MdbH) Add(t *dtm.Task) error {
	log.D("MdbH add task by id(%v)", t.Id)
	t.Time = util.Now()
	return m.C().Insert(t)
}

// update task to db
func (m *MdbH) Update(t *dtm.Task) error {
	t.Time = util.Now()
	return m.C().Update(bson.M{"_id": t.Id}, t)
}

// delete task to db
func (m *MdbH) Del(t *dtm.Task) error {
	log.D("MdbH delete task by id(%v)", t.Id)
	return m.C().RemoveId(t.Id)
}

func (m *MdbH) ClearSyncTask() error {
	_, err := m.C().UpdateAll(bson.M{"mid": util.MID()}, bson.M{"$set": bson.M{"mid": ""}})
	return err
}

// list task from db
func (m *MdbH) List(mid string, running []string, status string, skip, limit int) (int, []*dtm.Task, error) {
	if len(mid) > 0 {
		_, err := m.C().UpdateAll(
			bson.M{
				"mid": mid,
				"_id": bson.M{"$nin": running},
			},
			bson.M{
				"$set": bson.M{"mid": ""},
			})
		if err != nil {
			return 0, nil, err
		}
	}
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
	var rts []*dtm.Task
	if len(mid) > 0 {
		for _, t := range ts {
			err = m.C().Find(bson.M{
				"$and": []bson.M{
					bson.M{"_id": t.Id},
					bson.M{"$or": []bson.M{
						bson.M{"mid": ""},
						bson.M{"mid": bson.M{"$exists": false}},
					}},
				},
			}).Apply(mgo.Change{
				Update: bson.M{
					"$set": bson.M{
						"mid": mid,
					},
				},
			}, nil)
			if err == nil {
				rts = append(rts, t)
			} else if err == mgo.ErrNotFound {
				continue
			} else {
				return 0, nil, err
			}
		}
	} else {
		rts = ts
	}
	var total int = 0
	if err == nil {
		total, err = m.C().Count()
	}
	return total, rts, err
}

// find task by id
func (m *MdbH) Find(id string) (*dtm.Task, error) {
	var ts []*dtm.Task
	var err = m.C().Find(bson.M{"_id": id}).All(&ts)
	var task *dtm.Task
	if err == nil {
		if len(ts) > 0 {
			task = ts[0]
		} else {
			err = util.NOT_FOUND
		}
	}
	return task, err
}

// database creator
func MdbH_dc(uri, name string) (dtm.DbH, error) {
	db, err := mgo.DialShared("mongodb://" + name + "@" + uri)
	if err == nil {
		return &MdbH{
			Name: "ffcm_task",
			MC:   db.Collection,
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
