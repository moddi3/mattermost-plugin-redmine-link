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

func parseLink(link string) (map[string]string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "https"
		u, _ = url.Parse(u.String())
	}

	return map[string]string{
		"Scheme": u.Scheme,
		"Host":   u.Hostname(),
		"Path":   u.Path,
		"Hash":   u.Fragment,
	}, nil
}

func extractTrackerLinks(input string, redmineHost string) []string {
	var matches []string

	pattern := `(?<!\]\()(?:https?:\/\/|(?<!\S)|(?<!\W))` + regexp.QuoteMeta(redmineHost) + `\/issues\/\d+(?:#note-\d+)?(?![^\[]*\])`
	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.FindStringMatch(input)

	for match != nil {
		matches = append(matches, match.String())
		match, _ = re.FindNextMatch(match)
	}
	return matches
}

// i need to get necessary data from issues response and transform it to to map
// issueID -> { issueName, issueStatus, issueNoteAnchor }
func processIssuesResponse(issuesResponse IssuesResponse) map[string]map[string]string {
	issuesMap := make(map[string]map[string]string)

	for _, issue := range issuesResponse.Issues {
		issueID := fmt.Sprintf("%d", issue.ID)
		issuesMap[issueID] = map[string]string{
			"Subject":    issue.Subject,
			"Status":     issue.Status.Name,
			"Tracker":    issue.Tracker.Name,
			"AssignedTo": issue.AssignedTo.Name,
		}
	}
	return issuesMap
}

func (p *Plugin) getRedmineInstanceURL() (string, string) {
	configuration := p.getConfiguration()
	if configuration.RedmineInstanceURL == "" {
		return "", ""
	}
	parsedURL, err := parseLink(configuration.RedmineInstanceURL)
	if err != nil {
		return "", ""
	}
	return fmt.Sprintf("%s://%s/", parsedURL["Scheme"], parsedURL["Host"]), parsedURL["Host"]
}

func (p *Plugin) getIssuesData(issueIDs []string) (map[string]map[string]string, error) {
	configuration := p.getConfiguration()

	idsParam := strings.Join(issueIDs, ",")
	// https://www.redmine.org/issues.json?issue_id=1,2,3&status_id=*
	redmineURL, _ := p.getRedmineInstanceURL()
	url := fmt.Sprintf("%s%s?issue_id=%s&status_id=*", redmineURL, "issues.json", idsParam)

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

	return processIssuesResponse(issuesResponse), nil
}

func (p *Plugin) transformMessageLinks(message string, links []string) string {
	if len(links) == 0 {
		return message
	}

	var transformedParts []string
	startIndex := 0
	issuesIDs := make([]string, 0, len(links))
	issuesHashes := make([]string, 0, len(links))

	// Collect issue IDs from links
	for _, link := range links {
		parsedLink, _ := parseLink(link)
		issueID := strings.TrimPrefix(parsedLink["Path"], "/issues/")

		issuesIDs = append(issuesIDs, issueID)
		issuesHashes = append(issuesHashes, parsedLink["Hash"])
	}

	// Get issue names for all issue IDs in a single API request
	issuesData, err := p.getIssuesData(issuesIDs)

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

		issueData := issuesData[issuesIDs[i]]

		if issueData["Subject"] == "" {
			// If issue subject is not found, use the original link
			transformedParts = append(transformedParts, link)
		} else {
			// Create transformed link with issue subject
			hash := ""
			if issuesHashes[i] != "" {
				hash = "#" + issuesHashes[i]
			}
			transformedLink := fmt.Sprintf("[%s%s](%s)", issueData["Subject"], hash, link)
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
	redmineURL, redmineHost := p.getRedmineInstanceURL()

	if redmineURL != "" {
		newPost.Message = p.transformMessageLinks(newPost.Message, extractTrackerLinks(newPost.Message, redmineHost))
	}
	return newPost, ""
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
