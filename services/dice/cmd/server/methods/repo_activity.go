package methods

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) RepoActivity(ctx context.Context, request *pb.RepoRequest) (*pb.RepoActivityResponse, error) {

	var response pb.RepoActivityResponse

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
	fullreponame := fmt.Sprintf("%s/%s", request.Orgname, request.Reponame)
	rowsEvents, err := db.Query("select created_at, type, actor, summaryline from githubevents where repo_name = ? order by created_at desc limit 5", fullreponame)
	if err != nil {
		db.Close()
		return nil, err
	}
	createdat := ""
	eventtype := ""
	actor := ""
	summaryline := ""
	for rowsEvents.Next() {
		err = rowsEvents.Scan(&createdat, &eventtype, &actor, &summaryline)
		if err != nil {
			db.Close()
			return nil, err
		}
		response.Activity = append(response.Activity, &pb.ActivityResponse{CreatedAt: createdat, Actor: actor, RepoName: fullreponame, Type: eventtype, Summaryline: summaryline})
	}
	rowsEvents.Close()
	db.Close()

	return &response, nil

}
