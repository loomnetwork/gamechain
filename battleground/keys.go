package battleground

import "strings"

func InitDataKey() []byte {
	return []byte("init")
}

func NewUserKeySpace(userId string) *UserKeySpace {
	return &UserKeySpace{userId: strings.TrimSpace(userId)}
}

type UserKeySpace struct {
	userId string
}

func (u *UserKeySpace) AccountKey() []byte {
	return []byte("user:" + u.userId)
}

func (u *UserKeySpace) DecksKey() []byte {
	return []byte("user:" + u.userId + ":deck")
}

func (u *UserKeySpace) CardCollectionKey() []byte {
	return []byte("user:" + u.userId + ":cards")
}
