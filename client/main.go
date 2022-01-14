package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func GetBody() {
	for {
		response, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Fatalln(err)
		}
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(bytes))
		response.Body.Close()
	}
	wg.Done()
}

func main() {
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go GetBody()
	}
	wg.Wait()
}
