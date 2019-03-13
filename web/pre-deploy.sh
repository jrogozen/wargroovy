#!/bin/bash

rm -rf ./src
mkdir ./src

rsync -av --exclude='web' --exclude='.git' ../ ./src
