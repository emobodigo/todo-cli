package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

type Todo struct {
	Id        int
	Task      string
	CreatedAt time.Time
	Done      bool
}

func Add(task string) {
	file, err := openFile("todo.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer closeFile(file)

	list, err := loadTodos(file)

	if err != nil {
		fmt.Println(err)
	}
	currentId := len(list)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	newTodo := Todo{
		Id:        currentId + 1,
		Task:      task,
		CreatedAt: time.Now().In(loc),
		Done:      false,
	}
	list = append(list, newTodo)
	err = saveTodo(list)
	if err != nil {
		fmt.Println(err)
	}
}

func List(isAll bool) {
	file, err := openFile("todo.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer closeFile(file)

	list, err := loadTodos(file)
	if err != nil {
		fmt.Println(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 10, 10, ' ', tabwriter.TabIndent)
	_, err = fmt.Fprintf(w, "ID\tTask\tCreated\tDone\n")

	for _, todo := range list {
		if isAll == false && todo.Done {
			continue
		}
		tDiff := timediff.TimeDiff(todo.CreatedAt)

		if err != nil {
			fmt.Println(err)
		}
		_, _ = fmt.Fprintf(w, "%d\t%s\t%s\t%t\n",
			todo.Id, todo.Task, tDiff, todo.Done)
	}
	_ = w.Flush()

}

func Complete(id string) {
	file, err := openFile("todo.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer closeFile(file)

	list, err := loadTodos(file)
	if err != nil {
		fmt.Println(err)
	}
	for i := range list {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		}
		if list[i].Id == idInt {
			list[i].Done = true
		}
	}
	err = saveTodo(list)
	if err != nil {
		fmt.Println(err)
	}
}

func Delete(id string) {
	file, err := openFile("todo.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer closeFile(file)

	list, err := loadTodos(file)
	if err != nil {
		fmt.Println(err)
	}
	for i := range list {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		}
		if list[i].Id == idInt {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	err = saveTodo(list)
	if err != nil {
		fmt.Println(err)
	}
}

func loadTodos(file *os.File) ([]Todo, error) {
	reader := csv.NewReader(file)

	_, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return []Todo{}, nil
		}
		return nil, err
	}

	var todos []Todo

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			return nil, err
		}
		createdAt, err := time.ParseInLocation("2006-01-02 15:04:05", record[2], loc)
		if err != nil {
			return nil, err
		}

		isDone, err := strconv.ParseBool(record[3])
		if err != nil {
			return nil, err
		}

		todo := Todo{
			Id:        id,
			Task:      record[1],
			CreatedAt: createdAt,
			Done:      isDone,
		}

		todos = append(todos, todo)
	}
	return todos, nil
}

func saveTodo(todos []Todo) error {
	file, err := os.Create("todo.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"No", "Task", "CreatedAt", "Done"})
	if err != nil {
		return err
	}

	for _, todo := range todos {
		record := []string{
			strconv.Itoa(todo.Id),
			todo.Task,
			todo.CreatedAt.Format("2006-01-02 15:04:05"),
			strconv.FormatBool(todo.Done),
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}
	return writer.Error()
}

func openFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}
	return f.Close()
}
