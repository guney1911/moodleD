package main

import (
	"fmt"
	"log"
)

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
func generateModuleId(sectionID, moduleID int) string {
	return fmt.Sprintf("%v|%v", sectionID, moduleID)
}
