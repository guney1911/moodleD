package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"time"
)

const saveLocation = "/moodleD"
const version = 1

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
type courseData struct {
	Total   int      `json:"total"`
	Courses []course `json:"courses"`
}
type course struct {
	ID        int    `json:"id"`
	Fullname  string `json:"fullname"`
	Shortname string `json:"shortname"`
}
type courseContent []struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Modules []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"modules"`
}
type eventData struct {
	description string
	throw       time.Time
	ID          int
}

type userData struct {
	AppVersion     int             `json:"app_version"`
	ApiData        Api             `json:"api_data"`
	Courses        courseData      `json:"courses"`
	ContentIDStore map[string]bool `json:"content_id_store"`
}

type Api struct {
	Token        string `json:"token"`
	PrivateToken string `json:"privatetoken"`
	Base         string `json:"base"`
}

func (d eventData) notify() {
	d.throw = d.throw.Add(-5 * time.Minute)
	fmt.Printf("created thread for event %v, will notify at %v:%v\n", d.ID, d.throw.Hour(), d.throw.Minute())
	dur := d.throw.Sub(time.Now())
	time.Sleep(dur)
	beeep.Notify("MoodleD", d.description, "")
}
