#!/bin/bash

ssh -N -f -g -R 5000:localhost:5000 -p 32756 root@localhost
