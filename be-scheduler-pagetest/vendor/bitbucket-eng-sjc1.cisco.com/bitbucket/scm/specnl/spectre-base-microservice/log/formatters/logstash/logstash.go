package logstash

import (
	"encoding/json"
	"fmt"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log" // spectre
)

// Formatter generates json in logstash format.
// Logstash site: http://logstash.net/
type LogstashFormatter struct {
	Type string // if not empty use for logstash type field.

	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

func (f *LogstashFormatter) Format(entry *log.Entry) ([]byte, error) { // spectre
	fields := make(log.Fields) // spectre
	for k, v := range entry.Data {
		fields[k] = v
	}

	fields["@version"] = 1

	if f.TimestampFormat == "" {
		f.TimestampFormat = log.DefaultTimestampFormat // spectre
	}

	fields["@timestamp"] = entry.Time.Format(f.TimestampFormat)

	// set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}
	fields["message"] = entry.Message

	// set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}
	fields["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			fields["fields.type"] = v
		}
		fields["type"] = f.Type
	}

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
