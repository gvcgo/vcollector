#!/bin/bash

condaUrl="https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-s390x.sh"

if [ "$(arch)" == "x86_64" ]; then
    condaUrl="https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh"
elif [ "$(arch)" == "i386" ] && [ "$(uname)" == "Darwin" ]; then
    condaUrl="https://repo.anaconda.com/miniconda/Miniconda3-latest-MacOSX-x86_64.sh"
fi

echo "download url: $condaUrl"

cd ~
curl -L "$condaUrl" -o "miniconda.sh"

if [ -f "~/miniconda.sh" ];then
    sh ~/miniconda.sh
fi
