package master

import (
	"encoding/json"
	"github.com/Goular/crontab/master/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	// 单例对象
	G_apiServer *ApiServer
)

// 任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

// 保存任务接口
// post job = {"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	// 任务保存到ETCD中
	var (
		err     error
		postJob string
		job     *common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	// 1.解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 2.获取表单中的job字段值
	postJob = req.PostForm.Get("job")
	// 3.反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 4.保存到ETCD
	if oldJob, err = G_jobMgr.SaveJob(job); err != nil {
		goto ERR
	}
	// 5.返回正常的应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6.返回异常的应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 删除任务接口
// POST /job/delete name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	// 任务保存到ETCD中
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	// 1.post参数获取
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 删除任务名
	name = req.PostForm.Get("name")
	//去删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	// 5.返回正常的应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6.返回异常的应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 初始化服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 开启异步协程HTTP服务
	go httpServer.Serve(listener)

	return
}
