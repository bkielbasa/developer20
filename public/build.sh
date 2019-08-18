git pull

bundle install
bundle exec jekyll build --drafts
scp -r ./_site bklimczak@kodcast.pl:~/domains/dev.developer20.com/public_html

