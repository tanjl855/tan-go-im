package log

import "testing"

func TestNewLog(t *testing.T) {
	InitLog("D:\\Goland_project\\tan-go-im\\log\\im_log.log", "INFO", "json", true)
	Panic("lalala", Log)
}
