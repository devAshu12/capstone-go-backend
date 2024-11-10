package handlers

import (
	"encoding/json"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"
	"sync"
	"time"
)

var (
	queue      []types.VideoProgress
	queueMutex sync.Mutex
	Quit       chan struct{}
)

func batchInsert(progressBatch []types.VideoProgress) {
	// query := `INSERT INTO user_progress (user_id, video_id, progress, time_spent, completion, updated_at)
	//           VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	//           ON CONFLICT (user_id, video_id)
	//           DO UPDATE SET progress = $3, time_spent = $4, completion = $5, updated_at = CURRENT_TIMESTAMP`
	for _, progress := range progressBatch {
		fmt.Println("Processing and saving to DB:", progress.VideoID)
	}
}

func ProcessQueue(interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			dispatchBatch(batchSize)

		case <-Quit: // Stop processing when quit signal is received
			fmt.Println("Shutting down queue processor...")
			dispatchBatch(0) // Flush remaining items
			return
		}
	}
}

func dispatchBatch(batchSize int) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(queue) >= batchSize || (batchSize == 0 && len(queue) > 0) {
		progressBatch := make([]types.VideoProgress, len(queue))
		copy(progressBatch, queue)
		queue = []types.VideoProgress{}
		go batchInsert(progressBatch)
	}
}

func UpdateProgress(w http.ResponseWriter, r *http.Request) {
	var progress types.VideoProgress
	err := json.NewDecoder(r.Body).Decode(&progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	queueMutex.Lock()
	queue = append(queue, progress)
	fmt.Println("Added to queue:", progress.VideoID)
	queueMutex.Unlock()

	w.WriteHeader(http.StatusAccepted)
}
