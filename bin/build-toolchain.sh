#!/bin/bash

echo "Building clang and binutils based cross compilation toolchain for raspberry pis"

set -e

download() {
  mkdir -p downloads
  for url in $*; do
    filename=$(basename $url)
    if [[ ! -f downloads/$filename ]]; then
      echo "Downloading $url"
      (cd downloads && curl -# -O $url)
    fi
    dirname=$(basename $filename .tar.gz)
    dirname=$(basename $dirname .tar.xz)
    if [[ ! -d $dirname  ]]; then
      echo "Unpacking $filename"
      tar zxf downloads/$filename
    fi
  done
}

CLANG_SRC=cfe-3.6.2.src
LLVM_SRC=llvm-3.6.2.src
BINUTILS_SRC=binutils-2.25

download \
http://ftp.gnu.org/gnu/binutils/$BINUTILS_SRC.tar.gz \
http://llvm.org/releases/3.6.2/$CLANG_SRC.tar.xz \
http://llvm.org/releases/3.6.2/$LLVM_SRC.tar.xz

HERE=$(pwd)
INST_DIR=$HERE/install
LLVM_INST=${INST_DIR}/llvm
BINUTILS_INST=${INST_DIR}/binutils

[[ ! -d $LLVM_INST ]] && mkdir -p $LLVM_INST
[[ ! -d $BINUTILS_INST ]] && mkdir -p $BINUTILS_INST

build_binutils() {
  time (
    cd $BINUTILS_SRC
    ./configure --target=arm-linux-gnueabi --program-prefix= --prefix=$BINUTILS_INST --with-sysroot=yes
    make -j8
    make install
  )
}

build_clang() {
  time (
    bdir=build.$LLVM_SRC
    mkdir -p $bdir
    cd $bdir
    cmake -G"Unix Makefiles" \
-DCMAKE_INSTALL_PREFIX=$LLVM_INST \
-DLLVM_TARGETS_TO_BUILD="ARM" \
-DLLVM_EXTERNAL_CLANG_SOURCE_DIR=../$CLANG_SRC \
../$LLVM_SRC
    make -j8
    make install
  )
}

#build_binutils
#build_clang

echo "LLVM installed in $LLVM_INST"
echo "binutils installed in $BINUTILS_INST"

write_config() {
  echo ""
  echo "PIXC_INSTALL=$HERE"
  echo "PIXC_LLVM_BIN=$LLVM_INST/bin"
  echo "PIXC_BINUTILS_BIN=$BINUTILS_INST/bin"
}

write_config > .toolchain.config


