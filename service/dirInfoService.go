package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"orealtest/config"
	"path/filepath"
	"sync"
)

var numCh = make(chan struct{}, config.MaxRoutineNum)

type dirInfo struct {
	Path      string `json:"path"`
	DirCount  int    `json:"dirCount"`
	FileCount int    `json:"fileCount"`
	TotalSize int64  `json:"totalSize"`
}

func getDirInfo(path string, ch chan<- *fileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	numCh <- struct{}{}
	fileInfos, _ := ioutil.ReadDir(path)
	<-numCh

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			wg.Add(1)
			go getDirInfo(filepath.Join(path, fileInfo.Name()), ch, wg)
		} else {
			ch <- getDirItems(path)
		}
	}
}

func DirInfoHandlder(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	path := filepath.Join(config.Root, p)

	var dirCount, fileCount int
	var totalSize int64

	var wg sync.WaitGroup
	ch := make(chan *fileInfo)

	wg.Add(1)
	go getDirInfo(path, ch, &wg)

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for v := range ch {
		dirCount += v.dirCount
		fileCount += v.fileCount
		totalSize += v.totalSize
	}

	res, _ := json.Marshal(&dirInfo{
		Path:      path,
		DirCount:  dirCount,
		FileCount: fileCount,
		TotalSize: totalSize,
	})

	w.Write(res)
}
