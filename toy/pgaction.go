package toy

import (
	"gopkg.in/pg.v5"
	"log"
)

type User struct {
	id   int64
	name string
}

func FetchUser(db *pg.DB, userid int64) (*User, error) {
	user := User{
		id:   userid,
		name: "uncle wang",
	}
	log.Printf("Get user %d successful")
	return &user, nil
}
