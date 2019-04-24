package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	data, err := read("data.txt")
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}

	err = plotData("out.png", data)
	if err != nil {
		log.Fatalf("Could not plot data: %v", err)
	}
}

type xy struct {
	X, Y float64
}

func read(filename string) (plotter.XYs, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data plotter.XYs

	s := bufio.NewScanner(f)
	for s.Scan() {
		var x, y float64
		_, err := fmt.Sscanf(s.Text(), "%f,%f", &x, &y)
		if err != nil {
			log.Printf("Discarding data point: %q: %v", s.Text(), err)
			continue
		}
		data = append(data, struct{ X, Y float64 }{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("Could not scan: %v", err)
	}
	return data, nil
}

func plotData(path string, d plotter.XYs) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Could not create %s: %v", path, err)
	}

	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("Could not create plot: %v", err)
	}
	s, err := plotter.NewScatter(d)

	if err != nil {
		return fmt.Errorf("Could not create scatter: %v", err)
	}

	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	x, c := linearRegression(d)
	// fake linear regression resutl
	l, err := plotter.NewLine(plotter.XYs{
		{0, c}, {20, 20*x + c},
	})
	if err != nil {
		return fmt.Errorf("Coult not create line: %v", err)
	}

	p.Add(l)

	w, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("Could not create writer: %v", err)
	}

	_, err = w.WriteTo(f)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v", path, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Could not close %s: %v", path, err)
	}
	return nil
}

func linearRegression(d plotter.XYs) (m, c float64) {
	const (
		min   = -100.0
		max   = 100.0
		delta = 0.1
	)
	minCost := math.MaxFloat64
	for im := min; im < max; im += delta {
		for ic := min; ic < max; ic += delta {
			cost := computeCost(d, im, ic)
			if cost < minCost {
				minCost = cost
				m, c = im, ic
			}
		}
	}

	return m, c
}

func computeCost(xys plotter.XYs, m, c float64) float64 {
	// 1/N * sum((y - (m*x+c))^2)
	s := 0.0
	for _, xy := range xys {
		d := xy.Y - (xy.X*m + c)
		s += d * d
	}
	return s / float64(len(xys))
}
