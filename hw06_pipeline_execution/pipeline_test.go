package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)

			// Из способов закрыть горутину знаю только этот
			// Позволить ей слушать сигнальный канал, операция чтенния на нем будет блокирующей потому что туда никто не пишет
			// В момент закрытия сигнального канала - чтение из него разблокируется,
			//- срабатывает select -> return -> defer -> close(out)
			// горутина завершает работу и закрывает свой канал с данными
			// в следующем stage инструкция for v := range in { прерывается так как канал in закрыт в предыдущем stage

			// получаем глобальный сигнальный канал
			done := DoneChannel()

			// Если это решение неверно - пробовал также слушать и закрывать доступный канал out в качестве сигнального
			// но это опасно - так как можно попасть на ситуацию когда будем пробовать писать в закрытый канал out

			// Даже с таким  решением тест "done case" срабатывает через раз :(
			// Время закрытия пайплана через раз получается больше ожидаемого
			// Что я принципиально делаю не так ? ...

			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					//require.Eventually(t, func() bool { return true }, sleepPerStage+time.Second, sleepPerStage)
					select {
					case <-done:
						return
					case out <- f(v):
					}
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}
