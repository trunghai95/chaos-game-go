package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

/**
A transformation is in form of
	x_(i+1) = a*x_i + b*y_i + e
	y_(i+1) = c*x_i + d*y_i + f
and is chosen with a probability of p
(All probabilities will be standardized to range [0-1])
*/
type Transformation struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
	D float64 `json:"d"`
	E float64 `json:"e"`
	F float64 `json:"f"`
	P float64 `json:"p"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DrawConfig struct {
	PointSize  float64   `json:"point_size"`
	PointColor []float64 `json:"point_color"`
	Scale      float64   `json:"scale"`
}

type Config struct {
	Transformations []Transformation `json:"transformations"`
	StartPoint      Point            `json:"start_point"`
	LoopCount       int              `json:"loop_count"`
	Draw            DrawConfig       `json:"draw"`

	// Cumulative probability, range [0-1]
	cumProb []float64
}

func (c *Config) LoadFromJsonFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}

	c.afterLoad()
	return nil
}

func (c *Config) afterLoad() {
	if c.Draw.PointSize <= 0 {
		c.Draw.PointSize = 1.0
	}
	for i := range c.Draw.PointColor {
		if c.Draw.PointColor[i] < 0 {
			c.Draw.PointColor[i] = 0
		} else if c.Draw.PointColor[i] > 1 {
			c.Draw.PointColor[i] = 1
		}
	}
	if len(c.Draw.PointColor) < 3 {
		log.Println("config point color is invalid")
		c.Draw.PointColor = []float64{0, 0, 0}
	}
	if c.Draw.Scale <= 0 {
		c.Draw.Scale = 1
	}

	sumP := float64(0)
	for _, t := range c.Transformations {
		sumP += t.P
	}
	c.cumProb = make([]float64, len(c.Transformations))
	last := float64(0)
	for i, t := range c.Transformations {
		t.P = t.P / sumP
		last = last + t.P
		c.cumProb[i] = last
	}
}

func (c *Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
