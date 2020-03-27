package models

type AccessList struct {
	Id          int `orm:"pk"`
	Code        string
	Title       string
	Private_key string
	Private_iv  string
	Push_url    string
	Status      int
	Domain      string
	Back_url    string
}

type AccessIp struct {
	Id        string `orm:"pk"`
	Access_id int
	Ip        string
}

type PayBankList struct {
	Id             int `orm:"pk"`
	WebsiteType    int
	Cash_type_code string
	Class_id       int
	Bank_code      string
	Bank_title     string
	Jump_status    int
	Style_code     string
}

type PayThirdpayConfig struct {
	Id             int `orm:"pk"`
	Status         int
	Cash_type_id   string
	Cash_type_code string
	Api_config     string
	Back_domain    string
}

type PayList struct {
	Id              string `orm:"pk"`
	Status          int    //订单状态(1=完成；0=处理中；-1=删除；2=非法订单；9=订单过期)
	Cash_type_code  string
	Pay_id          int
	Access_id       int
	Access_code     string
	Push_status     int
	Push_num        int
	Amout           float64
	Real_amout      float64
	Ctime           int64
	Finish_time     int64
	Ordernumber     string
	Web_ordernumber string
	Pay_ordernumber string
	Pay_class       int
	Bank_code       string
	Account         string
	Source_url      string
	Note            string
	Web_type        int
}

type CashList struct {
	Id              string `orm:"pk"`
	Status          int
	Bank_code       string
	Bank_title      string
	Card_number     string
	Card_name       string
	Account         string
	Cash_type       int
	Amout           float64
	Real_amout      float64
	Fee_amout       float64
	Create_time     string
	Finish_time     string
	Web_ordernumber string
	Order_number    string
	Access_id       int
	Pay_number      string
	Pay_class       int
	Pay_id          int
	Push_status     int
	Push_num        int
	Note            string
}

type PayCashType struct {
	Id      string `orm:"pk"`
	Status  int
	Code    string
	Title   string
	Is_auto int
	Is_df   int
}

type SysBank struct {
	Id         int `orm:"pk"`
	Status     int
	Bank_code  string
	Bank_title string
}

type PayClass struct {
	Id    int
	Code  int
	Title string
}
