#
#  Copyright 2018 Nalej
# 

# Name of the target applications to be built
APPS=system-model system-model-cli

# Use global Makefile for common targets
export
%:
	$(MAKE) -f Makefile.golang $@
