# Darwin hosted go cross compiler for raspberry pi's with cgo enabled.

This note outlines how to set up a Darwin hosted go cross compiler for arm-based systems with support for c-go. If you don’t need to use cgo then the normal go 1.5 cross compilation support is very easy to use and you needn't bother with any of this.

## The C/CXX Cross Compiler

After some failed attempts to use crosstools-ng and embtoolkit I settled on llvm and clang as the cross compiler. This proved relatively easy once I’d instrumented clang to see what it was actually doing with its myriad command line options. Building an instance of clang and llvm with the ability to generate arm code is straightforward, however the standard Darwin linker cannot be used, rather an appropriately configured version of gnu ld (from binutils) is required. The combination of clang and gnu ld can generate and link binaries that will run on arm based systems. The only thing missing is access to system headers and libraries (things like crt0.o, libgcc.a) which provide access to the runtime on the target system. crosstools attempts to cross compile these from source which doesn’t appear to be possible on Darwin since crosstools depends on glibc which is unsupported (and seems deliberately configured to resist being used) on Darwin. I took the easy way out of copying these from a carefully configured target system. These ‘images’ need to be carefully maintained to remain in sync with the intended target systems. In practice this is no different to using crosstools-ng - i.e. regardless of how these runtime libraries/includes are created they must be in sync with the target. Ideally a hermetic build and compatibility test are the ideal solution, but this is beyond me for now.

### Binutils configuration

Given a directory with binutils (I used 2.25), the following configuration can be used:

```
INSTDIR=<example> ./configure --target=arm-linux-gnueabi --program-prefix= --prefix=$INSTDIR/binutils --with-sysroot=yes
```

The parameters are as follows:
* --target is appropriate for a linux arm system with elf formats
* --program-prefix is set to the empty string so that ld is installed as “ld” and not as “arm-ld”
* --prefix is the location to install into
* --with-sysroot=yes enables the --sysroot parameter for ld (which is used by clang).

### Clang/llvm configuration

```
INSTDIR=<example>
LLVM_SRC=$LLVM/llvm-3.6.2.src
CLANG_SRC=$LLVM/cfe-3.6.2.src

cmake -G"Unix Makefiles" \
-DCMAKE_INSTALL_PREFIX=$INSTDIR/llvm \
-DLLVM_TARGETS_TO_BUILD="ARM" \
-DLLVM_EXTERNAL_CLANG_SOURCE_DIR=$CLANG_SRC \
$LLVM_SRC
```

This is a simple configuration for clang/llvm which enables code generation for ARM. Note that the default ‘target’ for clang will still be host that it is built on - i.e. Darwin.

### Scripts and installation locations

The script, build-toolchain.sh, will download and build this clang based cross compilation tool chain. Run it in an initially empty directory. It will create a directory *install* with subdirectories *llvm* and *binutils*. You'll need to refer to these installation directories when running the toolchain as explained below. In addition the build-toolchain.sh writes out a file called *.toolchain.config* that contains shell variable assignments that can be sourced by the other scripts provided.

### Using the cross compiler

Running the cross compiler consists of running clang with the appropriate arguments:

```
INSTDIR=<example>
SYSROOT=<example>
TARGET=arm-linux-gnueabihf
clang \
--target=$TARGET \
-B$INSTDIR/binutils/bin \
--sysroot=$SYSROOT \
-isysroot $SYSROOT/usr/lib/gcc/$TARGET \
```

The command line flags are as follows:

* --target, the target system type represented as a llvm ‘triple’ - architecture, operating system, abi. The values used here are appropriate for a raspberry pi 1 or 2. In theory more specific options may be used (e.g. armv6m) though I have not tested the; similarly -mcpu, mfpu etc are theoretically usable.
* -B points to the directory containing the gnu ld binary built above
* --sysroot and -isysroot (the -- and - matter!) point to the locations of the runtime headers and libraries for the target system. --sysroot points to the system as a whole, whereas -isysroot to the gcc runtime headers+libraries within it.

### Obtaining a 'sysroot'

Obtaining usable copies of the sysroot and isysroot directories from a target system is annoying but at least possible! For the purposes of testing, I simply copied the following directories from a raspberry pi to a local directory:

* /lib
* /usr/lib
* /usr/include

These provided suifficient for my initial testing. Obviously it would good to automate this process.

### Scripts to run the compiler



## The Go Cross Compiler

Cross Compiling With Go and Cgo

Cross compiling go with cgo dependencies is more complicated than for pure go. In particular, the go toolchain must be built with a cgo cross compiler preconfigured for it. This is achieved by setting: CGO_ENABLED=1, GO_EXTLINK_ENABLED=1, CC_FOR_TARGET and CXX_FOR_TARGET when invoking make.bash for the target architecture. For now,

