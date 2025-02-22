# gogrpcpeer
Simple GRPC server-client written in Go

It uses mTLS encryption, so you need firstly to generate certs.
Edit `cert/sshgen.sh` file and run it. It will prepare all necessary certs and keys.
File `sshrem.sh` is for clearing all certificates, be careful.

After install dependencies with `go get` command.
Then you have to compile proto file(s) with Protocol Buffer Compiler.
So, firstly you need to install it.
For Linux:
```
$ apt install -y protobuf-compiler
$ protoc --version
```

For MacOS:
```
$ brew install protobuf
$ protoc --version
```

After installation finished, run `protocomp.sh` script, it will run `protoc` compiler
and create all .go files you need to deal with protobuf.

So, you are ready to go.

Start server with `go run .`
Start client with `go run client/client.go`

Good luck!

P.S. If you copy prepared `cert` folder to the [TypeScript version of GRPC server-client](https://github.com/dioxine/tsgrpcpeer) and use same `.proto` file,
it will work with TypeScript as server and Go as client or visa-versa.
