package methods

import (
	"database/sql"

	"golang.org/x/net/context"

	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) UserRepositories(ctx context.Context, request *pb.UserRequest) (*pb.UserRepositoriesResponse, error) {
	var response pb.UserRepositoriesResponse

	// Database
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// TODO: Fetch repositories?

	rowsRepos, err := db.Query("select distinct repo_name from githubevents where actor = ? order by repo_name", request.Username)
	if err != nil {
		db.Close()
		return nil, err
	}
	reponame := ""
	username := ""
	displayname := ""
	for rowsRepos.Next() {
		err = rowsRepos.Scan(&reponame)
		if err != nil {
			db.Close()
			return nil, err
		}
		repositoryResponse := pb.RepositoryResponse{Name: reponame}
		// Get team members
		rowsMembers, err := db.Query("select distinct u.login username, u.name displayname from githubusers u inner join githubevents e on u.login = e.actor where e.actor <> ? and e.repo_name = ?", request.Username, reponame)
		if err != nil && err != sql.ErrNoRows {
			db.Close()
			return nil, err
		}
		for rowsMembers.Next() {
			err = rowsMembers.Scan(&username, &displayname)
			if err != nil {
				db.Close()
				return nil, err
			}
			repositoryResponse.TeamMembers = append(repositoryResponse.TeamMembers, &pb.UserResponse{Username: username, DisplayName: displayname})
		}
		rowsMembers.Close()
		response.Repository = append(response.Repository, &repositoryResponse)
	}
	rowsRepos.Close()
	db.Close()

	return &response, nil

}
