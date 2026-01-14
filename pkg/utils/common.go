package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/utils/constant"
)

func PrintJson(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("marshal error: %v", err.Error())
	}
	return string(b)
}

func FormatBJTime(t time.Time) string {
	location := time.FixedZone("CST", 8*3600)
	return t.In(location).Format("2006-01-02 15:04:05")
}

func FormatPrice(price float64) string {
	if price >= 1 {
		return fmt.Sprintf("%.2f", price)
	}
	return fmt.Sprintf("%.4f", price)
}

func FormatRelativeTime(t time.Time) string {
	bjTime := FormatBJTime(t)
	duration := time.Since(t)

	if duration < time.Minute {
		return fmt.Sprintf("%s (%d s ago)", bjTime, int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%s (%d min ago)", bjTime, int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%s (%d hours ago)", bjTime, int(duration.Hours()))
	} else {
		return bjTime
	}
}

// QuietHoursParams defines the quiet hours configuration for a task
type QuietHoursParams struct {
	Enabled            bool
	StartHour          int
	EndHour            int
	Behavior           string // "pause" or "throttle"
	ThrottleMultiplier int    // e.g. 5
}

// ShouldExecTask determines if a task should run based on Quiet Hours configuration.
// It returns true if the task should run.
func ShouldExecTask(params QuietHoursParams, lastRun time.Time, interval time.Duration) bool {
	if !params.Enabled {
		return true // Feature disabled, always run
	}

	location := time.FixedZone("CST", 8*3600)
	nowBJ := time.Now().In(location)
	hour := nowBJ.Hour()

	inQuietHours := false
	if params.StartHour < params.EndHour {
		// e.g., 0 to 7
		if hour >= params.StartHour && hour < params.EndHour {
			inQuietHours = true
		}
	} else if params.StartHour > params.EndHour {
		// e.g., 22 to 7 (crosses midnight)
		if hour >= params.StartHour || hour < params.EndHour {
			inQuietHours = true
		}
	} else {
		// startH == endH, effectively disabled or empty range
		return true
	}

	if !inQuietHours {
		return true
	}

	// In Quiet Hours
	if params.Behavior == constant.QUIET_HOURS_BEHAVIOR_THROTTLE {
		if params.ThrottleMultiplier <= 1 {
			return true // No throttling
		}
		// Check if enough time has passed: interval * multiplier
		return time.Since(lastRun) >= interval*time.Duration(params.ThrottleMultiplier)
	}

	// Default behavior is "pause"
	return false
}
