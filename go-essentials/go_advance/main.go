package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Car struct {
	Manufacturer string
	Model        string
	Sales        float64
	VehicleType  string
	Price        float64
	Horsepower   int
}

type Metrics struct {
	TotalSales         float64
	SalesByRegion      map[string]float64
	SalesByType        map[string]float64
	PriceDistribution  map[string]int
	HorsepowerByRegion map[string][]int
	mutex              sync.RWMutex
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Open the CSV file
	file, err := os.Open("../csv/Car_sales.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}
	// Initialize metrics
	metrics := &Metrics{
		SalesByRegion:      make(map[string]float64),
		SalesByType:        make(map[string]float64),
		PriceDistribution:  make(map[string]int),
		HorsepowerByRegion: make(map[string][]int),
	}

	// Create processing pipeline
	rawRecords := readCSV(ctx, file)
	validRecords := validateRecords(ctx, rawRecords)
	processedCars := processRecords(ctx, validRecords, 3) // 3 workers

	// Start metrics aggregation
	metricsDone := make(chan struct{})
	go aggregateMetrics(ctx, processedCars, metrics, metricsDone)

	// Wait for completion
	select {
	case <-metricsDone:
		fmt.Println("Processing completed successfully")
	case <-ctx.Done():
		fmt.Println("Processing timed out")
	}

	// Print final metrics
	printMetrics(metrics)
}

func readCSV(ctx context.Context, file *os.File) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1 // Allow variable fields

		// Skip header
		if _, err := reader.Read(); err != nil {
			log.Printf("Error reading header: %v", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				record, err := reader.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("CSV read error: %v", err)
					continue
				}
				out <- record
			}
		}
	}()
	return out
}

func validateRecords(ctx context.Context, in <-chan []string) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)
		for record := range in {
			if len(record) < 15 {
				log.Printf("Invalid record: %v", record)
				continue
			}
			select {
			case out <- record:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func processRecords(ctx context.Context, in <-chan []string, workers int) <-chan *Car {
	out := make(chan *Car)
	var wg sync.WaitGroup

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for record := range in {
				car, err := parseCar(record)
				if err != nil {
					log.Printf("Parse error: %v", err)
					continue
				}
				select {
				case out <- car:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func aggregateMetrics(ctx context.Context, in <-chan *Car, metrics *Metrics, done chan<- struct{}) {
	defer close(done)

	for car := range in {
		metrics.mutex.Lock()

		// Update total sales
		metrics.TotalSales += car.Sales

		// Update regional sales
		metrics.SalesByRegion[car.Manufacturer] += car.Sales

		// Update vehicle type sales
		metrics.SalesByType[car.VehicleType] += car.Sales

		// Update price distribution
		priceRange := fmt.Sprintf("%d-%d", int(car.Price/10)*10, int(car.Price/10)*10+10)
		metrics.PriceDistribution[priceRange]++

		// Update horsepower data
		metrics.HorsepowerByRegion[car.Manufacturer] = append(
			metrics.HorsepowerByRegion[car.Manufacturer],
			car.Horsepower,
		)

		metrics.mutex.Unlock()
	}
}

func parseCar(record []string) (*Car, error) {
	car := &Car{
		Manufacturer: strings.TrimSpace(record[0]),
		Model:        strings.TrimSpace(record[1]),
		VehicleType:  strings.TrimSpace(record[4]),
	}

	var err error

	salesStr := strings.TrimSpace(record[2])
	if salesStr == "." {
		salesStr = "0.0"
	}
	if car.Sales, err = strconv.ParseFloat(salesStr, 64); err != nil {
		return nil, fmt.Errorf("invalid sales: %w", err)
	}

	priceStr := strings.TrimSpace(record[5])
	if priceStr == "." {
		priceStr = "0.0"
	}
	if car.Price, err = strconv.ParseFloat(priceStr, 64); err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}
	
    hrsStr:= strings.TrimSpace(record[7])
	if hrsStr == "." {
		hrsStr = "0"
	}
	if car.Horsepower, err = strconv.Atoi(hrsStr); err != nil {
		return nil, fmt.Errorf("invalid horsepower: %w", err)
	}

	return car, nil
}

func printMetrics(metrics *Metrics) {
	metrics.mutex.RLock()
	defer metrics.mutex.RUnlock()

	fmt.Printf("\n=== Total Sales: %.2fK ===\n", metrics.TotalSales)

	fmt.Println("\nRegional Sales:")
	for region, sales := range metrics.SalesByRegion {
		fmt.Printf("- %-15s: %.2fK\n", region, sales)
	}

	fmt.Println("\nVehicle Type Sales:")
	for vType, sales := range metrics.SalesByType {
		fmt.Printf("- %-15s: %.2fK\n", vType, sales)
	}

	fmt.Println("\nPrice Distribution ($K):")
	for priceRange, count := range metrics.PriceDistribution {
		fmt.Printf("- %-10s: %d cars\n", priceRange, count)
	}

	fmt.Println("\nAverage Horsepower by Region:")
	for region, hpValues := range metrics.HorsepowerByRegion {
		sum := 0
		for _, hp := range hpValues {
			sum += hp
		}
		avg := math.Round(float64(sum) / float64(len(hpValues)))
		fmt.Printf("- %-15s: %.0f hp\n", region, avg)
	}
}
