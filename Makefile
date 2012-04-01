ifeq ($(MYTARGDIR),)
	MYTARGDIR:=target
endif

CONF_NAME=kurz.conf

ifeq ($(PREFIX),)
	PREFIX:=/usr
endif

ifeq ($(STATIC_DIR),)
	STATIC_DIR:=$(PREFIX)/share/kurz
endif

TARG=src/kurz

CLEANFILES=$(MYTARGDIR)


all: bin-dist


$(TARG):
	@go build -o $(TARG) src/*.go


clean:
	@rm -r $(CLEANFILES)
	@rm -r $(TARG)

bin-dist: $(TARG) assets
	@cp -rf conf/kurz.conf $(MYTARGDIR)/etc/kurz/
	@cp -rf stuff/init-script/kurz $(MYTARGDIR)/etc/rc.d/init.d/
	@sed 's?=static?=$(STATIC_DIR)?' conf/$(CONF_NAME) > $(MYTARGDIR)/etc/kurz/$(CONF_NAME)
	@cp $(TARG) $(MYTARGDIR)/$(PREFIX)/bin
	@git log --pretty=format:"kurz.go %H" -1 > $(MYTARGDIR)/$(STATIC_DIR)/_version


assets: directories
	@cp -r stuff/assets/* $(MYTARGDIR)/$(STATIC_DIR)


directories:
	@mkdir -p $(MYTARGDIR)/$(STATIC_DIR)
	@mkdir -p $(MYTARGDIR)/etc/kurz/
	@mkdir -p $(MYTARGDIR)/etc/rc.d/init.d/
	@mkdir -p $(MYTARGDIR)/$(PREFIX)/bin

.PHONY: clean bin-dist assets directories kurz
