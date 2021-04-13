package loop_test

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"guitest/goey/internal/nopanic"
	"guitest/goey/loop"
)

func ExampleRun() {
	// This init function will be used to create a window on the GUI thread.
	init := func() error {
		// Create an empty window.  Note that most user code should instead
		// create a window, which will handle lock counting.
		loop.AddLockCount(1)

		go func() {
			// Because of goroutine, we are now off the GUI thread.
			// Schedule an action.
			err := loop.Do(func() error {
				loop.AddLockCount(-1)
				fmt.Println("...like tears in rain")
				return nil
			})
			if err != nil {
				fmt.Println("Error:", err)
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
	// ...like tears in rain
}

func ExampleDo() {
	err := loop.Do(func() error {
		// Inside this closure, we will be executing only on the GUI thread.
		_, err := fmt.Println("Hello.")
		// Return the error (if any) back to the caller.
		return err
	})

	// Report on the success or failure
	fmt.Println("Previous call to fmt.Println resulted in ", err)
}

func TestMain(m *testing.M) {
	loop.TestMain(m)
}

func TestRun(t *testing.T) {
	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			// Try running the main loop again, but in parallel.  We should get an error.
			err := loop.Run(func() error {
				return nil
			})
			if err != loop.ErrAlreadyRunning {
				t.Errorf("Expected ErrAlreadyRunning, got %s", err)
			}

			// Close the window.  This should stop the GUI loop.
			err = loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				t.Errorf("Unexpected error in call to Do")
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Unexpected error in call to Run")
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestRunWithError(t *testing.T) {
	const errorString = "No luck"

	// Make sure that error is passed through to caller
	init := func() error {
		return errors.New(errorString)
	}

	err := loop.Run(init)
	if err == nil {
		t.Errorf("Unexpected success, no error returned")
	} else if s := err.Error(); errorString != s {
		t.Errorf("Unexpected error, want %s, got %s", errorString, s)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestRunWithPanic(t *testing.T) {
	const errorString = "No luck"

	defer func() {
		r := recover()
		if r != nil {
			if pe, ok := r.(nopanic.PanicError); ok {
				r = pe.Value()
			}

			if s, ok := r.(string); !ok {
				t.Errorf("Unexpected recover, %v", r)
			} else if s != errorString {
				t.Errorf("Unexpected recover, %s", s)
			}
		} else {
			t.Errorf("Missing panic")
		}

		// Make sure that window count is properly maintained.
		if c := loop.LockCount(); c != 0 {
			t.Errorf("Want lockCount==0, got lockCount==%d", c)
		}
	}()

	// Make sure that error is passed through to caller
	init := func() error {
		panic(errorString)
	}

	// Make sure that window count is properly maintained.
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}

	err := loop.Run(init)
	if err == nil {
		t.Errorf("Unexpected success, no error returned")
	} else if s := err.Error(); errorString != s {
		t.Errorf("Unexpected error, want %s, got %s", errorString, s)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestRunWithWindowClose(t *testing.T) {
	// Make sure that error is passed through to caller
	init := func() error {
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		loop.AddLockCount(-1)
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Unexpected error in call to Run")
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestDo(t *testing.T) {
	count := uint32(0)

	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			// Run the actions, which are counted.
			for i := 0; i < 10; i++ {
				err := loop.Do(func() error {
					atomic.AddUint32(&count, 1)
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Check for an error return
			err := loop.Do(func() error {
				// Some example error
				return loop.ErrNotRunning
			})
			if err != loop.ErrNotRunning {
				t.Errorf("Error in Do, expected %v, got %v", loop.ErrNotRunning, err)
			}

			// Close the window
			err = loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
	if c := atomic.LoadUint32(&count); c != 10 {
		t.Errorf("Want count=10, got count==%d", c)
	}
}

func BenchmarkRunNoInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := loop.Run(func() error {
			return nil
		})
		if err != nil {
			b.Errorf("Call to Run failed: %s", err)
		}
	}
}

func BenchmarkRun(b *testing.B) {
	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			b.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			b.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			// Close the window
			err := loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				b.Errorf("Error in Do, %s", err)
			}
		}()

		return nil
	}

	for i := 0; i < b.N; i++ {
		err := loop.Run(init)
		if err != nil {
			b.Errorf("Call to Run failed: %s", err)
		}
	}
}

func BenchmarkDo(b *testing.B) {
	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			b.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			b.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := loop.Do(func() error {
					return nil
				})
				if err != nil {
					b.Errorf("Error in Do, %s", err)
				}
			}
			b.StopTimer()

			// Close the window
			err := loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				b.Errorf("Error in Do, %s", err)
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		b.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := loop.LockCount(); c != 0 {
		b.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestDoFailure(t *testing.T) {
	err := loop.Do(func() error {
		return nil
	})

	if err != loop.ErrNotRunning {
		t.Errorf("Unexpected success in call to Do")
	}
}

func TestDoWithError(t *testing.T) {
	const errorString = "No luck"

	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			// Run the actions, which are counted.
			err := loop.Do(func() error {
				return errors.New(errorString)
			})
			if err == nil {
				t.Errorf("Failed to return error in Do")
			} else if err.Error() != errorString {
				t.Errorf("Incorrect error returned in Do, %v != %v", err.Error(), errorString)
			}

			// Close the window
			err = loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestDoWithPanic(t *testing.T) {
	const errorString = "No luck"

	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Failed to recover the expected panic")
				} else if _, ok := r.(error); !ok {
					t.Errorf("Unexpected value for the panic")
				}

				// Need to head back to the GUI thread to stop the event loop.
				_ = loop.Do(func() error {
					// Close the window
					loop.AddLockCount(-1)
					return nil
				})
			}()

			// Run the actions, which are counted.
			_ = loop.Do(func() error {
				panic(errorString)
			})
			t.Errorf("unreachable")
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
}

func TestThunderHerdOfDo(t *testing.T) {
	count := uint32(0)
	// This is the limit for number of simultaneous goroutines under the race
	// detector.
	const herdSize = 8128

	init := func() error {
		// Verify that the test is starting in the correct state.
		if c := loop.LockCount(); c != 1 {
			t.Errorf("Want lockCount==1, got lockCount==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
		loop.AddLockCount(1)
		if c := loop.LockCount(); c != 2 {
			t.Fatalf("Want lockCount==2, got lockCount==%d", c)
		}

		go func() {
			// At least let the run loop start before we hammer it.
			err := loop.Do(func() error {
				atomic.AddUint32(&count, 1)
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}

			// Start the herd.
			wg := sync.WaitGroup{}
			for i := 0; i < herdSize; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					err := loop.Do(func() error {
						atomic.AddUint32(&count, 1)
						return nil
					})
					if err != nil {
						t.Errorf("Error in Do, %s", err)
					}
				}()
			}
			wg.Wait()

			// Close the window
			err = loop.Do(func() error {
				loop.AddLockCount(-1)
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}()

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := loop.LockCount(); c != 0 {
		t.Errorf("Want lockCount==0, got lockCount==%d", c)
	}
	if c := atomic.LoadUint32(&count); c != herdSize+1 {
		t.Errorf("Want count=10, got count==%d", c)
	}
}
