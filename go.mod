module github.com/yakumioto/hlf-deploy

go 1.13

require (
	github.com/Shopify/sarama v1.24.1 // indirect
	github.com/fsouza/go-dockerclient v1.6.0 // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hyperledger/fabric v1.4.3
	github.com/hyperledger/fabric-amcl v0.0.0-20190902191507-f66264322317 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20191121202242-f5500d5e3e85
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta1.0.20191231170015-e7b9b0dc1316
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/sykesm/zap-logfmt v0.0.3 // indirect
	go.uber.org/zap v1.13.0 // indirect
)

exclude (
	github.com/hyperledger/fabric v1.4.4
	github.com/hyperledger/fabric v2.0.0-alpha+incompatible
	github.com/hyperledger/fabric v2.0.0-beta+incompatible
)
