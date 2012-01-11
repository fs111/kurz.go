include $(GOROOT)/src/Make.inc
TARG=src/kurz

ifeq ($(MYTARGDIR),)
	MYTARGDIR:=target
endif

CONF_NAME=kurz.conf

ifeq ($(STATIC_DIR),)
	STATIC_DIR:=share/kurz
endif

CLEANFILES=$(MYTARGDIR)
GOFILES=\
	src/*.go\

all: bin-dist

.PHONY: clean bin-dist assets directories

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

include $(GOROOT)/src/Make.cmd
