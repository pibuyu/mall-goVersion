package payment

type AlipayWebPayReqStruct struct {
	OutTradeNo  string  `json:"outTradeNo"` //支付的流水单号
	Subject     string  `json:"subject"`
	TotalAmount float32 `json:"totalAmount"`
}
