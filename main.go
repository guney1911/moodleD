package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var _state userData

func main() {
	boolptr := flag.Bool("login", false, "-login : start the login wizard")
	refresh := flag.Int("refresh", 20, "-refresh=xx : refresh every xx minutes. Default is 20 min.")
	flag.Parse()
	if *boolptr {
		login()
		return
	}

	err := load(saveLocation, &_state)
	api := &_state.ApiData

	if err != nil {
		fmt.Println("Couldn't read save file!")
		beeep.Alert("MoodleD", "moodleD couldn't read or find save file! \n Please start the program from a terminal with the -login option to login to Moodle.", "")
	}
	//checkNewContent(api)
	//createNotifyThreads(api.getEvents(time.Now()))
	//createNotifyThreads(api.getEvents(time.Date(2021,1,11,0,0,0,0,time.Local)))
	go checkNew(api, *refresh)
	beeep.Notify("MoodleD", "moodleD started", "")
	select {}

}

func checkNewContent(api *Api) {
	if _state.Courses.Courses == nil {
		fmt.Println("No course-list found. Requesting a new one.")
		saveCourses(api)
	}
	if _state.ContentIDStore == nil {
		fmt.Println("No content-list found. Requesting a new one ")
		_state.ContentIDStore = make(map[string]bool)

		for _, cours := range _state.Courses.Courses {
			content := api.getContent(cours.ID)

			for _, s := range content {
				_state.ContentIDStore[strconv.Itoa(s.ID)] = true

				for _, module := range s.Modules {
					_state.ContentIDStore[generateModuleId(s.ID, module.ID)] = true

				}
			}
		}
		save(_state, saveLocation)
		return
	}
	var newItems bool
	changes := make(map[string]int)
	for _, cours := range _state.Courses.Courses {
		content := api.getContent(cours.ID)

		for _, s := range content {

			if !_state.ContentIDStore[strconv.Itoa(s.ID)] {
				_state.ContentIDStore[strconv.Itoa(s.ID)] = true
				newItems = true
				changes[cours.Shortname]++
			}
			for _, module := range s.Modules {
				if !_state.ContentIDStore[generateModuleId(s.ID, module.ID)] {
					_state.ContentIDStore[generateModuleId(s.ID, module.ID)] = true
					newItems = true
					changes[cours.Shortname]++
				}
			}
		}
	}
	if newItems {
		notifyNewContent(changes)
		go save(_state, saveLocation)
	} else {
		fmt.Println("No new content found")
	}
}

func notifyNewContent(changes map[string]int) {

	for s, i := range changes {
		beeep.Notify("MoodleD", fmt.Sprintf("There were %v new changes in %s.", i, s), "")
	}

}

func saveCourses(api *Api) {
	data := api.getCourseIDs()
	_state.Courses = data
	save(_state, saveLocation)
}

func createNotifyThreads(a []eventData) {
	for _, e := range a {
		go e.notify()
	}
	s := fmt.Sprintf("Created %v new notifiers.", len(a))
	beeep.Notify("MoodleD", s, "")
}

func checkNew(api *Api, refresh int) {
	for {
		beeep.Notify("MoodleD", "Refreshed", "")
		createNotifyThreads(api.getEvents(time.Now()))
		go checkNewContent(api)
		time.Sleep(time.Duration(refresh) * time.Minute)

	}
}
func save(state userData, location string) {
	c, _ := os.UserConfigDir()
	location = c + location
	d, err := json.Marshal(state)
	assertErr(err)
	_, err = os.Create(location)
	assertErr(err)
	err = ioutil.WriteFile(location, d, 0644)
	assertErr(err)
	return
}

func load(location string, data *userData) error {
	c, _ := os.UserConfigDir()
	location = c + location
	d, err := ioutil.ReadFile(location)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(d, data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if data.AppVersion == 0 {
		var a Api
		data.AppVersion = version
		err = json.Unmarshal(d, &a)
		if err != nil {
			return err
		}
		data.ApiData = a
		save(*data, saveLocation)
	}

	return nil
}
