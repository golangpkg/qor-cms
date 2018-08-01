package auth

import (
	"github.com/astaxie/beego/logs"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/session/manager"
	"github.com/golangpkg/qor-cms/models"
	"github.com/golangpkg/qor-cms/controllers"
)

// //########################## 定义admin 权限 ##########################
type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	logs.Info(" user not login ")
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	logs.Info(" user  logout ")
	return "/auth/logout"
}

//从session中获得当前用户。
func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	adminUserName := manager.SessionManager.Get(c.Request, controllers.USER_SESSION_NAME)
	logs.Info("########## adminUser %v", adminUserName)
	if adminUserName != "" {
		userInfo := models.UserInfo{}
		userInfo.UserName = adminUserName
		return &userInfo
	}
	return nil
}
