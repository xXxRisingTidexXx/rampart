FROM ubuntu:16.04

ARG CONDA_DIR=/opt/conda

ENV PATH $CONDA_DIR/bin:$PATH

COPY requirements.txt .
COPY twinkle twinkle

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates cmake build-essential gcc g++ git wget && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    /bin/bash Miniconda3-latest-Linux-x86_64.sh -f -b -p $CONDA_DIR && \
    export PATH="$CONDA_DIR/bin:$PATH" && \
    conda config --set always_yes yes --set changeps1 no && \
    conda install -q -y --file requirements.txt && \
    git clone --recursive --branch stable --depth 1 https://github.com/Microsoft/LightGBM && \
    cd LightGBM/python-package && python setup.py install && \
    apt-get autoremove -y && \
    apt-get clean && \
    conda clean -a -y && \
    rm -rf /usr/local/src/*
