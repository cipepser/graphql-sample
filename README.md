# graphql-sample

[Real-time Chat with GraphQL Subscriptions in Go](https://outcrawl.com/go-graphql-realtime-chat/)の写経

## gqlgenのインストール

```sh
❯ go get -u github.com/vektah/gqlgen
```

2018/07/16現在だと互換性が崩れている。
いったん以下でやったけど、まだ`generated.go`がだめ。

```sh
❯ git clone https://github.com/vektah/gqlgen.git
❯ git checkout 325c45a40b41dec948abd1138cc8f84ae815b285
```

### 実装

スキーマを定義する。

`server/schema.graphql`を追加
`server/graphql.go`を追加  

`server/graphql.go`は一旦以下だけあれば、generateできる。

```go
//go:generate gqlgen -schema ./schema.graphql
package server
```

## Docker

`docker-compose.yaml`は変更なし。
`Dockerfile`は`WORKDIR`だけ自分のものに変える。

```sh
❯ atom Dockerfile
```

```Dockerfile
WORKDIR /go/src/github.com/cipepser/graphql-sample
```

### vgo

`Dockerfile`に書いたようにパッケージを`vendor`に入れておく。

```sh
❯ vgo mod -vendor
```

`app`と`redis`を起動。

```sh
❯ docker-compose up -d
```

## frontend

今回は[もとのディレクトリ](https://github.com/tinrab/graphql-realtime-chat/tree/master/frontend)をコピー

```sh
❯ cd frontend
```

`yarn`をインストールしておく。

```sh
❯ brew install yarn
```

front側を立ち上げる。

```sh
❯ yarn serve
yarn run v1.7.0
$ vue-cli-service serve
 INFO  Starting development server...
 98% after emitting CopyPlugin

 DONE  Compiled successfully in 3556ms                                                                                                                         12:56:38


  App running at:
  - Local:   http://localhost:3000/
  - Network: http://192.168.100.115:3000/

  Note that the development build is not optimized.
  To create a production build, run yarn build.
```


最初うまくいかなかったので、以下あたりを実行した。

```sh
npm install vue
npm install -g npm
npm uninstall -g vue-cli
npm install
npm audit fix
```

## References
* [Real-time Chat with GraphQL Subscriptions in Go](https://outcrawl.com/go-graphql-realtime-chat/)