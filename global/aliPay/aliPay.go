package aliPay

import (
	"github.com/smartwalle/alipay/v3"
	"gomall/global/config"
	"log"
)

var alipayClient *alipay.Client

func ReturnsInstance() *alipay.Client {
	var alipayConfig = config.Config.AlipayConfig
	var err error
	if alipayClient, err = alipay.New(alipayConfig.AppId, alipayConfig.PrivateKey, false); err != nil {
		log.Fatalf("init alipay client failed :%v", err)
	}

	//加载应用公钥证书
	if err = alipayClient.LoadAppCertPublicKeyFromFile("./config/appPublicCert.crt"); err != nil {
		log.Fatalf("加载laipay公钥证书 failed :%v", err)
	}

	//加载支付宝根证书
	if err = alipayClient.LoadAliPayRootCertFromFile("./config/alipayRootCert.crt"); err != nil {
		log.Fatalf("加载alipay根证书 failed :%v", err)
	}

	//加载支付宝公钥证书
	if err = alipayClient.LoadAlipayCertPublicKeyFromFile("./config/alipayPublicCert.crt"); err != nil {
		log.Fatalf("加载alipay公钥证书 failed :%v", err)
	}

	return alipayClient
}
