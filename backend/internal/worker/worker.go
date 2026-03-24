package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"localVercel/db"
	"localVercel/internal/deployer"
	"localVercel/internal/queue"
	"localVercel/models"
	"log"
	"time"
)

type Worker struct {
	Queue    queue.Queue
	Deployer *deployer.Deployer
}

func NewWorker(q queue.Queue, d *deployer.Deployer) *Worker {
	return &Worker{
		Queue:    q,
		Deployer: d,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("Worker started, waiting for jobs...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping...")
			return
		default:
			// 1. Dequeue
			job, err := w.Queue.Dequeue(ctx)
			if err != nil {
				// Handle timeout/empty if Redis blocks or fails
				if err != nil && err.Error() != "redis: nil" {
					// Use specific error checks if needed, but BRPop returns redis.Nil on timeout if configured? 
					// BRPop with 0 timeout blocks indefinitely unless connection lost
					log.Printf("Error dequeuing: %v", err)
					time.Sleep(2 * time.Second) // backoff
				}
				continue
			}

			log.Printf("Processing job %s (type: %s)", job.ID, job.Type)
			
			// 2. Process
			if job.Type == "deploy" {
				w.processDeployJob(ctx, job)
			} else {
				log.Printf("Unknown job type: %s", job.Type)
			}
		}
	}
}

type DeployPayload struct {
	DeploymentID string `json:"deployment_id"`
	RepoURL      string `json:"repo_url"`
	Branch       string `json:"branch"`
	ProjectID    string `json:"project_id"`
	BuildCmd     string `json:"build_cmd"`
	OutputDir    string `json:"output_dir"`
}

func (w *Worker) processDeployJob(ctx context.Context, job *queue.Job) {
	var payload DeployPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		log.Printf("Invalid payload for job %s: %v", job.ID, err)
		return
	}

	// Update DB: Building
	var deployment models.Deployment
	if err := db.DB.First(&deployment, "id = ?", payload.DeploymentID).Error; err != nil {
		log.Printf("Deployment not found: %s", payload.DeploymentID)
		return
	}
	deployment.Status = "building"
	deployment.StartedAt = time.Now()
	db.DB.Save(&deployment)

	// 1. Build
	// Detect framework first if not provided?
	// For now assume Deployer handles logic or payload has basic info.
	// We'll pass "" as framework to let Deployer auto-detect
	logs, err := w.Deployer.BuildProject(ctx, payload.DeploymentID, payload.RepoURL, payload.Branch, "", payload.BuildCmd)
	
	deployment.Logs = logs
	if err != nil {
		deployment.Status = "failed"
		deployment.CompletedAt = time.Now()
		db.DB.Save(&deployment)
		log.Printf("Build failed for job %s: %v", job.ID, err)
		return
	}

	// 2. Deploy (Copy artifacts)
	finalPath, err := w.Deployer.DeployArtifacts(payload.DeploymentID, payload.OutputDir)
	if err != nil {
		deployment.Status = "failed"
		deployment.Logs += "\nDeployment error: " + err.Error()
		deployment.CompletedAt = time.Now()
		db.DB.Save(&deployment)
		log.Printf("Deploy failed for job %s: %v", job.ID, err)
		return
	}

	// Success
	deployment.Status = "ready" // or "success"
	deployment.PreviewURL = fmt.Sprintf("http://%s.localhost:3000", payload.ProjectID) // Mock URL
	deployment.CompletedAt = time.Now()
	// Storing finalPath optionally or inferring it
	db.DB.Save(&deployment)
	
	log.Printf("Job %s completed successfully. Deployed to %s", job.ID, finalPath)
}
