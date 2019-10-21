package sp

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config 内容站配置
type Config struct {
	Service map[string]ServiceInfo
}

type ServiceInfo struct {
	ServiceID     string `yaml:"service_id" orm:"pk;column(service_id)"`
	ServiceName   string `yaml:"service_name"`
	ShortCode     string `yaml:"short_code"`
	OperatorId    string `yaml:"operator_id"`
	ApiKey        string `yaml:"api_key"`
	ApiSecret     string `yaml:"api_secret"`
	RUrl          string `yaml:"rurl"`
	KeyWord       string `yaml:"key_word"`
	ProductName   string `yaml:"product_name"`
	Description   string `yaml:"description"`
	DescriptionAr string `yaml:"description_ar"`
	Service       string `yaml:"service"`
	Content       string `yaml:"content"`
	ContentAr     string `yaml:"content_ar"`
	LimitSubNum   int    `yaml:"limit_sub_num"`
	CampID        int    `yaml:"camp_id"`
	ApplicationId string `yaml:"application_id"`
	CountryId     string `yaml:"country_id"`
	CpId          string `yaml:"cp_id"`
	UrlPost       string `yaml:"url_post"`
	Price         string `yaml:"price"`
}

const (
	WapIdentifyUser int = iota + 1
	GetUser
	WapAuthorize
	GetSubscription
	CloseSubscription
)

var ServiceData = make(map[string]ServiceInfo)

func (server *ServiceInfo) TableName() string {
	return "server_info"
}

func InitServiceConfig() {
	filename, _ := filepath.Abs("resource/config/conf.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
	ServiceData = config.Service
}

type CommandParameter struct {
	Types          int
	TrackID        string
	IP             string
	Uid            string
	SessionID      string
	SubscriptionId string
}

//// GetDifferentCommandParameter NTH不同的步骤请求不同的接口  Command
//func GetDifferentCommandParameter(paramer CommandParameter, service ServiceInfo) (Parameters string) {
//	switch paramer.Types {
//	case WapIdentifyUser:
//		Parameters = fmt.Sprintf("command=wapIdentifyUser&username=%s&password=%s&serviceCode=%s&userIp=%s"+
//			"&callbackUrl=%s", service.Username, service.Password, service.ServiceID, paramer.IP,
//			url.QueryEscape(service.WapIdentifyUserCallbackURL+paramer.TrackID))
//	case GetUser:
//		Parameters = fmt.Sprintf("command=getUser&username=%s&password=%s&serviceCode=%s&uid=%s",
//			service.Username, service.Password, service.ServiceID, paramer.Uid)
//	case WapAuthorize:
//		Parameters = fmt.Sprintf("command=wapAuthorize&username=%s&password=%s&serviceCode=%s&price="+
//			"%s&callbackUrl=%s&uid=%s&serviceUrl=%s&notificationUrl=%s", service.Username, service.Password, service.ServiceID,
//			service.Price, url.QueryEscape(service.WapAuthorizeCallbackURL+paramer.TrackID), paramer.Uid, url.QueryEscape(service.ContentURL),
//			url.QueryEscape(service.NotificationURL))
//	case GetSubscription:
//		Parameters = fmt.Sprintf("command=getSubscription&username=%s&password=%s&serviceCode=%s&sessionId=%s",
//			service.Username, service.Password, service.ServiceID, paramer.SessionID)
//	case CloseSubscription:
//		Parameters = fmt.Sprintf("command=closeSubscription&username=%s&password=%s&serviceCode=%s&subscriptionId=%s",
//			service.Username, service.Password, service.ServiceID, paramer.SubscriptionId)
//	}
//	logs.Info(Parameters)
//	return
//}
