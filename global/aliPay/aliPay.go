package aliPay

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"gomall/global/config"
	"log"
	"os"
	"path/filepath"
)

var alipayClient *alipay.Client

func getConfigPath(filename string) string {
	//判断一下是在启动项目还是运行测试类，返回不同的配置文件的路径
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//判断当前是否在test目录下
	if filepath.Base(curDir) == "test" {
		path := filepath.ToSlash("../config/" + filename)
		fmt.Println("path=" + path)
		return path
	}
	path := filepath.ToSlash("./config/" + filename)
	fmt.Println("path=" + path)
	return path
}

func ReturnsInstance() *alipay.Client {

	_ = getConfigPath("appPublicCert.crt")

	var alipayConfig = config.Config.AlipayConfig
	var err error
	if alipayClient, err = alipay.New(alipayConfig.AppId, alipayConfig.PrivateKey, false); err != nil {
		log.Fatalf("init alipay client failed :%v", err)
	}

	//加载应用公钥证书
	//if err = alipayClient.LoadAppCertPublicKeyFromFile("./config/appPublicCert.crt"); err != nil {
	//	log.Fatalf("加载laipay公钥证书 failed :%v", err)
	//}
	if err = alipayClient.LoadAppCertPublicKeyFromFile(getConfigPath("appPublicCert.crt")); err != nil {
		log.Fatalf("加载laipay公钥证书 failed :%v", err)
	}

	//加载支付宝根证书
	//if err = alipayClient.LoadAliPayRootCertFromFile("./config/alipayRootCert.crt"); err != nil {
	//	log.Fatalf("加载alipay根证书 failed :%v", err)
	//}
	if err = alipayClient.LoadAliPayRootCertFromFile(getConfigPath("alipayRootCert.crt")); err != nil {
		log.Fatalf("加载alipay根证书 failed :%v", err)
	}

	//加载支付宝公钥证书
	//if err = alipayClient.LoadAlipayCertPublicKeyFromFile("./config/alipayPublicCert.crt"); err != nil {
	//	log.Fatalf("加载alipay公钥证书 failed :%v", err)
	//}
	if err = alipayClient.LoadAlipayCertPublicKeyFromFile(getConfigPath("alipayPublicCert.crt")); err != nil {
		log.Fatalf("加载alipay公钥证书 failed :%v", err)
	}

	return alipayClient
}
