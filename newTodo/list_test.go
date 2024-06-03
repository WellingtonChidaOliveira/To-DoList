package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateNewTodoList(t *testing.T) {
	got, err := CreateNewTodoList("test")
	if err != nil {
		t.Fatalf("CreateNewTodoList() error = %v", err)
	}
	//defer os.Remove(got.Filename)

	want := TodoList{Filename: fmt.Sprintf("todo/test-%s.csv", time.Now().Format("2006-01-02")), Todos: []Todo{}}

	if got.Filename != want.Filename {
		t.Errorf("CreateNewTodoList() = %v, want %v", got.Filename, want.Filename)
	}

	if _, err := os.Stat(got.Filename); os.IsNotExist(err) {
		t.Errorf("CreateNewTodoList() failed to create a new file, %v", err)
	}
}

func TestGetTodoList(t *testing.T) {
	fileName := fmt.Sprintf("todo/test-%s.csv", time.Now().Format("2006-01-02"))
	todoList := TodoList{Filename: fileName, Todos: []Todo{}}
	got, err := GetTodoList(todoList.Filename)
	if err != nil {
		t.Fatalf("GetTodoList() error = %v", err)
	}

	if got.Filename != todoList.Filename {
		t.Errorf("GetTodoList() = %v, want %v", got.Filename, todoList.Filename)
	}
}

func TestAddTodo(t *testing.T) {
	fileName := fmt.Sprintf("todo/test-%s.csv", time.Now().Format("2006-01-02"))
	todoList := TodoList{Filename: fileName, Todos: []Todo{}}
	got, err := AddTodo(todoList.Filename, Todo{Title: "test", Completed: false})
	if err != nil {
		t.Fatalf("AddTodo() error = %v", err)
	}

	if got.Filename != todoList.Filename {
		t.Errorf("AddTodo() = %v, want %v", got.Filename, todoList.Filename)
	}

	if len(got.Todos) != 1 {
		t.Errorf("AddTodo() = %v, want %v", len(got.Todos), 1)
	}
}

func TestRemoveTodo(t *testing.T) {
	fileName := fmt.Sprintf("todo/test-%s.csv", time.Now().Format("2006-01-02"))
	todoList := TodoList{Filename: fileName, Todos: []Todo{}}
	todo := Todo{Title: "test", Completed: false}
	got, err := AddTodo(todoList.Filename, todo)
	if err != nil {
		t.Fatalf("AddTodo() error = %v", err)
	}

	if got.Filename != todoList.Filename {
		t.Errorf("AddTodo() = %v, want %v", got.Filename, todoList.Filename)
	}

	if len(got.Todos) != 1 {
		t.Errorf("AddTodo() = %v, want %v", len(got.Todos), 1)
	}

	got, err = RemoveTodo(todoList.Filename, todo)
	if err != nil {
		t.Fatalf("RemoveTodo() error = %v", err)
	}

	if len(got.Todos) != 0 {
		t.Errorf("RemoveTodo() = %v, want %v", len(got.Todos), 0)
	}
}
