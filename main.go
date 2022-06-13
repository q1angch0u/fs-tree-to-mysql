package main

import (
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"my-tree/dao"
	"sync"
	"time"
)

type Path struct {
	ID       uint `gorm:"primarykey"`
	Path     string
	ParentId int
}

var db *gorm.DB
var wg sync.WaitGroup

func processData(parentId int, path ...string) int {

	if len(path) == 0 {
		return -1
	}

	var paths []Path

	for _, v := range path {
		paths = append(paths, Path{
			Path:     v,
			ParentId: parentId,
		})
	}

	result := db.Create(&paths)
	if result.Error != nil {
		log.Fatalf("insert db error, err: %s", result.Error)
		return -1
	}

	return int(paths[0].ID)

}

func getFileList(path string, n int) {
	defer wg.Done()
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if file.IsDir() {
			dataId := processData(n, file.Name())
			if dataId == -1 {
				return
			}
			wg.Add(1)
			go getFileList(path+"/"+file.Name(), dataId)
		} else {
			processData(n, file.Name())
		}
	}

}

func main() {

	startTime := time.Now().Unix()

	db = dao.GetDB()
	db.AutoMigrate(&Path{})

	wg.Add(1)
	getFileList("~", 0)
	wg.Wait()

	finishTime := time.Now().Unix()

	fmt.Printf("run time: %d\n", finishTime-startTime)

}
