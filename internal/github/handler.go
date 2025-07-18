package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitnotify/internal/config"

	"github.com/google/go-github/v63/github"
)

// Handler handles GitHub webhook events
type Handler struct {
	config *config.Config
}

// NewHandler creates a new webhook handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
	}
}

// HandleWebhook processes incoming GitHub webhook requests
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the webhook signature
	_, err = github.ValidatePayload(r, []byte(h.config.WebhookSecret))
	if err != nil {
		log.Printf("Invalid webhook signature: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the event type from headers
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		log.Printf("Missing X-GitHub-Event header")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process the event based on its type
	if err := h.processEvent(eventType, payload); err != nil {
		log.Printf("Error processing event: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success response to GitHub
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// processEvent handles different GitHub event types
func (h *Handler) processEvent(eventType string, payload []byte) error {
	switch eventType {
	case "issues":
		return h.handleIssuesEvent(payload)
	case "pull_request":
		return h.handlePullRequestEvent(payload)
	default:
		log.Printf("Ignoring event: %s", eventType)
		return nil
	}
}

// handleIssuesEvent processes GitHub issue events
func (h *Handler) handleIssuesEvent(payload []byte) error {
	var event github.IssuesEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal issues event: %w", err)
	}

	// Check if we should notify for this event
	if !h.shouldNotify("issues", event.GetAction(), event.GetRepo().GetName()) {
		return nil
	}

	// Print notification message
	action := event.GetAction()
	switch action {
	case "opened":
		log.Printf("New Issue Opened: #%d - %s by %s",
			event.GetIssue().GetNumber(),
			event.GetIssue().GetTitle(),
			event.GetIssue().GetUser().GetLogin())
	case "closed":
		log.Printf("Issue Closed: #%d - %s by %s",
			event.GetIssue().GetNumber(),
			event.GetIssue().GetTitle(),
			event.GetIssue().GetUser().GetLogin())
	case "reopened":
		log.Printf("Issue Reopened: #%d - %s by %s",
			event.GetIssue().GetNumber(),
			event.GetIssue().GetTitle(),
			event.GetIssue().GetUser().GetLogin())
	default:
		log.Printf("Issue %s: #%d - %s by %s",
			action,
			event.GetIssue().GetNumber(),
			event.GetIssue().GetTitle(),
			event.GetIssue().GetUser().GetLogin())
	}

	return nil
}

// handlePullRequestEvent processes GitHub pull request events
func (h *Handler) handlePullRequestEvent(payload []byte) error {
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal pull request event: %w", err)
	}

	// Check if we should notify for this event
	if !h.shouldNotify("pull_request", event.GetAction(), event.GetRepo().GetName()) {
		return nil
	}

	// Print notification message
	action := event.GetAction()
	switch action {
	case "opened":
		log.Printf("New Pull Request Opened: #%d - %s by %s",
			event.GetPullRequest().GetNumber(),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetUser().GetLogin())
	case "closed":
		if event.GetPullRequest().GetMerged() {
			log.Printf("Pull Request Merged: #%d - %s by %s",
				event.GetPullRequest().GetNumber(),
				event.GetPullRequest().GetTitle(),
				event.GetPullRequest().GetUser().GetLogin())
		} else {
			log.Printf("Pull Request Closed: #%d - %s by %s",
				event.GetPullRequest().GetNumber(),
				event.GetPullRequest().GetTitle(),
				event.GetPullRequest().GetUser().GetLogin())
		}
	case "reopened":
		log.Printf("Pull Request Reopened: #%d - %s by %s",
			event.GetPullRequest().GetNumber(),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetUser().GetLogin())
	default:
		log.Printf("Pull Request %s: #%d - %s by %s",
			action,
			event.GetPullRequest().GetNumber(),
			event.GetPullRequest().GetTitle(),
			event.GetPullRequest().GetUser().GetLogin())
	}

	return nil
}

// shouldNotify checks if we should send a notification based on configuration
func (h *Handler) shouldNotify(eventType, action, repoName string) bool {
	// If no notification rules are configured, notify for everything
	if len(h.config.Notifications) == 0 {
		return true
	}

	for _, notification := range h.config.Notifications {
		// Check if event type matches
		if notification.EventType != eventType {
			continue
		}

		// Check if action matches (if actions are specified)
		if len(notification.Actions) > 0 {
			actionMatches := false
			for _, allowedAction := range notification.Actions {
				if allowedAction == action {
					actionMatches = true
					break
				}
			}
			if !actionMatches {
				continue
			}
		}

		// Check if repo matches (if repos are specified)
		if len(notification.Repos) > 0 {
			repoMatches := false
			for _, allowedRepo := range notification.Repos {
				if allowedRepo == repoName {
					repoMatches = true
					break
				}
			}
			if !repoMatches {
				continue
			}
		}

		// If we get here, all conditions match
		return true
	}

	return false
}
