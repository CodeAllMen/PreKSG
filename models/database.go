package models

type Track struct {
	Id          int64  `orm:"pk;auto"`
	CampId      string `orm:"size(20);index"`
	AffName     string `orm:"size(20);index"`
	PubId       string `orm:"size(100);index"`
	ProId       string `orm:"size(10)"`
	ClickId     string `orm:"size(100)"`
	ShortCode   string `orm:"size(10)"`
	Keyword     string `orm:"size(10)"`
	ProductName string `orm:"size(30)"`
	Bd          string `orm:"size(10)"`
	Time        string `orm:"size(30);index"`
	Ip          string `orm:"size(30)"`
	Agent       string `orm:"size(300)"`
	ClickStatus string `orm:"size(10)"`
	ClickTime   string `orm:"size(30)"`
	Operator    string `orm:"size(10)"`
	AocUrl      string `orm:"size(300)"`
	AocError    string `orm:"size(100)"`
}

type Old_Postback struct {
	Id       int64  `orm:"pk;auto"`
	CampId   string `orm:"size(20)"`
	AffName  string `orm:"size(255)"`
	Url      string `orm:"size(255)"`
	Country  string `orm:"size(20)"`
	Operator string `orm:"size(30)"`
	Payout   string `orm:"size(20)"`
	Rate     int
	Cap      int
}

type PostbackRecord struct {
	Id      int64  `orm:"pk;auto"`
	CampId  string `orm:"size(20);index"`
	AffName string `orm:"size(20);index"`
	Clickid string `orm:"size(100)"`
	Proid   string `orm:"size(100)"`
	Pubid   string `orm:"size(100);index"`
	Time    string `orm:"size(50)"`
}

type MoStruct struct {
	Id        int64  `orm:"pk;auto"`
	Moid      string `orm:"size(50);index"`
	Shortcode string `orm:"size(20);index"`
	Keyword   string `orm:"size(20);index"`
	Msisdn    string `orm:"size(20);index"`
	Subtime   string `orm:"size(30);index"`
	Unsubtime string `orm:"size(30);index"`
	Operator  string `orm:"size(20);index"`
	Country   string `orm:"size(20);index"`
	Status    int    `orm:"index"`
	IP        string `orm:"size(30)"`

	ProductName string `orm:"size(30)"`
	TrackId     string `orm:"size(50)"`
	CampId      string `orm:"size(20);index"`
	AffName     string `orm:"size(30);index"`
	PubId       string `orm:"size(100);index"`
	ProId       string `orm:"size(10)"`
	ClickId     string `orm:"size(300)"`

	PostbackCode    string `orm:"size(10)"`
	PostbackStatus  int
	PostbackTime    string `orm:"size(30)"`
	PostbackMessage string `orm:"size(30)"`
}

type MtStruct struct {
	Id        int64  `orm:"pk;auto"`
	Moid      string `orm:"size(50);index"`
	Shortcode string `orm:"size(20);index"`
	Keyword   string `orm:"size(20);index"`
	Msisdn    string `orm:"size(20);index"`
	Operator  string `orm:"size(20);index"`
	Country   string `orm:"size(20);index"`
	Time      string `orm:"size(30)"`
	CampId    string `orm:"size(20);index"`
	Charge    string
	Subid     int64

	ResponseStatus    string `orm:"size(20)"`
	ResponseMessageid string `orm:"size(20)"`
	ResponseErrorcode string `orm:"size(20)"`
}

type DnStruct struct {
	Id      int64  `orm:"pk;auto"`
	Keyword string `orm:"size(20);index"`
	Price   string `orm:"size(10);;index"`
	Time    string `orm:"size(30);index"`
	SubId   int64
	AffName string `orm:"size(50);index"`
	PubId   string `orm:"size(100);index"`
	Charge  string

	DnStatus int `orm:"index"`

	RequestId       string `orm:"size(100)"`
	TransactionId   string `orm:"size(100)"`
	Shortcode       string `orm:"size(100)"`
	ChannelId       string `orm:"size(100)"`
	ApplicationId   string `orm:"size(100)"`
	Country         string `orm:"size(100)"`
	OperatorId      string `orm:"size(100)"`
	Msisdn          string `orm:"size(100)"`
	Mtid            string `orm:"size(100)"`
	ActivityTime    string `orm:"size(100)"`
	SubscriptionEnd string `orm:"size(100)"`
	Type            string `orm:"size(100)"`
	SubType         string `orm:"size(100)"`
	Status          string `orm:"size(100)"`
	Rate            string `orm:"size(100)"`
}

//func init() {
//	orm.RegisterModel(new(Track), new(Old_Postback), new(PostbackRecord), new(MoStruct), new(DnStruct), new(MtStruct), new(sp.Postback), new(sp.AffTrack))
//}
