package main

import (
	"fmt"
	"sync"
)

// structura de la funcion que se almacenara
type Function func(key int) (interface{}, error)

// structura que tendran los datos ya calculados
type FunctionResult struct {
	value interface{}
	err   error
}

// structura para guardar en memoria los datos almacenados
type Memory struct {
	f          Function                       // campo para almacenar la funcion
	cache      map[interface{}]FunctionResult // campo para almacenar los datos que ya se calcularon previamente
	InProgress map[interface{}]bool
	IsPending  map[interface{}][]chan FunctionResult
	lock       sync.RWMutex
}

// funcion que crea un nuevo sistema de cache que apunta a Memory
/*
	parametros: f: funcion a evaluar
	return: regresa un tipo Memory
*/
func NewCache(f Function) *Memory {
	return &Memory{
		f:          f,
		cache:      make(map[interface{}]FunctionResult),
		InProgress: make(map[interface{}]bool),
		IsPending:  make(map[interface{}][]chan FunctionResult),
	}
}

// funcion que regresa la funcion fibonacci que se pasa
// como parametro a la funcion NewCache
/*
	parametros: el numero a calcular
	return: un generico y un error
*/
func GetFibonnaci(n int) (interface{}, error) {
	// time.Sleep(1 * time.Second)
	return Fibonacci(n), nil
}

// funcion fibonacci
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// función que comprueba si otro proceso está calculando el mismo valor
/*
	parametros: el numero a calcular
	return: un generico y un error
*/
func (m *Memory) Work(key int) (interface{}, error) {
	m.lock.RLock()
	exists := m.InProgress[key]
	m.lock.RUnlock()
	if exists {
		response := make(chan FunctionResult)
		defer close(response)
		m.lock.Lock()
		m.IsPending[key] = append(m.IsPending[key], response)
		m.lock.Unlock()
		fmt.Printf("Waiting for response: %d\n", key)
		res := <-response
		fmt.Printf("Response done, received %v\n", res)
		return res.value, res.err
	}

	m.lock.Lock()
	m.InProgress[key] = true
	m.lock.Unlock()

	fmt.Printf("Calculate Expensive Fibonnaci to %d\n", key)
	result := FunctionResult{}
	result.value, result.err = m.f(key)
	if result.err != nil {
		return result.value, result.err
	}

	m.lock.RLock()
	pendingWorkers, exists := m.IsPending[key]
	m.lock.RUnlock()

	if exists {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result
		}
		fmt.Printf("Result sent - all pending workers ready job: %d\n", key)
	}
	m.lock.Lock()
	m.InProgress[key] = false
	m.IsPending[key] = make([]chan FunctionResult, 0)
	m.lock.Unlock()

	return result.value, result.err
}

// funcion que evalua si el valor ya fue calculado
/*
	parametros: el numero a evaluar
	return: un generico y un error
*/
func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.RLock()
	result, exist := m.cache[key]
	m.lock.RUnlock()
	if !exist {
		result.value, result.err = m.Work(key)
		m.lock.Lock()
		m.cache[key] = result
		m.lock.Unlock()
	}
	return result.value, result.err
}

func main() {
	memory := NewCache(GetFibonnaci)
	// cache := NewCache(GetFibonnaci)
	values := []int{45, 42, 41, 51, 42, 46, 47, 39, 45, 51, 50, 51, 42, 51, 41, 47}
	var wg sync.WaitGroup
	for _, v := range values {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			result, err := memory.Get(value)
			if err != nil {
				fmt.Printf("Number %d could not be calculated: Error: %v", value, err)
			} else {
				fmt.Printf("Result for value: %d is %v \n", value, result)
			}
		}(v)
	}
	wg.Wait()

	fmt.Println("-----Cheking cache-----")

	for _, v := range values {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			result, err := memory.Get(value)
			if err != nil {
				fmt.Printf("Number %d could not be calculated: Error: %v", value, err)
			} else {
				fmt.Printf("Result for value: %d is %v \n", value, result)
			}
		}(v)
	}
	wg.Wait()
}
