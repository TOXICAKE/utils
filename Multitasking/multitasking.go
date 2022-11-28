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

func (lrm *Multitasking) AddTask(taskInfo interface{}) {
	lrm.taskQueue <- taskInfo
	utils.DebugEcho("Join task successfully")
}

func (lrm *Multitasking) Register(taskFunc func(), execFunc func(interface{}, context.Context) interface{}) {
	lrm.taskFunc = taskFunc
	lrm.execFunc = execFunc
}

func (lrm *Multitasking) Run() ([]interface{}, error) {
	var result []interface{}
	utils.DebugEcho("Try run module '" + lrm.name + "'...")
	if lrm.taskFunc == nil || lrm.execFunc == nil {
		errors.New("Multitasking '" + lrm.name + "' must be registered")
	}

	go func() {
		lrm.taskFunc()
		close(lrm.taskQueue)
		utils.DebugEcho("Task-thread closed.")
	}() //启动任务写入进程
	utils.DebugEcho("Start task-thread")

	for i := 0; i < lrm.threads; { //启动任务执行进程
		lrm.wg.Add(1)
		go func(tid int) {
			defer lrm.wg.Done()
			for task := range lrm.taskQueue {
				utils.DebugEcho(tid, "::Execute task <"+reflect.TypeOf(task).String()+">"+fmt.Sprintf("%s", task))
				lrm.resultChan <- lrm.execFunc(task, lrm.ctx)
				utils.DebugEcho(tid, "::Task <"+reflect.TypeOf(task).String()+">"+fmt.Sprintf("%s", task)+" Done")
			}
			utils.DebugEcho(tid, "::Exec-thread closed.")
		}(i)
		utils.DebugEcho(i, "::New exec-thread")
		i += 1
	}
	stop := make(chan struct{})
	go func() { //启动执行结果读取进程
		for r := range lrm.resultChan {
			result = append(result, r)
		}
		utils.DebugEcho("Result-thread closed.")
		close(stop)
	}()
	utils.DebugEcho("Started result-thread")

	lrm.wg.Wait()
	close(lrm.resultChan)
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
