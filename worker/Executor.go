package worker

import (
	"context"
	"github.com/Goular/crontab/common"
	"os/exec"
	"time"
)

// 任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {

		var (
			cmd    *exec.Cmd
			err    error
			output []byte
			result *common.JobExecuteResult
		)
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		// 记录任务开始时间
		result.StartTime = time.Now()

		// 执行shell命令(Linux使用)
		// cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)

		// 执行shell命令(Windows使用)
		cmd = exec.CommandContext(context.TODO(), "E:\\cygwin64\\bin\\bash.exe", "-c", info.Job.Command)

		// 执行并捕获输出
		output, err = cmd.CombinedOutput()

		// 记录任务结束时间
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err

		// 任务执行完成后，把执行的结果返回给Scheduler，Scheduler会从executingTable中删除掉执行的记录
		G_scheduler.PushJobResult(result)
	}()
}

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}
