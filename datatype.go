package main

import (
	"encoding/json"
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const saveLocation = "/moodleD"

type events struct {
	Events []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Modulename  string `json:"modulename"`
		Instance    int    `json:"instance"`
		Eventtype   string `json:"eventtype"`
		Timestart   int    `json:"timestart"`

		Course struct {
			ID       int    `json:"id"`
			Fullname string `json:"fullname"`
		} `json:"course"`

		Normalisedeventtypetext string `json:"normalisedeventtypetext"`
		Action                  struct {
			Name       string `json:"name"`
			URL        string `json:"url"`
			Actionable bool   `json:"actionable"`
		} `json:"action"`
		URL string `json:"url"`
	} `json:"events"`

	Neweventtimestamp int `json:"neweventtimestamp"`
	Date              struct {
		Seconds   int    `json:"seconds"`
		Minutes   int    `json:"minutes"`
		Hours     int    `json:"hours"`
		Mday      int    `json:"mday"`
		Wday      int    `json:"wday"`
		Mon       int    `json:"mon"`
		Year      int    `json:"year"`
		Yday      int    `json:"yday"`
		Weekday   string `json:"weekday"`
		Month     string `json:"month"`
		Timestamp int    `json:"timestamp"`
	} `json:"date"`
}
type eventData struct {
	description string
	throw       time.Time
	ID          int
}

type Api struct {
	Token        string `json:"token"`
	PrivateToken string `json:"privatetoken"`
	Base         string `json:"base"`
}

func getApi(base, username, password string) *Api {
	var data Api
	data.Base = base
	data.authenticate(username, password)
	return &data
}

func (api *Api) authenticate(username, password string) {
	url := fmt.Sprintf("%s/login/token.php?username=%s&password=%s&service=moodle_mobile_app", api.Base, username, password)
	resp, err := http.Get(url)
	assertErr(err)
	buffer, err := ioutil.ReadAll(resp.Body)
	assertErr(err)
	resp.Body.Close()
	fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &api)
	assertErr(err)

}

var IDStore []int

func (api *Api) getEvents(t time.Time) []eventData {
	if IDStore == nil {
		IDStore = make([]int, 0)
	}
	arguments := map[string]string{
		"year":  strconv.Itoa(t.Year()),
		"month": strconv.Itoa(int(t.Month())),
		"day":   strconv.Itoa(t.Day()),
	}
	buffer := api.request("core_calendar_get_calendar_day_view", arguments)
	fmt.Println(string(buffer))
	var d events
	err := json.Unmarshal(buffer, &d)
	assertErr(err)
	eventDataSlice := make([]eventData, 0)
	for _, e := range d.Events {
		if !contains(IDStore, e.ID) {
			IDStore = append(IDStore, e.ID)
			tm := time.Unix(int64(e.Timestart), 0)
			s := fmt.Sprintf("You have a '%s-%s' event in 5 minutes: \n '%s' in %s at %v:%v \n", e.Eventtype, e.Modulename, e.Name, e.Course.Fullname, tm.Hour(), tm.Minute())
			eventDataSlice = append(eventDataSlice, eventData{
				description: s,
				throw:       tm,
				ID:          e.ID,
			})
		}
	}
	return eventDataSlice
}

func (api *Api) request(verb string, arguments map[string]string) []byte {
	s := fmt.Sprintf("%s/webservice/rest/server.php?wstoken=%s&moodlewsrestformat=json", api.Base, api.Token)
	u, err := url.Parse(s)
	assertErr(err)
	q := u.Query()
	q.Set("wsfunction", verb)
	for k, v := range arguments {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	assertErr(err)
	buffer, err := ioutil.ReadAll(resp.Body)
	assertErr(err)
	resp.Body.Close()
	return buffer
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func assertErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func (d eventData) notify() {
	d.throw = d.throw.Add(-5 * time.Minute)
	fmt.Printf("created thread for event %v, will notify at %v:%v\n", d.ID, d.throw.Hour(), d.throw.Minute())
	dur := d.throw.Sub(time.Now())
	time.Sleep(dur)
	beeep.Notify("Moodle Notification", d.description, "")
}
