 1426  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/rfielding/events
 1427  curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o "github-api-responses/users/rfielding/events.json" https://api.github.com/users/rfielding/events
 1428  mkdir -p github-api-responses/user_rfielding
 1429  curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o "github-api-responses/user_rfielding/events.json" https://api.github.com/users/rfielding/events
 1430  mkdir -p github-api-responses/user_rfielding/org_deciphernow
 1431  curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o "github-api-responses/user_rfielding/org_deciphernow/events.json" https://api.github.com/users/rfielding/events/orgs/deciphernow
 1432  mkdir -p github-api-responses/user_lucasmoten/org_deciphernow
 1433  curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o "github-api-responses/user_lucasmoten/org_deciphernow/events.json" https://api.github.com/users/lucasmoten/events/orgs/deciphernow
 1434  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-data-aac/commits
 1435  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-data-aac/events
 1436  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-data-aac/contributors
 1437  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-data/contributors
 1438  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-data-static-pages/contributors
 1439  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/contributors
 1440  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/rfielding/events
 1441  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/rfielding/events | grep DecipherNow
 1442  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/7ruth/events | grep DecipherNow
 1443  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/7ruth/events | grep name | grep DecipherNow
 1444  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/shanberg/events | grep name | grep DecipherNow
 1445  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/users/shanberg/events
 1446  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/events
 1447  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/commits
 1448  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/comments
 1449  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/assignees
 1450  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/tags
 1451  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/trees
 1452  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/git/tags
 1453  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/git/trees
 1454  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/statuses
 1455  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/contributors
 1456  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/labels
 1457  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/deployments
 1458  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/releases
 1459  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/object-drive-server/releases
 1460  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/object-drive-server/teams
 1461  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/object-drive-server/issues
 1462  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/object-drive-server/issues | more
 1463  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/issues | more
 1464  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/pulls | more
 1465  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/milestones | more
 1466  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/comments | more
 1467  curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/DecipherNow/gm-static-pages/branches | more
 1468  history