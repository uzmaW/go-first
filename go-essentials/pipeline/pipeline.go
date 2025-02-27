package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Data represents the data being processed
type Data struct {
    ID   int
    Name string
}

// ProcessData processes the data
func ProcessData(data Data, wg *sync.WaitGroup) {
    defer wg.Done()
    // Simulate some processing time
    time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
    fmt.Printf("Processed data: %+v\n", data)
}

// FanOutFanIn distributes the data to multiple workers
func FanOutFanIn(ctx context.Context, data []Data) {
    var wg sync.WaitGroup
    // Create a worker pool
    workerPool := make(chan Data, len(data))
    for _, d := range data {
        workerPool <- d
    }
    close(workerPool)
    // Start the workers
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            for d := range workerPool {
                ProcessData(d, &wg)
            }
        }()
    }
    // Wait for the workers to complete
    wg.Wait()
}

func main() {
    ctx := context.Background()
    // Generate some sample data
    data := []Data{
        {ID: 1, Name: "John"},
        {ID: 2, Name: "Jane"},
        {ID: 3, Name: "Bob"},
        {ID: 4, Name: "Alice"},
        {ID: 5, Name: "Mike"},
    }
    // Start the fan-out/fan-in pipeline
    FanOutFanIn(ctx, data)
    fmt.Println("All data processed")
}
