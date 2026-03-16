package conn

import (
	beckon "magic.pathao.com/beatbox/beckon"
	"magic.pathao.com/parcel/prism/internal/config"
	"magic.pathao.com/parcel/prism/internal/monitoring/beatbox"
)

func ConnectBeatbox(cfg *config.Config) *beatbox.BeatBox {
	conf := cfg.BeatBox
	beatbox := &beatbox.BeatBox{}
	beckon.SetStatsd(conf.Url)
	beatbox.Init()
	return beatbox
}
