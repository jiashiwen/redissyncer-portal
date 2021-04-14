module redissyncer-portal

go 1.15

require (
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc // indirect
	github.com/apache/calcite-avatica-go/v3 v3.2.0 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/cznic/fileutil v0.0.0-20181122101858-4d67cfea8c87 // indirect
	github.com/cznic/golex v0.0.0-20181122101858-9c343928389c // indirect
	github.com/cznic/internal v0.0.0-20181122101858-3279554c546e // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/ql v1.2.0 // indirect
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8 // indirect
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537 // indirect
	github.com/cznic/zappy v0.0.0-20181122101859-ca47d358d4b1 // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.3.0+incompatible // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/kshvakov/clickhouse v1.3.5 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/tidwall/gjson v1.6.8
	go.uber.org/zap v1.16.0
	google.golang.org/grpc/examples v0.0.0-20210311221743-f168a3cb3bf5 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/rana/ora.v4 v4.1.15 // indirect
	gopkg.in/yaml.v2 v2.3.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)

//replace github.com/coreos/etcd => go.etcd.io/etcd/v3 v3.5.0-alpha.0

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
