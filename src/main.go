package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	"github.com/clarkbains/counter"
)

//Test my sample package.
func timoutTest() {
	//Create a timer that takes 200 ms to increment
	counter := counter.NewCounter()
	counter.SetAddDelay(200 * time.Millisecond)

	//Keep track of how long our deadline should be
	t := time.Millisecond * 220
	
	for ;; {

		//Gradually decrement our deadline length
		if t.Milliseconds() > 2  {
			//Decrement the duration
			t = time.Duration((t.Milliseconds() - 1) * int64(time.Millisecond))
			fmt.Printf("New timeout is %s\n", t)

		}
		//Create a context that times out after the deadline length.
		ctx, cancel := context.WithTimeout(context.Background(), t) 

		//Add one to our counter, passing context. If the counter can increment it before the end of the context, then all works, otherwise we get an err
		err := counter.AddOneWithContext(ctx)
		if err != nil {
			//If we have an error, our context deadline passed.
			fmt.Printf("Error while adding one: %s\n", err)
			cancel()
			return
		} 
		cancel()
		//Display the value of the counter
		counter.LogValue()
		
		time.Sleep(time.Second)
	}
}

//Print the system hostname in an infinite loop
func hostnamePrinter () {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Could not get hostname\n")
		return
	}
	
	for ;; {
		fmt.Printf("My Hostname is %s\n", name)
		time.Sleep(time.Second * 10)
	}
}

func main () {
	//On any request with / as a prefix (all of them)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		//I store my html inside the /views dir in my contatiner
		fp := path.Join("/views", r.URL.Path)
		fmt.Printf("Trying to read file: %s\n", fp)
		//Read the bytes, and then proceed only if second clause is true (no error)
		if data, err := os.ReadFile(fp); err == nil {
			//Data is a byte array, which we can easily send back through the response writer
			w.Write(data)
			return
		}
		//Otherwise we can send a message back to the user
		fmt.Fprintf(w, "Could not find file! (\"%s\")", r.URL)
	})

	//Start a bunch of things as goroutines (the run in the background)
	var wg sync.WaitGroup

	//Start the timeout test.
	wg.Add(1)
	go func ()  {
		timoutTest()
		wg.Done()
	}()

	//Start the hostname goroutine
	wg.Add(1)
	go func ()  {
		hostnamePrinter()
		wg.Done()
	}()

	//Start the http listener goroutine
	wg.Add(1)
	go func ()  {
		http.ListenAndServe(":8080", nil)
		wg.Done()
	}()

	//Go exits when it gets to the end, so we use a wait group as a semaphore to wait until it has a value of 0. wg.Add added one, and wg.Done subracted one
	wg.Wait()
}