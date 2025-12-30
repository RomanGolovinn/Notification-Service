package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"Notification-Service/internal/service"
)

// I use mailpit for test

type MailpitResponse struct {
	Messages []struct {
		ID      string
		Snippet string
		To      []struct {
			Address string
		} `json:"To"`
		Bcc []struct {
			Address string
		} `json:"Bcc"`
		Subject string
	} `json:"messages"`
}

func TestEmailSender_SendNotification(t *testing.T) {

	message := fmt.Sprintf("Auto-test-%d", time.Now().UnixNano())

	notification := service.Notification{
		ID:       1,
		Address:  "test@gmail.com",
		Message:  message,
		Type:     service.Email,
		IsSended: false,

		CreatedAt:   time.Now(),
		ScheduledAt: time.Now(),
		SentAt:      nil,
	}

	email := EmailConfig{
		Username: "test",
		Password: "test",
		Host:     "localhost",
		Port:     1025,
		From:     "test",
	}

	ctx := context.Background()

	sender := NewEmailSender(email)

	err := sender.SendNotification(ctx, notification)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8025/api/v1/messages")
	var data MailpitResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Decode missed: %v", err)
	}
	defer resp.Body.Close()

	found := false
	for _, m := range data.Messages {
		addressInBcc := false
		if len(m.Bcc) > 0 && m.Bcc[0].Address == "test@gmail.com" {
			addressInBcc = true
		}

		if addressInBcc {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("email_sender_test fail: letter '%s' not found in Mailpit", message)
	} else {
		fmt.Println("email_sender_test pass")
	}
}
