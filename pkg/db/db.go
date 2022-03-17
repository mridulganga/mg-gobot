package db

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/mridulganga/mg-gobot/pkg/constants"
	"github.com/mridulganga/mg-gobot/pkg/models"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DB struct {
	db *leveldb.DB
}

func NewDB(path string) DB {
	db, err := leveldb.OpenFile(constants.DBPath, nil)
	if err != nil {
		panic(err)
	}
	return DB{
		db: db,
	}
}

func (d DB) GetUser(id string) *models.User {
	log.Debugf("get user %v", id)
	value, err := d.db.Get([]byte("user-"+id), nil)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil
	}
	user := models.User{}
	json.Unmarshal(value, &user)
	return &user
}

func (d DB) PutUser(u models.User) {
	log.Debugf("put user %v", u)
	jsonBytes, _ := json.Marshal(u)
	err := d.db.Put([]byte("user-"+u.ID), jsonBytes, nil)
	if err != nil {
		log.Errorf("error: %v", err)
	}
}

func (d DB) ListUsers() []models.User {
	log.Debug("list users")
	users := []models.User{}
	iter := d.db.NewIterator(util.BytesPrefix([]byte("user-")), nil)
	for iter.Next() {
		log.Debugf("iter %v", string(iter.Key()))
		u := models.User{}
		json.Unmarshal(iter.Value(), &u)
		users = append(users, u)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil
	}
	return users
}
