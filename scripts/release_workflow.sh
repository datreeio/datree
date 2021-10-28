actor=$1
branch_name=$2
commit_sha=$3
version=$4

echo "Triggered by: $actor!"

echo "branch name: $branch_name"
# if [ "$branch_name" != "main" ]; then
#     echo "Release should build only from main branch"
#     exit 1
# fi
if [ -z "$TRAVIS_API_TOKEN" ]; then
    echo "TRAVIS_API_TOKEN is empty"
    exit 1
fi

BODY="
{
  \"request\": {
    \"branch\": \"$branch_name\",
    \"sha\": \"$commit_sha\",
    \"merge_mode\": \"deep_merge_append\",
    \"config\": {
      \"env\": {
        \"global\": [
          \"RELEASE_DATREE_PROD=true\",
          \"RELEASE_CANDIDATE_VERSION=$version\"
        ]
      }
    }
  }
}
"
GH_ORGANIZATION_NAME="datreeio"
GH_REPOSITORY_NAME="datree"

curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -H "Travis-API-Version: 3" \
    -H "Authorization: token ${TRAVIS_API_TOKEN}" \
    -d "$BODY" \
    "https://api.travis-ci.com/repo/${GH_ORGANIZATION_NAME}%2F${GH_REPOSITORY_NAME}/requests"
