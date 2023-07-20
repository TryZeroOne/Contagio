#!/bin/bash
SUDO_UID=0
INFO="\e[42m\e[44m[INFO]\e[0m"
ERROR="\e[42m\e[40m[INFO]\e[0m"

all() {

    sed -i 's/1024/999999/g' /usr/include/bits/typesizes.h
    ulimit -n 99999
    ulimit -u 99999
    ulimit -e 99999

    echo -e "$INFO Golang installation..."
    wget https://storage.googleapis.com/golang/go1.20.linux-amd64.tar.gz

    sudo tar -C /usr/local -xzvf go1.20.linux-amd64.tar.gz >/dev/null 2>&1
    echo 'export PATH="/usr/local/go/bin:$PATH"' >>~/.bashrc
    echo 'export GOPATH=$HOME/go' >>~/.bashrc

    echo -e "$INFO Vlang installation..."

    rm -rf vlang
    git clone --depth=1 https://github.com/vlang/v vlang
    cd vlang
    make
    # mv ../v/ /opt
    ./v symlink
    cd ..
    # cd ..

    echo -e "$INFO Docker installation..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo groupadd docker >/dev/null 2>&1
    sudo usermod -aG docker $USER >/dev/null 2>&1
    systemctl start docker

    echo -e "$INFO Nim installation..."
    curl https://nim-lang.org/choosenim/init.sh -sSf | sh
    sleep 1

    echo 'export PATH=$HOME/.nimble/bin:$PATH' >>~/.bashrc

    echo -e "$INFO Garble installation..."
    echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >>~/.bashrc

    rm -rf assets >/dev/null 2>&1
    rm -rf temp >/dev/null 2>&1
    mkdir assets >/dev/null 2>&1
    mkdir temp >/dev/null 2>&1
    pip3 install pycryptodome >/dev/null 2>&1

    cd temp >/dev/null 2>&1
    echo -e "$INFO Upx installation..."
    wget https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz >/dev/null 2>&1
    tar -xvf upx-4.0.2-amd64_linux.tar.xz >/dev/null 2>&1

    cd upx-4.0.2-amd64_linux >/dev/null 2>&1
    sudo cp upx /usr/bin >/dev/null 2>&1

    cd ../../ >/dev/null 2>&1

    rm -rf v
    rm -rf temp >/dev/null 2>&1
}

pacmanPkg() {
    echo -e "$INFO Updating..."
    pacman -S sudo --noconfirm >/dev/null 2>&1
    sudo pacman -Fy --noconfirm
    sudo pacman -Syu --noconfirm

    echo -e "$INFO Snap installation..."
    sudo pacman -S snapd --noconfirm >/dev/null 2>&1
    sudo systemctl enable --now snapd.socket
    sudo ln -s /var/lib/snapd/snap /snap >/dev/null 2>&1

    if [ $? -ne 0 ]; then
        echo -e "$ERROR Snap error..."
    fi

    echo -e "$INFO Gcc installation..."
    sudo pacman -S gcc --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Snap error..."
    fi

    sudo pacman -S xz --noconfirm >/dev/null 2>&1

    echo -e "$INFO Curl installation..."
    sudo pacman -S curl --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Curl error..."
    fi

    echo -e "$INFO Make installation..."
    sudo pacman -S make --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Make error..."
    fi

    echo -e "$INFO Git installation..."
    sudo pacman -S git --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Git error..."
    fi

    echo -e "$INFO Tor installation..."
    sudo pacman -S tor --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tor error..."
    fi

    echo -e "$INFO OpenSSL installation..."
    sudo pacman -S openssl --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR OpenSSL error..."
    fi

    echo -e "$INFO Screen installation..."
    sudo pacman -S screen --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Screen error..."
    fi

    echo -e "$INFO Sqlite3 installation..."
    sudo pacman -S sqlite3 --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Sqlite3 error..."
    fi

    echo -e "$INFO Python3 installation..."
    sudo pacman -S python3 --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Python3 error..."
    fi

    echo -e "$INFO Wget installation..."
    sudo pacman -S wget --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Wget error..."
    fi

    echo -e "$INFO Tar installation..."
    sudo pacman -S tar --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tar error..."
    fi

    echo -e "$INFO Nano installation..."
    sudo pacman -S nano --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Nano error..."
    fi

}

dnfPkg() {
    echo -e "$INFO Updating..."
    sudo dnf update -y
    sudo dnf install -y epel-release >/dev/null 2>&1
    sudo dnf install epel-release >/dev/null 2>&1
    sudo yum install epel-release >/dev/null 2>&1
    sudo dnf upgrade

    echo -e "$INFO Gcc installation..."
    sudo dnf install -y gcc >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Gcc error..."
    fi

    echo -e "$INFO Curl installation..."
    sudo dnf install -y curl >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Curl error..."
    fi

    echo -e "$INFO Git installation..."
    sudo dnf install -y git >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Git error..."
    fi

    echo -e "$INFO Tor installation..."
    sudo dnf install -y tor >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tor error..."
    fi

    echo -e "$INFO Screen installation..."
    sudo dnf install -y screen >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Screen error..."
    fi

    echo -e "$INFO Sqlite3 installation..."
    sudo dnf install -y sqlite-devel >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Sqlite3 error..."
    fi

    echo -e "$INFO Python3 installation..."
    sudo dnf install -y python3 >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Python3 error..."
    fi

    echo -e "$INFO Wget installation..."
    sudo dnf install -y wget >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Wget error..."
    fi

    echo -e "$INFO OpenSSL installation..."
    sudo dnf install -y openssl >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR OpenSSL error..."
    fi

    echo -e "$INFO Tar installation..."
    sudo dnf install -y tar >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tar error..."
    fi

    echo -e "$INFO Make installation..."
    sudo dnf install -y make >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Make error..."
    fi
    echo -e "$INFO Nano installation..."
    sudo dnf install -y nano >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Nano error..."
    fi

    sudo dnf install -y xz >/dev/null 2>&1
}

aptPkg() {
    echo -e "$INFO Updating..."
    sudo apt-get upgrade -y
    sudo apt-get update -y

    echo -e "$INFO Gcc installation..."
    sudo apt install gcc -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Gcc error..."
    fi

    echo -e "$INFO Git installation..."
    sudo apt install git -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Git error..."
    fi

    echo -e "$INFO Tor installation..."
    sudo apt install tor -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tor error..."
    fi

    echo -e "$INFO Make installation..."
    sudo apt install make -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Make error..."
    fi
    echo -e "$INFO Curl installation..."
    sudo apt install curl -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Curl error..."
    fi

    echo -e "$INFO OpenSSL installation..."
    sudo apt install openssl -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR OpenSSL error..."
    fi

    echo -e "$INFO Screen installation..."
    sudo apt install screen -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Screen error..."
    fi
    echo -e "$INFO Sqlite3 installation..."
    sudo apt install sqlite3 -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Sqlite3 error..."
    fi

    systemctl restart systemd-logind.service -y >/dev/null 2>&1

    echo -e "$INFO Python3 installation..."
    sudo apt install python3 -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Python3 error..."
    fi
    echo -e "$INFO Wget installation..."
    sudo apt install wget -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Wget error..."
    fi

    echo -e "$INFO Tar installation..."
    sudo apt install tar -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Tar error..."
    fi

    echo -e "$INFO Nano installation..."
    sudo apt install nano -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "$ERROR Nano error..."
    fi

}

launch() {
    if [ "$UID" -eq "$SUDO_UID" ]; then
        echo -e "\e[42m\e[37m[INFO]\e[0m Sudo: Yes"
    else
        sudo -v
        echo -e "\e[42m\e[37m[INFO]\e[0m Sudo: Yes"
    fi

    case "$(uname -s)" in
    Linux)
        if command -v apt-get >/dev/null 2>&1; then
            echo -e "\e[42m\e[37m[INFO]\e[0m Package manager found: apt"
            sleep 0.5
            aptPkg
            all

        elif command -v dnf >/dev/null 2>&1; then
            echo -e "\e[42m\e[37m[INFO]\e[0m Package manager found: dnf"
            sleep 0.5
            dnfPkg
            all

        elif
            command -v pacman >/dev/null 2>&1
        then
            echo -e "\e[42m\e[37m[INFO]\e[0m Package manager found: pacman"
            sleep 0.5
            pacmanPkg
            all

        else
            echo "$ERROR Unknown package manager"
            exit
        fi
        ;;
    *)
        echo "Unsupported operating system"
        ;;
    esac
}

clean() {
    echo -e "$INFO Contagio setup..."
    git clone https://github.com/TryZeroOne/Contagio >/dev/null 2>&1
    cd Contagio
    v installer/update.v -o cup
    sudo cp cup /bin/

    rm -rf themes config.toml assets sqlite tests README.md .gitignore setup.txt go.mod go.sum installer cup

    mkdir themes
    curl https://raw.githubusercontent.com/TryZeroOne/Contagio/main/installer/empty_config.toml -o config.toml
    curl https://raw.githubusercontent.com/TryZeroOne/Contagio/main/themes/empty.toml -o themes/theme.toml
}

default() {
    echo -e "$INFO Contagio setup..."
    git clone https://github.com/TryZeroOne/Contagio >/dev/null 2>&1
    cd Contagio
    v installer/update.v -o cup
    sudo cp cup /bin/

    rm -rf installer go.mod go.sum cup
}

if [[ $1 == "-clean" ]]; then
    launch
    clean

elif [[ $1 == "-default" ]]; then
    launch
    default
else
    echo """
    sudo bash installer.sh [-default/-clean]
    see https://github.com/TryZeroOne/Contagio#installation
    """
fi
