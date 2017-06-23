package fluent

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

const (
	tagField     = "tag"
	messageField = "message"
)

var levels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

// Hook implements a Logrus hook for Fluent.
type Hook struct {
	Fluent     *fluent.Fluent
	Tag        string
	DefaultTag string
}

// Levels reutrns a list of levels to fire this hook for.
func (h *Hook) Levels() []logrus.Level {
	return levels
}

// Fire is invoked by logrus and sends logs to Fluent.
func (h *Hook) Fire(entry *logrus.Entry) error {
	data := make(map[string]interface{})

	// Loop through entry data to avoid modifications.
	for k, v := range entry.Data {
		data[k] = v
	}

	tag := h.getTag(data)

	if _, ok := data[messageField]; !ok {
		data[messageField] = entry.Message
	}

	return h.Fluent.PostWithTime(tag, entry.Time, data)
}

// getTag finds the appropriate tag.
// Order of detection:
// 1. Value of Tag field
// 2. Tag defined in context
// 3. Value of DefaultTag field
func (h *Hook) getTag(data map[string]interface{}) string {
	if h.Tag != "" {
		return h.Tag
	}

	if tag, ok := data[tagField].(string); ok {
		delete(data, tagField)

		return tag
	}

	return h.DefaultTag
}
