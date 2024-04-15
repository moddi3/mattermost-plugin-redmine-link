package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"sync"

	"github.com/dlclark/regexp2"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

type IssueResponse struct {
	Issue Issue `json:"issue"`
}
type IssuesResponse struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	ID int `json:"id"`
	// Project             Project  `json:"project"`
	// Tracker             Tracker  `json:"tracker"`
	// Status              Status   `json:"status"`
	// Priority            Priority `json:"priority"`
	// Author              Person   `json:"author"`
	// AssignedTo          Person   `json:"assigned_to"`
	// Parent              *Parent  `json:"parent,omitempty"` // Optional field
	Subject             string   `json:"subject"`
	Description         string   `json:"description"`
	StartDate           string   `json:"start_date"`
	DueDate             *string  `json:"due_date,omitempty"` // Optional field
	DoneRatio           int      `json:"done_ratio"`
	IsPrivate           bool     `json:"is_private"`
	EstimatedHours      *float64 `json:"estimated_hours,omitempty"`       // Optional field
	TotalEstimatedHours *float64 `json:"total_estimated_hours,omitempty"` // Optional field
	SpentHours          float64  `json:"spent_hours"`
	TotalSpentHours     float64  `json:"total_spent_hours"`
	CreatedOn           string   `json:"created_on"`
	UpdatedOn           string   `json:"updated_on"`
	ClosedOn            *string  `json:"closed_on,omitempty"` // Optional field
}

// func getIssueName(issueID string) (string, error) {
// 	url := fmt.Sprintf("%s%s%s", p.getRedmineInstanceURL(), issueID, ".json")
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("X-Redmine-API-Key", "5225a5f42e854fca558358866d7d253631189cb8")
// 	resp, err := http.DefaultClient.Do(req) //nolint

// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var issueResponse IssueResponse
// 	err = json.NewDecoder(resp.Body).Decode(&issueResponse)
// 	if err != nil {
// 		return "", err
// 	}

//		return issueResponse.Issue.Subject, nil
//	}

func (p *Plugin) getRedmineInstanceURL() string {
	configuration := p.getConfiguration()
	if configuration.RedmineInstanceURL == "" {
		return ""
	}
	url, err := url.Parse(configuration.RedmineInstanceURL)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s/", url.Scheme, url.Host)
}

func (p *Plugin) getIssueNames(issueIDs []string) (map[string]string, error) {
	configuration := p.getConfiguration()

	idsParam := strings.Join(issueIDs, ",")
	// https://www.redmine.org/issues.json?issue_id=1,2,3&status_id=*
	url := fmt.Sprintf("%s%s?issue_id=%s&status_id=*", p.getRedmineInstanceURL(), "issues.json", idsParam)

	req, _ := http.NewRequest("GET", url, nil)

	if configuration.RedmineAPIKey != "" {
		req.Header.Set("X-Redmine-API-Key", configuration.RedmineAPIKey)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}

	defer resp.Body.Close()

	var issuesResponse IssuesResponse
	err = json.NewDecoder(resp.Body).Decode(&issuesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	issuesSubjectMap := make(map[string]string)
	for _, issue := range issuesResponse.Issues {
		issueIDStr := fmt.Sprintf("%d", issue.ID)

		issuesSubjectMap[issueIDStr] = issue.Subject
	}
	return issuesSubjectMap, nil
}

func (p *Plugin) extractTrackerLinks(input string) []string {
	var matches []string
	escapedURL := strings.ReplaceAll(regexp.QuoteMeta(p.getRedmineInstanceURL()), "/", "\\/")
	pattern := regexp2.MustCompile(`(?<!\]\()`+escapedURL+`issues\/\d+(?![^\[]*\])`, 0)

	match, _ := pattern.FindStringMatch(input)

	for match != nil {
		matches = append(matches, match.String())
		match, _ = pattern.FindNextMatch(match)
	}

	return matches
}

func getIssueIDFromLink(link string, url string) string {
	return strings.TrimPrefix(link, url+"issues/")
}

func (p *Plugin) transformMessageLinks(message string, links []string) string {
	if len(links) == 0 {
		return message
	}

	var transformedParts []string
	startIndex := 0
	issueIDs := make([]string, 0, len(links))

	// Collect issue IDs from links
	for _, link := range links {
		issueID := getIssueIDFromLink(link, p.getRedmineInstanceURL())
		issueIDs = append(issueIDs, issueID)
	}

	// Get issue names for all issue IDs in a single API request
	issuesSubjectMap, err := p.getIssueNames(issueIDs)

	if err != nil {
		// If there is an error fetching issue names, return the original message
		return message
	}

	// Transform message links based on the fetched issue names
	for i, link := range links {
		linkIndex := strings.Index(message[startIndex:], link)
		if linkIndex == -1 {
			continue
		}

		linkIndex += startIndex
		transformedParts = append(transformedParts, message[startIndex:linkIndex])

		issueName := issuesSubjectMap[issueIDs[i]]

		if issueName == "" {
			// If issue name is not found, use the original link
			transformedParts = append(transformedParts, link)
		} else {
			// Create transformed link with issue name
			transformedLink := fmt.Sprintf("[%s](%s)", issueName, link)
			transformedParts = append(transformedParts, transformedLink)
		}

		// Update start index for the next iteration
		startIndex = linkIndex + len(link)
	}

	// Append remaining part of the message
	transformedParts = append(transformedParts, message[startIndex:])

	return strings.Join(transformedParts, "")
}

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	newPost := post.Clone()

	if p.getRedmineInstanceURL() != "" {
		newPost.Message = p.transformMessageLinks(newPost.Message, p.extractTrackerLinks(newPost.Message))
	}
	return newPost, ""
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
