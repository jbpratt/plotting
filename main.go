package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gonum.org/v1/plot"
)

func main() {
	data, err := read("data.txt")
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}

	_ = data
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %v", err)
	}
	w, err := p.WriterTo(512, 512, "png")
	if err != nil {
		log.Fatalf("Could not create writer: %v", err)
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Fatalf("Could not create out file: %v", err)
	}

	_, err = w.WriteTo(f)
	if err != nil {
		log.Fatalf("Could not write to out file: %v", err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("Could not close out file: %v", err)
	}
}

type xy struct {
	x, y float64
}

func read(filename string) ([]xy, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []xy

	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("Discarding data point: %q: %v", s.Text(), err)
		}
		data = append(data, xy{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}
	return data, nil
}
