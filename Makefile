HOSTNAME=github.com
NAMESPACE=roperto
NAME=hue
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64

default: install

clean:
	rm -rv ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH} || true

build: clean
	go install -gcflags="all=-N -l"

install: build
	rm -fr definitions/.terraform/providers
	rm -f definitions/.terraform.lock.hcl

apply: install
	cd definitions && terraform apply
