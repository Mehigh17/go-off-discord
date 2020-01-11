package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// MessageAuthor represents the author of a message.
type MessageAuthor struct {
	ID string `json:"id"`
}

// ChannelMessage represents the message owned by a user sent on a channel.
type ChannelMessage struct {
	ID      string        `json:"id"`
	Content string        `json:"content"`
	Author  MessageAuthor `json:"author"`
	Hit     bool          `json:"hit"`
}

// ChannelMessagesResponse represents the answer given by the API when a search fetched.
type ChannelMessagesResponse struct {
	TotalResults uint32             `json:"total_results"`
	Messages     [][]ChannelMessage `json:"messages"`
}

// AccountConfiguration represents the user credentials.
type AccountConfiguration struct {
	AuthenticationToken string `json:"authToken"`
	UserID              string `json:"userId"`
}

// Client represents a GOD client on which commands are executed upon.
type Client struct {
	baseURL       string
	Configuration AccountConfiguration
}

func (client Client) startDeletion(channelID string) {
	resp, err := client.loadMessages(channelID)
	if err != nil {
		panic(err)
	}

	initialMsgCount := resp.TotalResults
	removedMesages := 0

	for resp.TotalResults != 0 {
		const maxRetries int = 15

		for _, messages := range resp.Messages {
			for _, msg := range messages {
				if msg.Author.ID == client.Configuration.UserID && msg.Hit {
					waitTime := 250
					retries := 0
					for retries < maxRetries {
						time.Sleep(time.Duration(waitTime) * time.Millisecond)

						err = client.deleteMessage(channelID, msg)
						if err != nil {
							retries++
							waitTime += 50
							fmt.Printf("(%s) %s (retry no. %d/%d, waiting %dms) ['%s']\n", channelID, err, retries, maxRetries, waitTime, msg.Content)
						} else {
							removedMesages++
							fmt.Printf("(%s) [%d/%d] Removed message: \"%s\" \n", channelID, removedMesages, initialMsgCount, msg.Content)
							break
						}
					}
				}
			}
		}

		resp, err = client.loadMessages(channelID)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("No more messages have been found on channel %s.\n", channelID)
}

func (client Client) loadMessages(channelID string) (ChannelMessagesResponse, error) {
	var serverResponse ChannelMessagesResponse
	httpClient := http.Client{}

	url := fmt.Sprintf("%s/%s/messages/search?author_id=%s", client.baseURL, channelID, client.Configuration.UserID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return serverResponse, err
	}
	request.Header.Set("Authorization", client.Configuration.AuthenticationToken)

	serverIndexed := false
	var resp *http.Response
	for !serverIndexed {
		resp, err = httpClient.Do(request)
		if err != nil {
			return serverResponse, err
		}

		serverIndexed = resp.StatusCode == http.StatusOK
		if serverIndexed {
			break
		}

		fmt.Printf("The server has not been yet indexed. Trying again in 1s... (Received code: %d) \n", resp.StatusCode)
		time.Sleep(1 * time.Second)
	}

	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return serverResponse, err
	}

	err = json.Unmarshal(bodyData, &serverResponse)
	if err != nil {
		return serverResponse, err
	}

	return serverResponse, nil
}

func (client Client) deleteMessage(channelID string, message ChannelMessage) error {
	delurl := fmt.Sprintf("%s/%s/messages/%s", client.baseURL, channelID, message.ID)

	request, err := http.NewRequest("DELETE", delurl, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Authorization", client.Configuration.AuthenticationToken)

	httpClient := http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if (resp.StatusCode >= 200 && resp.StatusCode <= 299) || (resp.StatusCode == 404) {
		return nil
	}

	return fmt.Errorf("failed to remove the message with id %s (Response code: %d)", message.ID, resp.StatusCode)
}
