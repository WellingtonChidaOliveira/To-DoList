package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Todo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoList struct {
	Filename string `json:"filename"`
	Todos    []Todo `json:"todos"`
}

var todos []Todo

func CreateNewTodoList(filename string) (TodoList, error) {
	// Create a new file with the given filename
	fileName := fmt.Sprintf("todo/%s-%s.csv", filename, time.Now().Format("2006-01-02"))
	file, err := os.Create(fileName)
	if err != nil {
		return TodoList{}, err
	}
	defer file.Close()

	_, err = file.WriteString("Title,Completed\n")
	if err != nil {
		return TodoList{}, err
	}

	// Return a new TodoList with the filename and an empty list of todos
	return TodoList{Filename: fileName, Todos: []Todo{}}, nil
}

func GetTodoList(filename string) (TodoList, error) {

    file, err := os.Open(filename)
    if err != nil {
        return TodoList{}, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    // Skip the header line
    _, err = reader.Read()
    if err != nil {
        return TodoList{}, err
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return TodoList{}, err
        }

        completed := false
        if strings.ToLower(record[1]) == "true" {
            completed = true
        }

        todos = append(todos, Todo{Title: record[0], Completed: completed})
    }

    // Return a new TodoList with the filename and the todos
    return TodoList{Filename: filename, Todos: todos}, nil
}

func AddTodo(filename string, todo Todo) (TodoList, error) {
	path := fmt.Sprintf("todo/%s", filename)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return TodoList{}, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{todo.Title, fmt.Sprintf("%t", todo.Completed)})
	if err != nil {
		return TodoList{}, err
	}

	// Return a new TodoList with the filename and the todos
	return TodoList{Filename: filename, Todos: append(todos, todo)}, nil
}

func RemoveTodo(filename string, todo Todo) (TodoList, error) {
	todos, err := GetTodoList(filename)
	if err != nil {
		return TodoList{}, err
	}

	var newTodos []Todo
	for _, t := range todos.Todos {
		if t.Title != todo.Title {
			newTodos = append(newTodos, t)
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return TodoList{}, err
	}
	defer file.Close()

	_, err = file.WriteString("Title,Completed\n")
	if err != nil {
		return TodoList{}, err
	}


	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, t := range newTodos {
		err = writer.Write([]string{t.Title, fmt.Sprintf("%t", t.Completed)})
		if err != nil {
			return TodoList{}, err
		}
	}
	
	return TodoList{Filename: filename, Todos: newTodos}, nil
}


func main() {
	action := flag.String("action", "", "The action to perform (create, get, add, remove)")
	filename := flag.String("filename", "", "The base filename of the TodoList")
	title := flag.String("title", "", "The title of the Todo")
	completed := flag.Bool("completed", false, "Whether the Todo is completed")

	flag.Parse()

	// Append the date to the filename
	fullFilename := fmt.Sprintf("%s-%s.csv", *filename, time.Now().Format("2006-01-02"))

	switch *action {
	case "create":
		todoList, err := CreateNewTodoList(*filename)
		if err != nil {
			fmt.Println("Failed to create a new TodoList")
			return
		}
		fmt.Println("Created a new TodoList:", todoList.Filename)

	case "get":
		todoList, err := GetTodoList(fullFilename)
		if err != nil {
			fmt.Println("Failed to get the TodoList")
			return
		}
		fmt.Println("Got the TodoList:", todoList)

	case "add":
		_, err := AddTodo(fullFilename, Todo{Title: *title, Completed: *completed})
		if err != nil {
			fmt.Println("Failed to add a new Todo")
			return
		}
		fmt.Println("Added a new Todo:", *title)

	case "remove":
		_, err := RemoveTodo(fullFilename, Todo{Title: *title, Completed: *completed})
		if err != nil {
			fmt.Println("Failed to remove the Todo")
			return
		}
		fmt.Println("Removed the Todo:", *title)

	default:
		fmt.Println("Unknown action:", *action)
	}
}