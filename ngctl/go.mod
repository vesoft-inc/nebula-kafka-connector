module github.com/vesoft-inc/nebula-ng-tools/ngctl

go 1.21

toolchain go1.22.4

replace github.com/vesoft-inc/nebula-ng-tools/golang => ../golang

replace github.com/vesoft-inc/nebula-ng-tools/ngadm => ../ngadm

replace github.com/vesoft-inc/nebula-ng-tools/agent => ../agent

require (
	github.com/jedib0t/go-pretty/v6 v6.5.9
	github.com/manifoldco/promptui v0.9.0
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.9.0
	github.com/vesoft-inc/nebula-ng-tools/golang v0.0.0
	github.com/vesoft-inc/nebula-ng-tools/ngadm v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/ScaleFT/sshkeys v1.2.0 // indirect
	github.com/appleboy/easyssh-proxy v1.5.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/bcrypt_pbkdf v0.0.0-20150205184540-83f37f9c154a // indirect
	github.com/go-resty/resty/v2 v2.10.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/unknwon/goconfig v1.0.0 // indirect
	github.com/vesoft-inc/go-pkg v0.0.0-20231117110005-307b542ecb31 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/grpc v1.66.2 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
