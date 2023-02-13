package waste

import (
	"fmt"
	mrand "math/rand"
	"runtime"
	"time"
)

func CPU(interval time.Duration, runinterval time.Duration) {
	//var buffer []byte
	//if len(Buffers) > 0 {
	//	buffer = Buffers[0].B[:8*MiB]
	//} else {
	//	buffer = make([]byte, 8*MiB)
	//}
	//rand.Read(buffer)
	//
	//// construct XChaCha20 stream cipher
	//cipher, err := chacha20.NewUnauthenticatedCipher(buffer[:32], buffer[:24])
	//if err != nil {
	//	panic(cipher)
	//}

	for {
		//RunCPULoad(2, 5, mrand.Intn(10)+10)
		RunCPULoad(2, runinterval, mrand.Intn(10)+10)
		fmt.Println("Waiting")
		time.Sleep(interval + time.Duration(mrand.Intn(10))*time.Second)
	}
}

func RunCPULoad(coresCount int, timeSeconds time.Duration, percentage int) {
	runtime.GOMAXPROCS(coresCount)
	//fmt.Println("Current percentage is %d", percentage)

	//limiter := &cpulimit.Limiter{
	//	MaxCPUUsage:     50.0,                   // throttle if current cpu usage is over 50%
	//	MeasureInterval: time.Millisecond * 333, // measure cpu usage in an interval of 333 milliseconds
	//	Measurements:    3,                      // use the average of the last 3 measurements for cpu usage calculation
	//}
	//limiter.Start()
	//defer limiter.Stop()

	// second     ,s  * 1
	// millisecond,ms * 1000
	// microsecond,μs * 1000 * 1000
	// nanosecond ,ns * 1000 * 1000 * 1000

	// every loop : run + sleep = 1 unit

	// 1 unit = 100 ms may be the best
	unitHundresOfMicrosecond := 50
	runMicrosecond := unitHundresOfMicrosecond * percentage
	sleepMicrosecond := unitHundresOfMicrosecond*100 - runMicrosecond
	//	time.Sleep(time.Duration(timeSeconds) * time.Second)
	stopch := make(chan bool)
	for i := 0; i < coresCount; i++ {
		go func(stop chan bool) {
			runtime.LockOSThread()
			// endless loop
			for {
				select {
				case _ = <-stop:
					return
				case <-time.After(time.Millisecond):
					break
				}
				begin := time.Now()
				for {
					// run 100%
					//limiter.Wait()
					if time.Now().Sub(begin) > time.Duration(runMicrosecond)*time.Microsecond {
						break
					}
				}
				// sleep
				time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
			}
		}(stopch)
	}
	// how long
	time.Sleep(timeSeconds)
	stopch <- true
	stopch <- true
	close(stopch)
}
