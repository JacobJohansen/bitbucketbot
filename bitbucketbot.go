package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
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

var config = Configuration{}
var users = map[string]string{}

func main() {

	slackClient := slack.New(config.SlackAPIKey)

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
func loadConfiguration() {
	file, e := ioutil.ReadFile("./configuration.json")
	if e != nil {
		log.Fatal("Could not read ./configuration.json: " + e.Error())
	}

	var config Configuration
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatal("Failed load configuration: " + err.Error())
	}
}

func loadUsers(slack *slack.Client) (err error) {
	slackUsers, err := slack.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range slackUsers {
		users[user.Profile.Email] = user.ID
	}

	return nil
}

func sendUserMessage(email, message string, client *slack.Client) {
	_, _, channelID, err := client.OpenIMChannel(users[email])
	if err != nil {
		log.Println("Couldn't find user to send message: " + email)
		return
	}

	client.PostMessage(channelID, message, slack.PostMessageParameters{})
}
