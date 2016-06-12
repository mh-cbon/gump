
go build -o build/gump gump.go

vagrant ssh -c "sh /vagrant/vagrant/git.sh && cd ~/gow/src/github.com/mh-cbon/gump && GO15VENDOREXPERIMENT=1 GOPATH=/home/vagrant/gow GOROOT=/home/vagrant/go ~/go/bin/go test"
