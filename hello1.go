package main

import (
	"math"
	"fmt"
)

func main() {
	rect := Rectangle{height:2, width:4}
	circ := Circle{2}

	fmt.Println("Area of rectangle = " , getArea(rect), " and area of circle= " , getArea(circ))
}

type Shape interface{
	area() float64
}

type Rectangle struct{
	height float64
	width  float64
}

type Circle struct{
	radius float64
}

func (r Rectangle) area() float64{
	return r.height * r.width
}

func (c Circle) area() float64{
	return math.Pi * math.Pow(c.radius, 2)
}

func getArea(shape Shape) float64{
	return shape.area()
}
