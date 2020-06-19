package main
import "net/http"
import "fmt"
import "bytes"
import "log"
import "sync"
import "time"

func publish (count int) {
	http.Post("http://127.0.0.1:4151/topic/delete?topic=test", "", nil)
	for x := 0; x<count; x++ {
		var data string
		for y := 0; y<publish_multi; y++ {
			data = data+"\n"+string(y)+"-"+string(x)
		}
		_, err := http.Post("http://127.0.0.1:4151/mpub?topic=test", "application/text", bytes.NewBuffer([]byte(data)))
		if err != nil {
			fmt.Printf("http send: error %v\n", err)
		}
	}
	stats := "http://127.0.0.1:4151/stats"
	resp, err := http.Get(stats)
	if err != nil {
		fmt.Printf("http send: error %v\n", err)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		stats := buf.String()
		fmt.Printf("http get nsqd stats: %s\n", stats)
	}
	fmt.Printf("%s\n", stats)
	fmt.Printf("done publishing: message_count:%d\n", publish_multi*count)
}

var publish_multi int = 1000
var publish_count int = 50000
var message_count int = 0
var serve_mutex = &sync.Mutex{}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	serve_mutex.Lock()
	message_count++
	i := message_count
	serve_mutex.Unlock()
	if (i == publish_count || i % (publish_count/20) == 0) {
		t := time.Now()
		x := (i*100)/publish_count
		fmt.Printf("[%s] message_count:%d missing:%d %d%%\n", t.Format(time.RFC1123), i, publish_count-i, x)
	}
}

func serve() {
	http.HandleFunc("/", httpHandler)
	address := "192.168.0.111:8080"
	fmt.Printf("serve: listening on \"%s\"\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}



func main () {
	go signalHandling()
	publish(publish_count/publish_multi)
	serve()
}
