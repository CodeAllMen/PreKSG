package sp

import (
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/postback"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(AffTrack), new(mo.Mo), new(ChargeNotification),
		new(postback.Postback), new(notification.Notification), new(click.HourClick))
}
