//go:generate protoc.exe database.proto --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative

package protos
