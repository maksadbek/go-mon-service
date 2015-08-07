#!/bin/bash
case "$1" in
    start)
        /usr/local/bin/gomon -conf=/etc/gomon/conf.toml
        ;;
    stop)
        /usr/local/bin/gomon -s stop
        ;;
    *)
      echo "Usage: /etc/init.d/blah {start|stop}"
      exit 1
esac
exit 0
