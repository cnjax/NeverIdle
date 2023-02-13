package waste

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000} //kB
//type Myspeedtest speedtest.Server

func Network(interval time.Duration) {
	for {
		user, err := speedtest.FetchUserInfo()
		if err != nil {
			fmt.Println("[NETWORK] Error when fetching user info:", err)
			time.Sleep(time.Minute)
			continue
		}
		serverList, err := speedtest.FetchServers(user)
		if err != nil {
			fmt.Println("[NETWORK] Error when fetching servers:", err)
			time.Sleep(time.Minute)
			continue
		}
		//targets, err := serverList.FindServer([]int{})
		//if err != nil {
		//	fmt.Println("[NETWORK] Error when finding target:", err)
		//	time.Sleep(time.Minute)
		//	continue
		//}

		// pick random

		//s := targets[rand.Int31n(int32(len(targets)))]
		s := serverList[rand.Int31n(int32(len(serverList)))]

		err = s.PingTest()
		if err != nil {
			s.Latency = -1
		}

		dlURL := strings.Split(s.URL, "/upload.php")[0]
		for i := 1; i <= 50; i++ {
			err := downloadRequest(context.Background(), dlURL, rand.Intn(9))
			if err != nil {
				break
			}
			time.Sleep(time.Second)
			err = uploadRequest(context.Background(), s.URL, rand.Intn(9))
			if err != nil {
				break
			}
			time.Sleep(time.Second)
		}
		//
		//err = s.DownloadTest(false)
		//if err != nil {
		//	s.DLSpeed = -1
		//}
		//
		//err = s.UploadTest(false)
		//if err != nil {
		//	s.ULSpeed = -1
		//}

		//		fmt.Println("[NETWORK] SpeedTest Ping:", s.Latency, ",", s.DLSpeed, ",", "Upload:", s.ULSpeed, "via", s.String())

		runtime.GC()
		time.Sleep(interval + time.Second*time.Duration(rand.Int31n(200)))
	}

}

func downloadRequest(ctx context.Context, dlURL string, w int) error {
	size := dlSizes[w]
	xdlURL := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"
	fmt.Println(xdlURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, xdlURL, nil)
	if err != nil {
		return err
	}
	doer := &http.Client{}
	resp, err := doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func uploadRequest(ctx context.Context, ulURL string, w int) error {
	size := ulSizes[w]
	reader := speedtest.NewRepeatReader((size*100 - 51) * 10)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ulURL, reader)
	req.ContentLength = reader.ContentLength
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	doer := &http.Client{}
	resp, err := doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}
