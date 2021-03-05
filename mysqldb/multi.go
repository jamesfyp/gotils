package mysqldb

import (
	"gorm.io/gorm"
	"sync"
)

var dsns = make(map[string]string, 0)
var instances sync.Map

func Add(name, dsn string) {

	dsns[name] = dsn

}

func GetInstance(name string) *gorm.DB {

	ins, ok := instances.Load(name)

	if ok {
		return ins.(*gorm.DB)
	}

	dsn, ok1 := dsns[name]

	if !ok1 {

		return nil
	}

	ins = initDB(dsn)

	instances.Store(name, ins)

	return ins.(*gorm.DB)
}
