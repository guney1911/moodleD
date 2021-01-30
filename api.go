package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getApi(base, username, password string) (a Api) {
	a.Base = base
	a.authenticate(username, password)
	return
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

var _eventIDStore []int

func (api *Api) getEvents(t time.Time) []eventData {
	if _eventIDStore == nil {
		_eventIDStore = make([]int, 0)
	}
	arguments := map[string]string{
		"year":  strconv.Itoa(t.Year()),
		"month": strconv.Itoa(int(t.Month())),
		"day":   strconv.Itoa(t.Day()),
	}
	buffer := api.request("core_calendar_get_calendar_day_view", arguments)
	var d events
	err := json.Unmarshal(buffer, &d)
	assertErr(err)
	eventDataSlice := make([]eventData, 0)
	for _, e := range d.Events {
		if !contains(_eventIDStore, e.ID) {
			_eventIDStore = append(_eventIDStore, e.ID)
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

func (api *Api) getCourseIDs() courseData {
	data := api.request("core_course_search_courses", map[string]string{"criterianame": "search", "criteriavalue": " "})
	var c courseData
	err := json.Unmarshal(data, &c)
	assertErr(err)
	return c
}

func (api *Api) getContent(id int) (data courseContent) {

	res := api.request("core_course_get_contents", map[string]string{"courseid": strconv.Itoa(id), "options[0][name]": "excludecontents", "options[0][value]": "true"})
	err := json.Unmarshal(res, &data)
	assertErr(err)
	return
}
func (api *Api) request(verb string, arguments map[string]string) []byte {
	s := fmt.Sprintf("%s/webservice/rest/server.php?wstoken=%s&moodlewsrestformat=json", api.Base, api.Token)
	u, err := url.Parse(s)
	assertErr(err)
	q := u.Query()
	q.Set("wsfunction", verb)

	if arguments != nil {
		for k, v := range arguments {
			q.Set(k, v)
		}
	}

	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	assertErr(err)
	buffer, err := ioutil.ReadAll(resp.Body)
	assertErr(err)
	resp.Body.Close()
	return buffer
}
