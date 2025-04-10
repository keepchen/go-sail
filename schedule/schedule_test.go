package schedule

import (
	"fmt"
	"testing"
)

func TestGenerateJobNameKey(t *testing.T) {
	t.Log(generateJobNameKey("tester"))
}

func TestNewJob(t *testing.T) {
	t.Run("NewJob", func(t *testing.T) {
		scheduler := NewJob("NewJob", func() {
			fmt.Println("NewJob...")
		})
		t.Log(scheduler)
	})
}

func TestJob(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		scheduler := Job("Job", func() {
			fmt.Println("Job...")
		})
		t.Log(scheduler)
	})
}

func TestJobIsRunning(t *testing.T) {
	t.Run("JobIsRunning", func(t *testing.T) {
		t.Log(JobIsRunning("tester"))
	})
}

func TestCall(t *testing.T) {
	t.Run("Call", func(t *testing.T) {
		Call("tester", false)
	})
}

func TestMustCall(t *testing.T) {
	t.Run("MustCall", func(t *testing.T) {
		scheduler := NewJob("MustCall", func() {
			fmt.Println("MustCall...")
		})
		t.Log(scheduler)
		MustCall("MustCall", true)
	})
}
