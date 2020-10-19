/**
 * Created on 2020/10/20.
 */

package tasktimer_test

import (
	"fmt"
	"testing"
	"time"

	"tasktimer"
)

func TestRegisterFixedTimer(t *testing.T) {
	taskTimerMgr := tasktimer.NewTaskTimerMgr()
	tm := time.Now().Unix() + 13
	timerId := taskTimerMgr.RegisterFixedTimer(tm, func(now time.Time) {
		fmt.Println("###########定时器回调", now)
	})
	fmt.Println("###", timerId)
	taskTimerMgr.Start()
	time.Sleep(time.Second * 16)
	taskTimerMgr.Stop()
}

func TestRegisterEveryDayHourTimer(t *testing.T) {
	taskTimerMgr := tasktimer.NewTaskTimerMgr()
	timerId := taskTimerMgr.RegisterEveryDayHourTimer(17, func(now time.Time) {
		fmt.Println("###########定时器回调", now)
	})
	fmt.Println("###", timerId)
	taskTimerMgr.Start()
	time.Sleep(time.Hour)
}

func TestRegisterWeeklyHourTimer(t *testing.T) {
	taskTimerMgr := tasktimer.NewTaskTimerMgr()
	timerId := taskTimerMgr.RegisterWeeklyHourTimer(time.Monday, 17, func(now time.Time) {
		fmt.Println("###########定时器回调", now)
	})
	fmt.Println("###", timerId)
	taskTimerMgr.Start()
	time.Sleep(time.Hour)
}
