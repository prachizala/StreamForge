package queue

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type JobStatus string

const (
	StatusPending    JobStatus = "PENDING"
	StatusProcessing JobStatus = "PROCESSING"
	StatusCompleted  JobStatus = "COMPLETED"
	StatusFailed     JobStatus = "FAILED"
)

type Job struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Status    JobStatus `json:"status"`
	Progress  int       `json:"progress"`
	Stage     string    `json:"stage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JobDispatcher struct {
	mu         sync.RWMutex
	jobs       map[string]*Job
	workChan   chan *Job
	listeners  map[string][]chan Job
	listenerMu sync.Mutex
}

func NewDispatcher(workerCount int) *JobDispatcher {
	d := &JobDispatcher{
		jobs:      make(map[string]*Job),
		workChan:  make(chan *Job, 100),
		listeners: make(map[string][]chan Job),
	}

	for i := 0; i < workerCount; i++ {
		go d.workerLoop()
	}

	return d
}

func (d *JobDispatcher) CreateJob(filename string) *Job {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	jobID := "job_" + hex.EncodeToString(b)

	job := &Job{
		ID:        jobID,
		Filename:  filename,
		Status:    StatusPending,
		Progress:  0,
		Stage:     "Queued in Gateway",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	d.mu.Lock()
	d.jobs[jobID] = job
	d.mu.Unlock()

	d.workChan <- job
	return job
}

func (d *JobDispatcher) GetJob(jobID string) (*Job, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	j, exists := d.jobs[jobID]
	return j, exists
}

func (d *JobDispatcher) workerLoop() {
	for job := range d.workChan {
		d.processJob(job)
	}
}

func (d *JobDispatcher) processJob(job *Job) {
	stages := []struct {
		progress int
		stage    string
		delay    time.Duration
	}{
		{progress: 25, stage: "Rust: Extracting Audio Waveform Vector", delay: 40 * time.Millisecond},
		{progress: 50, stage: "Rust: Generating Animated WebP Thumbnails", delay: 40 * time.Millisecond},
		{progress: 75, stage: "Rust: Embedding StreamForge Watermark", delay: 40 * time.Millisecond},
		{progress: 100, stage: "Pipeline Complete", delay: 20 * time.Millisecond},
	}

	d.updateJobStatus(job.ID, StatusProcessing, 5, "Dispatched to Rust Media Kernel")

	for _, s := range stages {
		time.Sleep(s.delay)
		d.updateJobStatus(job.ID, StatusProcessing, s.progress, s.stage)
	}

	d.updateJobStatus(job.ID, StatusCompleted, 100, "Ready for Stream Playback")
}

func (d *JobDispatcher) updateJobStatus(jobID string, status JobStatus, progress int, stage string) {
	d.mu.Lock()
	job, exists := d.jobs[jobID]
	if exists {
		job.Status = status
		job.Progress = progress
		job.Stage = stage
		job.UpdatedAt = time.Now()
	}
	d.mu.Unlock()

	if exists {
		d.notifyListeners(jobID, *job)
	}
}

func (d *JobDispatcher) notifyListeners(jobID string, job Job) {
	d.listenerMu.Lock()
	defer d.listenerMu.Unlock()

	chans := d.listeners[jobID]
	for _, ch := range chans {
		select {
		case ch <- job:
		default:
		}
	}
}
