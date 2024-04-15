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
	plugin := Plugin{}

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
		Message: "This is a test message with a tracker link: https://tracker.sendpulse.com/issues/49033",
	}
	expectedPost := &model.Post{
		Id:      "post2",
		Message: "This is a test message with a tracker link: [Дизайн для UA вебінара з Владиславом Комаревичем 04.04.24](https://tracker.sendpulse.com/issues/49033)",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 3: Multiple tracker links
	postModel = &model.Post{
		Id:      "post3",
		Message: "This is a test message with multiple tracker links: https://tracker.sendpulse.com/issues/49033 and https://tracker.sendpulse.com/issues/49034",
	}
	expectedPost = &model.Post{
		Id:      "post3",
		Message: "This is a test message with multiple tracker links: [Дизайн для UA вебінара з Владиславом Комаревичем 04.04.24](https://tracker.sendpulse.com/issues/49033) and [Заменить png на svg на странице велком пейдж CRM](https://tracker.sendpulse.com/issues/49034)",
	}
	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 4: Multiple tracker links and markdown links
	postModel = &model.Post{
		Id:      "post4",
		Message: "This is a test message with multiple tracker links: [a link](https://tracker.sendpulse.com/issues/49033) and https://tracker.sendpulse.com/issues/49034 and https://tracker.sendpulse.com/issues/999999 and https://tracker.sendpulse.com/issues/49034",
	}
	expectedPost = &model.Post{
		Id:      "post4",
		Message: "This is a test message with multiple tracker links: [a link](https://tracker.sendpulse.com/issues/49033) and [Заменить png на svg на странице велком пейдж CRM](https://tracker.sendpulse.com/issues/49034) and https://tracker.sendpulse.com/issues/999999 and [Заменить png на svg на странице велком пейдж CRM](https://tracker.sendpulse.com/issues/49034)",
	}
	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(expectedPost, newPost)

	// Test case 5: Error retrieving issue name
	postModel = &model.Post{
		Id:      "post5",
		Message: "This is a test message with a tracker link: https://tracker.sendpulse.com/issues/999999",
	}

	newPost, _ = plugin.MessageWillBePosted(nil, postModel)
	assert.Equal(postModel, newPost)
}
