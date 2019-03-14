#!/bin/bash

rm -rf ./src
mkdir ./src

apt-get install rsync

rsync -av --exclude='web' --exclude='.git' ../ ./src
