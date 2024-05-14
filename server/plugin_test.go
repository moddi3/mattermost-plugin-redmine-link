package main

import (
	"fmt"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Description     string
	InputMessage    string
	ExpectedMessage string
	OldMessage      string
}

func TestMessagehooks(t *testing.T) {
	plugin := &Plugin{
		configuration: &configuration{
			RedmineInstanceURL: "https://www.redmine.org",
		},
	}

	// Nested map to group test cases for each hook
	testCases := map[string][]TestCase{
		"MessageWillBePosted": {
			{
				Description:     "No tracker links",
				InputMessage:    "This is a test message without any tracker links.",
				ExpectedMessage: "This is a test message without any tracker links.",
			},
			{
				Description:  "Single tracker link",
				InputMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/40556",
				ExpectedMessage: fmt.Sprintf("This is a test message with a tracker link: %s", createTransformedLink(
					"Focus on the textarea after clicking the Edit Journal button",
					"https://www.redmine.org/issues/40556",
					"",
					map[string]string{
						"Subject":    "Focus on the textarea after clicking the Edit Journal button",
						"Status":     "Closed",
						"Priority":   "Normal",
						"UpdatedOn":  "2024-04-29T19:23:49Z",
						"AssignedTo": "Marius BĂLTEANU",
						"Tracker":    "Feature",
						"ID":         "40556",
						"Author":     "Yasu Saku",
					},
				)),
			},
			{
				Description:  "Multiple tracker links",
				InputMessage: "This is a test message with multiple tracker links: https://www.redmine.org/issues/40556 and https://www.redmine.org/issues/40559",
				ExpectedMessage: fmt.Sprintf("This is a test message with multiple tracker links: %s and %s",
					createTransformedLink(
						"Focus on the textarea after clicking the Edit Journal button",
						"https://www.redmine.org/issues/40556",
						"",
						map[string]string{
							"Subject":    "Focus on the textarea after clicking the Edit Journal button",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-29T19:23:49Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Feature",
							"ID":         "40556",
							"Author":     "Yasu Saku",
						},
					),
					createTransformedLink(
						"Fix incorrect icon image paths for Wiki help pages",
						"https://www.redmine.org/issues/40559",
						"",
						map[string]string{
							"Subject":    "Fix incorrect icon image paths for Wiki help pages",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-16T19:26:17Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Patch",
							"ID":         "40559",
							"Author":     "Katsuya HIDAKA",
						},
					),
				),
			},
			{
				Description:  "Multiple tracker links with markdown links",
				InputMessage: "This is a test message with multiple tracker links: [a link](https://www.redmine.org/issues/40556) and https://www.redmine.org/issues/40559 and https://www.redmine.org/issues/999999 and https://www.redmine.org/issues/40559",
				ExpectedMessage: fmt.Sprintf("This is a test message with multiple tracker links: [a link](https://www.redmine.org/issues/40556) and %s and %s and %s",
					createTransformedLink(
						"Fix incorrect icon image paths for Wiki help pages",
						"https://www.redmine.org/issues/40559",
						"",
						map[string]string{
							"Subject":    "Fix incorrect icon image paths for Wiki help pages",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-16T19:26:17Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Patch",
							"ID":         "40559",
							"Author":     "Katsuya HIDAKA",
						},
					),
					"https://www.redmine.org/issues/999999",
					createTransformedLink(
						"Fix incorrect icon image paths for Wiki help pages",
						"https://www.redmine.org/issues/40559",
						"",
						map[string]string{
							"Subject":    "Fix incorrect icon image paths for Wiki help pages",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-16T19:26:17Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Patch",
							"ID":         "40559",
							"Author":     "Katsuya HIDAKA",
						},
					),
				),
			},
			{
				Description:  "Tracker links with http protocol and without protocol",
				InputMessage: "This is a test message with an http tracker link: http://www.redmine.org/issues/40556 and without protocol www.redmine.org/issues/40556",
				ExpectedMessage: fmt.Sprintf("This is a test message with an http tracker link: %s and without protocol %s",
					createTransformedLink(
						"Focus on the textarea after clicking the Edit Journal button",
						"http://www.redmine.org/issues/40556",
						"",
						map[string]string{
							"Subject":    "Focus on the textarea after clicking the Edit Journal button",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-29T19:23:49Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Feature",
							"ID":         "40556",
							"Author":     "Yasu Saku",
						},
					),
					createTransformedLink(
						"Focus on the textarea after clicking the Edit Journal button",
						"www.redmine.org/issues/40556",
						"",
						map[string]string{
							"Subject":    "Focus on the textarea after clicking the Edit Journal button",
							"Status":     "Closed",
							"Priority":   "Normal",
							"UpdatedOn":  "2024-04-29T19:23:49Z",
							"AssignedTo": "Marius BĂLTEANU",
							"Tracker":    "Feature",
							"ID":         "40556",
							"Author":     "Yasu Saku",
						},
					),
				),
			},
			{
				Description:     "Error retrieving issue name",
				InputMessage:    "This is a test message with a tracker link: https://www.redmine.org/issues/999999",
				ExpectedMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/999999",
			},
			{
				Description:  "Tracker link with a note anchor",
				InputMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/40538#note-4",
				ExpectedMessage: fmt.Sprintf("This is a test message with a tracker link: %s", createTransformedLink(
					"Hi, can you help me with a Version Extended?",
					"https://www.redmine.org/issues/40538#note-4",
					"#note-4",
					map[string]string{
						"Subject":   "Hi, can you help me with a Version Extended?",
						"Status":    "Reopened",
						"Priority":  "Normal",
						"UpdatedOn": "2024-04-09T10:42:41Z",
						"Tracker":   "Patch",
						"ID":        "40538",
						"Author":    "Enzo Pellecchia",
					},
				)),
			},
			{
				Description:  "Tracker link with query params",
				InputMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/40538?issue_count=453&issue_position=2&next_issue_id=40506",
				ExpectedMessage: fmt.Sprintf("This is a test message with a tracker link: %s", createTransformedLink(
					"Hi, can you help me with a Version Extended?",
					"https://www.redmine.org/issues/40538?issue_count=453&issue_position=2&next_issue_id=40506",
					"",
					map[string]string{
						"Subject":   "Hi, can you help me with a Version Extended?",
						"Status":    "Reopened",
						"Priority":  "Normal",
						"UpdatedOn": "2024-04-09T10:42:41Z",
						"Tracker":   "Patch",
						"ID":        "40538",
						"Author":    "Enzo Pellecchia",
					},
				)),
			},
			{
				Description:  "Tracker link with query params and note anchor",
				InputMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/40538?issue_count=453&issue_position=2&next_issue_id=40506#note-4",
				ExpectedMessage: fmt.Sprintf("This is a test message with a tracker link: %s", createTransformedLink(
					"Hi, can you help me with a Version Extended?",
					"https://www.redmine.org/issues/40538?issue_count=453&issue_position=2&next_issue_id=40506#note-4",
					"#note-4",
					map[string]string{
						"Subject":   "Hi, can you help me with a Version Extended?",
						"Status":    "Reopened",
						"Priority":  "Normal",
						"UpdatedOn": "2024-04-09T10:42:41Z",
						"Tracker":   "Patch",
						"ID":        "40538",
						"Author":    "Enzo Pellecchia",
					},
				)),
			},
		},
		"MessageWillBeUpdated": {
			{
				Description:  "Message updated with a single tracker link",
				OldMessage:   "this is an old message",
				InputMessage: "This is a test message with a tracker link: https://www.redmine.org/issues/40556",
				ExpectedMessage: fmt.Sprintf("This is a test message with a tracker link: %s", createTransformedLink(
					"Focus on the textarea after clicking the Edit Journal button",
					"https://www.redmine.org/issues/40556",
					"",
					map[string]string{
						"Subject":    "Focus on the textarea after clicking the Edit Journal button",
						"Status":     "Closed",
						"Priority":   "Normal",
						"UpdatedOn":  "2024-04-29T19:23:49Z",
						"AssignedTo": "Marius BĂLTEANU",
						"Tracker":    "Feature",
						"ID":         "40556",
						"Author":     "Yasu Saku",
					},
				)),
			},
		},
	}

	for hook, testCases := range testCases {
		hook := hook
		testCases := testCases
		t.Run(hook, func(t *testing.T) {
			runTestCases(t, plugin, testCases, hook)
		})
	}
}

func runTestCases(t *testing.T, plugin *Plugin, testCases []TestCase, hook string) {
	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Description, func(t *testing.T) {
			t.Parallel()
			var newPost *model.Post

			if hook == "MessageWillBePosted" {
				newPost, _ = plugin.MessageWillBePosted(nil, &model.Post{Message: testCase.InputMessage})
				assert.Equal(t, testCase.ExpectedMessage, newPost.Message, fmt.Sprintf("Failed test case: %s", testCase.Description))
			} else if hook == "MessageWillBeUpdated" {
				newPost, _ = plugin.MessageWillBeUpdated(nil, &model.Post{Message: testCase.InputMessage}, &model.Post{Message: testCase.OldMessage})
				assert.Equal(t, testCase.ExpectedMessage, newPost.Message, fmt.Sprintf("Failed test case: %s", testCase.Description))
			}
		})
	}
}
