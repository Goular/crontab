package master

import (
	"context"
	"github.com/Goular/crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Mongo打包日志管理
type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)
	// 1.建立连接
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri)); err != nil {
		return
	}

	// 2.选择db和collection
	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("test_db").Collection("log"),
	}

	return
}

// 查看任务日志
func (logMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		filter  *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor  *mongo.Cursor
		jobLog  *common.JobLog
		skip64  int64
		limit64 int64
	)
	logArr = make([]*common.JobLog, 0)
	// 过滤条件
	filter = &common.JobLogFilter{JobName: name}
	// 按照任务开始时间倒排
	logSort = &common.SortLogByStartTime{
		SortOrder: -1,
	}
	skip64 = int64(skip)
	limit64 = int64(limit)
	// 查询
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, &options.FindOptions{Sort: logSort, Skip: &skip64, Limit: &limit64}); err != nil {
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		// 反序列化bson
		if err = cursor.Decode(jobLog); err != nil {
			continue // 有日志但格式不合法
		}
		logArr = append(logArr, jobLog)
	}
	return
}
