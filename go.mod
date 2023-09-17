module partyframe

go 1.20

require (
	github.com/go-micro/plugins/v4/logger/zap v1.2.1
	go-micro.dev/v4 v4.10.2
	go.uber.org/zap v1.26.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace go-micro.dev/v4 => github.com/cherry-cup/go-micro/v4 v4.10.3-0.20230917064437-dd998711792c
