package walker

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/CarlosEduardoL/concurrent_commands/channelUtils"

	"github.com/CarlosEduardoL/concurrent_commands/model"
)

func WalkSizeChannel(path string) (<-chan model.File, <-chan model.Folder, <-chan int64) {
	cpuNumber := runtime.NumCPU()
	files, folders := Walk(path)
	sizeChanel := make(chan int64, cpuNumber)
	filesOut := make(chan model.File, cpuNumber)

	snifferIn, snifferOut := channelUtils.InfiniteChanel()

	go func() {
		for file := range files {
			sizeChanel <- file.Size()
			snifferIn <- file
		}
		close(snifferIn)
		close(sizeChanel)
	}()

	go func() {
		for file := range snifferOut {
			filesOut <- file.(model.File)
		}
		close(filesOut)
	}()

	return filesOut, folders, sizeChanel
}

func Walk(path string) (<-chan model.File, <-chan model.Folder) {
	cpuNumber := runtime.NumCPU()
	fileChannel := make(chan model.File, cpuNumber)
	folderChannel := make(chan model.Folder, cpuNumber)

	filesIn, filesOut := channelUtils.InfiniteChanel()
	foldersIn, foldersOut := channelUtils.InfiniteChanel()

	go func() {
		for file := range filesOut {
			fileChannel <- file.(model.File)
		}
		close(fileChannel)
	}()
	go func() {
		for folder := range foldersOut {
			folderChannel <- folder.(model.Folder)
		}
		close(folderChannel)
	}()

	go func() {
		walk(filepath.Dir(path), filepath.Base(path), filesIn, foldersIn)
		close(filesIn)
		close(foldersIn)
	}()

	return fileChannel, folderChannel
}

func walk(base, path string, fileChannel chan<- interface{}, folderChannel chan<- interface{}) {

	fullPath := filepath.Join(base, path)

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	folderChannel <- model.NewFolder(base, path)
	for _, f := range files {
		if f.IsDir() {
			walk(base, filepath.Join(path, f.Name()), fileChannel, folderChannel)
		} else {
			fileChannel <- model.NewFile(base, filepath.Join(path, f.Name()), f.Size())
		}
	}
}
