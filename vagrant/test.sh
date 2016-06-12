DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..
vagrant init ubuntu/trusty64 || echo "has already init"
vagrant up
vagrant ssh -c "sh /vagrant/vagrant/setup.sh"
sh $DIR/run.sh
