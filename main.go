package updater

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type GitHubRefsPayload struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

type GitHubGitPayload struct {
	Sha     string `json:"sha"`
	NodeID  string `json:"node_id"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
	Author  struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"committer"`
	Tree struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"tree"`
	Message string `json:"message"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
	Verification struct {
		Verified  bool        `json:"verified"`
		Reason    string      `json:"reason"`
		Signature interface{} `json:"signature"`
		Payload   interface{} `json:"payload"`
	} `json:"verification"`
}

func getRefs(repo string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/git/refs/heads/master", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyText
}

func getGit(repo, commitSha string) []byte {
	client := &http.Client{}
	// commitSha = "d81240518f62a06e1e49dfdba09b1eb4b54cba42"
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/git/commits/"+commitSha, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return bodyText
}

func getCurrentStoredSha(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

// func main() {
// 	argsWithoutProg := os.Args[1:]

// 	if len(argsWithoutProg) != 2 {
// 		log.Fatal("Incorrect number of arguments. Try [repo, comand]")
// 	}

// 	repo := argsWithoutProg[0]
// 	command := argsWithoutProg[1]

// 	storePath := "log"

// 	switch command {
// 	case "check":
// 		currentParentSha := getCurrentStoredSha(storePath)
// 		fmt.Println(currentParentSha)
// 	case "parent":
// 		sha := getParentCommit(repo)
// 		fmt.Println(sha)
// 	case "compare":
// 		currentParentSha := getCurrentStoredSha(storePath)
// 		sha := getParentCommit(repo)
// 		shouldUpdate := compareVersionsBySha(sha, currentParentSha)
// 		fmt.Println(shouldUpdate)
// 	default:
// 		fmt.Println("No Command Recognized")
// 	}
// }

func compareVersionsBySha(sha, currentSha string) bool {
	if sha == currentSha {
		return false
	}
	return true
}

func getParentCommit(repo string) string {

	bodyText := getRefs(repo)

	var payload GitHubRefsPayload
	err := json.Unmarshal(bodyText, &payload)
	if err != nil {
		log.Fatal(err)
	}

	commitSha := payload.Object.Sha
	data := getGit(repo, commitSha)

	var gitPayload GitHubGitPayload
	err = json.Unmarshal(data, &gitPayload)
	if err != nil {
		log.Fatal(err)
	}

	if len(gitPayload.Parents) == 0 {
		return "root"
	}
	return gitPayload.Parents[0].Sha

}

// Public functions

// CheckCurrentStoredVersion checks the stored value
func CheckCurrentStoredVersion(storePath string) string {
	currentParentSha := getCurrentStoredSha(storePath)
	return currentParentSha
}

// CheckCurrentGithubParent checks the Github value
func CheckCurrentGithubParent(repo string) string {
	sha := getParentCommit(repo)
	return sha
}

// CompareStoredVerionAndGihubVersion compares the stored and fetched values
func CompareStoredVerionAndGihubVersion(repo, storePath string) bool {
	currentParentSha := getCurrentStoredSha(storePath)
	sha := getParentCommit(repo)
	shouldUpdate := compareVersionsBySha(sha, currentParentSha)
	return shouldUpdate
}
