package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
)

type Passwords []Password

type Password struct {
	Id bson.ObjectId	`bson:"_id"`
	Site string			`bson:"site"`
	Username string		`bson:"username"`
	Secret string		`bson:"secret"`
	OldSecret string	`bson:"oldsecret"`
}

func (u *User) AddPassword(site, username, secret string) error {
	var p Password
	p.Id = bson.NewObjectId()
	p.Site = site
	p.Username = username
	p.Secret = secret
	return m.db.C("user").UpdateId(u.Id, bson.M {
		"$push": bson.M{
			"passwords": p,
		},
	})
}

func (u *User) RemovePassword(site string) error {
	return m.db.C("user").UpdateId(u.Id, bson.M {
		"$pull": bson.M {
			"passwords": bson.M {
				"site": site,
			},
		},
	})
}

func (u *User) FindPassword(site string) *Password {
	for _, p := range(u.Passwords) {
		if p.Site == site {
			return &p
		}
	}
	return nil
}

func (u *User) UpdatePassword(id bson.ObjectId, site, username, secret string) error {
	return m.db.C("user").Update(bson.M{
		"_id":u.Id,
		"passwords._id": id,
	}, bson.M {
		"$set": bson.M {
			"passwords.$.site": site,
			"passwords.$.username": username,
			"passwords.$.secret": secret,
		},
	})
}

func (u *User) UpdateAllPasswords(passwords Passwords) error {
	if len (u.Passwords) != len(passwords) {
		return errors.New("error passwords")
	}
	err := m.db.C("user").Update(bson.M {
		"_id":u.Id,
	}, bson.M {
		"$unset": bson.M{
			"passwords": -1,
		},
	})
	if err != nil {
		return err
	}
	err = m.db.C("user").Update(bson.M {
		"_id":u.Id,
	}, bson.M {
		"$set": bson.M {
			"passwords": passwords,
		},
	})
	return nil
}