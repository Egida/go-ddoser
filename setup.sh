sudo apt install golang git
git clone https://github.com/kenegdane/go-ddoser
go get "github.com/gamexg/proxyclient"
go mod init go-ddoser
go mod tidy
go build go-ddoser.go
