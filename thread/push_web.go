package thread

import (
	"paySystem/models"
	"time"
)

func ThreadPush() {
	PushWeb()
	//每2分钟执行一次
	bc_timer := time.NewTicker(time.Duration(1) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			PushWeb()
		}
	}
}

func PushWeb() {
	//每次100条数据，半个小时之前的数据
	now_unix := time.Now().Unix()
	now_unix = now_unix - 1800
	push_res := models.Gorm.PushPayList(now_unix)
	if len(push_res) > 0 {
		for _, p_list := range push_res {
			Push(p_list)
		}
	}
}

func ThreadPushPayFor() {
	PushPayForWeb()
	//每2分钟执行一次
	bc_timer := time.NewTicker(time.Duration(1) * time.Minute)
	for {
		select {
		case <-bc_timer.C:
			PushPayForWeb()
		}
	}
}

func PushPayForWeb() {
	dt, _ := time.ParseDuration("-30m")
	//每次100条数据，半个小时之前的数据
	now_time := time.Now().Add(dt)
	now_data := now_time.Format(format_date)
	push_res := models.Gorm.PushCashList(now_data)
	if len(push_res) > 0 {
		for _, c_list := range push_res {
			PushPayFor(c_list)
		}
	}
}
