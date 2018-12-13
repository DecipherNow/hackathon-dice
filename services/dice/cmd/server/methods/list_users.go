package methods

import (
	"golang.org/x/net/context"

	gc "github.com/lucasmoten/project-2502/services/dice/apis/github"
	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var response pb.ListUsersResponse

	// TODO: See if we already have the users in local data cache

	// The else condition is to ... Use the github client to fetch members of the org
	githubclient, err := gc.NewClient()
	if err != nil {
		return nil, err
	}
	members, err := githubclient.GetMembers()
	if err != nil {
		return nil, err
	}
	// TODO: After retrieving these, we should save to local data cache

	// Iterate the members (from github, or local cache)
	for _, m := range members {
		username := m.Login
		displayName := "No Display Name!"
		// TODO: Determine if we already have this user in the local data cache

		// The else condition is to ... Use the github client to fetch the details for a user
		userinfo, err := githubclient.GetUser(username)
		if err != nil {
			return nil, err
		}
		if userinfo.Name != username {
			displayName = userinfo.Name
		}

		// TODO: After retrieving, save to local data cache

		// Add to the response
		response.Users = append(response.Users, &pb.UserResponse{Username: username, DisplayName: displayName})
	}

	return &response, nil
}
