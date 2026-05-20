package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/common-nighthawk/go-figure"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Task struct {
	ID       string `json:"uuid"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

type Todos struct {
	Tasks []Task `json:"tasks"`
}

var (
	status []string
)

func main() {
	figure.NewColorFigure("ToDo", "", "Cyan", true).Print()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(color.Green + "\n ================ TO-DO ================ " + color.Reset)
		fmt.Println("\nOptions : \n 1.add Tasks \n 2.view todo lists \n 3.list of complete tasks \n 4.delete tasks \n 5.Exit")
		fmt.Print("\nchoose an Option : ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			addTasks(scanner)
		case "2":
			lists()
		case "3":
			compeletedTasks()
		case "4":
			deleteTasks()
		case "5":
			os.Exit(0)
		default:
			fmt.Println(color.Red + "Invalid option , choose between the options." + color.Red)
			time.Sleep(1 * time.Second)
		}

	}
}

func GenerateHash(name string) string {
	if name == "" {
		return "unknown task name"
	}

	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])[:10]
}

func addTasks(scanner *bufio.Scanner) {
	figure.NewColorFigure("Add Task", "", "Purple", true).Print()
	fmt.Println(color.Yellow + "\nprocessing ..." + color.Reset)
	time.Sleep(1 * time.Second)

	status = []string{"In Progress", "Completed"}

	tskName := readRequireInfo(*scanner, "Task Name : ")

	tskDate := readRequireInfo(*scanner, "Task Date : ")

	tskPriority := readRequireInfo(*scanner, "Task Priority [high , moderate , low]: ")

	// adding uuid to the task
	uuid := GenerateHash(tskName)

	fmt.Print("save task ? (y/n) ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	time.Sleep(1 * time.Second)

	task, err := loadTasks()
	if err != nil {
		fmt.Println(color.Yellow + "no task found for the first time!" + color.Reset)
	}

	time.Sleep(1 * time.Second)

	switch choice {
	case "y", "", "yes", "Y":
		tasks := Task{
			ID:       uuid,
			Name:     tskName,
			Date:     tskDate,
			Priority: tskPriority,
			Status:   status[0],
		}

		task.Tasks = append(task.Tasks, tasks)
		// call saveTask function for saving the tasks
		saveTasks(task)
	case "n", "no", "No", "NO", "0":
		fmt.Println(color.Yellow + "Canceled ..." + color.Reset)
		time.Sleep(1 * time.Second)
		main()
	}
}

func loadTasks() (Todos, error) {
	var readTask Todos
	binData, err := os.ReadFile(GetCurrentPath() + "data/tasks.json")
	if err != nil {
		fmt.Println("No data to read , failed")
		return Todos{}, err
	}

	if err := json.Unmarshal(binData, &readTask); err != nil {
		fmt.Println(color.Red + "Failed to decode the json data ." + color.Reset)
		return Todos{}, err
	}
	return readTask, nil
}

// readRequired function => give scanner and label of the I/O operation
//
// this is godoc comment for this function 
func readRequireInfo(scanner bufio.Scanner, label string) string {
	for {
		fmt.Printf("%s\n\u25CF %s%s", color.Cyan, label, color.Reset)
		scanner.Scan()
		ans := strings.TrimSpace(scanner.Text())

		if strings.Contains(label, "Priority") {
			switch ans {
			case "high":
				return ans
			case "moderate":
				return ans
			case "low":
				return ans
			default:
				fmt.Println(color.Red + "Invalid choice, choose between [high , moderate , low]." + color.Reset)
			}
		} else if strings.Contains(label, "Date") {

			result, err := checkDateValication(ans)
			if err != nil {
				fmt.Println(color.Red, err, color.Reset)
				time.Sleep(1 * time.Second)
				continue
			}
			return result
		} else {
			return ans
		}

	}
}

func checkDateValication(usrDate string) (string, error) {
	pattern := `^\d{4}/\d{2}/\d{2}`

	regex := regexp.MustCompile(pattern)

	for {
		valid := regex.MatchString(usrDate)
		fmt.Println(color.Yellow + "checking date format ..." + color.Reset)
		time.Sleep(1 * time.Second)

		if !valid {
			fmt.Printf("%s%s → %v%s,", color.Red, usrDate, valid, color.Reset)
			return "", errors.New("invalid date Format.")
		} else {
			fmt.Println(color.Green + "\u2713 Format date Accepted!" + color.Reset)
			return usrDate, nil
		}
	}

}

func saveTasks(task Todos) {
	data := Todos{Tasks: task.Tasks}
	jsonbytes, _ := json.MarshalIndent(data, "", " ")

	// currentDir, _ := os.Getwd()
	checkDataDir, err := checkingDataDir()
	if err != nil {
		fmt.Println(color.Red + checkDataDir + color.Reset)
	} else {
		fmt.Println(color.Green + checkDataDir + color.Reset)
	}

	chmodCmd := exec.Command("chmod", "+x", GetCurrentPath()+"/data")
	if err := chmodCmd.Run(); err != nil {
		fmt.Println(err)
	}
	// writing file
	if err := os.WriteFile(GetCurrentPath()+"/data/"+"tasks.json", jsonbytes, 0644); err != nil {
		panic(err)
	} else {
		fmt.Println(color.Green + "task Successfully created." + color.Reset)
		main()
	}
}

func checkingDataDir() (string, error) {
	if err := os.MkdirAll(GetCurrentPath()+"/data", 0755); err != nil {
		return "Error creating directory", err
	} else {
		return "data directory created!", nil
	}
}

func GetCurrentPath() string {
	currentPath, _ := os.Getwd()
	return currentPath + "/"
}

func lists() {

}

func compeletedTasks() {

}

func deleteTasks() {

}
