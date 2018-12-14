package methods

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"time"

	gc "github.com/lucasmoten/project-2502/services/dice/apis/github"
	pb "github.com/lucasmoten/project-2502/services/dice/protobuf"
)

func (s *serverData) ListRepositories(ctx context.Context, request *pb.ListRepositoriesRequest) (*pb.ListRepositoriesResponse, error) {
	var response pb.ListRepositoriesResponse
	db, err := sql.Open("sqlite3", "./dice.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	dbRows, err := db.Query("select updated_at from repositories order by updated_at desc limit 1")

	if err != nil {
		db.Close()
		return nil, err
	}

	updated_at := "2018-12-13T00:00:00Z" // just a default in the past

	for dbRows.Next() {
		dbRows.Scan(&updated_at)
	}
	dbRows.Close()

	then, err := time.Parse(time.RFC3339, updated_at)

	if err != nil {
		db.Close()
		return nil, err
	}

	duration := time.Since(then)
	dataisold := (duration.Hours() > 2)

	if dataisold {
		updated_at = time.Now().Format(time.RFC3339)

		githubclient, err := gc.NewClient()
		if err != nil {
			db.Close()
			return nil, err
		}

		pageHasResults := true
		pageNumber := 0

		stmtRepoUpdate, err := db.Prepare("update repositories set node_id=?, name=?, " +
			"full_name=?, description=?, language=?, default_branch=?, created_at=?, pushed_at=?, updated_at=?, fork=?, " +
			"private=?, archived=?, forks_count=?, network_count=?, open_issues_count=?, " +
			"stargazers_count=?, subscribers_count=?, watchers_count=?, size=? where id=?")

		if err != nil {
			db.Close()
			return nil, err
		}

		stmtRepoInsert, err := db.Prepare("insert into repositories (id, node_id, name," +
			" full_name, description, language, default_branch, created_at, pushed_at, updated_at, fork," +
			" private, archived, forks_count, network_count, open_issues_count," +
			" stargazers_count, subscribers_count, watchers_count, size) values" +
			" (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

		if err != nil {
			db.Close()
			return nil, err
		}

		for pageHasResults {
			pageNumber++

			repositories, err := githubclient.GetRepositories(pageNumber)

			if err != nil {
				db.Close()
				return nil, err
			}

			pageHasResults = len(repositories) > 0

			for _, repo := range repositories {
				var rowsCount sql.NullString
				err = db.QueryRow("select count(0) from repositories where id=?", repo.ID).Scan(&rowsCount)

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
					_, err := stmtRepoUpdate.Exec(repo.NodeID, repo.Name, repo.FullName,
						repo.Description, repo.Language, repo.DefaultBranch, repo.CreatedAt, repo.PushedAt, repo.UpdatedAt, repo.Fork,
						repo.Private, repo.Archived, repo.ForksCount, repo.NetworkCount, repo.OpenIssuesCount,
						repo.StargazersCount, repo.SubscribersCount, repo.WatchersCount, repo.Size)
					if err != nil {
						db.Close()
						return nil, err
					}
				} else {
					// insert
					_, err := stmtRepoInsert.Exec(repo.ID, repo.NodeID, repo.Name, repo.FullName,
						repo.Description, repo.Language, repo.DefaultBranch, repo.CreatedAt, repo.PushedAt, repo.UpdatedAt, repo.Fork,
						repo.Private, repo.Archived, repo.ForksCount, repo.NetworkCount, repo.OpenIssuesCount,
						repo.StargazersCount, repo.SubscribersCount, repo.WatchersCount, repo.Size)
					if err != nil {
						db.Close()
						return nil, err
					}
				}

			}

			var rowsCount sql.NullString
			err = db.QueryRow("select count(0) from repositories where id=?", 0).Scan(&rowsCount)
			if rowsCount.String == "1" {
				// Update when last checked
				stmtCache, err := db.Prepare("UPDATE repositories SET updated_at = ? WHERE name = '__cachestate__'")
				if err != nil {
					db.Close()
					return nil, err
				}
				_, err = stmtCache.Exec(updated_at)
				if err != nil {
					db.Close()
					return nil, err
				}
			} else {
				_, err = stmtRepoInsert.Exec("0", "TEST", "__cachestate__", "TEST", "TEST", "TEST", "TEST", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z", "2018-12-13T00:00:00Z", false, false, false, 1, 1, 1, 1, 1, 1, 1)
				if err != nil {
					db.Close()
					return nil, err
				}
			}
		}
	}
	// Users in the database reflect the current state
	rowsRepositories, err := db.Query("select  name, full_name, description, " +
		"language, default_branch, created_at, pushed_at, updated_at, fork, private, archived, " +
		"forks_count, network_count, open_issues_count, stargazers_count, subscribers_count, " +
		"watchers_count, size from repositories order by updated_at")
	if err != nil {
		db.Close()
		return nil, err
	}
	name := ""
	fullName := ""
	description := ""
	language := ""
	defaultBranch := ""
	createdTime := ""
	pushedTime := ""
	updatedTime := ""
	fork := false
	private := false
	archived := false
	forksCount := int64(0)
	networkCount := int64(0)
	openIssuesCount := int64(0)
	stargazersCount := int64(0)
	subscribersCount := int64(0)
	watchersCount := int64(0)
	size := int64(0)

	openIssueCount := int64(0)
	for rowsRepositories.Next() {
		err = rowsRepositories.Scan(&name, &fullName, &description, &language, &defaultBranch,
			&createdTime, &pushedTime, &updatedTime, &fork, &private, &archived,
			&forksCount, &networkCount, &openIssuesCount,
			&stargazersCount, &subscribersCount, &watchersCount, &size)
		if err != nil {
			db.Close()
			return nil, err
		}
		if name != "__cachestate__" {
			response.Repositories = append(response.Repositories, &pb.RepositoryResponse2{
				Name:             name,
				FullName:         fullName,
				Description:      description,
				Language:         language,
				DefaultBranch:    defaultBranch,
				CreatedAt:        createdTime,
				PushedAt:         pushedTime,
				UpdatedAt:        updatedTime,
				Fork:             fork,
				Private:          private,
				Archived:         archived,
				ForksCount:       forksCount,
				NetworkCount:     networkCount,
				OpenIssuesCount:  openIssueCount,
				StargazersCount:  stargazersCount,
				SubscribersCount: subscribersCount,
				WatchersCount:    watchersCount,
				Size:             size,
			})
		}
	}

	return &response, nil
}
