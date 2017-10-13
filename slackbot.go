package bitbucketbot

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)

const (
	slackAPIKey     = "1"
	bitbucketAPIKey = "1"
)

type eventHandler func(payload interface{})

type webHook struct {
	eventTypes map[string]eventHandler
}

// Configuration used for slack api and bitbucket api
type Configuration struct {
	SlackAPIKey     string `json:"SlackAPIKey"`
	BitBucketAPIKey string `json:"BitBucketAPIKey"`
}

var users = map[string]string{}

func main() {
	slackClient := slack.New(slackAPIKey)

	var err = loadUsers(slackClient)
	if err != nil {
		log.Fatal("Failed to load slack users: " + err.Error())
	}

	r := mux.NewRouter()
	// r.HandleFunc("/bitbucket", bitbucketHandler)
	http.Handle("/", r)
}

// func bitbucketHandler(w http.ResponseWriter, r *http.Request) {

// }

// func handleNewPullRequest(payload payload.PullRequestCreatedPayload, slack slack.Client) {
// 	users, err := slack.GetUsers()
// }

func loadUsers(slack *slack.Client) (err error) {
	slackUsers, err := slack.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range slackUsers {
		users[user.ID] = user.Profile.Email
	}

	return nil
}
