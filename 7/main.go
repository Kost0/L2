package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	// Создаем небуферированный канал
	c := make(chan int)
	go func() {
		// Записываем в него данные
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		// Закрываем канал
		close(c)
	}()
	// Возвращаем канал
	return c
}

func merge(a, b <-chan int) <-chan int {
	// Создаем результирующий канал
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				// Если канал не закрыт
				if ok {
					// Передаем данные в результирующий канал
					c <- v
				} else {
					// Иначе обнуляем канал a
					a = nil
				}
			case v, ok := <-b:
				// Если канал не закрыт
				if ok {
					// Передаем данные в результирующий канал
					c <- v
				} else {
					// Иначе обнуляем канал b
					b = nil
				}
			}
			// Если оба канала обнулены
			if a == nil && b == nil {
				// Закрываем результирующий канал
				close(c)
				// Выходим из горутины
				return
			}
		}
	}()
	// Возвращаем канал с объединенными данными
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	// Получаем каналы и запускаем передачу в них информации
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	// Запускаем конвейер
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}
