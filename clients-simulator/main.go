package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ResponseApi int
const (
    OK ResponseApi = iota
    Timeout
    Refused
	EOF
	Unknow
)
type ResponseK8s struct{
	res ResponseApi
	t time.Duration 
}
const (
	routinesCount = 60
)
func main() {

	for i:=0; i<=9; i++ {
		fmt.Printf("cycle %d\n",i)
		cycle()
	}
}
func cycle() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	ch := make(chan ResponseK8s, routinesCount)


	go startApiCalls(ctxTimeout,ch)



	 func() {
		routineOK := make([]time.Duration,0)
		routineTimeout := make([]time.Duration,0)
		routineRefused := make([]time.Duration,0)
		routineEOF := make([]time.Duration,0)
		routineUnknow := make([]time.Duration,0)

		for {
			select {
			case k8sRes, ok := <-ch:
						if !ok {
							fmt.Println("Channel closed.")
							okStr := buildLogFromDurations("OK",routineOK, true)
							timeoutStr := buildLogFromDurations("Timeout",routineTimeout, true)
							refusedStr := buildLogFromDurations("Refused",routineRefused, false)
							EOFStr := buildLogFromDurations("EOF",routineEOF, true)
							unknowStr := buildLogFromDurations("Unknow",routineUnknow, false)				
							fmt.Printf("Timeout reached. %s %s %s %s %s\n", okStr, timeoutStr, refusedStr, EOFStr, unknowStr)
							return
						}

						switch k8sRes.res {
						case OK:
							routineOK = append(routineOK, k8sRes.t)
						case Timeout:
							routineTimeout = append(routineTimeout, k8sRes.t)
						case Refused:
							routineRefused = append(routineRefused, k8sRes.t)
						case EOF:
							routineEOF = append(routineEOF, k8sRes.t)
						default:
							routineUnknow = append(routineUnknow, k8sRes.t)
						}
				/*case <-ctxTimeout.Done():	
				return	*/		
			}
		}
	}()
}

func buildLogFromDurations(prefix string, durations []time.Duration, withAverage bool ) string {

	simple := fmt.Sprintf("%s %d/%d", prefix, len(durations),routinesCount)

	if !withAverage {
		return simple
	}

	return fmt.Sprintf("%s with average %d",simple , CalculateAverageDuration(durations))
}

func startApiCalls(ctx context.Context, ch chan ResponseK8s){
	// Start goroutines to fill the channel
	var wg sync.WaitGroup

	// Increment the wait group counter for the goroutines in fillCh
	wg.Add(routinesCount)
	for i:=0; i<= routinesCount-1; i++ {
				// go callHttp(ctx, ch, &wg)

				go func() {



					defer wg.Done()

	url := "http://0.0.0.0:2000/order/pizza"
	//url := "http://localhost:8080/order/pizza"

	// Make the POST request
	timePostStarted := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		if strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
			ch <- ResponseK8s{
				res: Refused,
				t:   time.Since(timePostStarted),
			}
			return
		}

		if strings.Contains(err.Error(), "EOF"){
			ch <- ResponseK8s{
				res: EOF,
				t:   time.Since(timePostStarted),
			}
			return
		}

		fmt.Println(err)

		ch <- ResponseK8s{
			res: Unknow,
			t:   time.Since(timePostStarted),
		}
		
		return
	}
	defer resp.Body.Close()

	//dead,_ := ctx.Deadline()
	//elapsedTime := time.Since(dead)

	select {
	case <-ctx.Done():
		ch <- ResponseK8s{
			res: Timeout,
			t:   time.Since(timePostStarted),
		}
	case ch <- ResponseK8s{
		res: OK,
		t:   time.Since(timePostStarted),
	}:
	}
				}()

						}

		wg.Wait()
		close(ch)

	}

	func CalculateAverageDuration(durations []time.Duration) time.Duration {
		if len(durations) == 0 {
			return 0
		}
	
		var total time.Duration
		for _, d := range durations {
			total += d
		}
	
		average := total / time.Duration(len(durations))
		return average
	}