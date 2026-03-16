package beckon

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	statsd_client "github.com/cactus/go-statsd-client/statsd"
)

type AuditWith string

const (
	Statsd   AuditWith = "statsd"
	Filebeat AuditWith = "filebeat"
)

type beatBoxer struct {
	key         string
	description string
	identifier  []string
	variable    []string
	identifyBy  []string
	auditWith   AuditWith
	loggerKey   string
	err         error
	statsd      statsd_client.Statter
}

var beatBoxerSync sync.Once

var beatBoxerMapper map[string]*beatBoxer

func init() {
	beatBoxerSync.Do(func() {
		beatBoxerMapper = make(map[string]*beatBoxer)
	})
}

type Setter interface {
	SetIdentifier(...string) Setter
	SetVariable(...string) Setter
	SetDescription(string) Setter
	AuditWith(AuditWith, ...string) Setter
	Done()
}

var _ Setter = &beatBoxer{}

func SetBeatBoxer(key string) Setter {
	if key == "" {
		return &beatBoxerError{
			key: key,
			err: errors.New("key invalid"),
		}
	}

	if _, found := beatBoxerMapper[key]; found {
		return &beatBoxerError{
			key: key,
			err: fmt.Errorf(`key "%s" already exists`, key),
		}
	}
	return &beatBoxer{key: key, description: "...", loggerKey: "BeatBox"}
}

func (b *beatBoxer) SetIdentifier(identifiers ...string) Setter {
	m := make(map[string]bool)

	b.identifier = make([]string, 0)
	for _, h := range identifiers {
		if _, found := m[h]; found {
			return &beatBoxerError{
				key: b.key,
				err: fmt.Errorf(`can't use duplicate identifier "%s" in BeatBox`, h),
			}
		}
		m[h] = true
		b.identifier = append(b.identifier, h)
	}
	return b
}

func (b *beatBoxer) SetVariable(variables ...string) Setter {
	m := make(map[string]bool)

	b.variable = make([]string, 0)
	for _, h := range variables {
		if _, found := m[h]; found {
			return &beatBoxerError{
				key: b.key,
				err: fmt.Errorf(`can't use duplicate variable "%s" in BeatBox`, h),
			}
		}
		m[h] = true
		b.variable = append(b.variable, h)
	}
	return b
}

func (b *beatBoxer) SetDescription(description string) Setter {
	b.description = description
	return b
}

func (b *beatBoxer) AuditWith(with AuditWith, url ...string) Setter {
	b.auditWith = with
	if with == Statsd {
		if len(url) > 1 {
			fmt.Printf("[BeatBoxError] # only supports single URL for statsd\n")
			b.auditWith = Filebeat
			b.loggerKey = "StatsD"
			return b
		}

		var statsdURL string
		if len(url) == 1 {
			statsdURL = url[0]
		} else {
			statsdURL = GetStatsdURL()
		}

		if statsdURL != "" {
			client, err := statsd_client.NewClientWithConfig(&statsd_client.ClientConfig{
				Address:       statsdURL,
				Prefix:        b.key,
				ResInterval:   0,
				UseBuffered:   false,
			})
			if err != nil {
				fmt.Printf("[BeatBoxError] # %s\n", err.Error())
				b.auditWith = Filebeat
				b.loggerKey = "StatsD"
			} else {
				b.statsd = client
			}
		} else {
			b.auditWith = Filebeat
			b.loggerKey = "StatsD"
		}
	}
	return b
}

func (b *beatBoxer) Done() {
	beatBoxerMapper[b.key] = b
}

type beatBoxerError struct {
	key string
	err error
}

func (b *beatBoxerError) SetIdentifier(_ ...string) Setter {
	return b
}

func (b *beatBoxerError) SetVariable(_ ...string) Setter {
	return b
}

func (b *beatBoxerError) SetDescription(description string) Setter {
	return b
}

func (b *beatBoxerError) AuditWith(_ AuditWith, _ ...string) Setter {
	return b
}

func (b *beatBoxerError) Done() {
	beatBoxerMapper[b.key] = &beatBoxer{
		key: b.key,
		err: b.err,
	}
	fmt.Printf("[BeatBoxError] # %s\n", b.err.Error())
}

type Sender interface {
	Describe(description string) Sender
	IdentifyBy(...string) Sender
	Gauge(...int64)
	Duration(...time.Duration)
	IncV(int64)
	Inc()
}

func BeatBox(key string) Sender {
	bb, found := beatBoxerMapper[key]
	if !found {
		return beatBoxerError{
			err: fmt.Errorf(`key "%s" doesn't exist`, key),
		}
	}
	if bb.err != nil {
		return beatBoxerError{
			key: key,
			err: fmt.Errorf("ignoring BeatBox logging. Reason: %v", bb.err),
		}
	}
	return *bb
}

func (b beatBoxer) Describe(description string) Sender {
	b.description = description
	return b
}

func (b beatBoxer) IdentifyBy(identifyBy ...string) Sender {
	if len(identifyBy) != len(b.identifier) {
		return beatBoxerError{
			err: fmt.Errorf(`length of values and identifier mismatch for key "%s"`, b.key),
		}
	}

	b.identifyBy = make([]string, 0)
	for _, h := range identifyBy {
		b.identifyBy = append(b.identifyBy, h)
	}
	return b
}

func (b beatBoxer) Gauge(values ...int64) {
	if len(values) != len(b.variable) {
		fmt.Println("[BeatBoxError] <mismatch> # length of values and variables mismatch")
		return
	}

	if b.auditWith == "" {
		fmt.Println("[BeatBoxError] <not_found> # nothing set to audit with")
		return
	}

	switch b.auditWith {
	case Filebeat:
		data := make([]string, 0)
		for i, v := range b.identifyBy {
			data = append(data, fmt.Sprintf("%s=%v", b.identifier[i], v))
		}
		for i, v := range values {
			data = append(data, fmt.Sprintf("%s=%v", b.variable[i], v))
		}
		fmt.Printf("[%s] [%v] <%s> :: %s :: # %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), b.loggerKey, b.key, strings.Join(data, " | "), b.description)
	case Statsd:
		var prefix string
		if len(b.identifyBy) > 0 {
			prefix = strings.Join(b.identifyBy, ".")
		}
		statsdClient := b.statsd.NewSubStatter(prefix)
		for i, v := range values {
			if err := statsdClient.Gauge(b.variable[i], v, 1.0); err != nil {
				fmt.Printf("[BeatBoxError] <error> # failed to push into stats. error: %v\n", err)
			}
		}
	}
}

func (b beatBoxer) Duration(values ...time.Duration) {
	if len(values) != len(b.variable) {
		fmt.Println("[BeatBoxError] <mismatch> # length of values and variables mismatch")
		return
	}

	if b.auditWith == "" {
		fmt.Println("[BeatBoxError] <not_found> # nothing set to audit with")
		return
	}

	switch b.auditWith {
	case Filebeat:
		data := make([]string, 0)
		for i, v := range b.identifyBy {
			data = append(data, fmt.Sprintf("%s=%v", b.identifier[i], v))
		}
		for i, v := range values {
			data = append(data, fmt.Sprintf("%s=%v", b.variable[i], v.Seconds()*1000.0))
		}
		fmt.Printf("[%s] [%v] <%s> :: %s :: # %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), b.loggerKey, b.key, strings.Join(data, " | "), b.description)
	case Statsd:
		var prefix string
		if len(b.identifyBy) > 0 {
			prefix = strings.Join(b.identifyBy, ".")
		}
		statsdClient := b.statsd.NewSubStatter(prefix)
		for i, v := range values {
			if err := statsdClient.TimingDuration(b.variable[i], v, 1.0); err != nil {
				fmt.Printf("[BeatBoxError] <error> # failed to push into stats. error: %v\n", err)
			}
		}
	}
}

func (b beatBoxer) Inc() {
	if b.auditWith == "" {
		fmt.Println("[BeatBoxError] <not_found> # nothing set to audit with")
		return
	}

	switch b.auditWith {
	case Filebeat:
		data := make([]string, 0)
		for i, v := range b.identifyBy {
			data = append(data, fmt.Sprintf("%s=%v", b.identifier[i], v))
		}
		for _, v := range b.variable {
			data = append(data, fmt.Sprintf("%s=%v", v, 1))
		}
		fmt.Printf("[%s] [%v] <%s> :: %s :: # %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), b.loggerKey, b.key, strings.Join(data, " | "), b.description)
	case Statsd:
		var prefix string
		if len(b.identifyBy) > 0 {
			prefix = strings.Join(b.identifyBy, ".")
		}
		statsdClient := b.statsd.NewSubStatter(prefix)
		for _, v := range b.variable {
			if err := statsdClient.Inc(v, 1, 1.0); err != nil {
				fmt.Printf("[BeatBoxError] <error> # failed to push into stats. error: %v\n", err)
			}
		}
	}
}

func (b beatBoxer) IncV(value int64) {
	if b.auditWith == "" {
		fmt.Println("[BeatBoxError] <not_found> # nothing set to audit with")
		return
	}

	switch b.auditWith {
	case Filebeat:
		data := make([]string, 0)
		for i, v := range b.identifyBy {
			data = append(data, fmt.Sprintf("%s=%v", b.identifier[i], v))
		}
		for _, v := range b.variable {
			data = append(data, fmt.Sprintf("%s=%v", v, value))
		}
		fmt.Printf("[%s] [%v] <%s> :: %s :: # %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), b.loggerKey, b.key, strings.Join(data, " | "), b.description)
	case Statsd:
		var prefix string
		if len(b.identifyBy) > 0 {
			prefix = strings.Join(b.identifyBy, ".")
		}
		statsdClient := b.statsd.NewSubStatter(prefix)
		for _, v := range b.variable {
			if err := statsdClient.Inc(v, value, 1.0); err != nil {
				fmt.Printf("[BeatBoxError] <error> # failed to push into stats. error: %v\n", err)
			}
		}
	}
}

func (b beatBoxerError) Describe(_ string) Sender {
	return b
}

func (b beatBoxerError) IdentifyBy(_ ...string) Sender {
	return b
}

func (b beatBoxerError) Gauge(_ ...int64) {
	fmt.Printf("[BeatBoxError] # %s\n", b.err.Error())
}

func (b beatBoxerError) Duration(values ...time.Duration) {
	fmt.Printf("[BeatBoxError] # %s\n", b.err.Error())
}

func (b beatBoxerError) Inc() {
	fmt.Printf("[BeatBoxError] # %s\n", b.err.Error())
}

func (b beatBoxerError) IncV(_ int64) {
	fmt.Printf("[BeatBoxError] # %s\n", b.err.Error())
}
