package application

import (
	"github.com/boomfunc/root/base/conf"
)

func JSON(router *conf.Router) Interface {
	return &Application{
		router: router,
		packer: new(jsonPacker),
	}
}

func HTTP(router *conf.Router) Interface {
	return &Application{
		router: router,
		packer: new(httpPacker),
	}
}
