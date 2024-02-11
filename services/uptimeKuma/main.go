package uptimekuma

import "log/slog"

func SendMonitorStatus(instance, token string, up bool) error {
	status := "up"
	if !up {
		status = "down"
	}

	slog.Debug("Sending monitor status", "instance", instance, "status", status)

	return nil
}
