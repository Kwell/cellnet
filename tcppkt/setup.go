package tcppkt

import (
	"github.com/davyxu/cellnet"
	"net"
)

func ProcSysMsg(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case cellnet.SessionConnectErrorEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnectError{}})
		case cellnet.SessionClosedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionClosed{}})
		case cellnet.SessionAcceptedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionAccepted{}})
		case cellnet.SessionConnectedEvent:
			userFunc(cellnet.RecvMsgEvent{Ses: ev.Ses, Msg: &SessionConnected{}})
		}

		return userFunc(raw)
	}
}

func ProcTLVPacket(userFunc cellnet.EventFunc) cellnet.EventFunc {

	return func(raw cellnet.EventParam) cellnet.EventResult {

		switch ev := raw.(type) {

		case cellnet.RecvEvent: // 接收数据事件

			if result := onRecvLTVPacket(ev.Ses, userFunc); result != nil {
				return result
			}

		case cellnet.SendMsgEvent: // 发送数据事件

			if result := onSendLTVPacket(ev.Ses, userFunc, ev.Msg); result != nil {
				return result
			}

		case cellnet.SessionCleanupEvent:
			// 取Socket连接
			conn, ok := ev.Ses.Raw().(net.Conn)

			// 转换错误，或者连接已经关闭时退出
			if ok && conn != nil {
				conn.Close()
			}
		}

		return userFunc(raw)
	}
}