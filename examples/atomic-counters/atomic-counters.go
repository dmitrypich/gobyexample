// Основной механизм управления состоянием в Go
// – связь через каналы. Вы видели это в примере
// с [набором обработчиков](worker-pools.html).
// Но есть несколько других вариантов для управления
// состоянием. В этом примере мы рассмотрим
// использование пакета `sync/atomic` для атомарных
// счетчиков, доступных нескольким горутинам.

package main

import "fmt"
import "time"
import "sync/atomic"
import "runtime"

func main() {

	// Мы будем использовать беззнаковое целое
	// число для представления нашего
	// (всегда положительного) счетчика.
	var ops uint64 = 0

	// Для эмуляции конкурентных обновлений,
	// мы запустим 50 горутин, которые будут
	// увеличивать счетчик примерно
	// каждую миллисекунду.
	for i := 0; i < 50; i++ {
		go func() {
			for {
				// Для атомарного увеличения счетчика `ops`
				// мы используем `AddUint64`, предоставляя
				// указатель на память объявленного беззнакового
				// целого числа с использование `&` в синтаксисе.
				atomic.AddUint64(&ops, 1)

				// Разрешим остальным горутинам продолжить работу.
				runtime.Gosched()
			}
		}()
	}

	// Подождем секунду для накопления
	// результатов работы.
	time.Sleep(time.Second)

	// В случае безопасного использования счетчика
	// во время продолжающихся обновлений из
	// иных горутин, мы извлекаем копию его текущего
	// значение в `opsFinal` с помощью `LoadUint64`.
	// Как и выше, мы должны передать эту функцию
	// адреса в памяти с использованием `&ops` в синтаксисе
	// для извлечения значения.
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}
