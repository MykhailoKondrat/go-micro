package main

import (
	"context"
	"github.com/MykhailoKondrat/go-micro/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "propduct-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gb := handlers.NewGoodBuy(l)
	sm := http.NewServeMux()

	sm.Handle("/", hh)
	sm.Handle("/buy", gb)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

	//myChan := make(chan int)
	////defer close(myChan)
	////
	////fmt.Println(<-myChan)
	//wg := new(sync.WaitGroup)
	//wg.Add(3)
	//go increaseByOne(myChan, wg)
	//go increaseByOne(myChan, wg)
	//go printChan(myChan, wg)
	//
	//wg.Wait()
	//
	//fmt.Println("program is over")
	////go increaseByOne(myChan, wg)
	////fmt.Println(x)
}

//func increaseByOne(c chan int, wg *sync.WaitGroup) {
//	defer wg.Done()
//	fmt.Println("Running increase routine")
//	sum := 0
//	for i := 0; i < 4; i++ {
//		sum += i
//	}
//	c <- sum
//}
//func printChan(c chan int, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	fmt.Println("Running print routine")
//	x, y := <-c, <-c
//	fmt.Println(x, y)
//}
