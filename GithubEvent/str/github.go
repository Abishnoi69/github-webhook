package str

import (
	"fmt"
	"github-webhook/GithubEvent/config"
	"github.com/google/go-github/v71/github"
	"log"
	"net/http"
	"strings"
)

// GitHubWebhook processes GitHub webhooks
func GitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		log.Printf("Error validating payload: %v\n", err)
		http.Error(w, "Invalid payload", http.StatusUnauthorized)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("Error parsing webhook: %v\n", err)
		http.Error(w, "Error parsing webhook", http.StatusInternalServerError)
		return
	}

	// Prioritize critical or frequent event types
	var message string
	switch e := event.(type) {
	case *github.PushEvent:
		message = handlePushEvent(e)
	case *github.PullRequestEvent:
		message = handlePullRequestEvent(e)
	case *github.IssuesEvent:
		message = handleIssuesEvent(e)
	case *github.PingEvent:
		message = handlePingEvent(e)

	// Handle review-related events
	case *github.PullRequestReviewEvent:
		message = handlePullRequestReviewEvent(e)
	case *github.PullRequestReviewCommentEvent:
		message = handlePullRequestReviewCommentEvent(e)

	// Handle repository and organization events
	case *github.RepositoryEvent:
		message = handleRepositoryEvent(e)
	case *github.RepositoryDispatchEvent:
		message = handleRepositoryDispatchEvent(e)
	case *github.OrganizationEvent:
		message = handleOrganizationEvent(e)
	case *github.OrgBlockEvent:
		message = handleOrgBlockEvent(e)

	// Handle CI/CD and deployment-related events
	case *github.CheckRunEvent:
		message = handleCheckRunEvent(e)
	case *github.CheckSuiteEvent:
		message = handleCheckSuiteEvent(e)
	case *github.WorkflowRunEvent:
		message = handleWorkflowRunEvent(e)
	case *github.WorkflowJobEvent:
		message = handleWorkflowJobEvent(e)
	case *github.DeploymentEvent:
		message = handleDeploymentEvent(e)
	case *github.DeploymentStatusEvent:
		message = handleDeploymentStatusEvent(e)

	// Handle advisory and security-related events
	case *github.SecurityAdvisoryEvent:
		message = handleSecurityAdvisoryEvent(e)
	case *github.MembershipEvent:
		message = handleMembershipEvent(e)
	case *github.MilestoneEvent:
		message = handleMilestoneEvent(e)

	// Handle less frequent or low-priority events
	case *github.CommitCommentEvent:
		message = handleCommitCommentEvent(e)
	case *github.ForkEvent:
		message = handleForkEvent(e)
	case *github.ReleaseEvent:
		message = handleReleaseEvent(e)
	case *github.StarEvent:
		message = handleStarEvent(e)
	case *github.WatchEvent:
		message = handleWatchEvent(e)
	case *github.LabelEvent:
		message = handleLabelEvent(e)
	case *github.MarketplacePurchaseEvent:
		message = handleMarketplacePurchaseEvent(e)
	case *github.PageBuildEvent:
		message = handlePageBuildEvent(e)
	case *github.DeployKeyEvent:
		message = handleDeployKeyEvent(e)
	case *github.StarredRepository:
		message = handleStarredEvent(e)
	case *github.CreateEvent:
		message = handleCreateEvent(e)
	case *github.DeleteEvent:
		message = handleDeleteEvent(e)
	case *github.IssueCommentEvent:
		message = handleIssueCommentEvent(e)
	case *github.MemberEvent:
		message = handleMemberEvent(e)
	case *github.PublicEvent:
		message = handlePublicEvent(e)
	case *github.StatusEvent:
		message = handleStatusEvent(e)
	case *github.WorkflowDispatchEvent:
		message = handleWorkflowDispatchEvent(e)
	case *github.TeamAddEvent:
		message = handleTeamAddEvent(e)
	case *github.TeamEvent:
		message = handleTeamEvent(e)
	case *github.PackageEvent:
		message = handlePackageEvent(e)
	case *github.GollumEvent:
		message = handleGollumEvent(e)
	case *github.MetaEvent:
		message = handleMetaEvent(e)
	// Catch-all fallback for unhandled events
	default:
		log.Printf("Unhandled event type: %s\n", github.WebHookType(r))
		message = fmt.Sprintf("Unhandled event type: %s", github.WebHookType(r))
	}

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id query parameter", http.StatusBadRequest)
		return
	}

	err = sendToTelegram(chatID, message)
	if err != nil {
		http.Error(w, strings.ReplaceAll(err.Error(), config.BotToken, "$Bot"), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(message))
}
