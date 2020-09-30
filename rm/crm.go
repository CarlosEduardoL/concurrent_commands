package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/CarlosEduardoL/concurrent_commands/model"
	. "github.com/CarlosEduardoL/concurrent_commands/util"
	"github.com/CarlosEduardoL/concurrent_commands/walker"
	tm "github.com/buger/goterm"
	"github.com/nathan-fiscaletti/consolesize-go"
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
	objetivePath := os.Args[1]
	path, err := filepath.Abs(objetivePath)
	if err != nil {
		tm.Println(fmt.Errorf("Invalid Path %s", path))
		tm.Flush()
		return
	}
	cpuNumber := runtime.NumCPU()

	start := time.Now()
	files, folders, sizeChanel := walker.WalkSizeChannel(path)

	var total int64
	var removed int64

	removeChan := make(chan int64, cpuNumber)

	banner := make([]string, cpuNumber)

	var wg sync.WaitGroup

	wg.Add(cpuNumber)

	for i := 0; i < cpuNumber; i++ {
		go func(i int) {
			defer wg.Done()
			for file := range files {
				banner[i] = fmt.Sprintf("Removing file %s with size %s", file.Name(), ByteSize(file.Size()))
				if file.Remove() != nil {
					banner[i] = "Error"
				} else {
					banner[i] = "Removed"
					removeChan <- file.Size()
				}
			}
		}(i)
	}

	printting := true

	go func() {
		for {
			select {
			case size, ok := <-sizeChanel:
				if !ok {
					sizeChanel = nil
				} else {
					total += size
				}
			case remove, ok := <-removeChan:
				if !ok {
					removeChan = nil
				} else {
					removed += remove
				}
			}
			if sizeChanel == nil && removeChan == nil {
				break
			}
		}
	}()

	cols, _ := consolesize.GetConsoleSize()
	fixSize := func(text string) string {
		pathers := fmt.Sprintf("%s%vv", "%-", cols-1)
		return fmt.Sprintf(pathers, text)
	}

	onFiles := true

	go func() {
		for printting {
			time.Sleep(500 * time.Millisecond)
			//tm.Clear()
			tm.MoveCursor(0, 0)
			tm.Println(fixSize("Welcome to Concurrent rm command by CarlosEduardoL"))
			tm.Println(fixSize(fmt.Sprintf("Time elapsed: %s", time.Since(start))))
			tm.Println(fixSize(fmt.Sprintf("%s / %s", ByteSize(removed), ByteSize(total))))
			if onFiles {
				for i := 0; i < cpuNumber; i++ {
					tm.Println(fixSize(fmt.Sprint("Goroutine[", i, "]: ", banner[i])))
				}
			} else {
				tm.Println(fixSize("Removing Folders"))
				for i := 0; i < cpuNumber-1; i++ {
					tm.Println(fixSize(""))
				}
			}
			tm.Flush()

		}
	}()

	wg.Wait()
	close(removeChan)

	onFiles = false

	dirStack := model.NewStack()

	for folder := range folders {
		dirStack.Push(folder)
	}

	for i := dirStack.Size(); i >= 0; i-- {
		dir, err := dirStack.Pop()
		if err != nil {
			continue
		}
		if err := dir.(model.Folder).Remove(); err != nil {
			fmt.Println(err)
		}
	}
	printting = false
	tm.MoveCursorUp(cpuNumber)
	tm.Println(fixSize("Finish"))
	tm.Flush()
}
