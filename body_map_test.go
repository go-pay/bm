package bm

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/go-pay/xlog"
)

func TestBodyMapSetBodyMap(t *testing.T) {
	xlog.SetLevel(xlog.DebugLevel)
	bm := make(BodyMap)
	// 1、配合map使用
	sceneInfo := make(map[string]map[string]string)
	h5Info := make(map[string]string)
	h5Info["type"] = "Wap"
	h5Info["wap_url"] = "https://www.fmm.ink"
	h5Info["wap_name"] = "H5测试支付"
	sceneInfo["h5_info"] = h5Info
	bm.Set("scene_info", sceneInfo)
	xlog.Debug("配合map使用：", bm) // map[scene_info:map[h5_info:map[type:Wap wap_name:H5测试支付 wap_url:https://www.fmm.ink]]]

	bm.Reset()
	xlog.Debug(bm) // []

	// 2、基础用法
	bm.Set("1key", "1value")
	bm.Set("2key", "2value")
	bm.Set("3key", "3value")
	xlog.Debug("基础用法：", bm)

	bm.Reset()
	xlog.Debug(bm) // []

	// 3、链式用法
	bm.Set("4key", "4value").
		Set("5key", "5value").
		Set("6key", "6value")
	xlog.Debug("链式用法：", bm)

	bm.Reset()
	xlog.Debug(bm) // []

	// 4、高级用法
	bm.SetBodyMap("scene_info", func(bm BodyMap) {
		bm.SetBodyMap("h5_info", func(bm BodyMap) {
			bm.Set("type", "Wap").
				Set("wap_url", "https://www.fmm.ink").
				Set("wap_name", "H5测试支付")
		})
	}).Set("7key", "7value").
		Set("8key", "8value")
	xlog.Debug("高级用法：", bm) // map[scene_info:map[h5_info:map[type:Wap wap_name:H5测试支付 wap_url:https://www.fmm.ink]]]
	xlog.Debug("高级用法 JsonBody：", bm.JsonBody())
}

func TestBodyMapMarshal(t *testing.T) {
	xlog.SetLevel(xlog.DebugLevel)
	bm := make(BodyMap)
	bm.Set("4key", "4value").
		Set("6key", "6value").
		Set("5key", "5value")
	jb := bm.JsonBody()
	xlog.Debug("jb:", jb)

	bm.Reset()

	bm.SetBodyMap("scene_info", func(bm BodyMap) {
		bm.SetBodyMap("h5_info", func(bm BodyMap) {
			bm.Set("type", "Wap").
				Set("wap_url", "https://www.fmm.ink").
				Set("wap_name", "H5测试支付")
		})
	}).Set("7key", "7value").
		Set("8key", "8value")
	jb2 := bm.JsonBody()
	xlog.Debug("jb2:", jb2)

	bm.Reset()

	bm.SetBodyMap("partner", func(bm BodyMap) {
		bm.Set("type", "APPID").
			Set("appid", "wx123456").
			Set("merchant_id", "88888")
	}).SetBodyMap("authorized_data", func(bm BodyMap) {
		bm.Set("business_type", "BUSIFAVOR_STOCK").
			Set("stock_id", "66666")
	}).Set("limit", 5).
		Set("offset", 10)

	urlParams := bm.EncodeURLParams()
	xlog.Debug("urlParams:", urlParams)
}

func TestBodyMapMarshalSlice(t *testing.T) {
	xlog.SetLevel(xlog.DebugLevel)
	type Receiver struct {
		Type        string `json:"type"`
		Account     string `json:"account"`
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}
	var rs []*Receiver
	item := &Receiver{
		Type:        "MERCHANT_ID",
		Account:     "190001001",
		Amount:      100,
		Description: "分到商户",
	}
	rs = append(rs, item)
	item2 := &Receiver{
		Type:        "PERSONAL_OPENID",
		Account:     "86693952",
		Amount:      888,
		Description: "分到个人",
	}
	rs = append(rs, item2)
	bs, _ := json.Marshal(rs)

	bm := make(BodyMap)
	bm.Set("nonce_str", "abc123").
		Set("transaction_id", "4208450740201411110007820472").
		Set("out_order_no", "P20150806125346")

	bm.Set("receivers", string(bs))

	xlog.Debug("JsonBody:", bm.JsonBody())
	//receiver := make(BodyMap)
	//receiver.Set("receiver", string(bs))
	//
	//body := receiver.JsonBody()
	bss, _ := xml.Marshal(bm)
	xlog.Debug("body:", string(bss))
}

func TestBodyUnmarshal(t *testing.T) {
	xlog.SetLevel(xlog.DebugLevel)
	jsonStr := `{"name":"jerry","age":18}`
	bm := make(BodyMap)
	err := bm.UnmarshalString(jsonStr)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("bm:", bm)
}

func TestSetSlice(t *testing.T) {
	xlog.SetLevel(xlog.DebugLevel)
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var us []*User
	for i := 0; i < 3; i++ {
		us = append(us, &User{
			Name: "Jerry",
			Age:  i,
		})
	}
	bm := make(BodyMap)
	bm.Set("slice", us)
	jb := bm.JsonBody()
	xlog.Debug("bm: ", jb)
}
