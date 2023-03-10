package avitotask

import "reflect"

type User struct {
	Id       int    `json:"-"`
	UName    string `json:"uname" binding:"required"`
	Balance  int    `json:"balance" binding:"required"`
	Reserved int    `json:"reserved"`
}

type P2p struct {
	SId    int `json:"sid" binding:"required"`
	DId    int `json:"did" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

type Balance struct {
	UserID        int `json:"uid" binding:"required"`
	ChangeBalance int `json:"change"`
}

type Service struct {
	Id    int `json:"sid"`
	Price int `json:"price"`
}

type Order struct {
	SId int `json:"sid" binding:"required"`
	UId int `json:"uid" binding:"required"`
}

type Accounting struct {
	UsersID   int `json:"auid"`
	ServiceID int `json:"asid"`
}

func IsEmptyStruct(s interface{}) bool {
	return reflect.DeepEqual(s, reflect.Zero(reflect.TypeOf(s)).Interface())
}
