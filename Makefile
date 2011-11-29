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


bin-dist: $(TARG)
	@mkdir -p $(MYTARGDIR)/etc/kurz/
	@mkdir -p $(MYTARGDIR)/bin
	@mkdir -p $(MYTARGDIR)/$(STATIC_DIR)
	@cp -rf conf/kurz.conf $(MYTARGDIR)/etc/kurz/
	@sed 's?=static?=$(STATIC_DIR)?' conf/$(CONF_NAME) > $(MYTARGDIR)/etc/kurz/$(CONF_NAME)
	@cp stuff/assets/* $(MYTARGDIR)/$(STATIC_DIR)
	@cp $(TARG) $(MYTARGDIR)/bin
	@git log --pretty=format:"kurz.go %H" -1 > $(MYTARGDIR)/$(STATIC_DIR)/_version


include $(GOROOT)/src/Make.cmd
