package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	tm "github.com/buger/goterm"

	"github.com/CarlosEduardoL/concurrent_commands/model"
	"github.com/CarlosEduardoL/concurrent_commands/walker"
)

func main() {

	tm.Clear()
	tm.MoveCursor(0, 0)
	tm.Println("Welcome to Concurrent rm command by CarlosEduardoL")
	tm.Flush()
	if len(os.Args) < 2 {
		tm.Println(fmt.Errorf("Not enough arguments\ntry crm <folder path>"))
		tm.Flush()
		return
	}

	sourcePath := os.Args[1]
	src, err := filepath.Abs(sourcePath)
	if err != nil {
		fmt.Println(fmt.Errorf("Invalid Path %s", src))
		return
	}
	cpuNumber := runtime.NumCPU()
	files := make(chan model.File, cpuNumber)
	folders := make(chan model.Folder, cpuNumber)

	walker.Walk(src, files, folders)

	var total int64 = 0
	var count int64 = 0

	dirStack := model.NewStack()

	for folder := range folders {
		count++
		dirStack.Push(folder)
	}

	fmt.Println(count)

	size := dirStack.Size()

	for i := 0; i < size; i++ {
		_, err := dirStack.Pop()
		if err != nil {
			continue
		}
		total++
	}

	fmt.Println(total)

}
