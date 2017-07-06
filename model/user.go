package model

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type User struct {
	Id bson.ObjectId	`bson:"_id"`
	Name string			`bson:"name"`
	Pass []byte			`bson:"pass"`

	Passwords Passwords	`bson:"passwords"`
}

func UserValidate(name string, pass []byte) (bool, error) {
	c, err := m.db.C("user").Find(bson.M{
		"name": name,
		"pass": pass,
	}).Count()
	if err != nil {
		return false, err
	}
	return c == 1, err
}

func UserAdd(name string, pass []byte) (error) {
	var u User
	u.Id = bson.NewObjectId()
	u.Name = name
	u.Pass = pass
	return m.db.C("user").Insert(&u)
}

func UserDelete(name string) (error) {
	var result User
	err := m.db.C("user").Find(bson.M{"name": name}).One(&result)
	if err != nil {
		return err
	}
	return m.db.C("user").Remove(result)
}

func UserFind(name string) (*User, error) {
	var user User
	err := m.db.C("user").Find(bson.M{
		"name": name,
	}).One(&user)
	return &user, err
}

func UserPassUpdate(name string, pass []byte) error {
	return m.db.C("user").Update(bson.M {
		"name": name,
	}, bson.M {
		"$set": bson.M{
			"pass": pass,
		},
	})
}