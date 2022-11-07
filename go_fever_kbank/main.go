package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	//Declare example struct
	type inputStruct struct {
		Value      int
		Wg         *sync.WaitGroup
		RetryCount int
	}

	//Initiate channel to push example records into go routines
	inputChannel := make(chan inputStruct)
	retryChannel := make(chan inputStruct)
	defer close(inputChannel)
	defer close(retryChannel)

	//Declare example go routine workers.
	goRoutineThreads := 200

	//Loop to extract threads to receive working records from channel
	for threads := 0; threads < goRoutineThreads; threads++ {
		//Declare func to handle retry input work in same go routine
		go func() {
			for rtc := range retryChannel {
				inputChannel <- rtc
			}
		}()

		//Declare func to accept input work
		go func(th int) {
			for workStruct := range inputChannel {

				//example time sleep to validate threads
				time.Sleep(20 * time.Millisecond)

				//if response is ok use wg Done to accept response
				//example case: if input % 15 ==0, retry 3 times
				if workStruct.Value%15 == 0 && workStruct.RetryCount < 3 {
					fmt.Printf("print retry count %d ,input from thread : %d, value : %d\n", workStruct.RetryCount, th, workStruct.Value)
					//put same work with add retry count to retry channel
					workStruct.RetryCount = workStruct.RetryCount + 1
					retryChannel <- workStruct
				} else if workStruct.Value%2 == 0 && workStruct.RetryCount == 3 {
					//if retry 3 times, force done
					fmt.Printf("print retry count %d ,input from thread : %d, value : %d force done\n", workStruct.RetryCount, th, workStruct.Value)
					workStruct.Wg.Done()
				} else {
					//simple case without conditions
					fmt.Printf("print input from thread : %d, value : %d, retry %d\n", th, workStruct.Value, workStruct.RetryCount)
					workStruct.Wg.Done()
				}
			}

		}(threads)
	}

	//Declare example records
	exampleRecord := 20000

	//Declare wait group to count all records to prevent error or leak action in go routines
	var wg sync.WaitGroup
	wg.Add(exampleRecord)

	for i := 0; i < exampleRecord; i++ {

		//fmt.Println("call api :", i+1)
		inputChannel <- inputStruct{
			Value:      i,
			Wg:         &wg,
			RetryCount: 0,
		}
	}
	//Use wg wait for all records return response from go routines
	wg.Wait()
	fmt.Println("exit from thread")

}
