package main

import (
	"fmt"
	"os"
)

func Foo() error {
	// Создается переменная типа *os.PathError, которая равна nil
	var err *os.PathError = nil
	// Возвращается интерфейс error
	// Он хранит информацию о:
	// Типе равном *os.PathError
	// Значении равном nil
	return err
}

func main() {
	// Переменная err получает тип error
	err := Foo()
	// Видим что err хранит значение nil
	fmt.Println(err)
	// Но не является nil
	fmt.Println(err == nil)
}

// Вывод
// <nil>
// false

// Интерфейс равен nil, только в том случае,
// если отсутствуют информация о типе и данные,
// точнее, они равны nil
// В нашем же случае, err хранила информацию о типе не равном nil

// Пустой интерфейс не требует никаких методов
