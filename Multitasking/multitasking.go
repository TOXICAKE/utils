package Multitasking

import (
	"context"
	"errors"
	"fmt"
	"gitlab.huaun.com/xuzy/utils"
	"reflect"
	"sync"
)

type Multitasking struct {
	name       string
	threads    int
	taskFunc   func()
	execFunc   func(interface{}, context.Context) interface{}
	taskQueue  chan interface{}
	resultChan chan interface{}
	wg         sync.WaitGroup
	ctx        context.Context
}

func (m *Multitasking) AddTask(taskInfo interface{}) {
	m.taskQueue <- taskInfo
	utils.DebugEcho("Join task successfully")
}

func (m *Multitasking) Register(taskFunc func(), execFunc func(interface{}, context.Context) interface{}) {
	m.taskFunc = taskFunc
	m.execFunc = execFunc
}

func (m *Multitasking) Run() ([]interface{}, error) {
	var result []interface{}
	utils.DebugEcho("Try run module '" + m.name + "'...")
	if m.taskFunc == nil || m.execFunc == nil {
		errors.New("Multitasking '" + m.name + "' must be registered")
	}

	go func() {
		m.taskFunc()
		close(m.taskQueue)
		utils.DebugEcho("Task-thread closed.")
	}() //启动任务写入进程
	utils.DebugEcho("Start task-thread")

	for i := 0; i < m.threads; { //启动任务执行进程
		m.wg.Add(1)
		go func(tid int) {
			defer m.wg.Done()
			for task := range m.taskQueue {
				utils.DebugEcho(tid, "::Execute task <"+reflect.TypeOf(task).String()+">"+fmt.Sprintf("%s", task))
				m.resultChan <- m.execFunc(task, m.ctx)
				utils.DebugEcho(tid, "::Task <"+reflect.TypeOf(task).String()+">"+fmt.Sprintf("%s", task)+" Done")
			}
			utils.DebugEcho(tid, "::Exec-thread closed.")
		}(i)
		utils.DebugEcho(i, "::New exec-thread")
		i += 1
	}
	stop := make(chan struct{})
	go func() { //启动执行结果读取进程
		for r := range m.resultChan {
			result = append(result, r)
		}
		utils.DebugEcho("Result-thread closed.")
		close(stop)
	}()
	utils.DebugEcho("Started result-thread")

	m.wg.Wait()
	close(m.resultChan)
	<-stop
	return result, nil
}

func NewLRModule(name string, threads int) *Multitasking {
	lrm := &Multitasking{
		name:       name,
		threads:    threads,
		taskQueue:  make(chan interface{}, threads),
		resultChan: make(chan interface{}, threads),
		wg:         sync.WaitGroup{},
	}
	lrm.ctx = context.Background()
	return lrm
}
