GOBIN ?= ~/go/bin

all:
	yarn build
	(cd pkg && make)
	${GOBIN}/mage -v
