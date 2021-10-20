package main

import (
	"fmt"
	"log"
	"sync"
	"time"
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
	f     Function               // campo para almacenar la funcion
	cache map[int]FunctionResult // campo para almacenar los datos que ya se calcularon previamente
	lock  sync.Mutex
}

// funcion que crea un nuevo sistema de cache que apunta a Memory
/*
	parametros: f: funcion a evaluar
	return: regresa un tipo Memory
*/
func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FunctionResult),
	}
}

// funcion que evalua si el valor ya fue calculado
/*
	parametros: el numero a evaluar
	return: un generico y un error
*/
func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.Lock()
	result, exist := m.cache[key]
	m.lock.Unlock()
	if !exist {
		m.lock.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		m.lock.Unlock()
	}
	return result.value, result.err
}

// funcion que regresa la funcion fibonacci que se pasa
// como parametro a la funcion NewCache
/*
	parametros: el numero a calcular
	return: un generico y un error
*/
func GetFibonnaci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

// funcion fibonacci
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func main() {
	cache := NewCache(GetFibonnaci)
	fibo := []int{15, 23, 31, 7, 36, 15, 43, 8, 17, 23}
	var wg sync.WaitGroup
	for _, v := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			start := time.Now()
			value, err := cache.Get(index)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%d, %s, %d\n", index, time.Since(start), value)
		}(v)
	}
	wg.Wait()
}
