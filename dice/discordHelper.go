package dice

import (
	"github.com/google/uuid"
	"time"
)

// NewDiscordConnItem 本来没必要写这个的，但是不知道为啥依赖出问题
func NewDiscordConnItem(token string) *EndPointInfo {
	conn := new(EndPointInfo)
	conn.Id = uuid.New().String()
	conn.Platform = "DISCORD"
	conn.ProtocolType = ""
	conn.Enable = false
	conn.RelWorkDir = "extra/discord-" + conn.Id
	conn.Adapter = &PlatformAdapterDiscord{
		EndPoint: conn,
		Token:    token,
	}
	return conn
}

// ServeDiscord gocqhttp_helper 中有一个相同的待重构方法，为了避免阻碍重构，先不写在一起了
func ServeDiscord(d *Dice, ep *EndPointInfo) {
	defer CrashLog()
	if ep.Platform == "DISCORD" {
		conn := ep.Adapter.(*PlatformAdapterDiscord)
		d.Logger.Infof("DiscordGo 尝试连接")
		if conn.Serve() != 0 {
			d.Logger.Errorf("连接Discord服务失败")
			ep.State = 3
			ep.Enable = false
			d.LastUpdatedTime = time.Now().Unix()
			d.Save(false)
		}
	}
}
