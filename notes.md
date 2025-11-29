 Input Types (Request body for mutations):what data the user must send during mutations.
input PaginationInput {
  skip: Int
  take: Int
}
Used for pagination:
skip → how many items to skip
take → how many items to return(Like LIMIT/OFFSET)



Steps for grpc file generation -

    wget https://github.com/protocolbuffers/protobuf/releases/download/v23.0/protoc-23.0-linux-x86_64.zip
    unzip protoc-23.0-linux-x86_64.zip -d protoc
    sudo mv protoc/bin/protoc /usr/local/bin/
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    echo $PATH
    export PATH="$PATH:$(go env GOPATH)/bin"
    source ~/.bashrc
    create the pb folder in your project and go to parent of this pb folder and then follow the below steps
    add this to account.proto - option go_package = "./pb";
    finally run this command - protoc --go_out=./pb --go-grpc_out=./pb account.proto
