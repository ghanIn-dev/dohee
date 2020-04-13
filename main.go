package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/tealeg/xlsx"
)

var errFailedFileOpen = errors.New("Faild File Open")

var resultFileName string = "Result.xlsx"
var resultSheetName string = "결과"
var userIDButEventSheetName string = "회원에 있는데 이벤트에 없는 경우"
var eventButUserIDSheetName string = "이벤트에 있는데 회원에 없는 경우"

func main() {
	userFile := flag.String("userFile", "", "회원파일경로")
	eventFile := flag.String("eventFile", "", "이벤트파일경로")

	flag.Parse()

	if *userFile == "" {
		fmt.Println("userFile 파일주소입력해임마")
		return
	}

	if *eventFile == "" {
		fmt.Println("eventFile 파일주소입력해임마")
		return
	}

	dohee(*userFile, *eventFile)
}

func dohee(userFile string, eventFile string) {

	//파일 읽기
	users, err := xlsx.FileToSlice(userFile)
	checkError(err)
	events, err := xlsx.FileToSlice(eventFile)
	checkError(err)

	usersID := makeIDSlice(users, 5)
	eventsID := makeIDSlice(events, 1)

	//create a new xlsx file and write a struct
	//in a new row
	f := xlsx.NewFile()
	resultSheet, err := f.AddSheet(resultSheetName)
	checkError(err)

	userIDButEventSheet, err := f.AddSheet(userIDButEventSheetName)
	checkError(err)

	eventButUserIDSheet, err := f.AddSheet(eventButUserIDSheetName)
	checkError(err)

	for i, userID := range usersID {
		_, exist := Find(eventsID, userID)
		if exist {
			fmt.Println("회원에 있고, 이벤트에도 있는", usersID[i], i)
			row := resultSheet.AddRow()
			row.WriteSlice(&users[0][i], -1)
		} else {
			fmt.Println("회원에 있는데 이벤트에 없는 경우", usersID[i], i)
			row := userIDButEventSheet.AddRow()
			row.WriteSlice(&users[0][i], -1)
		}
	}

	for i, eventID := range usersID {
		_, exist := Find(eventsID, eventID)
		if !exist {
			fmt.Println("이벤트에 있는데 회원에 없는 경우", eventsID[i], i)
			row := eventButUserIDSheet.AddRow()
			row.WriteSlice(&events[0][i], -1)
		}
	}

	f.Save(resultFileName)
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func checkError(err error) {
	if err != nil {
		log.Fatal(errFailedFileOpen)
	}
}

func makeIDSlice(sheet [][][]string, idIndex int) (idSlice []string) {
	for _, rows := range sheet {
		for _, cols := range rows {
			idSlice = append(idSlice, cols[idIndex])
		}
	}
	return idSlice
}
