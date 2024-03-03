package youtube

import (
	"errors"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ApiUrl = "localhost"
)

const (
	attempts  = 10
	wait_time = 15
	base      = "/aiuzubot/v3/bot/"
)

func StartBot(botName string, liveId string) error {
	done := false
	n := 0
	var err error
	getUrl := ApiUrl + base + botName + "/start"
	if liveId != "" {
		getUrl = getUrl + "?liveId=" + liveId
	}
	for !done && n < attempts {
		_, err = doGet(getUrl)
		if err != nil {
			log.Warnf("[startBot]Attempt: %d failed: %s", n, err.Error())
			n = n + 1
			time.Sleep(wait_time * time.Second)
		} else {
			done = true
		}
	}
	return err
}

func StopBot(botName string) error {
	done := false
	n := 0
	var err error
	getUrl := ApiUrl + base + botName + "/stop"
	for !done && n < attempts {
		_, err = doGet(getUrl)
		if err != nil {
			log.Warnf("[startBot]Attempt: %d failed: %s", n, err.Error())
			n = n + 1
			time.Sleep(wait_time * time.Second)
		} else {
			done = true
		}
	}
	return err
}

func doGet(u string) (*http.Response, error) {
	r, err := http.Get(u)
	err = handleResponse(r, err)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func handleResponse(r *http.Response, e error) error {
	if e != nil {
		return e
	} else if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	} else {
		bs, _ := io.ReadAll(r.Body)
		return errors.New("Error " + r.Status + " " + string(bs))
	}
}
