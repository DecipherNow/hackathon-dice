package methods

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	gc "github.com/lucasmoten/project-2502/services/dice/apis/github"
	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) UserActivity(ctx context.Context, request *pb.UserRequest) (*pb.UserActivityResponse, error) {

	var response pb.UserActivityResponse

	if err := checkAndUpdateEvents(); err != nil {
		return nil, err
	}

	// Database
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Events in the database reflect the current state
	rowsEvents, err := db.Query("select created_at, type, repo_name, summaryline from githubevents where actor = ? order by created_at desc limit 5", request.Username)
	if err != nil {
		db.Close()
		return nil, err
	}
	createdat := ""
	eventtype := ""
	reponame := ""
	summaryline := ""
	for rowsEvents.Next() {
		err = rowsEvents.Scan(&createdat, &eventtype, &reponame, &summaryline)
		if err != nil {
			db.Close()
			return nil, err
		}
		response.Activity = append(response.Activity, &pb.ActivityResponse{CreatedAt: createdat, Actor: request.Username, RepoName: reponame, Type: eventtype, Summaryline: summaryline})
	}
	rowsEvents.Close()
	db.Close()

	return &response, nil

}

func checkAndUpdateEvents() error {
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Str("service", "dice").Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Database
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Check existing cache to determine when last updated
	dbRows, err := db.Query("select created_at from githubevents order by created_at desc limit 1")
	if err != nil {
		db.Close()
		return err
	}
	createdat := "2018-12-13T00:00:00Z"
	for dbRows.Next() {
		dbRows.Scan(&createdat)
	}
	dbRows.Close()
	// And see if its been over 10 minutes
	then, err := time.Parse(time.RFC3339, createdat)
	if err != nil {
		db.Close()
		return err
	}
	duration := time.Since(then)
	dataisold := (duration.Minutes() > float64(viper.GetInt("github_age_events")))
	summaryline := ""

	// If data is old, then fetch events from github and update table
	if dataisold {
		logger.Debug().Float64("duration", duration.Minutes()).Str("created_at", createdat).Msg("dataisold")
		createdat = time.Now().UTC().Format(time.RFC3339)

		// Get from GitHub
		githubclient, err := gc.NewClient()
		if err != nil {
			db.Close()
			return err
		}
		stmtInsertEvent, err := db.Prepare("insert into githubevents (id, created_at, actor, repo_name, type, action, payload_comment_body, payload_issue_number, payload_issue_title, payload_issue_body, payload_pr_number, payload_pr_title, payload_pr_body, payload_pr_merged_dt, payload_pr_merged, ref, ref_type, payload_size, summaryline) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) on conflict(id) do update set id=id where 1=0")
		if err != nil {
			db.Close()
			return err
		}
		pageHasResults := true
		pageNumber := 0
		for pageHasResults && pageNumber < viper.GetInt("github_maxpage_events") {
			pageNumber++
			events, err := githubclient.GetOrganizationEvents(pageNumber)
			if err != nil {
				db.Close()
				return err
			}
			pageHasResults = len(events) > 0
			for _, event := range events {
				// determine summaryline
				switch event.Type {
				case "CreateEvent":
					if event.Payload.RefType == "repository" {
						summaryline = fmt.Sprintf("created %s named %s", event.Payload.RefType, event.Repo.Name)
					} else {
						summaryline = fmt.Sprintf("created %s named %s in %s", event.Payload.RefType, event.Payload.Ref, event.Repo.Name)
					}
					break
				case "DeleteEvent":
					summaryline = fmt.Sprintf("deleted %s named %s in %s", event.Payload.RefType, event.Payload.Ref, event.Repo.Name)
					break
				case "IssuesCommentEvent":
					summaryline = fmt.Sprintf("commented on %s#%d %s", event.Repo.Name, event.Payload.Issue.Number, event.Payload.Issue.Title)
					break
				case "IssuesEvent":
					summaryline = fmt.Sprintf("%s %s#%d %s", event.Payload.Action, event.Repo.Name, event.Payload.Issue.Number, event.Payload.Issue.Body)
					break
				case "PullRequestEvent":
					if event.Payload.PullRequest.Merged {
						summaryline = fmt.Sprintf("merged PR %s#%d %s", event.Repo.Name, event.Payload.Number, event.Payload.PullRequest.Title)
					} else {
						summaryline = fmt.Sprintf("%s PR %s#%d %s", event.Payload.Action, event.Repo.Name, event.Payload.Number, event.Payload.PullRequest.Title)
					}
					break
				case "PullRequestReviewCommentEvent":
					summaryline = fmt.Sprintf("reviewed PR %s#%d", event.Repo.Name, event.Payload.Number)
					break
				case "PushEvent":
					summaryline = fmt.Sprintf("pushed commit(s) to %s in %s", event.Payload.Ref, event.Repo.Name)
					break
				default:
					summaryline = fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
					break
				}

				event.ID = event.ID + ""
				// insert ... on conflict(id) do update set id=id where 1=0
				merged := "false"
				if event.Payload.PullRequest.Merged {
					merged = "true"
				}
				_, err := stmtInsertEvent.Exec(event.ID, event.CreatedAt, event.Actor.Login, event.Repo.Name, event.Type, event.Payload.Action, event.Payload.Comment.Body, event.Payload.Issue.Number, event.Payload.Issue.Title, event.Payload.Issue.Body, event.Payload.Number, event.Payload.PullRequest.Title, event.Payload.PullRequest.Body, event.Payload.PullRequest.MergedAt, merged, event.Payload.Ref, event.Payload.RefType, event.Payload.Size, summaryline)
				if err != nil {
					db.Close()
					return err
				}

			}
		}

		// Update when last checked
		stmtCache, err := db.Prepare("UPDATE githubevents SET created_at = ? WHERE actor = '__cachestate__'")
		if err != nil {
			db.Close()
			return err
		}
		_, err = stmtCache.Exec(createdat)
		if err != nil {
			db.Close()
			return err
		}
	}
	return nil
}
