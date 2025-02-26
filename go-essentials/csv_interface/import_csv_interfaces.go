package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Car struct {
	Manufacturer        string
	Model               string
	SalesInThousands    float64
	FourYearResaleValue float64
	VehicleType         string
	PriceInThousands    float64
	EngineSize          float64
	Horsepower          int
	Wheelbase           float64
	Width               float64
	Length              float64
	CurbWeight          float64
	FuelCapacity        float64
	FuelEfficiency      int
	LatestLaunch        time.Time
}

// SalesTracker interface defines methods for tracking sales
type SalesTracker interface {
	AddSales(category string, amount int)
	GetSales() map[string]int
	Lock()
	Unlock()
}

// BaseSalesTracker provides common implementation
type BaseSalesTracker struct {
	sync.Mutex
	Sales map[string]int
}

// RegionSales implements SalesTracker for regions
type RegionSales struct {
	BaseSalesTracker
}

// VehicleTypeSales implements SalesTracker for vehicle types
type VehicleTypeSales struct {
	BaseSalesTracker
}

// NewSalesTracker creates a new sales tracker
func NewSalesTracker() *BaseSalesTracker {
	return &BaseSalesTracker{
		Sales: make(map[string]int),
	}
}

// AddSales adds sales for a category
func (st *BaseSalesTracker) AddSales(category string, amount int) {
	st.Lock()
	defer st.Unlock()
	st.Sales[category] += amount
}

// GetSales returns the sales map
func (st *BaseSalesTracker) GetSales() map[string]int {
	st.Lock()
	defer st.Unlock()
	salesCopy := make(map[string]int)
	for k, v := range st.Sales {
		salesCopy[k] = v
	}
	return salesCopy
}

func main() {
	// Open the CSV file
	file, err := os.Open("../csv/Car_sales.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use a buffered channel to prevent blocking
	carsChan := make(chan Car, 100)
	recordsChan := make(chan []string, 100)

	var wg sync.WaitGroup

	// Consumer goroutine for cars
	wg.Add(1)
	go func() {
		defer wg.Done()
		records := make([][]string, 0, 1000) // Pre-allocate capacity

		for car := range carsChan {
			record := []string{
				car.Manufacturer,
				car.Model,
				strconv.FormatFloat(car.SalesInThousands, 'f', -1, 64),
				strconv.FormatFloat(car.FourYearResaleValue, 'f', -1, 64),
				car.VehicleType,
				strconv.FormatFloat(car.PriceInThousands, 'f', -1, 64),
				strconv.FormatFloat(car.EngineSize, 'f', -1, 64),
				strconv.Itoa(car.Horsepower),
				strconv.FormatFloat(car.Wheelbase, 'f', -1, 64),
				strconv.FormatFloat(car.Width, 'f', -1, 64),
				strconv.FormatFloat(car.Length, 'f', -1, 64),
				strconv.FormatFloat(car.CurbWeight, 'f', -1, 64),
				strconv.FormatFloat(car.FuelCapacity, 'f', -1, 64),
				strconv.Itoa(car.FuelEfficiency),
				car.LatestLaunch.Format("2006-01-02"),
			}

			// Debug print to check SalesInThousands value
			fmt.Printf("SalesInThousands for %s %s: %f\n", car.Manufacturer, car.Model, car.SalesInThousands)

			jsonRecord, err := json.Marshal(record)
			if err != nil {
				log.Printf("Error marshalling record: %v", err)
			} else {
				fmt.Printf("Processing car: %s\n", jsonRecord)
			}
			records = append(records, record)
			recordsChan <- record
		}
		close(recordsChan)
		fmt.Println("\nAll cars processed")
	}()

	// Create a new CSV reader
	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	_, err = csvReader.Read() // Skip header
	if err != nil {
		log.Fatal(err)
	}

	// Create a new RegionSales and VehicleTypeSales
	regionSales := &RegionSales{BaseSalesTracker: *NewSalesTracker()}
	vehicleTypeSales := &VehicleTypeSales{BaseSalesTracker: *NewSalesTracker()}

	// Start worker goroutines for processing sales data
	const numWorkers = 4
	var processWg sync.WaitGroup
	processWg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer processWg.Done()
			for record := range recordsChan {
				ProcessCarSalesData([]string{record[0], record[4], record[2]}, regionSales, vehicleTypeSales)
			}
		}()
	}

	// Read and process CSV records
	go func() {
		defer close(carsChan)
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading record: %v", err)
				continue
			}

			car, err := parseCar(record)
			if err != nil {
				log.Printf("Error parsing record: %v", err)
				continue
			}

			carsChan <- car
		}
	}()

	// Wait for all processing to complete
	wg.Wait()
	processWg.Wait()

	// Print results
	printSalesResults(regionSales, vehicleTypeSales)
}

func ProcessCarSalesData(record []string, regionTracker, vehicleTracker SalesTracker) {
	if len(record) < 3 {
		log.Printf("Invalid record length: expected at least 3, got %d", len(record))
		return
	}

	salesInThousands, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		log.Printf("Warning: invalid sales value for record: %v, error: %v", record, err)
		return
	}

	sales := int(salesInThousands * 1000)
	manufacturer := record[0]
	vehicleType := record[1]

	regionTracker.AddSales(manufacturer, sales)
	vehicleTracker.AddSales(vehicleType, sales)
}

func printSalesResults(regionSales *RegionSales, vehicleTypeSales *VehicleTypeSales) {
	fmt.Println("\nTotal sales for each region:")
	regionSales.Lock()
	for region, sales := range regionSales.Sales {
		fmt.Printf("%s: %.2f thousand\n", region, float64(sales)/1000)
	}
	regionSales.Unlock()

	fmt.Println("\nTotal sales for each vehicle type:")
	vehicleTypeSales.Lock()
	for vehicleType, sales := range vehicleTypeSales.Sales {
		fmt.Printf("%s: %.2f thousand\n", vehicleType, float64(sales)/1000)
	}
	vehicleTypeSales.Unlock()
}

func parseCar(record []string) (Car, error) {
	var car Car
	if len(record) != 15 {
		return car, fmt.Errorf("invalid field count: %d", len(record))
	}

	car.Manufacturer = record[0]
	car.Model = record[1]
	car.SalesInThousands = parseFloat(record[2])
	car.FourYearResaleValue = parseFloat(record[3])
	car.VehicleType = record[4]
	car.PriceInThousands = parseFloat(record[5])
	car.EngineSize = parseFloat(record[6])
	car.Horsepower = parseInt(record[7])
	car.Wheelbase = parseFloat(record[8])
	car.Width = parseFloat(record[9])
	car.Length = parseFloat(record[10])
	car.CurbWeight = parseFloat(record[11])
	car.FuelCapacity = parseFloat(record[12])
	car.FuelEfficiency = parseInt(record[13])

	date, err := time.Parse("2-Jan-06", record[14])
	if err != nil {
		return car, fmt.Errorf("invalid date format: %v", err)
	}
	car.LatestLaunch = date

	return car, nil
}

func parseFloat(s string) float64 {
	if s == "." {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func parseInt(s string) int {
	if s == "." {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
