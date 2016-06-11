rm -fr ~/git
rm -fr ~/git_with_ok_scripts
rm -fr ~/git_with_ko_scripts
rm -fr ~/git_with_glide_scripts

mkdir ~/git
cd ~/git
git init
git config user.email "john@doe.com"
git config user.name "John Doe"
touch tomate
git add -A
git commit -m "re v1"
git tag "notsemvertag"
git tag "v1.0.2"
git tag "v1.0.0"

mkdir ~/git_with_ok_scripts
cd ~/git_with_ok_scripts
git init
git config user.email "john@doe.com"
git config user.name "John Doe"
touch tomate
cat <<EOT > .version
preversion:echo "hello"
postversion:echo "goodbye"
EOT
git add -A
git commit -m "re v1"

mkdir ~/git_with_ko_scripts
cd ~/git_with_ko_scripts
git init
git config user.email "john@doe.com"
git config user.name "John Doe"
touch tomate
cat <<EOT > .version
preversion:eccho "hello"
postversion:eccho "goodbye"
EOT
git add -A
git commit -m "re v1"

mkdir ~/git_with_glide_scripts
cd ~/git_with_glide_scripts
git init
git config user.email "john@doe.com"
git config user.name "John Doe"
touch tomate
cat <<EOT > glide.yaml
scripts:
  preversion: echo "hello"
  postversion: echo "goodbye"
EOT
git add -A
git commit -m "re v1"
