#!/usr/bin/env bash

sudo rm /var/lib/dpkg/lock
sudo dpkg --configure -a

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF

apt update

apt-get install -y docker.io
apt-get install -y apt-transport-https
apt-get install -y kubelet kubeadm kubectl