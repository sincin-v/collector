build:
	go build -buildvcs=false -o /home/sinicin/projects/practic/collector/cmd/agent/agent /home/sinicin/projects/practic/collector/cmd/agent/main.go
	go build -buildvcs=false -o /home/sinicin/projects/practic/collector/cmd/server/server /home/sinicin/projects/practic/collector/cmd/server/main.go

tests: build
	go test -v ./...

linter:
	go vet -vettool=/usr/local/go/bin/statictest ./...
	golangci-lint run ./...

iter1: linter tests
	metricstest -test.v -binary-path=cmd/server/server -test.run=TestIteration1$

iter2: iter1
	metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent  -test.run=TestIteration2

iter3: iter2
	metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -test.run=TestIteration3

iter4: iter3
	metricstest -test.v -source-path=. -server-port=8888 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -test.run=TestIteration4

iter5: iter4
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -test.run=TestIteration5

sprint1: iter1 iter2 iter3 iter4 iter5


iter6: linter tests
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -test.run=TestIteration6

iter7: iter6
	SERVER_PORT=8080 ADDRESS="localhost:8080" metricstest -test.v -source-path=. -server-port=8080 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -test.run=TestIteration7$

iter8: iter7
	SERVER_PORT=8888 ADDRESS="localhost:8888" metricstest -test.v -source-path=. -server-port=8888 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -test.run=TestIteration8$


iter9: iter8
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -file-storage-path=/tmp/metric_storage_file.json -test.run=TestIteration9

sprint2: sprint1 iter6 iter7 iter8 iter9

iter10: linter tests
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -database-dsn='postgres://db_user:N6KAKgEDa27aot@0.0.0.0:5432/praktikum?sslmode=disable' -file-storage-path=/tmp/metric_storage_file.json -test.run=TestIteration10

iter11: iter10
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -database-dsn='postgres://db_user:N6KAKgEDa27aot@0.0.0.0:5432/praktikum?sslmode=disable' -file-storage-path=/tmp/metric_storage_file.json -test.run=TestIteration11

iter12: iter11
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -database-dsn='postgres://db_user:N6KAKgEDa27aot@0.0.0.0:5432/praktikum?sslmode=disable' -file-storage-path=/tmp/metric_storage_file.json -test.run=TestIteration12

iter13: iter12
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -database-dsn='postgres://db_user:N6KAKgEDa27aot@0.0.0.0:5432/praktikum?sslmode=disable' -file-storage-path=/tmp/metric_storage_file.json -test.run=TestIteration13

iter14: iter13
	SERVER_PORT=8989 ADDRESS="localhost:8989" metricstest -test.v -source-path=. -server-port=8989 -agent-binary-path=cmd/agent/agent   -binary-path=cmd/server/server -database-dsn='postgres://db_user:N6KAKgEDa27aot@0.0.0.0:5432/praktikum?sslmode=disable' -file-storage-path=/tmp/metric_storage_file.json -key=/tmp/tmp_file -test.run=TestIteration14

sprint3: sprint1 sprint2 iter10 iter11 iter12 iter13 iter14

ci: build linter
