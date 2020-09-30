package resolver

import (
	"os"
	"sync"

	"github.com/CarlosEduardoL/concurrent_commands/walker"

	"github.com/CarlosEduardoL/concurrent_commands/model"
)

func DirToStack(dir string) (files model.FileStack, folders model.Stack, size int64) {
	files = model.NewFileStack()
	folders = model.NewStack()

	size = dirToStack(dir, files, folders)

	return
}

func dirToStack(dir string, stack model.FileStack, folders model.Stack) int64 {

	wg := sync.WaitGroup{}

	wg.Add(2)

	files := make(chan model.File, 8)
	foldersChan := make(chan string, 8)

	walker.Walk(dir, files, foldersChan)

	var size int64

	go func() {
		for file := range files {
			stack.Push(file)
			size += file.Size()
		}
		wg.Done()
	}()

	go func() {
		for file := range files {
			stack.Push(file)
			size += file.Size()
		}
		wg.Done()
	}()

	wg.Wait()
	return size
}

func isValid(fp string) bool {
	if info, err := os.Stat(fp); err == nil {
		return info.IsDir()
	}
	return false
}
