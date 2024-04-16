package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(result)
	defer result.Body.Close()
	bodyBytes, err := io.ReadAll(result.Body)
	assert.Nil(err)
	bodyString := string(bodyBytes)

	assert.Equal("Hello, world!", bodyString)
}

func TestMessageWillBePosted(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{
		configuration: &configuration{
			RedmineInstanceURL: "https://www.redmine.org",
		},
	}

	// Test case 1: No tracker links
	postModel := &model.Post{
		Id:      "post1",
		Message: "This is a test message without any tracker links.",
	}
	newPost, _ := plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(postModel, newPost)

	// Test case 2: Single tracker link
	postModel = &model.Post{
		Id:      "post2",
		Message: "This is a test message with a tracker link: https://www.redmine.org/issues/40556",
	}
	expectedPost := &model.Post{
		Id:      "post2",
		Message: "This is a test message with a tracker link: [Focus on the textarea after clicking the Edit Journal button](https://www.redmine.org/issues/40556)",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 3: Multiple tracker links
	postModel = &model.Post{
		Id:      "post3",
		Message: "This is a test message with multiple tracker links: https://www.redmine.org/issues/40556 and https://www.redmine.org/issues/40559",
	}
	expectedPost = &model.Post{
		Id:      "post3",
		Message: "This is a test message with multiple tracker links: [Focus on the textarea after clicking the Edit Journal button](https://www.redmine.org/issues/40556) and [Fix incorrect icon image paths for Wiki help pages](https://www.redmine.org/issues/40559)",
	}
	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 4: Multiple tracker links and markdown links
	postModel = &model.Post{
		Id:      "post4",
		Message: "This is a test message with multiple tracker links: [a link](https://www.redmine.org/issues/40556) and https://www.redmine.org/issues/40559 and https://www.redmine.org/issues/999999 and https://www.redmine.org/issues/40559",
	}
	expectedPost = &model.Post{
		Id:      "post4",
		Message: "This is a test message with multiple tracker links: [a link](https://www.redmine.org/issues/40556) and [Fix incorrect icon image paths for Wiki help pages](https://www.redmine.org/issues/40559) and https://www.redmine.org/issues/999999 and [Fix incorrect icon image paths for Wiki help pages](https://www.redmine.org/issues/40559)",
	}
	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 5: Tracker links with http protocol and without protocol
	postModel = &model.Post{
		Id:      "post5",
		Message: "This is a test message with a http tracker link: http://www.redmine.org/issues/40556 and without protocol www.redmine.org/issues/40556",
	}
	expectedPost = &model.Post{
		Id:      "post5",
		Message: "This is a test message with a http tracker link: [Focus on the textarea after clicking the Edit Journal button](http://www.redmine.org/issues/40556) and without protocol [Focus on the textarea after clicking the Edit Journal button](www.redmine.org/issues/40556)",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 6: Error retrieving issue name
	postModel = &model.Post{
		Id:      "post6",
		Message: "This is a test message with a tracker link: https://www.redmine.org/issues/999999",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(postModel, newPost)

	// Test case 7: Tracker link with a note anchor
	postModel = &model.Post{
		Id:      "post7",
		Message: "This is a test message with a tracker link: https://www.redmine.org/issues/40538#note-4",
	}

	expectedPost = &model.Post{
		Id:      "post7",
		Message: "This is a test message with a tracker link: [Hi, can you help me with a Version Extended?#note-4](https://www.redmine.org/issues/40538#note-4)",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)
}
