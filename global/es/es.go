package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sethvargo/go-retry"
	"gomall/global/config"
	"log"
	"os"
	"time"
)

var Es *elastic.Client

func ReturnsInstance() *elastic.Client {
	var err error
	var esConfig = config.Config.ElasticSearchConfig
	b := retry.NewFibonacci(10 * time.Second) //重试的斐波那契机制，最大重试间隔时间为10秒
	ctx := context.Background()

	//需要注意host的格式
	host := fmt.Sprintf("http://%s:%d/", esConfig.Host, esConfig.Port)

	// 自定义日志记录器，控制日志输出级别
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags|log.Lshortfile)

	if err = retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		Es, err = elastic.NewClient(
			elastic.SetErrorLog(errorlog),
			elastic.SetURL(host))
		if err != nil {
			return err
		}
		_, _, err = Es.Ping(host).Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("重试5次后仍然无法连接es，请排查es服务端是否启动/配置信息是否正确，错误信息为： %v \n", err)
	}
	return Es
}
