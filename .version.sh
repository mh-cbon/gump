PREBUMP=
666 git fetch --tags origin master
666 git pull origin master

PREVERSION=
philea -s "666 go vet %s" "666 go-fmt-fail %s"
666 go run gump.go -v
666 changelog finalize --version !newversion!
666 commit -q -m "changelog: !newversion!" -f change.log

POSTVERSION=
666 changelog md -o CHANGELOG.md --vars='{"name":"gump"}'
666 commit -q -m "changelog: !newversion!" -f CHANGELOG.md
666 git push
666 git push --tags
666 gh-api-cli create-release -n release -o mh-cbon -r gump \
 --ver !newversion! -c "changelog ghrelease --version !newversion!" \
--draft !isprerelease!
666 go install --ldflags "-X main.VERSION=!newversion!"
