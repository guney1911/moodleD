package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	boolptr := flag.Bool("login", false, "-login : start the login wizard")
	refresh := flag.Int("refresh", 1, "-refresh=xx : refresh every xx minutes. Default is 360min/6h.")
	flag.Parse()
	if *boolptr {
		login()
		return
	}
	api := &Api{}
	c, _ := os.UserConfigDir()
	err := load(c+saveLocation, api)

	if err != nil {
		fmt.Println("Couldn't read save file!")
		beeep.Alert("MoodleD", "moodleD couldn't read or find save file! \n Please start the program from a terminal with the -login option to login to Moodle.", "")
	}
	createNotifyThreads(api.getEvents(time.Now()))
	//createNotifyThreads(api.getEvents(time.Date(2021,1,11,0,0,0,0,time.Local)))
	go checkNew(api, *refresh)
	beeep.Notify("MoodleD", "moodleD started", "")
	select {}

}

func loginRec() *Api {
	data := getApi(scanBase(), scanUserName(), scanPasswd())
	if data.Token == "" {
		data = loginRec()
	}
	return data
}
func login() {
	data := loginRec()
	c, _ := os.UserConfigDir()
	save(*data, c+saveLocation)
	fmt.Println("Thank you for your cooperation. You may start MoodleD now without the login option.")
}

func scanBase() string {
	fmt.Print("Please enter the Moodle server URL! It should look like 'https://moodle.example.com/moodle' \nURL:  ")
	var base string
	_, err := fmt.Scanln(&base)
	if err != nil {
		base = scanBase()
	}
	return base
}

func scanUserName() string {
	fmt.Print("Please enter your Moodle username! \nUsername:  ")
	var userName string
	_, err := fmt.Scan(&userName)
	if err != nil {
		userName = scanUserName()
	}
	return userName
}

func scanPasswd() string {
	fmt.Print("Please enter your Moodle password! It will not be saved locally but transmitted to the moodle server for an api token. \nPassword:  ")
	var passwd string
	_, err := fmt.Scan(&passwd)
	if err != nil {
		passwd = scanPasswd()
	}
	return passwd
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
		time.Sleep(time.Duration(refresh) * time.Minute)
		createNotifyThreads(api.getEvents(time.Now()))
		beeep.Notify("MoodleD", "Refreshed", "")
	}
}
func save(api Api, location string) {
	d, err := json.Marshal(api)
	assertErr(err)
	_, err = os.Create(location)
	assertErr(err)
	err = ioutil.WriteFile(location, d, 0644)
	assertErr(err)
	return
}

func load(location string, data *Api) error {
	d, err := ioutil.ReadFile(location)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(d))
	err = json.Unmarshal(d, data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
