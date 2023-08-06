package stat

import (
	"fmt"
	"testing"
	"time"
)

func TestNewSlidingWindow(t *testing.T) {
	// 创建一个滑动窗口实例，窗口大小为5，数据点过期时间为1000毫秒（1秒）
	window := NewSlidingWindow(5, 1000*time.Millisecond)

	// 添加一些数据点
	now := time.Now().UnixNano() / int64(time.Millisecond)
	dataPoints := []*DataPoint{
		{Timestamp: now, Value: 10},
		{Timestamp: now + 500, Value: 20},
		{Timestamp: now + 1000, Value: 30},
		{Timestamp: now + 1500, Value: 40},
		{Timestamp: now + 2000, Value: 50},
		{Timestamp: now + 2500, Value: 60},
	}

	for _, dp := range dataPoints {
		window.Add(dp.Value)
		time.Sleep(500 * time.Millisecond) // 模拟数据点到达的时间间隔
	}

	fmt.Println("Current sum:", window.currSum) // 输出当前窗口数据点的和
}
