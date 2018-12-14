package methods

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	gc "github.com/lucasmoten/project-2502/services/dice/apis/github"
	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var response pb.ListUsersResponse

	// Database
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check existing cache to determine when last updated
	dbRows, err := db.Query("select updated_at from githubusers order by updated_at desc limit 1")
	if err != nil {
		db.Close()
		return nil, err
	}
	updated_at := "2018-12-13T00:00:00Z" // just a default in the past
	for dbRows.Next() {
		dbRows.Scan(&updated_at)
	}
	dbRows.Close()
	// And see if its been over 2 hours
	then, err := time.Parse(time.RFC3339, updated_at)
	if err != nil {
		db.Close()
		return nil, err
	}
	duration := time.Since(then)
	dataisold := (duration.Minutes() > float64(viper.GetInt("github_age_users")))

	// If data is old, then fetch pages from github and update table
	if dataisold {
		updated_at = time.Now().UTC().Format(time.RFC3339)
		// Get from GitHub
		githubclient, err := gc.NewClient()
		if err != nil {
			db.Close()
			return nil, err
		}
		pageHasResults := true
		pageNumber := 0
		stmtUserUpdate, err := db.Prepare("update githubusers set avatar_url=?, name=?, email=?, updated_at=? where login=?")
		if err != nil {
			db.Close()
			return nil, err
		}
		stmtUserInsert, err := db.Prepare("insert into githubusers (login, avatar_url, name, email, created_at, updated_at) values (?,?,?,?,?,?)")
		if err != nil {
			db.Close()
			return nil, err
		}
		for pageHasResults && pageNumber < viper.GetInt("github_maxpage_users") {
			pageNumber++
			members, err := githubclient.GetMembers(pageNumber)
			if err != nil {
				db.Close()
				return nil, err
			}
			pageHasResults = len(members) > 0
			for _, m := range members {
				username := m.Login
				displayName := "No Display Name!"
				userinfo, err := githubclient.GetUser(username)
				if err != nil {
					db.Close()
					return nil, err
				}
				if userinfo.Name != username {
					displayName = userinfo.Name
				}
				// Determine if we already have this user in the local data cache
				var rowsCount sql.NullString
				err = db.QueryRow("select count(0) from githubusers where login=?", username).Scan(&rowsCount)
				hasRow := false
				if err != nil && err != sql.ErrNoRows {
					db.Close()
					return nil, err
				}
				if rowsCount.Valid {
					hasRow = (rowsCount.String != "0")
				}
				// Insert or Update
				if hasRow {
					// update
					_, err := stmtUserUpdate.Exec(m.AvatarURL, displayName, userinfo.Email, userinfo.Updated, username)
					if err != nil {
						db.Close()
						return nil, err
					}
				} else {
					// insert
					_, err := stmtUserInsert.Exec(username, m.AvatarURL, displayName, userinfo.Email, userinfo.Created, userinfo.Updated)
					if err != nil {
						db.Close()
						return nil, err
					}
				}
			}
		}

		// Update when last checked
		stmtCache, err := db.Prepare("UPDATE githubusers SET updated_at = ? WHERE login = '__cachestate__'")
		if err != nil {
			db.Close()
			return nil, err
		}
		_, err = stmtCache.Exec(updated_at)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	// Users in the database reflect the current state
	rowsUsers, err := db.Query("select login, name, email, avatar_url from githubusers order by name asc")
	if err != nil {
		db.Close()
		return nil, err
	}
	login := ""
	name := ""
	email := ""
	avatar_url := ""
	for rowsUsers.Next() {
		err = rowsUsers.Scan(&login, &name, &email, &avatar_url)
		if err != nil {
			db.Close()
			return nil, err
		}
		if login != "__cachestate__" {
			response.Users = append(response.Users, &pb.UserResponse{Username: login, DisplayName: name, Email: email, AvatarUrl: avatar_url})
		}
	}
	rowsUsers.Close()
	db.Close()

	return &response, nil
}
