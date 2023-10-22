/*
 * @Author: Guo Yaoze, Qu Zhixuan, Sun Longyu, Pu Zhengyu, Zhang Yuntao, Peng Huanran, Yang Qihao
 * @License: AGPLv3
 * @Date: 2023-10-19 23:13:43
 * @LastEditTime: 2023-10-22 14:40:26
 * @Description: The core of the program.
 * @FilePath: /MikuMikuLive/main.go
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
)

var argsParserInfo = []string{} // Info from args_parser.

// Parase os.Args.
func args_parser() {
	ProgramName = os.Args[0]
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--no-welcome", "-nw":
			NoWelcome = true
		case "--version", "-v":
			argsParserInfo = append(argsParserInfo, "MikuMikuLive - The Libre Live Toolkit", "Version: "+VER+" ("+VER_CODENAME+")")
		case "--http-port", "-hp":
			tmp, err := strconv.Atoi(os.Args[i+1])
			if !Err_handle(err) {
				HttpPort = tmp
			}
			i++
		default:
			argsParserInfo = append(argsParserInfo, color.HiRedString("Unrecognized arg \""+os.Args[i]+"\", is it correct? IGNORED it and keep going..."))
		}
	}
}

// Welcome info.
func welcome() {
	if !NoWelcome {
		fmt.Println("MikuMikuLive - The Libre Live Toolkit")
		fmt.Println("This is a libre software under the AGPLv3")
		fmt.Println("Version: " + VER + " (" + VER_CODENAME + ")")
	}
}

// As its name.
func print_args_parser_info() {
	for i := range argsParserInfo {
		fmt.Println(argsParserInfo[i])
	}
}

// HTTP handler func for "/".
func http_default_handler(w http.ResponseWriter, r *http.Request) {
	htmlLines := Read_lines("root.html")
	for i := range htmlLines {
		fmt.Fprintln(w, htmlLines[i])
	}
}

// Http server.
func http_server() {
	http.HandleFunc("/", http_default_handler)
	http.HandleFunc("/api/", Api_handler)
	http.ListenAndServe(HttpAddr+":"+strconv.Itoa(HttpPort), nil)
}

// LiveDescFile loader.
func loader() {
	fmt.Println("Input \"!EXIT!\" to exit MikuMikuLive instead of load a LDF.")
	for !File_exist(LiveDescFile) {
		fmt.Print("(LOAD) ")
		fmt.Scanln(&LiveDescFile)
		if LiveDescFile == "!EXIT!" {
			os.Remove(Unify_path(TmpDir, false))
			os.Exit(0)
		}
		if !File_exist(LiveDescFile) {
			color.HiRed("Error: LiveDescFile \"" + LiveDescFile + "\" is not exist.")
		}
	}
}

// The really important function main.
func main() {
	args_parser()
	welcome()
	print_args_parser_info()
	go http_server()
	Prepare_tmp(true)
	for {
		for !File_exist(LiveDescFile) {
			fmt.Println("Seems there is no valid LiveDescFile, entering the LOADER...")
			loader()
		}
		fmt.Println("Expanding the LDF...")
		LDF_expand(LiveDescFile)
		if !Check_conf_files_gened() {
			color.HiRed("Expanded conf files didn't pass the test, please check your input.")
			LiveDescFile = ""
		} else {
			Conf_loader()
			LatestScene = "(nil)"
			Controller()
			welcome()
		}
	}
}
