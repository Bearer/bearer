package main

import (
	ormlib "github.com/beego/beego/v2/client/orm"
	"github.com/example/mydb/v2"
)

func init() {
	ormlib.RegisterDriver("mydb", mydb.MyDB)

	ormlib.RegisterDataBase(`default`, `mydb`, "root:root@/orm_test?charset=utf8")
}
