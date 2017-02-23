PREBUMP=
  666 git fetch --tags origin master
  666 git pull origin master

PREVERSION=
  philea -s "666 go vet %s" "666 go-fmt-fail %s"
  666 go run gump.go -v
  666 go test
  666 go test config/*
  666 go test stringexec/*
  666 changelog finalize --version !newversion!
  666 commit -q -m "changelog: !newversion!" -f change.log

POSTVERSION=
  666 changelog md -o CHANGELOG.md --guess
  666 commit -q -m "changelog: !newversion!" -f CHANGELOG.md
  666 go install --ldflags "-X main.VERSION=!newversion!"
  666 emd gen -out README.md
  666 commit -q -m "README: !newversion!" -f README.md
  666 git push
  666 git push --tags
  666 gh-api-cli create-release -n release -o mh-cbon --guess \
   --ver !newversion! -c "changelog ghrelease --version !newversion!" \
  --draft !isprerelease!
