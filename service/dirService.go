package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"orealtest/config"
	"path/filepath"
)

type dir struct {
	Path string    `json:"path"`
	Dirs []dirItem `json:"dirs"`
}

type dirItem struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type fileInfo struct {
	dirItems  []dirItem
	dirCount  int
	fileCount int
	totalSize int64
}

func getDirItems(path string) *fileInfo {
	var dirItems []dirItem
	var dirCount, fileCount int
	var totalSize int64

	fileInfos, _ := ioutil.ReadDir(path)

	for _, info := range fileInfos {
		dirItems = append(dirItems, dirItem{
			Name:  info.Name(),
			IsDir: info.IsDir(),
			Size:  info.Size(),
		})
		if info.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
		totalSize += info.Size()
	}
	return &fileInfo{dirItems, dirCount, fileCount, totalSize}
}

func DirHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	path := filepath.Join(config.Root, p)

	fileInfoCh := getDirItems(path)
	res, _ := json.Marshal(&dir{
		Path: path,
		Dirs: fileInfoCh.dirItems,
	})

	w.Write(res)
}
