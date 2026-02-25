package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	slackapi "github.com/rneatherway/slack"
	"github.com/slack-go/slack"
)

// convertTimestampToSlackFormat converts a Slack timestamp (e.g., "1734567890.123456")
// to the format used in Slack URLs (e.g., "p1734567890123456")
func convertTimestampToSlackFormat(timestamp string) string {
	// Remove the decimal point and prepend with 'p'
	return "p" + strings.Replace(timestamp, ".", "", 1)
}

func main() {
	// Parse command line arguments
	channelID := flag.String("channel", "", "Slack channel link")
	teamDomain := flag.String("domain", "", "Slack team domain")
	lookback := flag.Int("lookback", 7, "Number of days to look back for messages")
	flag.Parse()

	client := slackapi.NewClient(*teamDomain)
	err := client.WithCookieAuth()

	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Get threads from the last 7 days
	threads, err := Threads(client, *channelID, time.Now().AddDate(0, 0, -1*(*lookback)))
	if err != nil {
		fmt.Printf("Error retrieving threads: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d threads in the last %d days\n\n", len(threads), *lookback)

	// Print thread URLs
	fmt.Println("Thread URLs:")
	for _, thread := range threads {
		threadURL := fmt.Sprintf("https://%s.slack.com/archives/%s/%s",
			*teamDomain, *channelID, convertTimestampToSlackFormat(thread.Timestamp))
		fmt.Printf("%s\n", threadURL)
	}
	fmt.Println()
}

// Replies retrieves conversation replies for a thread. See https://api.slack.com/methods/conversations.replies
func Replies(client *slackapi.Client, channelID string, threadTS string) (*HistoryResponse, error) {
	params := map[string]string{
		"channel": channelID,
		"ts":      threadTS,
	}

	body, err := client.API(context.Background(), "POST", "conversations.replies", params, nil)
	if err != nil {
		return nil, err
	}

	historyResponse := &HistoryResponse{}
	err = json.Unmarshal(body, historyResponse)
	if err != nil {
		return nil, err
	}

	if !historyResponse.Ok {
		return nil, fmt.Errorf("conversations.replies response not OK: %s", body)
	}

	return historyResponse, nil
}

// Threads retrieves all threads since the given time in a specific channel.
func Threads(client *slackapi.Client, channelID string, oldest time.Time) ([]slack.Message, error) {
	// Calculate timestamp for 7 days ago
	var allThreads []slack.Message
	cursor := ""

	fmt.Printf("Retrieving messages from channel %s since %s\n", channelID, oldest.Format("2006-01-02"))

	for {
		// Get conversation history with pagination
		// See https://api.slack.com/methods/conversations.history
		params := map[string]string{
			"channel": channelID,
			"oldest":  strconv.FormatInt(oldest.Unix(), 10),
			"limit":   "200", // Maximum allowed by Slack API
			"cursor":  cursor,
		}

		resp, err := client.API(context.Background(), "POST", "conversations.history", params, nil)
		if err != nil {
			return nil, fmt.Errorf("error getting conversation history: %w", err)
		}
		var history slack.GetConversationHistoryResponse
		if err := json.Unmarshal(resp, &history); err != nil {
			return nil, fmt.Errorf("error unmarshalling conversation history: %w", err)
		}

		// Filter messages that have thread_ts (are part of a thread)
		for _, message := range history.Messages {
			// Check if this message is a thread parent (has replies)
			if message.ReplyCount > 0 {
				allThreads = append(allThreads, message)
				fmt.Printf("Found thread: %s (replies: %d)\n", message.Timestamp, message.ReplyCount)
			}
		}

		// Check if there are more pages
		if !history.HasMore {
			break
		}
		cursor = history.ResponseMetaData.NextCursor
	}

	return allThreads, nil
}
