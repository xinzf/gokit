package storage

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"sync"
)

var Mongo *mongo
var mongoOnce sync.Once

type mongo struct {
	session *mgo.Session
	options MongoOptions
}

func (this *mongo) Init(opt ...MongoOption) error {
	var err error

	mongoOnce.Do(func() {

		Mongo = &mongo{options: newMongoOptions(opt...)}
		if Mongo.options.Addr == "" {
			err = errors.New("mongo config's addr is empty")
			return
		}

		if Mongo.options.Name == "" {
			err = errors.New("mongo config's name is empty")
			return
		}

		Mongo.session, err = mgo.Dial(Mongo.options.Addr)
		if err != nil {
			err = fmt.Errorf("mongodb init failed, err:" + err.Error())
			return
		}

		mgo.SetDebug(Mongo.options.Debug)
		if Mongo.options.logger != nil {
			mgo.SetLogger(Mongo.options.logger)
		}

		Mongo.session.SetMode(mgo.Monotonic, true)
		Mongo.options.logger.Info("Mongo","Mongo inited.")
	})
	return err
}

func (this *mongo) Use(name string) *mgo.Collection {
	s := this.session.Copy()
	return s.DB(this.options.Name).C(name)
}

func (this *mongo) Close() {
	this.session.Close()
}
