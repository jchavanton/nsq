package main
import "net/http"
import "fmt"
import "bytes"
import "log"
import "sync"
import "time"
import "strings"
import "strconv"
import "io/ioutil"

var ip string = "192.168.0.111"

func publish (count int) {
	http.Post("http://127.0.0.1:4151/topic/delete?topic=test", "", nil)
	counter := 0
	for x := 0; x<count; x++ {
		var data string = ""
		for y := 0; y<publish_multi; y++ {
			data = data+"\n"+strconv.Itoa(counter)
			counter++
		}

		var client = &http.Client{
			Timeout: time.Second * 10,
		}
		req, err := http.NewRequest("POST", "http://127.0.0.1:4151/mpub?topic=test", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/text")
		if err != nil {
			fmt.Printf("error setting header: %v\n", err)
		}
		_, err = client.Do(req)

		if err != nil {
			fmt.Printf("error http send: %v\n", err)
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
var messages [50000]int
var message_count int = 0
var serve_mutex = &sync.Mutex{}


func httpHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Millisecond)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	serve_mutex.Lock()
	message_count++
	i := message_count
	num, _ := strconv.Atoi(string(b))
	messages[num]++
	// c := messages[num]
	serve_mutex.Unlock()

	// fmt.Printf("updated:[%d]=[%d]\n", num, c)
	if (i == publish_count || i % (publish_count/20) == 0) {
		t := time.Now()
		x := (i*100)/publish_count
		fmt.Printf("[%s] message_count:%d missing:%d %d%%\n", t.Format(time.RFC1123), i, publish_count-i, x)
	}
}

func serve() {
	http.HandleFunc("/", httpHandler)
	address := ip + ":8080"
	fmt.Printf("serve: listening on \"%s\"\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}



func main () {
	go signalHandling()
	publish(publish_count/publish_multi)
	serve()
}
