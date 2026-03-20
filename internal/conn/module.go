package conn

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		ConnectPostgres,
		ConnectRedis,
	),
)