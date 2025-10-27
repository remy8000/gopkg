package logstash

import (
	"fmt"
	"reflect"
	"strings"
)

// this struct allows to build a standard message to include in logs "msg"
// this standardized mgs can be proceesed in Logstash with GROK, and allows to fill specific extra fields in logs

// warning the order of fields in struct have to match the grok pattern !
type Msg struct {
	// Kraken logs specific
	Issue string
	Name  string // source name
	Url   string
	// Standard
	Count int
	Error error
	Id    string
	Time  float64 // sec
	// misc message 'msg' or 'message' are already used by logstash
	// must be the last one for GREEDYDATA interception
	Content string
}

// note: with Grok
// count: text: total:O => will save only total=0, empty valy bypassed by grok
// so with the below function it stores <nil> and 0
// but we bypasss nil
func (m Msg) BuildMsg() string {
	var final []string

	// when logstash format json
	// it needs unescaped, proper strings
	m.Issue = fmt.Sprintf("%v", m.Issue)
	m.Content = fmt.Sprintf("%v", m.Content)

	// Manual copy to avoid merge errors
	output := Msg{
		Issue:   m.Issue,
		Name:    m.Name,
		Url:     m.Url,
		Count:   m.Count,
		Error:   m.Error,
		Id:      m.Id,
		Time:    m.Time,
		Content: m.Content,
	}

	t := reflect.TypeOf(output)
	v := reflect.ValueOf(output)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		// bypass '<nil>'
		if value == nil {
			value = fmt.Errorf("")
		}
		final = append(final, fmt.Sprintf("%v:%v", strings.ToLower(field.Name), value))
	}
	return strings.Join(final, " ")
}
