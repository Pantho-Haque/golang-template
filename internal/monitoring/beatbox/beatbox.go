package beatbox

import (
	"fmt"
	"time"

	beatbox "magic.pathao.com/beatbox/beckon"
)

type BeatBox struct {
}

func (bx BeatBox) Init() {
	beatbox.SetBeatBoxer("prism.redis_query_time").SetIdentifier("type", "success").
		SetVariable("time").AuditWith(beatbox.Statsd).Done()

	beatbox.SetBeatBoxer("prism.database_query_time").SetIdentifier("table", "method", "success").
		SetVariable("time").AuditWith(beatbox.Statsd).Done()

	beatbox.SetBeatBoxer("prism.service_call_time").SetIdentifier("name", "path", "success").
		SetVariable("time").AuditWith(beatbox.Statsd).Done()

	beatbox.SetBeatBoxer("prism.service_status_count").SetIdentifier("name", "path", "status_code").
		SetVariable("count").AuditWith(beatbox.Statsd).Done()
}

func (bx BeatBox) TimeTakenInRedisQuery(_type string, isSuccess bool, t time.Time) {
	success := "0"
	if isSuccess {
		success = "1"
	}

	beatbox.BeatBox("prism.redis_query_time").
		IdentifyBy(_type, success).
		Duration(time.Since(t))
}

func (bx BeatBox) TimeTakenInDatabaseQuery(table, method string, isSuccess bool, t time.Time) {
	success := "0"
	if isSuccess {
		success = "1"
	}

	beatbox.BeatBox("prism.database_query_time").
		IdentifyBy(table, method, success).
		Duration(time.Since(t))
}

func (bx BeatBox) TimeTakenInServiceCall(name, path string, isSuccess bool, t time.Time) {
	success := "0"
	if isSuccess {
		success = "1"
	}

	beatbox.BeatBox("prism.service_call_time").
		IdentifyBy(name, path, success).
		Duration(time.Since(t))
}

func (bx BeatBox) StatusCodeCountInServiceCall(name, path string, status int) {
	beatbox.BeatBox("prism.service_status_count").
		IdentifyBy(
			fmt.Sprintf("%v", name),
			fmt.Sprintf("%v", path),
			fmt.Sprintf("%v", status),
		).
		Inc()
}
