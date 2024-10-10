package loops

import (
	"context"
	"testing"
	"time"
)

// TestStopBeforeStart checks that tasks are not running when Stop is called before they start.
func TestStopBeforeStart(t *testing.T) {
	loops := New()
	loops.AddFunc(time.Hour, func(ctx context.Context) {
		// This job should not run since we will stop before it starts
		t.Error("Job should not have started running")
	})
	loops.Start()
	time.Sleep(50 * time.Millisecond)

	// Call Stop before starting the jobs
	ctx := loops.Stop()

	// Check that the context is done immediately since jobs haven't started
	select {
	case <-ctx.Done():
		// Context is done, which is expected
	case <-time.After(50 * time.Millisecond):
		t.Error("Expected context to be done immediately, but it wasn't")
	}
}

// TestStopWithRunningJobs tests stopping the Loops while jobs are actively running.
func TestStopWithRunningJobs(t *testing.T) {
	// This flag will indicate when the job is running.
	jobRunning := make(chan struct{})
	jobDone := make(chan struct{})

	l := New()
	// Adding a job that signals when it starts running and waits for a short period.
	l.AddFunc(50*time.Millisecond, func(ctx context.Context) {
		close(jobRunning)                  // Signal that the job is running.
		time.Sleep(200 * time.Millisecond) // Simulate long-running job.
		close(jobDone)                     // Signal that the job is done.
	})

	// Start the jobs.
	l.Start()

	// Wait for the job to start running.
	<-jobRunning

	// Stop the Loops while the job is running.
	ctx := l.Stop()

	// Wait for the job to finish.
	select {
	case <-jobDone:
		// Job finished successfully.
	case <-time.After(300 * time.Millisecond):
		// Context should be done, but we should allow enough time for the job to complete.
		t.Error("expected job to stop, but it is still running")
	}

	// Ensure the context is canceled.
	select {
	case <-ctx.Done():
		if ctx.Err() != context.Canceled {
			t.Error("expected context to be canceled, but it wasn't")
		}
	default:
		t.Error("expected context to be canceled, but it wasn't")
	}
}

// TestStopDuringTask tests stopping the task while it is in progress
func TestStopDuringTask(t *testing.T) {
	// Initialize a new loop manager instance
	loops := New()

	// Add a task that runs every 50 milliseconds, simulating work
	loops.AddFunc(50*time.Millisecond, func(ctx context.Context) {
		for {
			select {
			// Exit the loop when the context is canceled (task is stopped)
			case <-ctx.Done():
				return
			default:
			}
			// Simulate the task by sleeping for 50 milliseconds
			time.Sleep(50 * time.Millisecond)
		}
	})

	// Start the loops (tasks begin execution)
	loops.Start()

	// Allow the task to start by waiting for 100 milliseconds
	time.Sleep(100 * time.Millisecond)

	// Stop the loops while the task is running and get the stop context
	ctx := loops.Stop()

	// Verify that the context is immediately canceled after stopping
	select {
	case <-ctx.Done():
		// Context was canceled as expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected context to be canceled, but it wasn't")
	}
}
