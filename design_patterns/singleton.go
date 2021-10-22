package main

import (
	"fmt"
	"sync"
	"time"
)

type DataBase struct{}

func (DataBase) CreateSingleConnection() {
	fmt.Println("Creating Singleton for DB")
	time.Sleep(2 * time.Second)
}

var db *DataBase
var mutex sync.Mutex

func getDatabaseInstance() *DataBase {
	mutex.Lock()
	defer mutex.Unlock()
	if db == nil {
		fmt.Println("Creating DB Connection")
		db = &DataBase{}
		db.CreateSingleConnection()
	} else {
		fmt.Println("DB Already created")
	}
	return db
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getDatabaseInstance()
		}()
	}
	wg.Wait()
}
