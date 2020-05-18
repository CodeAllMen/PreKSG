/**
  create by yy on 2020/5/17
*/

package sp

import "github.com/astaxie/beego"

type CountController struct {
	beego.Controller
}

func (c *CountController) Count() {
	c.Data["json"] = "sad"

	c.ServeJSON()
}
