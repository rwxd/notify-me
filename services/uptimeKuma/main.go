package uptimekuma

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sagikazarmark/slog-shim"
)

func SendMonitorStatus(instance, token string, up bool, message string, ping string) error {
	status := "up"
	if !up {
		status = "down"
	}

	req, err := http.NewRequest("GET", strings.TrimSuffix(instance, "/")+"/api/push/"+token, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("status", status)
	if message != "" {
		q.Add("msg", message)
	}
	if ping != "" {
		q.Add("ping", ping)
	}
	req.URL.RawQuery = q.Encode()

	slog.Debug("Sending request to uptime-kuma", "url", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("failed to send status to uptime-kuma, status: %s, body: %s", string(resp.Status), string(body[:n]))
	}

	return nil
}
