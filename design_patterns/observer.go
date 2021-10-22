package main

import "fmt"

type Observer interface {
	getId() string
	updateValue(string)
}

type Topic interface {
	register(observer Observer)
	broadcast()
}

type Item struct {
	observers []Observer
	name      string
	available bool
}

func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

func (i *Item) register(observer Observer) {
	i.observers = append(i.observers, observer)
}

func (i *Item) broadcast() {
	for _, observer := range i.observers {
		observer.updateValue(i.name)
	}
}

func (i *Item) UpdateAvailable() {
	fmt.Printf("Item %s is available\n", i.name)
	i.available = true
	i.broadcast()
}

type EmailClient struct {
	id string
}

func (ec *EmailClient) getId() string {
	return ec.id
}

func (ec *EmailClient) updateValue(value string) {
	fmt.Printf("Send Email - %s available from client %s\n", value, ec.id)
}

func NewEmail(id string) Observer {
	return &EmailClient{
		id: id,
	}
}

func main() {
	nvidiaItem := NewItem("RTX 3080")
	firstObserver := NewEmail("12ab")
	secondObserver := NewEmail("32ba")

	nvidiaItem.register(firstObserver)
	nvidiaItem.register(secondObserver)
	nvidiaItem.UpdateAvailable()
}
