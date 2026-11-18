package queue

import (
	"testing"
	"time"
)

func TestJobDispatcherLifecycle(t *testing.T) {
	dispatcher := NewDispatcher(2)

	job := dispatcher.CreateJob("test_sample.mp4")
	if job == nil {
		t.Fatalf("expected job to be created, got nil")
	}

	if job.Status != StatusPending && job.Status != StatusProcessing {
		t.Errorf("expected initial status PENDING or PROCESSING, got %s", job.Status)
	}

	// Poll until completed
	timeout := time.After(3 * time.Second)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	completed := false
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for job to complete")
		case <-ticker.C:
			j, exists := dispatcher.GetJob(job.ID)
			if !exists {
				t.Fatalf("job not found in dispatcher")
			}
			if j.Status == StatusCompleted && j.Progress == 100 {
				completed = true
				break
			}
		}
		if completed {
			break
		}
	}

	if !completed {
		t.Errorf("job failed to reach completed state")
	}
}
