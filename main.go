package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

var (
	configFile = flag.String("configFile", "config.json", "Config JSON file name")
)

func main() {
	flag.Parse()
	cfg := &Config{}
	err := cfg.LoadFromJsonFile(*configFile)
	if err != nil {
		log.Fatalf("Error loading config file [%v]: %v", *configFile, err)
	}

	if len(cfg.Transformations) == 0 {
		log.Println("No transformation")
		return
	}
	transform(cfg)
}

func transform(cfg *Config) {
	rand.Seed(time.Now().UnixNano())
	// Prepare the list of all trace points
	points := []Point{cfg.StartPoint}
	x, y := cfg.StartPoint.X, cfg.StartPoint.Y
	maxX, maxY := x, y
	minX, minY := x, y

	for i := 0; i < cfg.LoopCount; i++ {
		r := rand.Float64()
		t := cfg.Transformations[len(cfg.Transformations)-1]
		for j := 0; j < len(cfg.Transformations); j++ {
			if cfg.cumProb[j] >= r {
				t = cfg.Transformations[j]
				break
			}
		}
		newX := t.A*x + t.B*y + t.E
		newY := t.C*x + t.D*y + t.F

		x, y = newX, newY
		maxX = max(maxX, x)
		maxY = max(maxY, y)
		minX = min(minX, x)
		minY = min(minY, y)

		points = append(points, Point{x, y})
	}

	log.Printf("Finish preparing points, x[%v;%v], y[%v;%v]", minX, maxX, minY, maxY)

	szX := (maxX - minX) * cfg.Draw.Scale
	szY := (maxY - minY) * cfg.Draw.Scale
	if szX > 2000 || szY > 2000 {
		log.Printf("Too large, can't draw: %vx%v", szX, szY)
		return
	}

	ctx := gg.NewContext(rnd(szX), rnd(szY))
	for _, p := range points {
		ctx.DrawPoint((p.X-minX)*cfg.Draw.Scale, (p.Y-minY)*cfg.Draw.Scale, cfg.Draw.PointSize)
	}
	ctx.SetRGB(cfg.Draw.PointColor[0], cfg.Draw.PointColor[1], cfg.Draw.PointColor[2])
	ctx.Fill()
	outputFile := fmt.Sprintf("%v.png", time.Now().UnixNano())
	err := ctx.SavePNG(outputFile)
	if err != nil {
		log.Fatal(err)
	}
}
