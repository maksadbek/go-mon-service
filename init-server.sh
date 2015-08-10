#!/bin/bash
mkdir /etc/gomon
mkdir /var/log/gomon
cp conf.toml.example /etc/gomon/conf.toml
cp gomond /etc/init.d/
cp go-mon-service /usr/local/bin/go-mon-service
update-rc.d gomond defaults
exit 0
