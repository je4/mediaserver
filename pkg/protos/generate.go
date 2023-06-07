//go:generate protoc.exe ../../data/protobuf/database.proto -I ../../data/protobuf/ --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative

package protos
