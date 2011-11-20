include $(GOROOT)/src/Make.inc
TARG=src/kurz

ifeq ($(MYTARGDIR),)
	MYTARGDIR:=target
endif

CLEANFILES=$(MYTARGDIR)
GOFILES=\
	src/*.go\

include $(GOROOT)/src/Make.cmd

myinstall: $(TARG)
	@mkdir -p $(MYTARGDIR)/etc/kurz/
	@mkdir -p $(MYTARGDIR)/bin
	@cp -rf conf/kurz.conf $(MYTARGDIR)/etc/kurz/
	@cp $(TARG) $(MYTARGDIR)/bin
