```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
```

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

```
go get -u google.golang.org/grpc
```


大体、protobufで作られるメソッドは決まってる


クライアント側
```go
client := NewHogeServiceClient 
client.Bar
```


サーバー側
```go
RegistorHogeServiceServer 
GetHoge
```

サービス系がhoge_grpc.pb.goでそれ以外がhoge.pb.goかな

関数の実装は"server interface"で検索すればおk

stream系は基本ループで送受信する。io.EOFでbreak