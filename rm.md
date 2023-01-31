 
 protoc --go_out=server/proto --go_opt=paths=source_relative \   
    --go-grpc_out=server/proto --go-grpc_opt=paths=source_relative \   
    api.proto

    //client
    protoc --go_out=client/proto --go_opt=paths=source_relative \   
    --go-grpc_out=client/proto --go-grpc_opt=paths=source_relative \   
    api.proto
