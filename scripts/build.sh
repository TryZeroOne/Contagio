#!/bin/bash

cmds=("GOOS=linux GOARCH=mips CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips.bin ." "GOOS=linux GOARCH=mipsle CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips32le.bin ." "GOOS=linux GOARCH=mips64le CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/mips64le.bin ." " GOOS=linux GOARCH=ppc64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/ppc64.bin ." "  GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/riscv.bin ." "GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm7.bin ." "GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm6.bin ." "GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/arm5.bin ." " GOOS=linux GOARCH=amd64 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/x86_64.bin ." "GOOS=linux GOARCH=386 CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/x32.bin ." "GOOS=linux GOARCH=s390x CGO_ENABLED=0 garble -seed=random -literals -tiny build -o ../bin/s390x.bin .")

upxcompress() {

    if ! command -v upx &>/dev/null; then
        echo "UPX Packer has not been installed. Skipping..."
        echo "Build success. Binary files in the bin folder"
        exit
    fi

    cd bin

    upxcmds=("upx --brute --lzma x86_64.bin"
        "upx --brute --lzma x32.bin"
        "upx --brute --lzma mips.bin"
        "upx --brute --lzma arm6.bin"
        "upx --brute --lzma arm5.bin"
        "upx --brute --lzma arm7.bin"
        "upx --brute --lzma mips32le.bin"
        "upx --brute --lzma ppc64.bin")

    total_cmds=${#upxcmds[@]}
    current_cmd=0

    for cmd in "${upxcmds[@]}"; do
        ((current_cmd++))
        clear
        echo ====================================
        echo "        [$current_cmd/$total_cmds] Compressing..."
        echo ====================================
        execute_command "$cmd"
    done

}

execute_command() {
    local command="$1"
    local output

    # echo $1

    output="$(eval "$command" 2>&1)"
    local exit_code=$?

    if [ $exit_code -ne 0 ]; then
        echo "Command failed with exit code $exit_code."
        echo "Error: $output"
        exit
    fi

}

launch() {
    rm -rf bin
    mkdir bin
    cd bot
    clear

    if [ "$2" != "" ]; then
        sleeptime="$2"
    else
        sleeptime=1 # If you set it to 0, the pc will die 💀
    fi

    total_cmds=${#cmds[@]}
    current_cmd=0

    if [ "$1" = "true" ]; then
        for cmd in "${cmds[@]}"; do
            ((current_cmd++))
            execute_command "$cmd" &
            clear
            echo ====================================
            echo "        [$current_cmd/$total_cmds] Building..."
            echo ====================================
            sleep "$sleeptime"
        done
        wait
    else
        for cmd in "${cmds[@]}"; do
            ((current_cmd++))
            clear
            echo ====================================
            echo "        [$current_cmd/$total_cmds] Building..."
            echo ====================================
            execute_command "$cmd"
        done
    fi

    cd ..
}

# go clean -cache

if ! command -v garble &>/dev/null; then
    echo "Garble has not been installed"
    exit
fi

if [[ $1 == "upx" ]]; then
    launch $2 $3

    upxcompress
elif [[ $1 == "default" ]]; then
    launch $2 $3
else
    echo "Invalid args"
fi
