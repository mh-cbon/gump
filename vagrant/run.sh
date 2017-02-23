
go build -o build/gump gump.go

GUMPDIR="~/gow/src/github.com/mh-cbon/gump"
GOBIN="GOPATH=/home/vagrant/gow GOROOT=/home/vagrant/go ~/go/bin/go"

vagrant ssh -c "sh /vagrant/vagrant/git.sh && cd $GUMPDIR && $GOBIN test && cd $GUMPDIR/config && $GOBIN test"
