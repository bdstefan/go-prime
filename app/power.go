package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/bdstefan/go-deploy-poc/nosql"
)

const exp = 5

var computeChan = make(chan string)
var redis = nosql.GetRedisClient()
var logFile, _ = os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

func publish(n int) {
	for i := 2; i <= n; i++ {
		go computePower(i, exp)
	}
}

func computePower(n int, exp int) {
	value := int(math.Pow(float64(n), float64(exp)))
	key := fmt.Sprintf("%v:%v", n, exp)

	if redis.Ping() != nil {
		log.Println("Saved to redis ", key, value)
		redis.Set(key, value, 30000)
	} else {
		os.Exit(255)
	}

	computeChan <- fmt.Sprintf("%v ^ %v = %v", n, exp, value)
}

func displayOutput(n int, w http.ResponseWriter) {
	for i := 2; i <= n; i++ {
		result := <-computeChan
		fmt.Fprintln(w, result)
	}
}

func compute(n int, w http.ResponseWriter) {
	log.SetOutput(logFile)
	publish(n)
	displayOutput(n, w)
}
