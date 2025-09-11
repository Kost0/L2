package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	// Создаем переменную с типом интерфейса error
	var err error
	// Получаем значение nil, упаковываем его в интерфейс,
	// из-за чего сам интерфейс не будет nil,
	// так как хранит информацию о типе error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

// Вывод: error
