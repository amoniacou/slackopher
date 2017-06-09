package slackopher

import (
	"github.com/nlopes/slack"
	"github.com/patrickmn/go-cache"
	"time"
	"errors"
)

type User struct {
	Info *slack.User
}

type Users struct {
	Cache *cache.Cache
}

func NewUsers() *Users {
	return &Users{
		Cache:cache.New(60 * time.Minute, 120 * time.Minute),
	}
}

func (u *Users) GetUser(id string) (User, error) {
	user, found := u.Cache.Get(id)
	if found {
		return user.(User), nil
	} else {
		return User{}, errors.New("User not found")
	}
}

func (u *Users) AddUser(user *slack.User) User {
	usr := User{
		Info: user,
	}
	u.Cache.Set(user.ID, usr, cache.DefaultExpiration)
	return usr
}

