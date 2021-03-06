/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

var DockerHealthCheck = `#!/bin/bash

# Copyright 2019 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script is intended to be run periodically, to check the health
# of docker.  If it detects a failure, it will restart docker using systemctl.

healthcheck() {
   if output=` + "`timeout 60 docker network ls`" + `; then
       echo "$output" | fgrep -qw host || {
          echo "docker 'host' network missing"
          return 1
       }
   else
       echo "docker returned $?"
       return 1
   fi
}

if healthcheck; then
  echo "docker healthy"
  exit 0
fi

echo "docker failed"
echo "Giving docker 30 seconds grace before restarting"
sleep 30

if healthcheck; then
  echo "docker recovered"
  exit 0
fi

echo "docker still unresponsive; triggering docker restart"
systemctl stop docker

echo "wait all tcp sockets to close"
sleep ` + "`cat /proc/sys/net/ipv4/tcp_fin_timeout`" + `

sleep 10
systemctl start docker

echo "Waiting 120 seconds to give docker time to start"
sleep 60

if healthcheck; then
  echo "docker recovered"
  exit 0
fi

echo "docker still failing"
`
