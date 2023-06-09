package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

func save(object any, session *mgo.Session, collection string) {
	err := session.DB("game_stats").C(collection).Insert(object)
	if err != nil {
		fmt.Printf("Could not save data in collection %s: %s", collection, err)
	}
}

func initializeMongoSession() (*mgo.Session, error) {
	return mgo.Dial("mongodb://maciek1651g:V7q3qYGppE2eaflc@ac-trubrvu-shard-00-02.pmjetfp.mongodb.net:27017,ac-trubrvu-shard-00-00.pmjetfp.mongodb.net:27017,ac-trubrvu-shard-00-01.pmjetfp.mongodb.net:27017,mdcluster0.pmjetfp.mongodb.net")
}
