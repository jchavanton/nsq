package main
import "os"
import "os/signal"
import "fmt"
import "time"
import "sync"

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			handle()
		}
	}
}

func (t *Task) Stop() {
	close(t.closed)
	t.wg.Wait()
}

func handle() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 200)
	}
}

func signalHandling() {
	task := &Task{
		ticker: time.NewTicker(time.Second * 2),
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			unique := 0
			duplicate := 0
			missing := 0
			for x:=0;x<publish_count;x++ {
				if messages[x] == 1 {
					unique++
				} else if messages[x] == 0 {
					missing++
				} else if messages[x] > 1 {
					duplicate = duplicate + messages[x] - 1
				}
			}
			completed := message_count*100/publish_count
			fmt.Printf("completed[%d%%] published[%d] message[%d]\n", completed, publish_count, message_count)
			fmt.Printf("unique[%d] missing[%d] duplicate[%d]\n", unique, missing, duplicate)
			os.Exit(1)
		}
	}()

	task.Run()
}
