package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rjeczalik/gh/webhook"
)

var (
	secret  = flag.String("secret", "", "GitHub webhook secret")
	token   = flag.String("token", "", "Slack API token")
	channel = flag.String("channel", "", "Slack channel name")
	base    = "https://slack.com/api/chat.postMessage?token=%s&channel=%s&text=%s"
)

type slack struct{}

func (s slack) queryURL(message string) string {
	escaped := url.QueryEscape(message)
	return fmt.Sprintf(base, *token, *channel, escaped)
}

func (s slack) Push(e *webhook.PushEvent) {
	message := fmt.Sprintf("%s pushed to <%s|%s>", e.Pusher.Email, e.Repository.URL, e.Repository.Name)

	_, err := http.Get(s.queryURL(message))
	if err != nil {
		log.Println(err)
	}
}

func (s slack) PullRequest(e *webhook.PullRequestEvent) {
	var message string

	switch e.Action {
	case "opened":
		message = fmt.Sprintf("[%s] opened a new pull request <%s|>", e.PullRequest.Head.Repo.FullName, e.PullRequest.HTMLURL)
	case "closed":
		message = fmt.Sprintf("[%s] closed a new pull request <%s|>", e.PullRequest.Head.Repo.FullName, e.PullRequest.HTMLURL)
	default:
		message = fmt.Sprintf("[%s] new action (%s) on pull request <%s|>", e.PullRequest.Head.Repo.FullName, e.Action, e.PullRequest.HTMLURL)
	}

	_, err := http.Get(s.queryURL(message))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(":8080", webhook.New(*secret, slack{})))
}
