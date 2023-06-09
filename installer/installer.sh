#!/bin/bash
PkgManager=''
SUDO_UID=0

all() {

    echo -e "\e[42m\e[44m[INFO]\e[0m Golang installation..."
    wget https://storage.googleapis.com/golang/go1.20.linux-amd64.tar.gz

    sudo tar -C /usr/local -xzvf go1.20.linux-amd64.tar.gz >/dev/null 2>&1
    echo export PATH="/usr/local/go/bin:$PATH" >>~/.bashrc

    echo -e "\e[42m\e[44m[INFO]\e[0m Docker installation..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo groupadd docker >/dev/null 2>&1
    sudo usermod -aG docker $USER >/dev/null 2>&1
    systemctl start docker

    echo -e "\e[42m\e[44m[INFO]\e[0m Nim installation..."
    curl https://nim-lang.org/choosenim/init.sh -sSf | sh
    sleep 1

    echo export "PATH=$HOME/.nimble/bin:$PATH" >>~/.bashrc

    echo -e "\e[42m\e[44m[INFO]\e[0m Garble installation..."
    echo export "PATH="$PATH:$(go env GOPATH)/bin"" >>~/.bashrc
    go install mvdan.cc/garble@latest

    rm -rf assets >/dev/null 2>&1
    rm -rf temp >/dev/null 2>&1
    mkdir assets >/dev/null 2>&1
    mkdir temp >/dev/null 2>&1
    pip3 install pycryptodome >/dev/null 2>&1

    cd temp >/dev/null 2>&1
    echo -e "\e[42m\e[44m[INFO]\e[0m Upx installation..."
    wget https://github.com/upx/upx/releases/download/v4.0.2/upx-4.0.2-amd64_linux.tar.xz >/dev/null 2>&1
    tar -xvf upx-4.0.2-amd64_linux.tar.xz >/dev/null 2>&1

    cd upx-4.0.2-amd64_linux >/dev/null 2>&1
    sudo cp upx /usr/bin >/dev/null 2>&1

    cd ../../ >/dev/null 2>&1

    rm -rf temp >/dev/null 2>&1
}

pacmanPkg() {
    echo -e "\e[42m\e[44m[INFO]\e[0m Updating..."
    pacman -S sudo --noconfirm >/dev/null 2>&1
    sudo pacman -Fy --noconfirm
    sudo pacman -Syu --noconfirm

    echo -e "\e[42m\e[44m[INFO]\e[0m Snap installation..."
    sudo pacman -S snapd --noconfirm >/dev/null 2>&1
    sudo systemctl enable --now snapd.socket
    sudo ln -s /var/lib/snapd/snap /snap >/dev/null 2>&1

    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Snap error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Gcc installation..."
    sudo pacman -S gcc --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Snap error..."
    fi

    sudo pacman -S xz --noconfirm >/dev/null 2>&1

    echo -e "\e[42m\e[44m[INFO]\e[0m Curl installation..."
    sudo pacman -S curl --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Curl error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Git installation..."
    sudo pacman -S git --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Git error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Screen installation..."
    sudo pacman -S screen --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Screen error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Sqlite3 installation..."
    sudo pacman -S sqlite3 --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Sqlite3 error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Python3 installation..."
    sudo pacman -S python3 --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Python3 error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Wget installation..."
    sudo pacman -S wget --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Wget error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Tar installation..."
    sudo pacman -S tar --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Tar error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Nano installation..."
    sudo pacman -S nano --noconfirm >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Nano error..."
    fi

}

dnfPkg() {
    echo -e "\e[42m\e[44m[INFO]\e[0m Updating..."
    sudo dnf update -y
    sudo dnf install -y epel-release >/dev/null 2>&1
    sudo dnf install epel-release >/dev/null 2>&1
    sudo yum install epel-release >/dev/null 2>&1
    sudo dnf upgrade

    echo -e "\e[42m\e[44m[INFO]\e[0m Gcc installation..."
    sudo dnf install -y gcc >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Gcc error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Curl installation..."
    sudo dnf install -y curl >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Curl error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Git installation..."
    sudo dnf install -y git >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Git error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Screen installation..."
    sudo dnf install -y screen >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Screen error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Sqlite3 installation..."
    sudo dnf install -y sqlite-devel >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Sqlite3 error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Python3 installation..."
    sudo dnf install -y python3 >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Python3 error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Wget installation..."
    sudo dnf install -y wget >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Wget error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Tar installation..."
    sudo dnf install -y tar >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Tar error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Nano installation..."
    sudo dnf install -y nano >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Nano error..."
    fi

    sudo dnf install -y xz >/dev/null 2>&1
}

aptPkg() {
    echo -e "\e[42m\e[44m[INFO]\e[0m Updating..."
    sudo apt-get upgrade -y
    sudo apt-get update -y

    echo -e "\e[42m\e[44m[INFO]\e[0m Gcc installation..."
    sudo apt install gcc -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Gcc error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Git installation..."
    sudo apt install git -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Git error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Curl installation..."
    sudo apt install curl -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Curl error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Screen installation..."
    sudo apt install screen -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Screen error..."
    fi
    echo -e "\e[42m\e[44m[INFO]\e[0m Sqlite3 installation..."
    sudo apt install sqlite3 -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Sqlite3 error..."
    fi

    systemctl restart systemd-logind.service -y >/dev/null 2>&1

    echo -e "\e[42m\e[44m[INFO]\e[0m Python3 installation..."
    sudo apt install python3 -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Python3 error..."
    fi
    echo -e "\e[42m\e[44m[INFO]\e[0m Wget installation..."
    sudo apt install wget -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Wget error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Tar installation..."
    sudo apt install tar -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Tar error..."
    fi

    echo -e "\e[42m\e[44m[INFO]\e[0m Nano installation..."
    sudo apt install nano -y >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo -e "\e[42m\e[40m[INFO]\e[0m Nano error..."
    fi

}

rebootsystem() {
    echo "Reboot is required. Reboot? (y/n)"
    read answer
    if [ "$answer" == "y" ]; then
        sleep 0.5
        sudo reboot
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
            echo "Unknown package manager"
        fi
        ;;
    *)
        echo "Unsupported operating system"
        ;;
    esac
}

clean() {
    echo -e "\e[42m\e[44m[INFO]\e[0m Contagio setup..."
    git clone https://github.com/TryZeroOne/Contagio
    cd Contagio >/dev/null 2>&1
    rm -rf themes >/dev/null 2>&1
    rm config.toml >/dev/null 2>&1
    touch config.toml >/dev/null 2>&1
    rm -rf assets >/dev/null 2>&1
    rm -rf sqlite >/dev/null 2>&1
    rm -rf tests >/dev/null 2>&1
    rm README.md >/dev/null 2>&1
    rm .gitignore >/dev/null 2>&1
    rm setup.txt >/dev/null 2>&1
    rm -rf installer
    rm go.mod >/dev/null 2>&1
    rm go.sum >/dev/null 2>&1
    go mod init contagio >/dev/null 2>&1
    go mod tidy >/dev/null 2>&1

}

if [[ $1 == "-clean" ]]; then
    launch
    clean
    rebootsystem

elif
    [[ $1 == "-default" ]]
then
    launch
    echo -e "\e[42m\e[44m[INFO]\e[0m Contagio setup..."

    git clone https://github.com/TryZeroOne/Contagio
    rm -rf installer
    rm go.mod >/dev/null 2>&1
    rm go.sum >/dev/null 2>&1
    go mod init contagio >/dev/null 2>&1
    go mod tidy >/dev/null 2>&1
    rebootsystem

else
    echo "Invalid args"
fi
