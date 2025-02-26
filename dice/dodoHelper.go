package dice

import (
	"github.com/google/uuid"
	"time"
)

func NewDodoConnItem(clientID string, token string) *EndPointInfo {
	conn := new(EndPointInfo)
	conn.Id = uuid.New().String()
	conn.Platform = "DODO"
	conn.ProtocolType = ""
	conn.Enable = false
	conn.RelWorkDir = "extra/dodo-" + conn.Id
	conn.Adapter = &PlatformAdapterDodo{
		EndPoint: conn,
		ClientID: clientID,
		Token:    token,
	}
	return conn
}

func ServeDodo(d *Dice, ep *EndPointInfo) {
	defer CrashLog()
	if ep.Platform == "DODO" {
		conn := ep.Adapter.(*PlatformAdapterDodo)
		d.Logger.Infof("Dodo 尝试连接")
		if conn.Serve() != 0 {
			d.Logger.Errorf("连接Dodo失败")
			ep.State = 3
			ep.Enable = false
			d.LastUpdatedTime = time.Now().Unix()
			d.Save(false)
		}
	}
}
