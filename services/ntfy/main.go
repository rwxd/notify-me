package ntfy

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Priority string

var (
	priorityMin     Priority = "min"
	priorityLow     Priority = "low"
	priorityDefault Priority = "default"
	priorityHigh    Priority = "high"
	priorityMax     Priority = "max"
)

type Notification struct {
	Topic    string
	Title    string
	Message  string
	Priority Priority
	Tags     []string
	Url      string
	Actions  string
	Delay    string
	Markdown bool
	Icon     string
}

func SendNotification(n *Notification, instance, user, pass string, token string) error {
	msg := strings.NewReader(n.Message)
	req, err := http.NewRequest("POST", strings.TrimSuffix(instance, "/")+"/"+n.Topic, msg)
	if err != nil {
		return err
	}

	if n.Title != "" {
		req.Header.Set("Title", n.Title)
	}
	if n.Priority != "" {
		req.Header.Set("Priority", string(n.Priority))
	}
	if len(n.Tags) > 0 {
		req.Header.Set("Tags", strings.Join(n.Tags, ","))
	}
	if n.Url != "" {
		req.Header.Set("Click", n.Url)
	}
	if n.Actions != "" {
		req.Header.Set("Actions", n.Actions)
	}
	if n.Delay != "" {
		req.Header.Set("Delay", n.Delay)
	}
	if n.Icon != "" {
		req.Header.Set("Icon", n.Icon)
	}
	if n.Markdown {
		req.Header.Set("Markdown", "true")
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	} else if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	slog.Debug("Sending request to ntfy", "url", req.URL.String(), "headers", req.Header)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("failed to send notification to ntfy, status: %s, body: %s", string(resp.Status), string(body[:n]))
	}

	return nil
}

func NewNotification(topic, title, message string, prio Priority, tags []string, url, actions, delay, icon string, markdown bool) *Notification {
	return &Notification{
		Topic:    topic,
		Title:    title,
		Message:  message,
		Priority: prio,
		Tags:     tags,
		Url:      url,
		Actions:  actions,
		Delay:    delay,
		Markdown: markdown,
		Icon:     icon,
	}
}
