package alog

import (
	"fmt"
	"time"
)

type Timer time.Time

func NewTimer() Timer {
	return Timer(time.Now())
}

func (t Timer) Elapsed() time.Duration {
	return time.Since(time.Time(t))
}

func (t Timer) FormatElapsed() string {
	return FormatDuration(t.Elapsed())
}

var elapsedTimeFormatLow = Colorify("@(green:%s)")
var elapsedTimeFormatMedium = Colorify("@(yellow:%s)")
var elapsedTimeFormatHigh = Colorify("@(red:%s)")

func (t Timer) FormatElapsedColor(mediumTime time.Duration, longTime time.Duration) string {
	return FormatDurationColor(t.Elapsed(), mediumTime, longTime)
}

func FormatDurationColor(duration time.Duration, mediumTime time.Duration, longTime time.Duration) string {
	elapsedStr := FormatDuration(duration)
	if duration < mediumTime {
		return fmt.Sprintf(elapsedTimeFormatLow, elapsedStr)
	}
	if duration < longTime {
		return fmt.Sprintf(elapsedTimeFormatMedium, elapsedStr)
	}
	return fmt.Sprintf(elapsedTimeFormatHigh, elapsedStr)
}

func FormatDuration(duration time.Duration) string {
	secs := duration.Seconds()
	if secs >= 600 {
		if secs >= 10*3600 {
			hours := duration.Hours()
			if hours > 9999 {
				return fmt.Sprintf("%4.0f", hours) + "h"
			} else if hours >= 99.95 {
				return fmt.Sprintf("%4.0f", hours)[:4] + "h"
			} else {
				return fmt.Sprintf("%4.1fh", hours)[:4] + "h"
			}
		} else {
			mins := duration.Minutes()
			if mins >= 99.95 {
				return fmt.Sprintf("%4.0f", mins)[:4] + "m"
			} else {
				return fmt.Sprintf("%4.1f", mins)[:4] + "m"
			}
		}
	} else {
		secs := duration.Seconds()
		if secs >= 0.9995 {
			if secs >= 99.95 {
				return fmt.Sprintf("%4.0f", secs)[:4] + "s"
			} else {
				return fmt.Sprintf("%4.2f", secs)[:4] + "s"
			}
		} else {
			if secs >= 0.00995 {
				return fmt.Sprintf("%3.0f", 1000*secs)[:3] + "ms"
			} else {
				return fmt.Sprintf("%3.1f", 1000*secs)[:3] + "ms"
			}
		}
	}
}
