package system

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DBApi
	JwtApi
	BaseApi
	SystemApi
	CasbinApi
	AutoCodeApi
	SystemApiApi
	AuthorityApi
	DictionaryApi
	AuthorityMenuApi
	OperationRecordApi
	DictionaryDetailApi
	AuthorityBtnApi
	SysExportTemplateApi
	AutoCodePluginApi
	AutoCodePackageApi
	AutoCodeHistoryApi
	AutoCodeTemplateApi
	SysParamsApi
	UserRegisterApi
}

var (
	apiService              = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService              = service.ServiceGroupApp.SystemServiceGroup.JwtService
	menuService             = service.ServiceGroupApp.SystemServiceGroup.MenuService
	userService             = service.ServiceGroupApp.SystemServiceGroup.UserService
	userExtendService       = service.ServiceGroupApp.SystemServiceGroup.UserExtendService
	userRegisterService     = service.ServiceGroupApp.SystemServiceGroup.UserRegisterService
	initDBService           = service.ServiceGroupApp.SystemServiceGroup.InitDBService
	casbinService           = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	baseMenuService         = service.ServiceGroupApp.SystemServiceGroup.BaseMenuService
	authorityService        = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	dictionaryService       = service.ServiceGroupApp.SystemServiceGroup.DictionaryService
	authorityBtnService     = service.ServiceGroupApp.SystemServiceGroup.AuthorityBtnService
	systemConfigService     = service.ServiceGroupApp.SystemServiceGroup.SystemConfigService
	sysParamsService        = service.ServiceGroupApp.SystemServiceGroup.SysParamsService
	operationRecordService  = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
	dictionaryDetailService = service.ServiceGroupApp.SystemServiceGroup.DictionaryDetailService
	autoCodeService         = service.ServiceGroupApp.SystemServiceGroup.AutoCodeService
	autoCodePluginService   = service.ServiceGroupApp.SystemServiceGroup.AutoCodePlugin
	autoCodePackageService  = service.ServiceGroupApp.SystemServiceGroup.AutoCodePackage
	autoCodeHistoryService  = service.ServiceGroupApp.SystemServiceGroup.AutoCodeHistory
	autoCodeTemplateService = service.ServiceGroupApp.SystemServiceGroup.AutoCodeTemplate
)
