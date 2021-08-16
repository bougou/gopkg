package screenshot

import "testing"

func TestTakeScreenShotToImageFile(t *testing.T) {
	// url := "http://xpmonitor-grafana.bj.govcloud.tencent.com/d/ttnbzATZz2/zhu-ji-ji-chu-jian-kong?orgId=1&kiosk=tv"
	url := "https://v.qq.com/"
	filename := "baidu.png"

	TakeScreenShotToImageFile(url, filename, "png")

	t.Error()

}
