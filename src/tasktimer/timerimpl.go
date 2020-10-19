/**
 * Created on 2020/10/20.
 */

package tasktimer

import "time"

type timerBase struct {
	timerId TimerId
	cb      TimerCB
}

//设置回调
func (t *timerBase) SetCallback(cb TimerCB) {
	t.cb = cb
}

//调用回调函数
func (t *timerBase) Callback(now time.Time) {
	if t.cb != nil {
		t.cb(now)
	}
}

//设置计时器ID
func (t *timerBase) SetTimerId(id TimerId) {
	t.timerId = id
}

//获取计时器ID
func (t *timerBase) GetTimerId() TimerId {
	return t.timerId
}

//固定时间点计时器
type FixedTimer struct {
	timerBase
	lastTs int64
	dstTs  int64
}

func (t *FixedTimer) init() {
	t.lastTs = time.Now().Unix()
}

//时间是否到达
func (t *FixedTimer) IsArrive(now time.Time) bool {
	nowTs := now.Unix()
	defer func() {
		t.lastTs = nowTs
	}()
	return t.dstTs > t.lastTs && t.dstTs <= nowTs
}

//计时器类型
func (t *FixedTimer) GetTimerType() TimerType {
	return TimerTypeFixed
}

//每天几点的计时器
type EveryDayHourTimer struct {
	timerBase
	lastHour int
	dstHour  int
}

func (t *EveryDayHourTimer) init() {
	t.lastHour = time.Now().Hour()
}

//时间是否到达
func (t *EveryDayHourTimer) IsArrive(now time.Time) bool {
	hour := now.Hour()
	defer func() {
		t.lastHour = hour
	}()
	return t.dstHour == hour && t.dstHour != t.lastHour
}

//计时器类型
func (t *EveryDayHourTimer) GetTimerType() TimerType {
	return TimerTypeEveryDayHour
}

//每周几点的计时器
type WeeklyHourTimer struct {
	timerBase
	lastWeek time.Weekday
	lastHour int
	dstWeek  time.Weekday
	dstHour  int
}

func (t *WeeklyHourTimer) init() {
	t.lastWeek = time.Now().Weekday()
	t.lastHour = time.Now().Hour()
}

//时间是否到达
func (t *WeeklyHourTimer) IsArrive(now time.Time) bool {
	week := now.Weekday()
	hour := now.Hour()
	defer func() {
		t.lastWeek = week
		t.lastHour = hour
	}()
	return week == t.dstWeek && t.dstHour == hour && t.dstHour != t.lastHour
}

//计时器类型
func (t *WeeklyHourTimer) GetTimerType() TimerType {
	return TimerTypeWeeklyHour
}

//每周几点几分的计时器
type WeeklyHourMinuteTimer struct {
	timerBase
	lastWeek   time.Weekday
	lastHour   int
	lastMinute int
	dstWeek    time.Weekday
	dstHour    int
	dstMinute  int
}

func (t *WeeklyHourMinuteTimer) init() {
	now := time.Now()
	t.lastWeek = now.Weekday()
	t.lastHour = now.Hour()
	t.lastMinute = now.Minute()
}

//时间是否到达
func (t *WeeklyHourMinuteTimer) IsArrive(now time.Time) bool {
	week := now.Weekday()
	hour := now.Hour()
	minute := now.Minute()
	defer func() {
		t.lastWeek = week
		t.lastHour = hour
		t.lastMinute = minute
	}()
	return t.dstWeek == week && t.dstHour == hour && t.dstMinute == minute && t.dstMinute != t.lastMinute
}

//计时器类型
func (t *WeeklyHourMinuteTimer) GetTimerType() TimerType {
	return TimerTypeWeeklyHourMinute
}
