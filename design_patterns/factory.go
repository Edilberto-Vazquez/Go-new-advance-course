package main

import "fmt"

type IProduct interface {
	setStock(stock int)
	getStock() int
	setName(name string)
	getName() string
}

type Computer struct {
	name  string
	stock int
}

func (c *Computer) setStock(stock int) {
	c.stock = stock
}

func (c *Computer) setName(name string) {
	c.name = name
}

func (c *Computer) getStock() int {
	return c.stock
}

func (c *Computer) getName() string {
	return c.name
}

type Laptop struct {
	Computer
}

func newLaptop() IProduct {
	return &Laptop{
		Computer: Computer{
			name:  "HP",
			stock: 25,
		},
	}
}

type Desktop struct {
	Computer
}

func (d *Desktop) MultiplyStock(n int) {
	d.stock = n
}

func newDesktop() IProduct {
	return &Desktop{
		Computer: Computer{
			name:  "DELL",
			stock: 25,
		},
	}
}

func GetComputerFactory(computerType string) (IProduct, error) {
	if computerType == "Laptop" {
		return newLaptop(), nil
	}
	if computerType == "Desktop" {
		return newDesktop(), nil
	}
	return nil, fmt.Errorf("invalid computer type")
}

func printNameAndStock(p IProduct) {
	fmt.Printf("Product name: %s, with stock %d\n", p.getName(), p.getStock())
}

func MultiplyStock(p IProduct, n int) {
	p.setStock(p.getStock() * 20)
}

func main() {
	laptop, _ := GetComputerFactory("Laptop")
	desktop, _ := GetComputerFactory("Desktop")

	printNameAndStock(laptop)
	printNameAndStock(desktop)

	MultiplyStock(laptop, 20)

	printNameAndStock(laptop)
	printNameAndStock(desktop)
}
