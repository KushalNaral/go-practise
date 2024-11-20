package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const FILE_CHANGE_SECOND = 5

// General user struct
type User struct {
	Id          string
	Name        string
	Email       string
	Contact     string
	Experiences []Experience
	Income      uint64
	Earnings    []Earning
	Notes       []Note
	JoinedAt    time.Time
	IsPresent   bool
}

// Note for user, can be created and attached to the user based on priority
type Note struct {
	CreatedAt time.Time
	CreatedBy string
	Note      string
	Priority  Priority
}

// Priority based on levels
type Priority struct {
	Name  string
	Level uint
}

// Work experience
type Experience struct {
	Title    string
	Duration int
}

// Details regarding salary and more ( has deduction and contributions will be made as receipt later)
type Earning struct {
	Date           string
	Amount         uint64
	Deductible     uint64
	Bonus          uint64
	SkillIncentive uint64
	NetPayable     uint64
	Contribution   Contribution
}

// Contribution represents various contributions or deductions.
type Contribution struct {
	ProvidentFund       uint64
	SST                 uint64
	RT                  uint64
	CIT                 uint64
	AttendanceDeduction uint64
	WelfareFund         uint64
}

func GenerateUniqueId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewUser() User {
	return User{
		Id:      GenerateUniqueId(),
		Name:    "John Doe",
		Email:   "john@doe.com",
		Contact: "98681002100",
		Experiences: []Experience{
			{
				Title:    "Intern At XYZ",
				Duration: 3,
			},
			{
				Title:    "Junior Engineer",
				Duration: 24,
			},
			{
				Title:    "Mid Level Analyst",
				Duration: 11,
			},
		},
		Income: 200000,
		Earnings: []Earning{
			{
				Date:           "2024-11-20",
				Amount:         47400,
				Deductible:     5610,
				Bonus:          0,
				SkillIncentive: 3000,
				NetPayable:     41790,
				Contribution: Contribution{
					ProvidentFund:       4800,
					SST:                 416,
					RT:                  293,
					CIT:                 0,
					AttendanceDeduction: 0,
					WelfareFund:         100,
				},
			},
		},
		Notes: []Note{
			{
				CreatedAt: time.Now(),
				CreatedBy: "Adam",
				Note:      "John has been early quite nicely",
				Priority: Priority{
					Name:  "moderate",
					Level: 3,
				},
			},
			{
				CreatedAt: time.Now(),
				CreatedBy: "HR",
				Note:      "Need to attend his leave counts and requests",
				Priority: Priority{
					Name:  "severe",
					Level: 5,
				},
			},
			{
				CreatedAt: time.Now(),
				CreatedBy: "Rubi",
				Note:      "Need to ask him for some tickets.",
				Priority: Priority{
					Name:  "minimal",
					Level: 1,
				},
			},
		},
		JoinedAt:  time.Now(),
		IsPresent: true,
	}
}

func CreateInitFile(user *User) {

	file, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	errr := os.WriteFile("./data.json", file, 0666)
	if errr != nil {
		log.Fatal(err)
	}
}

var wg sync.WaitGroup
var uC = make(chan User)

func thread1(user *User) {
	defer wg.Done()
	writeIntoFile(*user)
	count := 0
	for {
		count++
		fmt.Println(" id : \t ", count)

		user.IsPresent = !user.IsPresent
		fmt.Println(" id : \t ", count, " , status is ", user.IsPresent)
		writeIntoFile(*user)

		uC <- *user

		fmt.Println("waiting for 5 sec......")
		time.Sleep(time.Second * 5)
	}
}

func thread2() {
	for {
		select {
		case data, ok := <-uC:
			if !ok {
				fmt.Println("err occured in go routine")
			}

			fmt.Println(data)
			fmt.Println(reflect.TypeOf(data))
			writeIntoFile(data)
			ParseUserJson(data)
		}
	}
}

func main() {

	user := NewUser()
	wg.Add(1)

	go thread1(&user)
	go thread2()
	go Routes()

	wg.Wait()
}

func writeIntoFile(user User) {
	data, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("./test.json", []byte(data), 0666)

	file, err := os.OpenFile("./user_data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Err opeininig file:", err)
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(user)
	if err != nil {
		fmt.Println("err encoding json: ", err)
	}

}
