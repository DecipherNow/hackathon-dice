#!/bin/bash

orgname="deciphernow"
cachefolder="github-api-responses"

# Verify environment variable set
COUNT_GITHUB_TOKEN=$(printenv | grep "GITHUB_TOKEN" -c)
if [ $COUNT_GITHUB_TOKEN -eq 0 ]
then
  echo "envvar GITHUB_TOKEN is not set"
  exit 1
fi
# Verify jq is installed
r=($(which jq|wc -l))
if [ $r -eq 0 ]; then
    echo "This test script relies upon jq, a command line JSON processor to run."
    echo "See guidance here for installation in your environment"
    echo "https://stedolan.github.io/jq/download/"
    exit 1
fi
mkdir -p ${cachefolder}

# ==========================================================
# FUNCTIONS
drawline() {
  printf '%70s\n' | tr ' ' -
}
checkfile() {
    fetchfile=0
    if [ ! -f $1 ]; then
        fetchfile=1
    elif test `file "$1" -mmin +120`; then
        fetchfile=1
    fi
}
getorgpage() {
    checkfile ${cachefolder}/org_${orgname}.json
    if [ $fetchfile -gt 0 ]; then
        r=($(curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o ${cachefolder}/org_${orgname}.json https://api.github.com/orgs/${orgname}))
    fi
}
getmemberpage() {
    a=${@}
    pagenumber=1
    if [ $a -gt 0 ]; then
        pagenumber=$1
    fi
    
    #checkfile ${cachefolder}/org_${orgname}/members_${pagenumber}.json
    fetchfile=0
    if [ ! -f ${cachefolder}/org_${orgname}/members_${pagenumber}.json ]; then
        fetchfile=1
    elif test `find "${cachefolder}/org_${orgname}/members_${pagenumber}.json" -mmin +120`; then
        fetchfile=1
    fi
    if [ $fetchfile -gt 0 ]; then
        r=($(curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o ${cachefolder}/org_${orgname}/members_${pagenumber}.json https://api.github.com/orgs/${orgname}/members?page=${pagenumber}))
    fi
}
getuserpage() {
    fetchfile=0
    if [ ! -f ${cachefolder}/users_$1.json ]; then
        fetchfile=1
    elif test `find "${cachefolder}/users_$1.json" -mmin +120`; then
        fetchfile=1
    fi
    if [ $fetchfile -gt 0 ]; then
        r=($(curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o ${cachefolder}/users_$1.json https://api.github.com/users/$1))
    fi
}
getrepopage() {
    a=${@}
    pagenumber=1
    if [ $a -gt 0 ]; then
        pagenumber=$1
    fi
    
    #checkfile ${cachefolder}/org_${orgname}/repos_${pagenumber}.json
    fetchfile=0
    if [ ! -f ${cachefolder}/org_${orgname}/repos_${pagenumber}.json ]; then
        fetchfile=1
    elif test `find "${cachefolder}/org_${orgname}/repos_${pagenumber}.json" -mmin +120`; then
        fetchfile=1
    fi
    if [ $fetchfile -gt 0 ]; then
        r=($(curl -s -H "Authorization: token ${GITHUB_TOKEN}" -o ${cachefolder}/org_${orgname}/repos_${pagenumber}.json https://api.github.com/orgs/${orgname}/repos?page=${pagenumber}))
    fi
}

# ==========================================================
# API REPORTING
# ==========================================================
# This gives the endpoints available
#   curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com


# Get Organization Info
mkdir -p ${cachefolder}/org_${orgname}
getorgpage
jq -r '.name' ${cachefolder}/org_${orgname}.json
jq -r '.description' ${cachefolder}/org_${orgname}.json
jq -r '.location' ${cachefolder}/org_${orgname}.json
jq -r '(.total_private_repos|tostring) + " / " + (.plan.private_repos|tostring) + " private repos"' ${cachefolder}/org_${orgname}.json
drawline

# Get Members of Organization
printf 'Members\n'
drawline
memberpage=1
getmemberpage ${memberpage}
jqc=1
while [ $jqc -gt 0 ]
do
    for row in $(cat ${cachefolder}/org_${orgname}/members_${memberpage}.json | jq -r '.[] | @base64'); do
        _jq() {
            echo $row | base64 --decode | jq -r ${1}
        }
        member_login=$(_jq '.login')
        member_type=$(_jq '.type')
        getuserpage ${member_login}
        member_name=$(jq -r '.name' ${cachefolder}/users_${member_login}.json)
        printf '%20s | %30s | %10s \n' "$member_login" "$member_name" "$member_type"
        #echo $(_jq '.login + \\" - \\" + .type')
        #jq -r '.login' ${cachefolder}/org_${orgname}/members_${memberpage}.json
    done
    # get data for next page
    memberpage=$((memberpage+1))
    getmemberpage ${memberpage}
    jqc=($(jq '. | length' ${cachefolder}/org_${orgname}/members_${memberpage}.json))
done
drawline

# Get Repositories of Organization
printf 'Repositories\n'
drawline
repopage=1
getrepopage ${repopage}
jqc=1
while [ $jqc -gt 0 ]
do
    for row in $(cat ${cachefolder}/org_${orgname}/repos_${repopage}.json | jq -r '.[] | @base64'); do
        _jq() {
            echo $row | base64 --decode | jq -r ${1}
        }
        repo_name=$(_jq '.name')
        repo_description=$(_jq '.description')
        repo_license=$(_jq '.license.name')
        repo_openissues=$(_jq '.open_issues')
        repo_private=$(_jq '.private')
        repo_fork=$(_jq '.fork')
        #getissuespage ${repo_name}
        printf 'Name:        %s\n' "$repo_name"
        printf 'Description: %s\n' "$repo_description"
        printf '             License: %20s   Is Private: %5s   Is Fork: %5s\n' "$repo_license" "${repo_private}" "${repo_fork}" 
        printf '             Open Issues: %d\n' "$repo_openissues"
        printf '\n\n'
    done
    # get data for next page
    repopage=$((repopage+1))
    getrepopage ${repopage}
    jqc=($(jq '. | length' ${cachefolder}/org_${orgname}/repos_${repopage}.json))
done
drawline


exit 0


# Issues for a repository
curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/deciphernow/object-drive-server/issues
curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/deciphernow/object-drive-server/issues?page=2
[{},{}]
.number: number
.title: string
.state: open, closed ? 
.created_at: datetime
.updated_at: datetime
.body: string

# Releases for a repository
curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/deciphernow/object-drive-server/releases
curl -s -H "Authorization: token ${GITHUB_TOKEN}" https://api.github.com/repos/deciphernow/object-drive-server/releases?page=2
[{},{}]
.tag_name: string
.name: string
.created_at: iso8601 datetime
.body: string