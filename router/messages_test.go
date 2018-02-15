package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/traPtitech/traQ/model"
)

var (
	sampleText = "popopo"
)

func TestGetMessageByID(t *testing.T) {
	e, cookie, mw := beforeTest(t)

	message := makeMessage()

	c, rec := getContext(e, t, cookie, nil)
	c.SetPath("/messages/:messageID")
	c.SetParamNames("messageID")
	c.SetParamValues(message.ID)

	requestWithContext(t, mw(GetMessageByID), c)

	if rec.Code != http.StatusOK {
		t.Log(rec.Code)
		t.Fatal(rec.Body.String())
	}
	t.Log(rec.Body.String())
}

func TestGetMessagesByChannelID(t *testing.T) {
	e, cookie, mw := beforeTest(t)

	for i := 0; i < 5; i++ {
		makeMessage()
	}

	post := requestCount{
		Limit: 3,
		Count: 1,
	}
	body, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("PUT", "http://test", bytes.NewReader(body))

	c, rec := getContext(e, t, cookie, req)
	c.SetPath("/channels/:channelID/messages")
	c.SetParamNames("channelID")
	c.SetParamValues(testChannelID)
	requestWithContext(t, mw(GetMessagesByChannelID), c)

	if rec.Code != http.StatusOK {
		t.Log(rec.Code)
		t.Fatal(rec.Body.String())
	}

	var responseBody []MessageForResponse
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	if len(responseBody) != 3 {
		t.Errorf("No found all messages: want %d, actual %d", 3, len(responseBody))
	}

}

func TestPostMessage(t *testing.T) {
	e, cookie, mw := beforeTest(t)

	post := requestMessage{
		Text: "test message",
	}

	body, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "http://test", bytes.NewReader(body))
	c, rec := getContext(e, t, cookie, req)
	c.SetPath("/channels/:channelID/messages")
	c.SetParamNames("channelID")
	c.SetParamValues(testChannelID)
	requestWithContext(t, mw(PostMessage), c)

	message := &MessageForResponse{}

	result, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(result, message)
	if err != nil {
		t.Fatal(err)
	}

	if message.Content != post.Text {
		t.Errorf("message text is wrong: want %v, actual %v", post.Text, message.Content)
	}
	if len(message.StampList) != 0 {
		t.Errorf("StampList length is wrong: want 0, actual %d", len(message.StampList))
	}

	if rec.Code != http.StatusCreated {
		t.Log(rec.Code)
		t.Fatal(rec.Body.String())
	}
}

func TestPutMessageByID(t *testing.T) {
	e, cookie, mw := beforeTest(t)

	message := makeMessage()

	post := requestMessage{
		Text: "test message",
	}
	body, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("PUT", "http://test", bytes.NewReader(body))

	c, rec := getContext(e, t, cookie, req)
	c.SetPath("/messages/:messageID")
	c.SetParamNames("messageID")
	c.SetParamValues(message.ID)
	requestWithContext(t, mw(PutMessageByID), c)

	message, err = model.GetMessage(message.ID)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Log(rec.Code)
		t.Fatal(rec.Body.String())
	}

	if message.Text != post.Text {
		t.Fatalf("message text is wrong: want %v, actual %v", post.Text, message.Text)
	}

}

func TestDeleteMessageByID(t *testing.T) {
	e, cookie, mw := beforeTest(t)

	message := makeMessage()

	req := httptest.NewRequest("DELETE", "http://test", nil)

	c, rec := getContext(e, t, cookie, req)
	c.SetPath("/messages/:messageID")
	c.SetParamNames("messageID")
	c.SetParamValues(message.ID)
	requestWithContext(t, mw(DeleteMessageByID), c)

	message, err := model.GetMessage(message.ID)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusNoContent {
		t.Log(rec.Code)
		t.Fatal(rec.Body.String())
	}

	if message.IsDeleted != true {
		t.Fatalf("message text is wrong: want %v, actual %v", true, message.IsDeleted)
	}

}
