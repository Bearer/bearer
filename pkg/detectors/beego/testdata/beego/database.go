package main

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	ignoreme "github.com/other"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
	ignoreme.RegisterDataBase("ignore", "ignore", "something different")
}
