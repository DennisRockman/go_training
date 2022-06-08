package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("PushFront check", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		require.Equal(t, l.Front().Value, l.Back().Value)

		middle := l.PushFront(20) // [20, 10]
		require.Equal(t, l.Front().Value, 20)
		require.Equal(t, l.Back().Value, 10)

		l.PushFront(30)                                             // [30, 20 , 10]
		require.Equal(t, l.Front().Value, 30)                       // head
		require.Equal(t, l.Back().Value, 10)                        // tail
		require.Equal(t, middle.Value, 20)                          // middle
		require.Equal(t, l.Back().Prev.Value, l.Front().Next.Value) // middle from tail == middle from head
	})

	t.Run("PushBack check", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		require.Equal(t, l.Front().Value, l.Back().Value)

		middle := l.PushBack(20) // [10, 20]
		require.Equal(t, l.Front().Value, 10)
		require.Equal(t, l.Back().Value, 20)

		l.PushBack(30)                                              // [10 , 20 , 30]
		require.Equal(t, l.Front().Value, 10)                       // head
		require.Equal(t, l.Back().Value, 30)                        // tail
		require.Equal(t, middle.Value, 20)                          // middle
		require.Equal(t, l.Back().Prev.Value, l.Front().Next.Value) // middle from tail == middle from head
	})

	t.Run("Remove check", func(t *testing.T) {
		l := NewList()

		i := l.PushBack(10) // [10]
		require.Equal(t, 1, l.Len())
		require.NotNil(t, l.Front())
		require.NotNil(t, l.Back())

		l.Remove(i) // empty
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		i = l.PushFront(50) // [50]
		l.PushBack(70)      // [50, 70]
		l.Remove(i)         // [70]
		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front().Value, l.Back().Value)

		i = l.PushFront(75) // [75 , 70]
		l.PushFront(85)     // [85, 75 , 70]
		l.Remove(i)         // [85, 70]
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Value, 85)
		require.Equal(t, l.Back().Value, 70)
		require.Equal(t, l.Front().Next.Value, l.Back().Value)
		require.Equal(t, l.Back().Prev.Value, l.Front().Value)
	})

	t.Run("MoveToFront check", func(t *testing.T) {
		l := NewList()
		l.MoveToFront(l.PushBack(10))
		require.Equal(t, 1, l.Len())
		require.Equal(t, l.Front().Value, 10)
		require.Equal(t, l.Back().Value, l.Front().Value)
		l.Remove(l.Front()) // empty

		for i := 4; i >= 0; i-- {
			l.PushFront(i)
		}

		for i := 0; i < 5; i++ {
			l.MoveToFront(l.Back())
		}

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{0, 1, 2, 3, 4}, elems)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
