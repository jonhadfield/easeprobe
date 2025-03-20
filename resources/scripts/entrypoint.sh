#!/bin/sh

# Support the following running mode
# 1) run guardianprobe without any arguments
# 2) run guardianprobe with guardianprobe arguments
# 3) run the command in guardianprobe container
PROBE_CONFIG=${PROBE_CONFIG:-/opt/config.yaml}

echo "Using config file: ${PROBE_CONFIG}"

# docker run o2ip/guardianprobe
if [ "$#" -eq 0 ]; then
   exec /opt/guardianprobe
# docker run o2ip/guardianprobe -f config.yaml
elif [ "$1" != "--" ] && [ "$(echo $1 | head -c 1)" == "-" ] ; then
  exec /opt/guardianprobe "$@"
# docker run -it --rm o2ip/guardianprobe /bin/sh
# docker run -it --rm o2ip/guardianprobe -- /bin/echo hello world
else
  exec "$@"
fi
