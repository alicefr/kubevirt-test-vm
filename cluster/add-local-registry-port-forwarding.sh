#!/bin/bash

ssh -N -f -g -R 5000:localhost:5001 -p 32756 root@localhost
