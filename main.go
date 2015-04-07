package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type JIRAResponse struct {
	Issues []*JIRAIssue `json:"issues"`
}

type JIRAIssue struct {
	Key    string                 `json:"key"`
	Fields map[string]*JIRAStatus `json:"fields"`
}

type JIRAStatus struct {
	Name string `json:"name"`
}

func main() {
	project := flag.String("project", "", "The project ID in JIRA")
	url := flag.String("url", "", "The URL to your JIRA install")
	username := flag.String("username", "", "JIRA username")
	password := flag.String("password", "", "JIRA password")
	status := flag.String("status", "Live", "JIRA status to look for")

	flag.Parse()

	if len(*project) == 0 {
		log.Fatal("Project is required!")
	}

	if len(*url) == 0 {
		log.Fatal("URL is required!")
	}

	if len(*username) == 0 {
		log.Fatal("Username is required!")
	}

	if len(*password) == 0 {
		log.Fatal("Password is required!")
	}

	if len(*status) == 0 {
		log.Fatal("Status is required!")
	}

	// Get a list of branches...
	output, err := exec.Command("git", "branch").Output()

	if err != nil {
		log.Fatal(err)
	}

	out := string(output)

	// Split by new line
	branches := strings.Split(out, "\n")

	branchChecker := regexp.MustCompile(*project + "\\-\\d+")

	branchesToCheck := map[string]string{}

	branchKeys := []string{}
	for _, b := range branches {
		match := branchChecker.FindString(b)
		if len(match) > 0 {
			branchesToCheck[match] = strings.Replace(b, " ", "", -1)
			branchKeys = append(branchKeys, match)
		}
	}

	log.Println(branchesToCheck)

	req, err := http.NewRequest("GET", *url+"/rest/api/2/search?jql=key%20IN%20("+strings.Join(branchKeys, ",")+")%20AND%20status="+*status+"&fields=status", nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(*username, *password)

	client := http.Client{}

	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		fmt.Println("Invalid response code: ", response.StatusCode)
		log.Fatal("Cannot continue.")
	}
	log.Println(response.StatusCode)

	jresponse := &JIRAResponse{}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(jresponse)
	if err != nil {
		log.Fatal(err)
	}

	branchesToDelete := []string{}
	for _, i := range jresponse.Issues {
		branchesToDelete = append(branchesToDelete, branchesToCheck[i.Key])

	}

	log.Println(branchesToDelete)

	fmt.Sprintf("Deleting %d branches\n", len(branchesToDelete))
	for _, toDelete := range branchesToDelete {
		// Get a list of branches...
		fmt.Println("Running git branch -D " + toDelete)
		output, err := exec.Command("git", "branch", "-D", toDelete).CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			log.Fatal(err)
		}
		fmt.Println(string(output))
	}

	fmt.Println("Done!")

}
