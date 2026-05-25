package contexttest

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	fmt.Println(ctx)

	ctxTodo := context.TODO()
	fmt.Println(ctxTodo)
}

func TestContextWithValue(t *testing.T) {
	ctxA := context.Background()

	ctxB := context.WithValue(ctxA, "b", "B")
	ctxC := context.WithValue(ctxA, "c", "C")

	ctxD := context.WithValue(ctxB, "d", "D")
	ctxE := context.WithValue(ctxB, "e", "E")

	ctxF := context.WithValue(ctxC, "f", "F")
	ctxG := context.WithValue(ctxF, "g", "G")

	fmt.Println(ctxA)
	fmt.Println(ctxB)
	fmt.Println(ctxC)
	fmt.Println(ctxD)
	fmt.Println(ctxE)
	fmt.Println(ctxF)
	fmt.Println(ctxG)

	fmt.Println(ctxF.Value("f"))
	fmt.Println(ctxF.Value("c"))
	fmt.Println(ctxF.Value("b"))

	fmt.Println(ctxA.Value("b"))
}

/*
				[ ctxA ] (Akar Kosong)
                        │
         ┌──────────────┴──────────────┐
         ▼                             ▼
      [ ctxB ]                      [ ctxC ]
     Key: "b"="B"                  Key: "c"="C"
         │                             │
   ┌─────┴─────┐                       ▼
   ▼           ▼                    [ ctxF ]
[ ctxD ]    [ ctxE ]               Key: "f"="F"
Key: "d"="D" Key: "e"="E"              │
                                       ▼
                                    [ ctxG ]
                                   Key: "g"="G"
*/

// dari gemini
type contextKey string

func TestContextWithValueYangBenar(t *testing.T) {
	ctxA := context.Background()

	// Gunakan tipe data kustom tersebut saat membuat value
	const traceKey contextKey = "trace_id"
	ctxB := context.WithValue(ctxA, traceKey, "TX-1002")

	// Cara mengambilnya juga harus menggunakan tipe data yang sama
	fmt.Println(ctxB.Value(traceKey)) // Output: TX-1002
}

// context with cancel
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done(): // 3. Menangkap sinyal cancel dari luar!
				return // 4. Goroutine langsung mati (selesai) di sini.
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulasi slow
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break // 1. Keluar dari loop pembacaan counter
		}
	}
	cancel() // 2. 🎯 TOMBOL DARURAT DIPENCET!

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel() // jika proses nya lebih cepat dari waktu timeout, maka cancel() akan langsung menghentikan go routine tanpa harus menunggu hinggal waktu timeout habis

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel() // jika proses nya lebih cepat dari waktu deadline, maka cancel() akan langsung menghentikan go routine tanpa harus menunggu hinggal waktu deadline habis

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
