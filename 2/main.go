package main

import "fmt"

// Так как эта функция имеет именованную возвращаемую переменную,
// defer выполнится перед return, таким образом алгоритм выполняется в таком порядке:
// x = 1
// defer
// return
func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

// Так как эта функция не имеет именованную возвращаемую переменную,
// defer выполнится после return, таким образом алгоритм выполняется в таком порядке:
// x = 1
// return
// defer
func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}

// Вывод:
// 2
// 1
