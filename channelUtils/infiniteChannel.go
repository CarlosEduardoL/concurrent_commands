package channelUtils

import (
	"github.com/CarlosEduardoL/concurrent_commands/model"
)

func InfiniteChanel() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		inQueue := model.NewQueue()
		outChannel := func() chan interface{} {
			if inQueue.IsEmpty() {
				return nil
			}
			return out
		}

		for !inQueue.IsEmpty() || in != nil {
			select {
			case value, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue.Enqueue(value)
				}
			case outChannel() <- inQueue.Front():
				inQueue.Dequeue()
			}
		}
		close(out)
	}()
	return in, out
}
