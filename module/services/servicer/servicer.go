package servicer

import (
	"Panda/module/models/member"
)

type Services interface {
	User() member.UserStore
}
