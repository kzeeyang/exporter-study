package main

type Tarffic_exporter struct {
	URI    string
	mutex  sync.Mutex
	client *http.Client

}

func main() {

}