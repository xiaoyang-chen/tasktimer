/**
 * Desc:该定时器通过每5秒判断定时器是否到达，如果到达的话会执行回调。调整系统时间可以不用重启服务器即可生效。
 * Created on 2020/10/20.
 */

package tasktimer

import (
	"context"
	"sort"
	"sync"
	"time"
)

//常量
const (
	DefaultInterval = 3 //默认定时器刷新频率
)

// TimerType 定时器类型
type TimerType int32

const (
	TimerTypeFixed            TimerType = 1 //固定的时间点
	TimerTypeEveryDayHour               = 2 //每天几点
	TimerTypeWeeklyHour                 = 3 //每周几点
	TimerTypeWeeklyHourMinute           = 4 //每周几点几分
)

type TimerId int32
type TimerIdSlice []TimerId

func (p TimerIdSlice) Len() int           { return len(p) }
func (p TimerIdSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p TimerIdSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type TimerCB func(now time.Time)

type ITimer interface {
	init()                       //初始化
	IsArrive(now time.Time) bool //时间是否到达
	GetTimerType() TimerType     //计时器类型
	SetTimerId(id TimerId)       //设置计时器ID
	GetTimerId() TimerId         //获取计时器ID
	SetCallback(cb TimerCB)      //设置回调
	Callback(now time.Time)      //调用回调函数
}

type Manager struct {
	sync.RWMutex
	sync.WaitGroup
	timerId    TimerId
	timerMgr   map[TimerId]ITimer
	subCtx     context.Context
	cancelFunc context.CancelFunc
	interval   int32 //默认5秒为定时器刷新频率
}

func NewTaskTimerMgr() *Manager {
	mgr := new(Manager)
	mgr.init()
	return mgr
}

func (m *Manager) init() {
	m.subCtx, m.cancelFunc = context.WithCancel(context.Background())
	m.timerMgr = make(map[TimerId]ITimer)
	m.interval = DefaultInterval
	m.timerId = 0
}

func (m *Manager) Start() {
	m.Add(1)
	go m.timerRoutine()
}

func (m *Manager) Stop() {
	m.cancelFunc()
	m.Wait()
}

func (m *Manager) Clear() {
	m.Lock()
	defer m.Unlock()
	m.timerMgr = make(map[TimerId]ITimer)
}

//注册定时器：固定时间
func (m *Manager) RegisterFixedTimer(ts int64, cb TimerCB) TimerId {
	if ts == 0 {
		return 0
	}
	return m.register(&FixedTimer{
		timerBase: timerBase{cb: cb},
		dstTs:     ts,
	}).GetTimerId()
}

//注册定时器：每天几点
func (m *Manager) RegisterEveryDayHourTimer(hour int32, cb TimerCB) TimerId {
	return m.register(&EveryDayHourTimer{
		timerBase: timerBase{cb: cb},
		dstHour:   int(hour),
	}).GetTimerId()
}

//注册定时器：每周几点
func (m *Manager) RegisterWeeklyHourTimer(week time.Weekday, hour int32, cb TimerCB) TimerId {
	return m.register(&WeeklyHourTimer{
		timerBase: timerBase{cb: cb},
		dstWeek:   week,
		dstHour:   int(hour),
	}).GetTimerId()
}

//注册定时器：每周几点几分
func (m *Manager) RegisterWeeklyHourMinuteTimer(week time.Weekday, hour int32, minute int32, cb TimerCB) TimerId {
	return m.register(&WeeklyHourMinuteTimer{
		timerBase: timerBase{cb: cb},
		dstWeek:   week,
		dstHour:   int(hour),
		dstMinute: int(minute),
	}).GetTimerId()
}

func (m *Manager) UnRegisterById(id TimerId) {
	m.Lock()
	defer m.Unlock()
	delete(m.timerMgr, id)
}

func (m *Manager) UnRegister(timer ITimer) {
	m.Lock()
	defer m.Unlock()
	delete(m.timerMgr, timer.GetTimerId())
}

func (m *Manager) SetInterval(interval int32) {
	m.interval = interval
}

func (m *Manager) GetInterval() int32 {
	return m.interval
}

func (m *Manager) register(timer ITimer) ITimer {
	m.Lock()
	defer m.Unlock()

	m.timerId++
	timer.SetTimerId(m.timerId)
	timer.init()
	m.timerMgr[timer.GetTimerId()] = timer
	return timer
}

func (m *Manager) timerRoutine() {
	defer m.Done()

	ticker := time.NewTicker(time.Duration(m.interval) * time.Second)
	for {
		select {
		case <-m.subCtx.Done():
			return
		case now := <-ticker.C:
			m.loopTimer(now)
		}
	}
}

func (m *Manager) loopTimer(now time.Time) {
	m.Lock()
	defer m.Unlock()

	ids := m.getSortedTimerIdArray()
	for _, id := range ids {
		timer := m.timerMgr[id]
		if timer.IsArrive(now) {
			timer.Callback(now)
			if timer.GetTimerType() == TimerTypeFixed {
				//固定时间点的定时器是一次性的，用完即删
				delete(m.timerMgr, timer.GetTimerId())
			}
		}
	}
}

//排序
func (m *Manager) getSortedTimerIdArray() []TimerId {
	ids := make([]TimerId, 0, len(m.timerMgr))
	for id := range m.timerMgr {
		ids = append(ids, id)
	}
	sort.Sort(TimerIdSlice(ids))
	return ids
}
