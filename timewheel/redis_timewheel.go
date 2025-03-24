package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"
	"time"
	"unsafe"

	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

var (
	//go:embed add_task.lua
	AddTaskLuaScript string

	//go:embed remove_task.lua
	RemoveTaskLuaScript string

	//go:embed executable_tasks.lua
	GetExecutableTasksLuaScript string
)

type (
	RTimeWheel struct {
		redisCli   *redis.Client
		httpClient *req.Client
		stopC      chan struct{}
		ticker     *time.Ticker
	}

	rTaskElement struct {
		Key     string            `json:"key"`
		Url     string            `json:"url"`
		Headers map[string]string `json:"header"`
		Body    any               `json:"body"`
	}
)

func NewRTaskElement(url string, headers map[string]string, body any) *rTaskElement {
	return &rTaskElement{Body: body, Url: url, Headers: headers}
}

func (s *RTimeWheel) AddTask(key string, task *rTaskElement, executeAtTime time.Time) error {

	task.Key = key

	taskBody, _ := json.Marshal(task)

	// 执行 lua 脚本,实现添加定时任务到 redis zset 中
	s.redisCli.Eval(context.Background(), AddTaskLuaScript, []string{
		s.getMinuteSlice(executeAtTime),
		s.getDeleteSetKey(executeAtTime),
	}, executeAtTime.Unix(), key, string(taskBody))
	return nil
}

func (s *RTimeWheel) RemoveTask(key string, executeAtTime time.Time) {
	if err := s.redisCli.Eval(context.Background(), RemoveTaskLuaScript, []string{
		s.getDeleteSetKey(executeAtTime),
	}, key).Err(); err != nil {
		fmt.Println(err)
	}
}

func (s *RTimeWheel) Run() {
	slog.Info("redis timewheel started...")
	for {
		select {
		case <-s.stopC:
			return
		case <-s.ticker.C:
			go s.executeTasks()
		}
	}
}

func (s *RTimeWheel) Stop() {
	sync.OnceFunc(func() {
		s.redisCli.Close()
		s.httpClient.CloseIdleConnections()
		s.ticker.Stop()
		close(s.stopC)
	})()
}

func (s *RTimeWheel) getMinuteSlice(atTime time.Time) string {
	return fmt.Sprintf("timewheel_task_{%s}", cast.ToString(atTime.Minute()))
}

func (s *RTimeWheel) getDeleteSetKey(atTime time.Time) string {
	return fmt.Sprintf("timewheel_delst_{%s}", cast.ToString(atTime.Minute()))
}

func (s *RTimeWheel) executeTasks() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	tasks, err := s.getExecutableTasks(timeoutCtx)
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(tasks))
	for _, task := range tasks {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()

			// run task
			if err := s.execute(task); err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	wg.Wait()
}

func (s *RTimeWheel) getExecutableTasks(ctx context.Context) ([]*rTaskElement, error) {
	slog.Info("find executable tasks...")

	now := time.Now()

	minuteSlice := s.getMinuteSlice(now)

	deleteSetKey := s.getDeleteSetKey(now)
	nowSecond := time.Now()

	score1, score2 := nowSecond.Unix(), nowSecond.Add(time.Second).Unix()

	rawReply, err := s.redisCli.Eval(ctx, GetExecutableTasksLuaScript, []string{
		minuteSlice, deleteSetKey,
	}, score1, score2).Result()
	if err != nil {
		return nil, err
	}

	replies := cast.ToSlice(rawReply)

	if len(replies) == 0 {
		return nil, errors.New("123123123")
	}

	deleteds := cast.ToStringSlice(replies[0])
	deletedSet := make(map[string]struct{}, len(deleteds))

	for deleted := range slices.Values(deleteds) {
		deletedSet[deleted] = struct{}{}
	}

	tasks := make([]*rTaskElement, 0, len(replies)-1)

	for _, raw := range replies[1:] {
		task := &rTaskElement{}
		v := cast.ToString(raw)
		if err := json.Unmarshal(unsafe.Slice(unsafe.StringData(v), len(v)), task); err != nil {
			continue
		}

		if _, ok := deletedSet[task.Key]; ok {
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *RTimeWheel) execute(task *rTaskElement) error {
	slog.Info("executing task...", "key", task.Key)
	_, err := s.httpClient.R().
		SetHeaders(task.Headers).
		SetBody(task.Body).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.IsSuccessState() {
				defer resp.Body.Close()
				all, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return err
				}
				fmt.Println(string(all))
				return nil
			}
			return resp.Err
		}).Get(task.Url)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func NewRTimeWheel(redisAddr string) *RTimeWheel {
	return &RTimeWheel{
		redisCli:   redis.NewClient(&redis.Options{Addr: redisAddr, DB: 0}),
		httpClient: req.C(),
		stopC:      make(chan struct{}),
		ticker:     time.NewTicker(time.Minute),
	}
}
