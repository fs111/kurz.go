ifeq ($(MYTARGDIR),)
	MYTARGDIR:=target
endif

CONF_NAME=kurz.conf

ifeq ($(STATIC_DIR),)
	STATIC_DIR:=share/kurz
endif

TARG=src/kurz

CLEANFILES=$(MYTARGDIR)

all: bin-dist

$(TARG): src/*.go
	@go build -o $(TARG) src/*.go

test: src/*.go
	go test src/*.go

clean:
	@rm -rf $(CLEANFILES)
	@rm -rf $(TARG)

bin-dist: $(TARG) assets
	@cp -rf conf/kurz.conf $(MYTARGDIR)/etc/kurz/
	@sed 's?=static?=$(STATIC_DIR)?' conf/$(CONF_NAME) > $(MYTARGDIR)/etc/kurz/$(CONF_NAME)
	@cp $(TARG) $(MYTARGDIR)/bin
	@git log --pretty=format:"kurz.go %H" -1 > $(MYTARGDIR)/$(STATIC_DIR)/_version


assets: directories
	@cp -r stuff/assets/* $(MYTARGDIR)/$(STATIC_DIR)


directories:
	@mkdir -p $(MYTARGDIR)/$(STATIC_DIR)
	@mkdir -p $(MYTARGDIR)/etc/kurz/
	@mkdir -p $(MYTARGDIR)/bin

.PHONY: clean bin-dist assets directories test
