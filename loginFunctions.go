package main

import "fmt"

func loginRec() Api {
	a := getApi(scanBase(), scanUserName(), scanPasswd())
	if a.Token == "" {
		a = loginRec()
	}
	return a
}

func login() {
	_state.ApiData = loginRec()
	save(_state, saveLocation)
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
