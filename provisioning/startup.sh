 #! /bin/bash
 
sudo su -
apt update
apt install git
git clone https://github.com/mittz/role-play-webapp.git
cd ~/role-play-webapp
bash provisioning/os-setup-script.sh
source ~/.bashrc
source /etc/profile
systemctl start docker
cd ~/role-play-webapp/webapp
make upgrade-compose
make all

