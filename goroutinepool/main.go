/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"math/rand"
)

type Job struct {
	Id      int
	RandNum int
}

type Result struct {
	*Job
	sum int
}

//go协程池
func main() {
	//需要两个管道
	jobChan := make(chan *Job, 128)
	resultChan := make(chan *Result, 128)

	createPool(64, jobChan, resultChan)

	//打印协程数据
	go func() {
		for result := range resultChan {
			fmt.Printf("job id:%v randNum:%v sum:%d\n", result.Job.Id, result.RandNum, result.sum)
		}
	}()

	//创建数据
	var id int
	for {
		id++
		rNum := rand.Intn(10000)
		job := &Job{
			Id:      id,
			RandNum: rNum,
		}
		jobChan <- job
	}
}

func createPool(num int, jobChan chan *Job, resultChan chan *Result) {
	for i := 0; i < num; i++ {
		go func(jobChan chan *Job, resultChan chan *Result) {
			for job := range jobChan {
				randNum := job.RandNum
				var sum int
				for randNum != 0 {
					sum += randNum % 10
					randNum = randNum / 10
				}
				r := &Result{
					Job: job,
					sum: sum,
				}
				resultChan <- r
			}
		}(jobChan, resultChan)
	}
}
