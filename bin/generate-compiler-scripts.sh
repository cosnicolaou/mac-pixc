#!/bin/bash

set -e

# write script to run a cc+cxx cross compiler
write_xc_script() {
  echo "#!/bin/bash"
  cat .toolchain.config
  echo "TARGET=arm-linux-gnueabihf"
  echo "SYSROOT=$1"
  echo 'ISYSROOT=$SYSROOT/usr/lib/gcc/$TARGET'
  echo 'export PATH=$PIXC_LLVM_BIN:$PATH:'
  echo 'clang --target=$TARGET --sysroot=$SYSROOT  -isysroot $ISYSROOT -B$PIXC_BINUTILS_BIN "$@"'
}

# write a script to build a go cross compiler
write_go_xc_script() {
  echo "export CGO_ENABLED=1"
  echo "export CC_FOR_TARGET=$1"
  echo "export CXX_FOR_TARGET=$2"
  echo "export GOOS=linux"
  echo "export GOARCH=arm"
  echo "export GOARM=$3"
  echo "bash ./all.bash"
}

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <location of sysroot> <name for scripts>"
  exit 1
fi

sysroot=$1
name=$2
if [[ -z "$sysroot" ]]; then
  echo "missing sysroot argument"
  exit 1
fi

if [[ -z "$name" ]]; then
  echo "missing name argument"
  exit 1
fi

here=$(pwd)
cc_script=$here/bin/cc-arm-$name
cxx_script=$here/bin/cxx-arm-$name
rm -f $cc_script $cxx_script
write_xc_script $sysroot > $cc_script
ln -s $(basename $cc_script) $cxx_script
chmod +x $cc_script $cxx_script

build_go_arm6_script=bin/build-go-arm6.sh
build_go_arm7_script=bin/build-go-arm7.sh
rm -f $build_go_arm6_script $build_go_arm7_script
write_go_xc_script $cc_script $cxx_script 6 > $build_go_arm6_script
write_go_xc_script $cc_script $cxx_script 7 > $build_go_arm7_script
chmod +x $build_go_arm6_script $build_go_arm7_script