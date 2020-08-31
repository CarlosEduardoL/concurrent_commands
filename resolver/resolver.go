package resolver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/CarlosEduardoL/thread_commands/structures"
)

//NoValidError is a not valid path
type NoValidError struct {
}

func (NoValidError) Error() string {
	return "Not valid Path"
}

//ResolvePATH return the path to the file or error
func ResolvePATH(path string) (string, error) {
	var finalPath string
	currentPath, err := getCurrentDir()
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(path, ".\\") {
			if err != nil {
				return "", err
			}
			finalPath = strings.Replace(path, ".", currentPath, 1)
		} else if path == "." {
			if err != nil {
				return "", err
			}
			finalPath = strings.Replace(path, ".", fmt.Sprintf("%s%s", currentPath, "\\"), 1)
		} else {
			finalPath = path
		}
	} else {
		if strings.HasPrefix(path, "./") {
			if err != nil {
				return "", err
			}
			finalPath = strings.Replace(path, ".", currentPath, 1)
		} else if path == "." {
			if err != nil {
				return "", err
			}
			finalPath = strings.Replace(path, ".", fmt.Sprintf("%s%s", currentPath, "/"), 1)
		} else {
			finalPath = path
		}
	}
	if isValid(finalPath) {
		return finalPath, nil
	}
	return finalPath, NoValidError{}
}

func DirToStack(dir string) (structures.Stack, int64) {
	stack := structures.NewStack()

	size := dirToStack(dir, stack)

	return stack, size
}

func dirToStack(dir string, stack structures.Stack) int64 {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var size int64

	//stack.Push(dir)

	for _, f := range files {
		if f.IsDir() {
			size += dirToStack(fmt.Sprintf("%s%s%s", dir, "\\", f.Name()), stack)
		} else {
			stack.Push(fmt.Sprintf("%s%s%s", dir, "\\", f.Name()))
			size += f.Size()
		}
	}
	return size
}

func isValid(fp string) bool {
	if info, err := os.Stat(fp); err == nil {
		return info.IsDir()
	}
	return false
}

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}
