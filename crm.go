package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/CarlosEduardoL/concurrent_commands/resolver"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func main() {
	fmt.Println("Welcome to Concurrent rm command by CarlosEduardoL")
	if len(os.Args) < 2 {
		fmt.Println(fmt.Errorf("Not enough arguments\ntry crm <folder path>"))
		return
	}
	objetivePath := os.Args[1]
	path, err := resolver.ResolvePATH(objetivePath)
	if err != nil {
		fmt.Println(fmt.Errorf("Invalid Path %s", path))
		return
	}
	items, folders, size := resolver.DirToStack(path)
	cpuNumber := runtime.NumCPU()

	var total int64 = size

	p := mpb.New(mpb.WithWidth(64))

	bar := p.AddBar(total,
		mpb.PrependDecorators(decor.Counters(decor.UnitKiB, "% .1f / % .1f")),
		mpb.AppendDecorators(decor.Percentage()),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			),
		),
	)

	chanel := make(chan int64, cpuNumber)

	var wg sync.WaitGroup

	for i := 0; i < cpuNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {

				element, err := items.Pop()
				if err != nil {
					return
				}

				info, _ := os.Stat(element)
				err = os.Remove(element)

				if err != nil {
					fmt.Printf("error on %s\nerror:%s", element, err)
				}
				chanel <- info.Size()
			}
			for {
				folder, err := folders.Pop()
				if err != nil {
					return
				}
				err = os.Remove(folder)

				if err != nil {
					fmt.Printf("error on %s\nerror:%s", folder, err)
				}
			}

		}()
	}

	go func() {
		defer close(chanel)
		wg.Wait()
	}()

	for size := range chanel {
		bar.IncrInt64(size)
	}

	bar.SetTotal(total, true)

	p.Wait()

}
