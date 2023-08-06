package stat

import (
	"sync"
	"time"
)

// DataPoint 代表滑动窗口中的数据点
type DataPoint struct {
	Timestamp int64 // 时间戳，以毫秒为单位
	Value     int   // 数据值
}

// SlidingWindow 是滑动窗口的结构体
type SlidingWindow struct {
	data       []*DataPoint // 存储数据点的数组
	size       int          // 窗口大小，即数据点的最大数量
	currSum    int          // 当前窗口数据点的和
	mutex      sync.RWMutex
	expireTime int64 // 数据点过期时间，单位为毫秒
}

// NewSlidingWindow 创建一个新的滑动窗口.
// size < 0 意味着数据可以无限存储，慎用，容易导致 oom
// size = 0 意味着不会存储原始数据，无法使用 pct 等统计功能
func NewSlidingWindow(size int, expireTime time.Duration) *SlidingWindow {
	return &SlidingWindow{
		data:       make([]*DataPoint, 0, size),
		size:       size,
		expireTime: int64(expireTime / time.Millisecond),
	}
}

// Add 向滑动窗口中添加一个新的数据点，并返回当前窗口数据点的和
func (sw *SlidingWindow) Add(value int) int {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	dataPoint := &DataPoint{
		Timestamp: time.Now().UnixMilli(),
		Value:     value,
	}
	// 移除过期的数据点
	sw.removeExpiredDataPoints()

	// 添加新的数据点到数组
	sw.data = append(sw.data, dataPoint)

	// 更新当前窗口数据点的和
	sw.currSum += dataPoint.Value

	// 如果数据点数量超过窗口大小，则移除最旧的数据点
	if len(sw.data) > sw.size {
		oldestData := sw.data[0]
		sw.currSum -= oldestData.Value
		sw.data = sw.data[1:]
	}

	return sw.currSum
}

// GetLatestValue 获取滑动窗口中最新的数据点的值
func (sw *SlidingWindow) GetLatestValue() int {
	sw.mutex.RLock()
	defer sw.mutex.RUnlock()
	if len(sw.data) > 0 {
		return sw.data[len(sw.data)-1].Value
	}

	return 0
}

func (sw *SlidingWindow) GetData() []int {
	sw.mutex.RLock()
	defer sw.mutex.RUnlock()
	var record []int
	for _, d := range sw.data {
		record = append(record, d.Value)
	}
	return record
}

func (sw *SlidingWindow) GetSum() int {
	return sw.currSum
}

// 类似 qps， 单位 second
func (sw *SlidingWindow) GetIncreaseRatio() float64 {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()
	sw.removeExpiredDataPoints()
	if len(sw.data) == 0 || sw.expireTime == 0 {
		return 0
	}
	return float64(sw.currSum) / float64(sw.expireTime*1000)
}

// removeExpiredDataPoints 移除过期的数据点
func (sw *SlidingWindow) removeExpiredDataPoints() {
	currTime := time.Now().UnixNano() / int64(time.Millisecond)

	for len(sw.data) > 0 {
		oldestData := sw.data[0]
		if currTime-oldestData.Timestamp > sw.expireTime {
			sw.currSum -= oldestData.Value
			sw.data = sw.data[1:]
		} else {
			break
		}
	}
}

// PriorityQueue 优先队列，用于维护滑动窗口的数据点按时间戳排序
type PriorityQueue []*DataPoint

// Len 获取优先队列的长度
func (pq PriorityQueue) Len() int { return len(pq) }

// Less 定义数据点比较函数
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Timestamp < pq[j].Timestamp
}

// Swap 交换两个数据点的位置
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

// Push 向优先队列中添加数据点
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*DataPoint)
	*pq = append(*pq, item)
}

// Pop 从优先队列中取出最旧的数据点
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
