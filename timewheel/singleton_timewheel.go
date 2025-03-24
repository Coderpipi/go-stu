package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type (
	SingletonTimeWheel struct {
		// 时间轮运行间隔
		interval time.Duration

		// 定时器
		ticker *time.Ticker

		// 停止时间轮的通道
		stopC chan struct{}

		// 添加任务
		addTaskC chan *taskElement

		// 移除任务
		removeTaskc chan string

		// 所有时间槽的记录
		slots []*list.List

		// 当前槽idx
		curSlot int

		// task 容器
		keyToETask map[string]*list.Element
	}

	taskElement struct {
		task *Task

		// 定时任务挂在在环状数组中的索引
		pos int

		// 表示需要经过几轮才有能够被执行的机会
		cycle int
	}
)

func (s *SingletonTimeWheel) Stop() {
	sync.OnceFunc(func() {
		s.ticker.Stop()
		close(s.stopC)
	})()
}
func (s *SingletonTimeWheel) AddTask(task *Task) bool {
	if task == nil {
		return false
	}

	pos, cycle := s.getPostAndCycle(task.ExecuteAtTime)

	taskE := &taskElement{
		task:  task,
		pos:   pos,
		cycle: cycle,
	}

	s.addTaskC <- taskE

	return true
}

func (s *SingletonTimeWheel) RemoveTask(key string) {
	s.removeTaskc <- key
}

func (s *SingletonTimeWheel) Run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	for {
		select {
		case <-s.stopC:
			return
		case <-s.ticker.C:
			s.tick()
		case task := <-s.addTaskC:
			s.addTask(task)
		case removeKey := <-s.removeTaskc:
			s.removeTask(removeKey)
		}
	}
}

func (s *SingletonTimeWheel) tick() {
	l := s.slots[s.curSlot]

	// 推进 curSlot位置, 进行环状遍历
	defer s.circularIncr()

	s.execute(l)
}

func (s *SingletonTimeWheel) getPostAndCycle(atTime time.Time) (int, int) {
	delay := int(time.Until(atTime))

	// 定时任务的延时轮次
	cycle := delay / (len(s.slots) * int(s.interval))

	// 定时任务所属环状数组下标的 position
	pos := (s.curSlot + delay/int(s.interval)) % len(s.slots)

	return pos, cycle
}

func (s *SingletonTimeWheel) addTask(task *taskElement) {
	list := s.slots[task.pos]
	if _, ok := s.keyToETask[task.task.Key]; ok {
		s.removeTask(task.task.Key)
	}
	eTask := list.PushBack(task)
	s.keyToETask[task.task.Key] = eTask
}

func (s *SingletonTimeWheel) removeTask(key string) {
	eTask, ok := s.keyToETask[key]
	if !ok {
		return
	}
	delete(s.keyToETask, key)
	task := eTask.Value.(*taskElement)
	_ = s.slots[task.pos].Remove(eTask)
}

func (s *SingletonTimeWheel) circularIncr() {
	s.curSlot = (s.curSlot + 1) % len(s.slots)
}

func (s *SingletonTimeWheel) execute(l *list.List) {
	for e := l.Front(); e != nil; {
		taskE := e.Value.(*taskElement)
		if taskE.cycle > 0 {
			taskE.cycle--
			e = e.Next()
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()

			taskE.task.CallFunc()
		}()
		next := e.Next()
		l.Remove(e)
		delete(s.keyToETask, taskE.task.Key)
		e = next
	}
}

func NewSingletonTimeWheel(slotNum int, interval time.Duration) *SingletonTimeWheel {
	if slotNum <= 0 {
		slotNum = 10
	}

	if interval <= 0 {
		interval = time.Second
	}

	t := &SingletonTimeWheel{
		interval:    interval,
		ticker:      time.NewTicker(interval),
		stopC:       make(chan struct{}),
		keyToETask:  make(map[string]*list.Element),
		addTaskC:    make(chan *taskElement),
		removeTaskc: make(chan string),
	}

	for range slotNum {
		t.slots = append(t.slots, list.New())
	}

	return t
}
