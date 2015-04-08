package main

import (
	"labix.org/v2/mgo"
	"log"
)

var (
	mgoSession *mgo.Session
)

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MongoAddr)
		//log.Println("mongo addr:", MongoAddr)
		if err != nil {
			log.Println(err) // no, not really
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}

func withCollection(collection string, safe *mgo.Safe, s func(*mgo.Collection) error) error {
	session, err := getSession()
	if session == nil {
		return err
	}
	defer session.Close()

	session.SetSafe(safe)
	c := session.DB("sports").C(collection)
	return s(c)
}

func search(collection string, query interface{}, selector interface{},
	skip, limit int, sortFields []string, total *int, result interface{}) error {

	q := func(c *mgo.Collection) error {
		qy := c.Find(query)
		var err error

		if selector != nil {
			qy = qy.Select(selector)
		}

		if total != nil {
			if *total, err = qy.Count(); err != nil {
				return err
			}
		}

		if result == nil {
			return err
		}

		if limit > 0 {
			qy = qy.Limit(limit)
		}
		if skip > 0 {
			qy = qy.Skip(skip)
		}
		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.All(result)
	}

	return withCollection(collection, nil, q)
}

func findOne(collection string, query interface{}, sortFields []string, result interface{}) error {
	q := func(c *mgo.Collection) error {
		var err error
		qy := c.Find(query)

		if result == nil {
			return err
		}

		if len(sortFields) > 0 {
			qy = qy.Sort(sortFields...)
		}

		return qy.One(result)
	}

	return withCollection(collection, nil, q)
}

func updateId(collection string, id interface{}, change interface{}, safe bool) error {
	update := func(c *mgo.Collection) error {
		return c.UpdateId(id, change)
	}

	if safe {
		return withCollection(collection, &mgo.Safe{}, update)
	}
	return withCollection(collection, nil, update)
}
