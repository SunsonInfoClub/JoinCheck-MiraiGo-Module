package memcheck

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/jinzhu/gorm"

	//mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type joinChecker struct {
}

var instance *joinChecker
var logger = utils.GetModuleLogger("sic.joincheck")

//StuInfo 学生信息
type StuInfo struct {
	Name  string
	Grade int
	Class int
	QQ    int64
	IP    string
}

//Member 社团成员
type Member struct {
	StuInfo
	gorm.Model
}

func init() {
	instance = &joinChecker{}
	bot.RegisterModule(instance)
}

func (c *joinChecker) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "sic.joincheck",
		Instance: instance,
	}
}
func (c *joinChecker) Init() {
	logger.Info("SunsonCheck初始化完成")
}

func (c *joinChecker) PostInit() {}

func (c *joinChecker) Serve(b *bot.Bot) {
	register(b)
}
func (c *joinChecker) Start(b *bot.Bot) {}
func (c *joinChecker) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func register(b *bot.Bot) {
	b.OnUserWantJoinGroup(CheckonJoin)
}

//CheckonJoin 验证身份
func CheckonJoin(qqClient *client.QQClient, event *client.UserJoinGroupRequest) {
	logger.Infoln("新加群申请:", *event)
	if checkUin(event.RequesterUin) && isTheGroup(event.GroupCode) {
		event.Accept()
	}
}

func checkUin(uin int64) bool {
	db, err := gorm.Open("mysql", config.GlobalConfig.GetString("sic.checker.DBSource")+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return false
	}
	mem := Member{}
	db.Where("qq=?", uin).First(&mem)
	if mem.Name == "" {
		return false
	}
	return true
}

func isTheGroup(uin int64) bool {
	return uin == config.GlobalConfig.GetInt64("sic.checker.GroupUin")
}
