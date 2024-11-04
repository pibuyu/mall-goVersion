package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

type Info struct {
	SqlConfig *SqlConfigStruct
	RConfig   *RConfigStruct

	ProjectUrl     string
	MongoDBConfig  *MongoDBConfigStruct
	RabbitMQConfig *RabbitMQConfigStruct
}

func init() {
	//避免全局重复导包
	ReturnsInstance()
}

var Config = new(Info)
var cfg *ini.File
var err error

type RabbitMQConfigStruct struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Vhost    string `ini:"vhost"`
}

type SqlConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Host     int    `ini:"host"`
	Password string `ini:"password"`
	Database string `ini:"database"`
}

type RConfigStruct struct {
	IP       string `ini:"ip"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
}

type MongoDBConfigStruct struct {
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Port     int    `ini:"port"`
	Database string `ini:"database"`
}

func getConfigPath() string {
	//判断一下是在启动项目还是运行测试类，返回不同的配置文件的路径
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//判断当前是否在test目录下
	if filepath.Base(curDir) == "test" {
		return filepath.ToSlash("../config/config.ini")
	}
	return filepath.ToSlash("./config/config.ini")
}

func ReturnsInstance() *Info {
	Config.SqlConfig = &SqlConfigStruct{}

	//统一配置文件的路径
	path := getConfigPath()

	cfg, err = ini.Load(path)
	if err != nil {
		log.Fatalf("配置文件不存在,请检查环境: %v \n", err)
	}
	err = cfg.Section("mysql").MapTo(Config.SqlConfig)
	if err != nil {
		log.Fatalf("Mysql读取配置文件错误: %v \n", err)
	}
	//mongodb config
	Config.MongoDBConfig = &MongoDBConfigStruct{}
	err = cfg.Section("mongodb").MapTo(Config.MongoDBConfig)
	if err != nil {
		log.Fatalf("mongodb读取配置文件错误: %v \n", err)
	}
	//redis config
	Config.RConfig = &RConfigStruct{}
	err = cfg.Section("redis").MapTo(Config.RConfig)
	if err != nil {
		log.Fatalf("Redis读取配置文件错误: %v \n", err)
	}
	//rabbitmq config
	Config.RabbitMQConfig = &RabbitMQConfigStruct{}
	err = cfg.Section("rabbitmq").MapTo(Config.RabbitMQConfig)
	if err != nil {
		log.Fatalf("rabbitmq读取配置文件错误: %v \n", err)
	}

	return Config
}
