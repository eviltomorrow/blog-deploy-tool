#!/bin/bash

cd blog
rm -rf public
git pull origin master
hexo generate
cd /usr/share/nginx/html
rm -rf *
cp -r /root/roigo.me/blog/public/* .