package main
import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
    "time"
    "log"
    "github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
var (
    todos    []Todo
    todoID   int
    todoLock sync.Mutex
)
var db *sql.DB
func initDB() {
	var err error

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get database credentials from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort. dbName)

    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    go func() {
        time.Sleep(2 * time.Second)
        //todoLock.Lock()
        //defer todoLock.Unlock()
		//todo.ID = todoID
		//todoID++
		//todos = append(todos, todo)
		result, err := db.Exec("INSERT INTO todo (title, completed) VALUES (?, ?)", todo.Title, todo.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todo.ID = int(id)
		
       // todo.Completed = true
    }()
    //todoLock.Lock()
    //todoLock.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	todosChan := make(chan []Todo)
    errChan := make(chan error)

    go func() {
        rows, err := db.Query("SELECT id, title, completed FROM todo")
        if err != nil {
            errChan <- err
            return
        }
        defer rows.Close()

        var todos []Todo
        for rows.Next() {
            var todo Todo
            if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
                errChan <- err
                return
            }
            todos = append(todos, todo)
        }

        todosChan <- todos
    }()

    select {
    case todos := <-todosChan:
        json.NewEncoder(w).Encode(todos)
    case err := <-errChan:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    // todoLock.Lock()
    // defer todoLock.Unlock()


	// rows, err := db.Query("SELECT id, title, completed FROM todo")
    // if err != nil {
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }
    // defer rows.Close()

    // var todos []Todo
    // for rows.Next() {
    //     var todo Todo
    //     if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
    //         http.Error(w, err.Error(), http.StatusInternalServerError)
    //         return
    //     }
    //     todos = append(todos, todo)
    // }

    // json.NewEncoder(w).Encode(todos)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    go func() {
        _, err := db.Exec("UPDATE todo SET title = ?, completed = ? WHERE id = ?", todo.Title, todo.Completed, id)
        if err != nil {
            log.Printf("Error updating todo: %v", err)
            return
        }

        log.Printf("Updated todo with ID: %d", id)
    }()
    // id, _ := strconv.Atoi(mux.Vars(r)["id"])
    // var updatedTodo Todo
    // if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
    //     http.Error(w, err.Error(), http.StatusBadRequest)
    //     return
    // }

    // todoLock.Lock()
    // defer todoLock.Unlock()
    // for i, todo := range todos {
    //     if todo.ID == id {
    //         todos[i] = updatedTodo
    //         todos[i].ID = id
    //         json.NewEncoder(w).Encode(todos[i])
    //         return
    //     }
    // }
    // http.Error(w, "Todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
    //id, _ := strconv.Atoi(mux.Vars(r)["id"])
	vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    go func() {
        _, err := db.Exec("DELETE FROM todo WHERE id = ?", id)
        if err != nil {
            log.Printf("Error deleting todo: %v", err)
            return
        }

        log.Printf("Deleted todo with ID: %d", id)
    }() 
	w.WriteHeader(http.StatusNoContent)

    // todoLock.Lock()
    // defer todoLock.Unlock()
    // for i, todo := range todos {
    //     if todo.ID == id {
    //         todos = append(todos[:i], todos[i+1:]...)
    //         w.WriteHeader(http.StatusNoContent)
    //         return
    //     }
    // }
    // http.Error(w, "Todo not found", http.StatusNotFound)
}

func testConn(){
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:33060)/todo_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
:   "root",          // MySQL username
        Passwd: "1234",      // MySQL password
        Net:    "tcp",           // Network type
        Addr:   "127.0.0.1:33060", // MySQL server address
        DBName: "todo_db",       // Database name
    }

    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to the database!")
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/todos", createTodo).Methods("POST")
    r.HandleFunc("/todos", getTodos).Methods("GET")
    r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
    r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

    http.Handle("/", r)
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}