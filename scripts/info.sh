#!/bin/sh

VERSION="$1"

echo "[DONE] Now mcc version is $1"

cat << EOF

 ▄▄       ▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄ 
▐░░▌     ▐░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
▐░▌░▌   ▐░▐░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀ 
▐░▌▐░▌ ▐░▌▐░▌▐░▌          ▐░▌          
▐░▌ ▐░▐░▌ ▐░▌▐░▌          ▐░▌          
▐░▌  ▐░▌  ▐░▌▐░▌          ▐░▌          
▐░▌   ▀   ▐░▌▐░▌          ▐░▌          
▐░▌       ▐░▌▐░▌          ▐░▌          
▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ 
▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
 ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀ 
                                       
mcc — Multi Cluster Changer
---------------------------

mcc is a simple CLI tool to manage multiple Kubernetes cluster configs.
It stores kubeconfig files in ~/.kube/multi and switches them into ~/.kube/config when needed.

Examples:
  mcc add --file ./myconfig --name cluster1   # Add a new cluster config
  mcc ch cluster1                             # Switch to the cluster
  mcc list                                    # List stored clusters
  mcc delete cluster1                         # Remove a cluster config

Simple Usage
  sudo cp mcc /usr/local/bin/

EOF