module main

import os
import net.http
import json

const success = '\e[42m\e[44m[INFO]\e[0m'

const error = '\e[42m\e[41m[INFO]\e[0m'

const input = '\e[42m\e[100m[INPUT]\e[0m'

const default_whitelist = ['logs/', 'bin/', 'themes/', 'config.toml', 'sqlite/', '.cupignore']

fn main() {
	mut whitelist := []string{}

	cmd_version := fn () string {
		for x, i in os.args {
			if i == '--version' {
				if os.args.len < x + 2 {
					return ''
				} else {
					return os.args[x + 1]
				}
			}
		}
		return 'no'
	}()

	whitelist << read_whitelist()
	whitelist << default_whitelist

	if cmd_version == '' {
		println('${error} Invalid args')
		return
	}

	download_version(cmd_version, whitelist) or {
		println(err)
		return
	}
}

fn get_version() !string {
	lines := os.read_lines('./config.toml') or {
		return error("${error} Can't read config file (config.toml)")
	}
	for line in lines {
		if !line.starts_with('Version') {
			continue
		}
		return line.trim_string_left('Version = ')
	}

	return error("${error} Can't get version")
}

struct GitUrl {
	tarball_url string [required]
	tag_name    string [required]
}

fn download_version(version string, whitelist []string) ! {
	mut get_download_link := http.Response{}
	temp_dir := './tmp'
	mut latest_version := ''

	os.mkdir(temp_dir) or {}

	mut v := get_version() or {
		println(err)
		return
	}
	v = v.replace('"', '')

	if version == 'no' {
		get_download_link = http.get('https://api.github.com/repos/TryZeroOne/Contagio/releases/latest') or {
			return error('${error} Github api error: ${err}')
		}
	} else {
		get_download_link = http.get('https://api.github.com/repos/TryZeroOne/Contagio/releases/tags/${version}') or {
			return error('${error} Github api error: ${err}')
		}
		if get_download_link.status_code == 404 {
			return error('${error} Version ( ${version} ) not found. See https://github.com/TryZeroOne/Contagio/releases')
		}
	}

	res := json.decode(GitUrl, get_download_link.body) or {
		return error('${error} Json decoding error: ${err}')
	}

	if res.tag_name == v.str() {
		println('${error} The version is already installed')
		return
	}
	latest_version = res.tag_name

	for {
		install := os.input('${input} Version found (${latest_version})! Install? [y/n]: ')
		if install != 'y' && install != 'n' {
			continue
		} else {
			if install != 'y' {
				return
			}

			break
		}
	}

	backup(v)!

	println('${success} Downloading version ${latest_version}...')

	http.download_file(res.tarball_url, temp_dir + '/res.tar.gz') or {
		return error("${error} Can't download new version: ${err}")
	}

	tar := os.execute('tar -zxf ' + temp_dir + '/res.tar.gz -C ' + temp_dir + '/')
	if tar.exit_code != 0 {
		return error('${error} Tar error: ${tar.output}')
	}

	mut dir := os.ls(temp_dir + '/')!
	mut upd_dir := ''

	for d in dir {
		if d.starts_with('TryZeroOne') {
			upd_dir = d
			for i in whitelist {
				os.execute('rm -rf ' + temp_dir + '/' + d + '/' + i) // rm whitelist
			}
		}
	}

	// remove cur dir
	files := os.ls('./')!
	for file in files {
		if file.starts_with('backup') || file.starts_with('tmp') {
			continue
		}
		is_whitelisted := fn (file string, whitelist []string) bool {
			for w in whitelist {
				if w.starts_with(file) {
					return true
				}
			}
			return false
		}(file, whitelist)

		if is_whitelisted {
			continue
		}

		os.execute('rm -rf ' + file)
	}

	// copy to cur dir
	dir = os.ls(temp_dir + '/' + upd_dir + '/')!
	for d in dir {
		os.cp_all(temp_dir + '/' + upd_dir + '/' + d, './' + d, true) or { continue }
	}
	os.execute('rm -rf ./tmp/')
}

fn backup(version string) ! {
	backup_name := 'backup_' + version

	ls := os.ls('./')!

	os.mkdir(backup_name) or {
	}

	for i in ls {
		if i != backup_name {
			os.cp_all(i, backup_name + '/' + i, true) or { continue }
		}
	}

	println('${success} Backup has been created with name "${backup_name}"')
}

fn read_whitelist() []string {
	mut custom_whitelist := []string{}

	ignore := os.read_lines('.cupignore') or {
		println('${success} ".cupignore" not found. Skipping...')
		return []string{}
	}

	for i in ignore {
		if i.starts_with('//') || i.len < 1 || i == '' || i == ' ' {
			continue
		}
		custom_whitelist << i
	}

	println('${success} Read .cupignore ( ${custom_whitelist.len} lines )')

	return custom_whitelist
}
