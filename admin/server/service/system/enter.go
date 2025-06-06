package system

type ServiceGroup struct {
	JwtService
	ApiService
	MenuService
	UserService
	UserExtendService   // 新增OA登录
	UserRegisterService // 用户自主注册
	CasbinService
	InitDBService
	AutoCodeService
	BaseMenuService
	AuthorityService
	DictionaryService
	SystemConfigService
	OperationRecordService
	DictionaryDetailService
	AuthorityBtnService
	SysExportTemplateService
	SysParamsService
	AutoCodePlugin   autoCodePlugin
	AutoCodePackage  autoCodePackage
	AutoCodeHistory  autoCodeHistory
	AutoCodeTemplate autoCodeTemplate
}
