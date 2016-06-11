sudo apt-get install -y curl make binutils bison gcc build-essential
sudo apt-get install -y git

if [ ! -f go1.6.2.linux-amd64.tar.gz ]; then
  curl -s -S -L https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -o go1.6.2.linux-amd64.tar.gz
  tar -xf go1.6.2.linux-amd64.tar.gz
fi

mkdir -p ~/gow/src/github.com/mh-cbon
ln -s /vagrant ~/gow/src/github.com/mh-cbon/gump
