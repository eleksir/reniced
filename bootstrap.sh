#!/bin/bash

# Защита от дурака
if [[ -z $HOME ]]; then
	echo "HOME variable is not defined, please set it up and re-run bootstrap script"
	exit 1
fi

cpanm --installdeps --local-lib "$HOME/perl5"

# vim: set ft=sh noet ai ts=4 sw=4 sts=4: