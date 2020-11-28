package client

import (
	"fmt"
	"gitee.com/cristiane/go-common/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	t.Run("发送验证码", TestVerifyCodeSend)
	t.Run("注册用户", TestRegisterUser)
	t.Run("登录用户-验证码", TestLoginUserWithVerifyCode)
	t.Run("登录用户-密码", TestLoginUserWithPwd)
	t.Run("重置密码", TestLoginUserPwdReset)
	t.Run("获取用户信息", TestGetUserInfo)
	t.Run("用户申请提交审核资料", TestMerchantsMaterial)
	t.Run("商户提交开店材料", TestShopBusinessApply)
	t.Run("店铺上架商品", TestSkuBusinessPutAway)
	t.Run("获取店铺上架商品列表", TestGetSkuList)
	t.Run("补充商品属性", TestSkuBusinessSupplement)
	t.Run("添加商品到购物车", TestSkuJoinUserTrolley)
	t.Run("从购物车移除商品", TestSkuRemoveUserTrolley)
	t.Run("获取用户购物车列表", TestGetUserTrolleyList)
	t.Run("生成唯一订单号", TestGenOrderCode)
	t.Run("创建交易订单", TestTradeCreateOrder)
	t.Run("交易订单支付", TestOrderTradePay)
	t.Run("申请物流", TestLogisticsApply)
	t.Run("用户设置-地址变更", TestUserSettingAddress)
	t.Run("用户设置-获取收货地址", TestUserSettingAddressGet)
	t.Run("搜索-商品库存", TestSearchSkuInventory)
	t.Run("搜索-店铺", TestSearchShop)
	t.Run("获取店铺订单报告", TestGetOrderReport)
	t.Run("用户账户充值", TestUserAccountCharge)
	t.Run("订单评价", TestCommentsOrderCreate)
	t.Run("获取店铺评论列表", TestGetShopCommentsList)
	t.Run("修改评论标签", TestCommentsModify)
	t.Run("获取评论标签", TestCommentsTagList)
}

const benchCount = 90000

func BenchmarkGateway(b *testing.B) {
	b.Run("批量注册用户", BenchmarkRegisterUser)
	b.Run("批量充值", BenchmarkUserAccountCharge)
	b.Run("批量创建订单", BenchmarkTradeCreateOrder)
	b.Run("批量创建订单并支付-用户1", BenchmarkTestOrderTrade_1)
	b.Run("批量创建订单并支付-用户2", BenchmarkTestOrderTrade_2)
	b.Run("批量创建订单并支付-用户3", BenchmarkTestOrderTrade_3)
}

func TestGetUserInfo(t *testing.T) {
	r := baseUrl + userInfo
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestSearchShop(t *testing.T) {
	r := baseUrl + searchShop + "?keyword=交个朋友"
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	//req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestSearchSkuInventory(t *testing.T) {
	r := baseUrl + searchSkuInventory + "?keyword=洗发水"
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	//req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestGetUserTrolleyList(t *testing.T) {
	r := baseUrl + skuUserTrolleyList
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

type UserDeliveryInfo struct {
	Id           int64    `form:"id" json:"id"`
	DeliveryUser string   `form:"delivery_user" json:"delivery_user"`
	MobilePhone  string   `form:"mobile_phone" json:"mobile_phone"`
	Area         string   `form:"area" json:"area"`
	DetailedArea string   `form:"detailed_area" json:"detailed_area"`
	Label        []string `form:"label" json:"label"`
	IsDefault    bool     `form:"is_default" json:"is_default"`
}

type UserSettingAddressPutArgs struct {
	UserDeliveryInfo
	// 0-新增，1-修改，2-删除
	OperationType int `form:"operation_type" json:"operation_type"`
}

func TestUserSettingAddress(t *testing.T) {
	r := baseUrl + userSettingAddress
	t.Logf("request url: %s", r)
	args := UserSettingAddressPutArgs{
		UserDeliveryInfo: UserDeliveryInfo{
			Id:           130,
			DeliveryUser: "李大刀",
			MobilePhone:  "187553543534",
			Area:         "河北省石家庄",
			DetailedArea: "石家庄十来路",
			Label:        []string{"公司", "住宅"},
			IsDefault:    true,
		},
		OperationType: 0,
	}
	data := json.MarshalToStringNoError(args)
	t.Logf("req data: \n%v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

type OrderShopGoods struct {
	SkuCode string `form:"sku_code" json:"sku_code"`
	Price   string `form:"price" json:"price"`
	Amount  int64  `form:"amount" json:"amount"`
	Name    string `form:"name" json:"name"`
	Version int64  `form:"version" json:"version"`
}

type OrderShopSceneInfo struct {
	StoreInfo *OrderShopStoreInfo `form:"store_info" json:"store_info"`
}

type OrderShopStoreInfo struct {
	Id       int64  `form:"id" json:"id"`
	Name     string `form:"name" json:"name"`
	AreaCode string `form:"area_code" json:"area_code"`
	Address  string `form:"address" json:"address"`
}

type OrderShopDetail struct {
	ShopId    int64               `form:"shop_id" json:"shop_id"`
	CoinType  int32               `form:"coin_type" json:"coin_type"`
	Goods     []*OrderShopGoods   `form:"goods" json:"goods"`
	SceneInfo *OrderShopSceneInfo `form:"scene_info" json:"scene_info"`
}

type CreateTradeOrderArgs struct {
	Uid            int64              `json:"uid"`
	ClientIp       string             `json:"client_ip"`
	Description    string             `form:"description" json:"description"`
	DeviceId       string             `form:"device_id" json:"device_id"`
	OrderTxCode    string             `form:"order_tx_code" json:"order_tx_code"`
	UserDeliveryId int32              `form:"user_delivery_id" json:"user_delivery_id"`
	Detail         []*OrderShopDetail `json:"detail"`
}

func TestOrderTradePay(t *testing.T) {
	r := baseUrl + tradeOrderPay
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("tx_code", "a3d48269-caad-497e-9d2d-5e9a0cdc9c5c")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestCommentsModify(t *testing.T) {
	r := baseUrl + commentsTagsModify
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("operation_type", "0")
	data.Set("tag_code", uuid.New().String())
	data.Set("classification_major", "物流")
	data.Set("classification_medium", "配送")
	data.Set("classification_minor", "服务")
	data.Set("content", "配送员没有送货上门")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestCommentsTagList(t *testing.T) {
	r := baseUrl + commentsTagsList + "?tag_code=&classification_major=店铺&classification_medium="
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

type LogisticsCommentsInfo struct {
	LogisticsCode        string   `form:"logistics_code" json:"logistics_code"`
	FedexPack            int8     `form:"fedex_pack_star" json:"fedex_pack"`
	FedexPackLabel       []string `form:"fedex_pack_label" json:"fedex_pack_label"`
	DeliverySpeed        int8     `form:"delivery_speed" json:"delivery_speed"`
	DeliverySpeedLabel   []string `form:"delivery_speed_label" json:"delivery_speed_label"`
	DeliveryService      int8     `form:"delivery_service" json:"delivery_service"`
	DeliveryServiceLabel []string `form:"delivery_service_label" json:"delivery_service_label"`
	Comment              string   `form:"comment" json:"comment"`
}

type OrderCommentsInfo struct {
	ShopId    int64    `form:"shop_id" json:"shop_id"`
	OrderCode string   `form:"order_code" json:"order_code"`
	Star      int8     `form:"star" json:"star"`
	Content   string   `form:"content" json:"content"`
	ImgList   []string `form:"img_list" json:"img_list"`
	CommentId string   `form:"comment_id" json:"comment_id"`
}

type CreateOrderCommentsArgs struct {
	Anonymity             bool `form:"anonymity" json:"anonymity"`
	OrderCommentsInfo     OrderCommentsInfo
	LogisticsCommentsInfo LogisticsCommentsInfo
}

func TestCommentsOrderCreate(t *testing.T) {
	r := baseUrl + commentsOrderCreate
	t.Logf("request url: %s", r)
	args := CreateOrderCommentsArgs{
		Anonymity: false,
		OrderCommentsInfo: OrderCommentsInfo{
			ShopId:    30072,
			OrderCode: "000be2f2-489c-4e19-8e2a-731319c98aab",
			Star:      1,
			Content:   "经常在这家店购买，没毛病",
			ImgList:   []string{"image1"},
			CommentId: "",
		},
		LogisticsCommentsInfo: LogisticsCommentsInfo{
			LogisticsCode:        uuid.New().String(),
			FedexPack:            3,
			FedexPackLabel:       []string{"打包不结实"},
			DeliverySpeed:        3,
			DeliverySpeedLabel:   []string{"送货速度慢"},
			DeliveryService:      3,
			DeliveryServiceLabel: []string{"配送服务不到位"},
			Comment:              "配送人员没送到家门口",
		},
	}
	data := json.MarshalToStringNoError(args)
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", token_10050)
	commonTest(r, req, t)
}

func TestGetShopCommentsList(t *testing.T) {
	r := baseUrl + commentsShopList + "?shop_id=30072"
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func BenchmarkUserAccountCharge(b *testing.B) {
	r := baseUrl + userAccountCharge
	b.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("device_code", "OPPO Find20")
	data.Set("device_platform", "Android 10")
	data.Set("account_type", "0")
	data.Set("coin_type", "0")
	data.Set("amount", "999999999999999.99")
	b.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		b.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	for i := 0; i < benchCount; i++ {
		commonBenchmarkTest(r, req, b)
	}
	b.ReportAllocs()
}

func TestUserAccountCharge(t *testing.T) {
	r := baseUrl + userAccountCharge
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("device_code", "OPPO Find20")
	data.Set("device_platform", "Android 10")
	data.Set("account_type", "0")
	data.Set("coin_type", "0")
	data.Set("amount", "9999999999999999999999999")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

type ApplyLogisticsArgs struct {
	Uid          int              `json:"uid"`
	OutTradeNo   string           `json:"out_trade_no" form:"out_trade_no"`
	Courier      string           `json:"courier" form:"courier"`
	CourierType  int              `json:"courier_type" form:"courier_type"`
	ReceiveType  int              `json:"receive_type" form:"receive_type"`
	SendUser     string           `json:"send_user" form:"send_user"`
	SendAddr     string           `json:"send_addr" form:"send_addr"`
	SendPhone    string           `json:"send_phone" form:"send_phone"`
	SendTime     string           `json:"send_time" form:"send_time"`
	ReceiveUser  string           `json:"receive_user" form:"receive_user"`
	ReceiveAddr  string           `json:"receive_addr" form:"receive_addr"`
	ReceivePhone string           `json:"receive_phone" form:"receive_phone"`
	Goods        []GoodsLogistics `json:"goods" form:"goods"`
}

type GoodsLogistics struct {
	SkuCode string `json:"sku_code" form:"sku_code"`
	Name    string `json:"name" form:"name"`
	Kind    string `json:"kind" form:"kind"`
	Count   int64  `json:"count" form:"count"`
}

func TestLogisticsApply(t *testing.T) {
	r := baseUrl + logisticsApply
	t.Logf("request url: %s", r)
	applyReq := ApplyLogisticsArgs{
		Uid:          0,
		OutTradeNo:   uuid.New().String(),
		Courier:      "微商城快递",
		CourierType:  1,
		ReceiveType:  1,
		SendUser:     "李云龙",
		SendAddr:     "河北省石家庄市丰县迎宾路123号",
		SendPhone:    "18319430520",
		SendTime:     "2020-10-09 12:12:12",
		ReceiveUser:  "马司令",
		ReceiveAddr:  "浙江省杭州市余杭区西湖南路111雅静别院",
		ReceivePhone: "18319430520",
		Goods: []GoodsLogistics{
			{
				SkuCode: "2131d-f111-45e1-b68a-d602c2f0f1b3",
				Name:    "怡宝矿泉水",
				Kind:    "饮用水",
				Count:   98,
			},
		},
	}
	data := json.MarshalToStringNoError(applyReq)
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

// 批量创建订单
func BenchmarkTradeCreateOrder(b *testing.B) {
	r := baseUrl + tradeCreateOrder
	b.Logf("request url: %s", r)
	goods1 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-4898-a2b5-bcfebccb713b",
		Price:   "2.9",
		Amount:  1,
		Name:    "清风抽纸",
		Version: 1,
	}
	goods2 := OrderShopGoods{
		SkuCode: "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09",
		Price:   "23.78",
		Amount:  1,
		Name:    "百威淡色拉格啤酒",
		Version: 1,
	}
	// b882a5c9-564a-4912-a5d4-ce77de71577c
	detail := OrderShopDetail{
		ShopId:   30071,
		CoinType: 0, // 0-rmb,1-usdt
		Goods:    []*OrderShopGoods{&goods1, &goods2},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30071,
				Name:     "福建交个朋友",
				AreaCode: "福建",
				Address:  "福建交个朋友",
			},
		},
	}
	goods3 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-a2b5-a2b5-bcfebccb4898",
		Price:   "19.9",
		Amount:  1,
		Name:    "三只松鼠无骨凤爪",
		Version: 1,
	}
	detail2 := OrderShopDetail{
		ShopId:   30072,
		CoinType: 0,
		Goods:    []*OrderShopGoods{&goods3},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30072,
				Name:     "深圳市交个朋友科技有限公司",
				AreaCode: "深圳市南山区",
				Address:  "深圳市交个朋友科技有限公司",
			},
		},
	}
	data := CreateTradeOrderArgs{
		Description:    "双12预热",
		DeviceId:       "Galaxy Note20 Ultra",
		UserDeliveryId: 133,
		Detail:         []*OrderShopDetail{&detail, &detail2},
	}
	for i := 0; i < benchCount; i++ {
		data.OrderTxCode = uuid.New().String()
		req, err := http.NewRequest("POST", r, strings.NewReader(json.MarshalToStringNoError(data)))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", qToken)
		commonBenchmarkTest(r, req, b)
	}
	b.ReportAllocs()
}

func BenchmarkTestOrderTrade_1(b *testing.B) {
	createOrderUrl := baseUrl + tradeCreateOrder
	orderTradeUrl := baseUrl + tradeOrderPay
	goods1 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-4898-a2b5-bcfebccb713b",
		Price:   "2.9",
		Amount:  1,
		Name:    "清风抽纸",
		Version: 1,
	}
	goods2 := OrderShopGoods{
		SkuCode: "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09",
		Price:   "23.78",
		Amount:  1,
		Name:    "百威淡色拉格啤酒",
		Version: 1,
	}
	// b882a5c9-564a-4912-a5d4-ce77de71577c
	detail := OrderShopDetail{
		ShopId:   30071,
		CoinType: 0, // 0-rmb,1-usdt
		Goods:    []*OrderShopGoods{&goods1, &goods2},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30071,
				Name:     "福建交个朋友",
				AreaCode: "福建",
				Address:  "福建交个朋友",
			},
		},
	}
	goods3 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-a2b5-a2b5-bcfebccb4898",
		Price:   "19.9",
		Amount:  1,
		Name:    "三只松鼠无骨凤爪",
		Version: 1,
	}
	detail2 := OrderShopDetail{
		ShopId:   30072,
		CoinType: 0,
		Goods:    []*OrderShopGoods{&goods3},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30072,
				Name:     "深圳市交个朋友科技有限公司",
				AreaCode: "深圳市南山区",
				Address:  "深圳市交个朋友科技有限公司",
			},
		},
	}
	data := CreateTradeOrderArgs{
		Description:    "双12预热",
		DeviceId:       "Galaxy Note20 Ultra",
		UserDeliveryId: 132,
		Detail:         []*OrderShopDetail{&detail, &detail2},
	}
	for i := 0; i < benchCount; i++ {
		data.OrderTxCode = uuid.New().String()
		req, err := http.NewRequest("POST", createOrderUrl, strings.NewReader(json.MarshalToStringNoError(data)))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token_10050)
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Error(err)
			return
		}
		b.Logf("req url: %v status : %v", createOrderUrl, rsp.Status)
		if rsp.StatusCode != http.StatusOK {
			b.Error("StatusCode != 200")
			return
		}
		body, err := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			b.Error(err)
			return
		}
		var obj CreateOrderRsp
		err = json.Unmarshal(string(body), &obj)
		if err != nil {
			b.Error(err)
			return
		}
		if obj.Code != SuccessBusinessCode {
			log.Printf("business code != %v", SuccessBusinessCode)
			log.Printf("obj ==%+v,obj", obj)
			continue
		}
		if obj.Data.TxCode == "" {
			b.Errorf("创建订单交易号为空")
			continue
		}
		orderTradeReq := url.Values{}
		orderTradeReq.Set("tx_code", obj.Data.TxCode)
		req, err = http.NewRequest("POST", orderTradeUrl, strings.NewReader(orderTradeReq.Encode()))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("token", token_10050)
		commonBenchmarkTest(orderTradeUrl, req, b)
	}
	b.ReportAllocs()
}

func BenchmarkTestOrderTrade_2(b *testing.B) {
	createOrderUrl := baseUrl + tradeCreateOrder
	orderTradeUrl := baseUrl + tradeOrderPay
	goods1 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-4898-a2b5-bcfebccb713b",
		Price:   "2.9",
		Amount:  1,
		Name:    "清风抽纸",
		Version: 1,
	}
	goods2 := OrderShopGoods{
		SkuCode: "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09",
		Price:   "23.78",
		Amount:  1,
		Name:    "百威淡色拉格啤酒",
		Version: 1,
	}
	// b882a5c9-564a-4912-a5d4-ce77de71577c
	detail := OrderShopDetail{
		ShopId:   30071,
		CoinType: 0, // 0-rmb,1-usdt
		Goods:    []*OrderShopGoods{&goods1, &goods2},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30071,
				Name:     "福建交个朋友",
				AreaCode: "福建",
				Address:  "福建交个朋友",
			},
		},
	}
	goods3 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-a2b5-a2b5-bcfebccb4898",
		Price:   "19.9",
		Amount:  1,
		Name:    "三只松鼠无骨凤爪",
		Version: 1,
	}
	detail2 := OrderShopDetail{
		ShopId:   30072,
		CoinType: 0,
		Goods:    []*OrderShopGoods{&goods3},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30072,
				Name:     "深圳市交个朋友科技有限公司",
				AreaCode: "深圳市南山区",
				Address:  "深圳市交个朋友科技有限公司",
			},
		},
	}
	data := CreateTradeOrderArgs{
		Description:    "双12预热",
		DeviceId:       "Galaxy Note20 Ultra",
		UserDeliveryId: 133,
		Detail:         []*OrderShopDetail{&detail, &detail2},
	}
	for i := 0; i < benchCount; i++ {
		data.OrderTxCode = uuid.New().String()
		req, err := http.NewRequest("POST", createOrderUrl, strings.NewReader(json.MarshalToStringNoError(data)))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token_10051)
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Error(err)
			return
		}
		b.Logf("req url: %v status : %v", createOrderUrl, rsp.Status)
		if rsp.StatusCode != http.StatusOK {
			b.Error("StatusCode != 200")
			return
		}
		body, err := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			b.Error(err)
			return
		}
		var obj CreateOrderRsp
		err = json.Unmarshal(string(body), &obj)
		if err != nil {
			b.Error(err)
			return
		}
		if obj.Code != SuccessBusinessCode {
			log.Printf("business code != %v", SuccessBusinessCode)
			log.Printf("obj ==%+v,obj", obj)
			continue
		}
		if obj.Data.TxCode == "" {
			b.Errorf("创建订单交易号为空")
			continue
		}
		orderTradeReq := url.Values{}
		orderTradeReq.Set("tx_code", obj.Data.TxCode)
		req, err = http.NewRequest("POST", orderTradeUrl, strings.NewReader(orderTradeReq.Encode()))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("token", token_10051)
		commonBenchmarkTest(orderTradeUrl, req, b)
	}
	b.ReportAllocs()
}

// 批量创建订单并支付
func BenchmarkTestOrderTrade_3(b *testing.B) {
	createOrderUrl := baseUrl + tradeCreateOrder
	orderTradeUrl := baseUrl + tradeOrderPay
	goods1 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-4898-a2b5-bcfebccb713b",
		Price:   "2.9",
		Amount:  1,
		Name:    "清风抽纸",
		Version: 1,
	}
	goods2 := OrderShopGoods{
		SkuCode: "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09",
		Price:   "23.78",
		Amount:  1,
		Name:    "百威淡色拉格啤酒",
		Version: 1,
	}
	// b882a5c9-564a-4912-a5d4-ce77de71577c
	detail := OrderShopDetail{
		ShopId:   30071,
		CoinType: 0, // 0-rmb,1-usdt
		Goods:    []*OrderShopGoods{&goods1, &goods2},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30071,
				Name:     "福建交个朋友",
				AreaCode: "福建",
				Address:  "福建交个朋友",
			},
		},
	}
	goods3 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-a2b5-a2b5-bcfebccb4898",
		Price:   "19.9",
		Amount:  1,
		Name:    "三只松鼠无骨凤爪",
		Version: 1,
	}
	detail2 := OrderShopDetail{
		ShopId:   30072,
		CoinType: 0,
		Goods:    []*OrderShopGoods{&goods3},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30072,
				Name:     "深圳市交个朋友科技有限公司",
				AreaCode: "深圳市南山区",
				Address:  "深圳市交个朋友科技有限公司",
			},
		},
	}
	data := CreateTradeOrderArgs{
		Description:    "双12预热",
		DeviceId:       "Galaxy Note20 Ultra",
		UserDeliveryId: 131,
		Detail:         []*OrderShopDetail{&detail, &detail2},
	}
	for i := 0; i < benchCount; i++ {
		data.OrderTxCode = uuid.New().String()
		req, err := http.NewRequest("POST", createOrderUrl, strings.NewReader(json.MarshalToStringNoError(data)))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", token_10048)
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Error(err)
			return
		}
		b.Logf("req url: %v status : %v", createOrderUrl, rsp.Status)
		if rsp.StatusCode != http.StatusOK {
			b.Error("StatusCode != 200")
			return
		}
		body, err := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			b.Error(err)
			return
		}
		var obj CreateOrderRsp
		err = json.Unmarshal(string(body), &obj)
		if err != nil {
			b.Error(err)
			return
		}
		if obj.Code != SuccessBusinessCode {
			log.Printf("business code != %v", SuccessBusinessCode)
			log.Printf("obj ==%+v,obj", obj)
			continue
		}
		if obj.Data.TxCode == "" {
			b.Errorf("创建订单交易号为空")
			continue
		}
		orderTradeReq := url.Values{}
		orderTradeReq.Set("tx_code", obj.Data.TxCode)
		req, err = http.NewRequest("POST", orderTradeUrl, strings.NewReader(orderTradeReq.Encode()))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("token", token_10048)
		commonBenchmarkTest(orderTradeUrl, req, b)
	}
	b.ReportAllocs()
}

type CreateOrderRsp struct {
	Code int `json:"code"`
	Data struct {
		TxCode string `json:"tx_code"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func TestTradeCreateOrder(t *testing.T) {
	r := baseUrl + tradeCreateOrder
	t.Logf("request url: %s", r)
	goods1 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-4898-a2b5-bcfebccb713b",
		Price:   "2.9",
		Amount:  1,
		Name:    "清风抽纸",
		Version: 1,
	}
	goods2 := OrderShopGoods{
		SkuCode: "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09",
		Price:   "23.78",
		Amount:  1,
		Name:    "百威淡色拉格啤酒",
		Version: 1,
	}
	// b882a5c9-564a-4912-a5d4-ce77de71577c
	detail := OrderShopDetail{
		ShopId:   30071,
		CoinType: 0, // 0-rmb,1-usdt
		Goods:    []*OrderShopGoods{&goods1, &goods2},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30071,
				Name:     "福建交个朋友",
				AreaCode: "福建",
				Address:  "福建交个朋友",
			},
		},
	}
	goods3 := OrderShopGoods{
		SkuCode: "dd13b4aa-4121-a2b5-a2b5-bcfebccb4898",
		Price:   "19.9",
		Amount:  1,
		Name:    "三只松鼠无骨凤爪",
		Version: 1,
	}
	detail2 := OrderShopDetail{
		ShopId:   30072,
		CoinType: 0,
		Goods:    []*OrderShopGoods{&goods3},
		SceneInfo: &OrderShopSceneInfo{
			StoreInfo: &OrderShopStoreInfo{
				Id:       30072,
				Name:     "深圳市交个朋友科技有限公司",
				AreaCode: "深圳市南山区",
				Address:  "深圳市交个朋友科技有限公司",
			},
		},
	}
	data := CreateTradeOrderArgs{
		Description:    "双12预热",
		DeviceId:       "Galaxy Note20 Ultra",
		OrderTxCode:    uuid.New().String(),
		UserDeliveryId: 133,
		Detail:         []*OrderShopDetail{&detail, &detail2},
	}
	//log.Println(json.MarshalToStringNoError(data))
	t.Logf("req data: %v", json.MarshalToStringNoError(data))
	req, err := http.NewRequest("POST", r, strings.NewReader(json.MarshalToStringNoError(data)))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestVerifyCodeSend(t *testing.T) {
	r := baseUrl + verifyCodeSend
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("country_code", "86")
	data.Set("phone", "18319430520")
	data.Set("business_type", "1")
	data.Set("receive_email", "mybaishati@gmail.com")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func BenchmarkRegisterUser(b *testing.B) {
	r := baseUrl + registerUser
	b.Logf("request url: %s", r)
	for i := 0; i < benchCount; i++ {
		data := url.Values{}
		userName := GetFullName()
		data.Set("user_name", userName)
		data.Set("password", "07030501310")
		data.Set("sex", "1")
		data.Set("age", "33")
		data.Set("country_code", "86")
		data.Set("phone", fmt.Sprintf("%d%d", rand.Intn(9), time.Now().UnixNano()))
		data.Set("email", "mybaishati@gmail.com")
		data.Set("verify_code", "606347")
		data.Set("id_card_no", fmt.Sprintf("10000000%d", time.Now().Unix()))
		data.Set("contact_addr", fmt.Sprintf("南京市%s大院", userName))
		data.Set("invite_code", "494f85aa3000065")
		b.Logf("req data: %v", data)
		req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
		if err != nil {
			b.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("token", qToken)
		commonBenchmarkTest(r, req, b)
	}
}

func TestRegisterUser(t *testing.T) {
	r := baseUrl + registerUser
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("user_name", GetFullName())
	data.Set("password", "07030501310")
	data.Set("sex", "1")
	data.Set("age", "33")
	data.Set("country_code", "86")
	data.Set("phone", "15501707783")
	data.Set("email", "mybaishati@gmail.com")
	data.Set("verify_code", "606347")
	data.Set("id_card_no", fmt.Sprintf("10000000%d", time.Now().Unix()))
	data.Set("contact_addr", "廊坊市淮南路清明河畔李家大院")
	data.Set("invite_code", "")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestLoginUserWithVerifyCode(t *testing.T) {
	r := baseUrl + loginUserWithVerifyCode
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("country_code", "86")
	data.Set("phone", "18319430520")
	data.Set("verify_code", "876306")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestLoginUserWithPwd(t *testing.T) {
	r := baseUrl + loginUserWithPwd
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("country_code", "86")
	data.Set("phone", "18319430520")
	data.Set("password", "07030501310")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func BenchmarkTestLoginUserWithPwd(b *testing.B) {
	r := baseUrl + loginUserWithPwd
	b.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("country_code", "86")
	data.Set("phone", "15501707783")
	data.Set("password", "07030501310")
	b.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		b.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	for i := 0; i < math.MaxInt32; i++ {
		//b.Logf("request token=%v", qToken)
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Error(err)
			return
		}
		//b.Logf("req url: %v status : %v", r, rsp.Status)
		if rsp.StatusCode != http.StatusOK {
			b.Error("StatusCode != 200")
			return
		}
		body, err := ioutil.ReadAll(rsp.Body)
		defer rsp.Body.Close()
		if err != nil {
			b.Error(err)
			return
		}
		//b.Logf("req url: %v body : \n%s", r, body)
		var obj HttpCommonRsp
		err = json.Unmarshal(string(body), &obj)
		if err != nil {
			b.Error(err)
			return
		}
		if obj.Code != SuccessBusinessCode {
			b.Errorf("business code != %v", SuccessBusinessCode)
			b.Errorf("obj ==%+v,obj", obj)
			return
		}
	}
}

func TestLoginUserPwdReset(t *testing.T) {
	r := baseUrl + userPwdReset
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("verify_code", "381825")
	data.Set("password", "12345678")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestMerchantsMaterial(t *testing.T) {
	r := baseUrl + merchantsMaterial
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("operation_type", "0")
	data.Set("register_addr", "深圳市宝安区兴业路宝源二区72栋-深圳星光无限实业有限责任公司")
	data.Set("health_card_no", "R8nJ65TDUGAlqrwerSdb9")
	data.Set("identity", "1")
	data.Set("tax_card_no", "qX2Mr545kznWrlvO4sIp7")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestShopBusinessApply(t *testing.T) {
	r := baseUrl + shopBusinessApply
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("operation_type", "0")
	data.Set("shop_id", "123")
	data.Set("nick_name", "深圳市交个朋友科技有限公司")
	data.Set("full_name", "深圳市交个朋友科技有限公司")
	data.Set("register_addr", "深圳市宝安区兴业路宝源二区72栋")
	data.Set("merchant_id", "1081")
	data.Set("business_addr", "深圳市宝安区宝源二区73栋111号")
	data.Set("business_license", "qX2MkznWrlvO4sIp7")
	data.Set("tax_card_no", "qX2MkznWrlvO4sIp7")
	data.Set("business_desc", "qX2MkznWrlvO4sIp7")
	data.Set("social_credit_code", "qX2MkznWrlvO4sIp7")
	data.Set("organization_code", "qX2MkznWrlvO4sIp7")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestSkuBusinessPutAway(t *testing.T) {
	r := baseUrl + skuBusinessPutAway
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("operation_type", "0")
	data.Set("shop_id", "30071")
	data.Set("sku_code", uuid.New().String())
	data.Set("name", "百威淡色拉格啤酒")
	data.Set("price", "23.78")
	data.Set("title", "百威（Budweiser）淡色拉格啤酒 550ml*15听 整箱装")
	data.Set("sub_title", "百威（Budweiser）淡色拉格啤酒 550ml*15听 整箱装")
	data.Set("desc", "百威（Budweiser）淡色拉格啤酒 550ml*15听 整箱装")
	data.Set("production", "百威啤酒")
	data.Set("supplier", "百威旗舰店")
	data.Set("category", "11010")
	data.Set("color", "白色")
	data.Set("color_code", "199")
	data.Set("specification", "一盒200抽")
	data.Set("desc_link", "https://item.jd.com/2877592.html")
	data.Set("state", "1")
	data.Set("amount", "99999999999999999")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestGetSkuList(t *testing.T) {
	r := baseUrl + skuBusinessGetSkuList + "?shop_id=30069"
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestUserSettingAddressGet(t *testing.T) {
	r := baseUrl + userSettingAddress + "?delivery_id="
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestGenOrderCode(t *testing.T) {
	r := baseUrl + tradeOrderCodeGen
	t.Logf("request url: %s", r)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestGetOrderReport(t *testing.T) {
	r := baseUrl + reportOrder
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("shop_id", "30079")
	data.Set("start_time", "2019-11-22 08:46:41")
	data.Set("end_time", "2020-12-04 18:46:41")
	data.Set("page_size", "20")
	data.Set("page_num", "1")
	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", token_10051)
	commonTest(r, req, t)
}

func TestSkuBusinessSupplement(t *testing.T) {
	r := baseUrl + skuBusinessSupplement
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("operation_type", "0")
	data.Set("shop_id", "30071")
	data.Set("sku_code", "a3e5da0a-d3aa-43e2-a7b8-2c5e264e2a09")
	data.Set("name", "百威淡色拉格啤酒")
	data.Set("size", "1.8m")
	data.Set("shape", "长方形")
	data.Set("production_country", "百威淡色拉格啤酒")
	data.Set("production_date", "2020/10/19 15:20")
	data.Set("shelf_life", "3.3年")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestSkuJoinUserTrolley(t *testing.T) {
	r := baseUrl + skuJoinUserTrolley
	t.Logf("request url: %s", r)
	data := url.Values{}
	data.Set("shop_id", "30070")
	data.Set("sku_code", "b9dc7dea-b188-45cf-b8e7-1503898bd7ea")
	data.Set("count", "2")
	data.Set("time", "2020-09-08 23:32:35")
	data.Set("selected", "false")
	t.Logf("req data: %v", data)
	req, err := http.NewRequest("PUT", r, strings.NewReader(data.Encode()))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

func TestSkuRemoveUserTrolley(t *testing.T) {
	r := baseUrl + skuRemoveUserTrolley
	t.Logf("request url: %s", r)
	skuCode := "ec4abc12-9836-4546-a587-f72e375f7884"
	shopId := "30069"
	r += "?sku_code=" + skuCode + "&shop_id=" + shopId
	req, err := http.NewRequest("DELETE", r, nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("token", qToken)
	commonTest(r, req, t)
}

const (
	SuccessBusinessCode = 200
)

type HttpCommonRsp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func commonTest(r string, req *http.Request, t *testing.T) {
	t.Logf("request token=%v", qToken)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %v status : %v", r, rsp.Status)
	if rsp.StatusCode != http.StatusOK {
		t.Error("StatusCode != 200")
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("req url: %v body : \n%s", r, body)
	var obj HttpCommonRsp
	err = json.Unmarshal(string(body), &obj)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.Code != SuccessBusinessCode {
		t.Errorf("business code != %v", SuccessBusinessCode)
		t.Errorf("obj ==%+v,obj", obj)
		return
	}
}

func commonBenchmarkTest(r string, req *http.Request, b *testing.B) {
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		b.Error(err)
		return
	}
	b.Logf("req url: %v status : %v", r, rsp.Status)
	if rsp.StatusCode != http.StatusOK {
		b.Error("StatusCode != 200")
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		b.Error(err)
		return
	}
	b.Logf("req url: %v body : \n%s", r, body)
	var obj HttpCommonRsp
	err = json.Unmarshal(string(body), &obj)
	if err != nil {
		b.Error(err)
		return
	}
	if obj.Code != SuccessBusinessCode {
		log.Printf("business code != %v", SuccessBusinessCode)
		log.Printf("obj ==%+v,obj", obj)
		return
	}
}
