#!/bin/bash

podman system connection add test --identity $(pwd)/vmi-test/test-key ssh://root@localhost:32756/run/user/0/podman/podman.sock
