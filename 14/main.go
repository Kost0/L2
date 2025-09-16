package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	// Обработка крайних случаев
	if len(channels) == 0 {
		c := make(chan interface{})
		close(c)
		return c
	} else if len(channels) == 1 {
		return channels[0]
	}

	// Создаем канал для объединения
	orDone := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(channels))

	// Для каждого канала запускаем горутину
	for _, channel := range channels {
		go func(c <-chan interface{}) {
			defer wg.Done()
			select {
			// Если закрывается один из каналов
			case <-c:
				select {
				// Проверяем не закрыт ли уже orDone
				case <-orDone:
				// Если не закрыт, закрываем
				default:
					close(orDone)
				}
			// После закрытия orDone, завершаем горутину
			case <-orDone:
			}
		}(channel)
	}

	// Горутина для ожидания закрытия всех горутин
	go func() {
		wg.Wait()
		// Убеждаемся, что orDone закрыт
		select {
		case <-orDone:
		default:
			close(orDone)
		}
	}()

	return orDone
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
