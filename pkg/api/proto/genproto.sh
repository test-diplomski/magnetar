protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go_opt=Mmagnetar_model.proto=github.com/c12s/magnetar/pkg/api \
        magnetar_model.proto

protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go-grpc_out=../ \
        --go-grpc_opt=paths=source_relative \
        --go_opt=Mmagnetar.proto=github.com/c12s/magnetar/pkg/api \
        --go-grpc_opt=Mmagnetar.proto=github.com/c12s/magnetar/pkg/api \
        -I ./magnetar_model.proto \
        magnetar.proto

protoc --proto_path=./ \
        --go_out=../ \
        --go_opt=paths=source_relative \
        --go_opt=Mmagnetar.proto=github.com/c12s/magnetar/pkg/api \
        -I ./magnetar_model.proto \
        registration.proto