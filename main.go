package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"time"
)

func main() {
	timeout := os.Getenv("WERCKER_WAIT_GITHUB_STATUSES_TIMEOUT")
	githubToken := os.Getenv("WERCKER_WAIT_GITHUB_STATUSES_GITHUB_TOKEN")
	statusContextsList := os.Getenv("WERCKER_WAIT_GITHUB_STATUSES_STATUS_CONTEXTS")
	gitOwner := os.Getenv("WERCKER_GIT_OWNER")
	gitRepository := os.Getenv("WERCKER_GIT_REPOSITORY")
	gitCommit := os.Getenv("WERCKER_GIT_COMMIT")

	if githubToken == "" {
		fmt.Fprintf(os.Stderr, "No GitHub personal authentication token is configured, please set github_token in wercker.yml\n")
		os.Exit(1)
	}

	if statusContextsList == "" {
		fmt.Fprintf(os.Stderr, "No statuses are configured, please set status_contexts in wercker.yml\n")
		os.Exit(1)
	}

	timeoutMinutes := 25
	if timeout != "" {
		fmt.Scanf(timeout, "%d", &timeoutMinutes)
	}

	fmt.Printf("Waiting for successful outcome on %s/%s@%s for status checks (will wait for up to %d minutes): ",
		gitOwner, gitRepository, gitCommit, timeoutMinutes)

	timer := time.NewTimer(time.Duration(timeoutMinutes) * time.Minute)
	go func() {
		<-timer.C
		fmt.Fprintf(os.Stderr, "\nTimed out!")
		os.Exit(1)
	}()

	statusContextsNeeded := make(map[string]string)
	for _, c := range strings.Split(statusContextsList, ",") {
		statusContextsNeeded[strings.TrimSpace(c)] = ""
		fmt.Printf("%s ", c)
	}
	fmt.Printf("\n")

	for len(statusContextsNeeded) > 0 {
		fmt.Printf(".")

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)

		// get the statuses for a commit
		opt := &github.ListOptions{}
		statuses, _, err := client.Repositories.GetCombinedStatus(gitOwner, gitRepository, gitCommit, opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%s\n", err)
			os.Exit(1)
		}

		statusMap := make(map[string]string)
		for _, s := range statuses.Statuses {
			statusMap[*s.Context] = *s.State
		}

		for k := range statusContextsNeeded {
			if statusMap[k] == "success" {
				fmt.Printf("\n%s succeeded", k)
				delete(statusContextsNeeded, k)
			}
			if statusMap[k] == "failure" {
				fmt.Fprintf(os.Stderr, "\n%s failed!\n", k)
				os.Exit(1)
			}
		}

		if len(statusContextsNeeded) > 0 {
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Printf("\n")
}
