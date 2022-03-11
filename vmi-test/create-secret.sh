#!/bin/bash -x

KEY_FILE=test-key
ssh-keygen -t ed25519 -C "test" -N "" -f $KEY_FILE
ssh-agent $(pwd)/$KEY_FILE
kubectl get secret $KEY_FILE
if [ $? -eq 0 ]; then 
	kubectl delete secret $KEY_FILE
fi
kubectl create secret generic $KEY_FILE --from-file=$(pwd)/$KEY_FILE.pub
