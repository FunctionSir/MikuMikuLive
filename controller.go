/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-10-20 14:03:25
 * @LastEditTime: 2023-10-22 14:36:21
 * @LastEditors: FunctionSir
 * @Description:
 * @FilePath: /MikuMikuLive/controller.go
 */
package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func print_head() {
	fmt.Println("MikuMikuLive Controller - " + VER + " (" + VER_CODENAME + ")")
	fmt.Println("Current LDF: " + LiveDescFile)
	fmt.Println("Latest Scene: " + LatestScene)
}

func test_ffx(ffx string) bool {
	ff := exec.Command(ffx)
	ffOut, _ := ff.CombinedOutput()
	if len(ffOut) < FFMPEG_CHK_TS {
		return false
	} else {
		return true
	}
}

func list_scene() {
	Clear()
	print_head()
	var tmp string
	sectionList := SceneConf.SectionStrings()
	for i := 1; i < len(sectionList); i++ {
		tmp += sectionList[i] + "; "
	}
	fmt.Println("Scenes: " + tmp)
	fmt.Println("Total: " + strconv.Itoa(len(sectionList)-1))
	Pause()
}

func jump_to() {
	Clear()
	print_head()
	var tmp string
	sectionList := SceneConf.SectionStrings()
	for i := 1; i < len(sectionList); i++ {
		tmp += sectionList[i] + "; "
	}
	fmt.Println("Scenes: " + tmp)
	fmt.Println("Input \":q\" and press Enter to exit with nothong done.")
	for !SceneConf.HasSection(tmp) && tmp != ":q" {
		fmt.Print("Jmp to: ")
		tmp = Scan_str()
		if SceneConf.HasSection(tmp) {
			QuitAutoPlay = true
			time.Sleep(10 * time.Millisecond)
			Show_scene(tmp)
		} else {
			color.HiRed("Error: No such scene.")
		}
	}
}

func manually() {
	Clear()
	print_head()
	var f, d string
	fmt.Print("File = ")
	f = Scan_str()
	fmt.Print("Disp = ")
	d = Scan_str()
	FFPlay(f, d, nil)
}

func debug_settings() {
	Clear()
	print_head()
	fmt.Println(">按下[方括号]内的数字对应的键然后按下Enter来执行指令<")
	fmt.Println("[0] 不要让FFPlay置顶")
	fmt.Println("[1] 让FFPlay置顶")
	fmt.Println("[2] 让FFPlay有边框")
	fmt.Println("[3] 不要让FFPlay有边框")
	var cmd string
	fmt.Print(">>> ")
	fmt.Scanln(&cmd)
	switch cmd {
	case "0":
		NoAlwaysOnTop = true
	case "1":
		NoAlwaysOnTop = false
	case "2":
		PlayWithBorder = true
	case "3":
		PlayWithBorder = false
	default:
		color.HiRed("Error: unrecognized input, do nothing.")
		Pause()
	}
}

func Controller() {
	var cmd string
	for LiveDescFile != "" {
		Clear()
		print_head()
		fmt.Println(">按下[方括号]内的字母对应的键然后按下Enter来执行指令<")
		fmt.Println(">只是按下Enter来刷新<")
		fmt.Println("[T]est FFPlay&FFProbe 测试FFPlay&FFProbe")
		fmt.Println("[I]nit Display 初始化显示(会打断AutoPlay)")
		fmt.Println("[F]irst Scene 第一个Scene(会打断AutoPlay)")
		fmt.Println("[N]ext Scene 下一Scene(会打断AutoPlay)")
		fmt.Println("[P]rev Scene 上一Scene(会打断AutoPlay)")
		fmt.Println("[J]ump To Scene 跳至Scene(会打断AutoPlay)")
		fmt.Println("[A]uto 从LatestScene开始自动播放")
		fmt.Println("[Q]uite Auto Play 退出自动播放(现有Scene不会关闭)")
		fmt.Println("[L]ist Scene 列出全部Scene")
		fmt.Println("[M]anually 手动播放文件(将会叠在现有的场景上层且会受DurController的影响)")
		fmt.Println("[C]lose Scene 关闭Scene(以及手动播放的文件)(相当于杀死全部未被保护的FFPlay)(也会同时RstDurController)")
		fmt.Println("[R]st Dur Controller 重置DurController(相当于让当前Scene的Dur变为inf)")
		fmt.Println("[K]ill All 杀死全部MML启动的FFPlay(包括被保护的)(也会重置DurController)(也会退出AutoPlay)(也会重置DispInited状态)")
		fmt.Println("[D]ebug Settings 调试设置")
		fmt.Println("[U]nload Current LDF 卸载当前LDF")
		fmt.Print(">>> ")
		fmt.Scanln(&cmd)
		switch cmd {
		case "T", "t":
			Clear()
			print_head()
			pass := false
			if test_ffx(FFConf.Section("FFPLAY").Key("Exec").String()) {
				if test_ffx(FFConf.Section("FFPROBE").Key("Exec").String()) {
					pass = true
				}
			}
			if pass {
				color.HiGreen("The simple check returned \"true\", test PASS.")
				fmt.Println("In next test, you should see a large picture on you screen after you press Enter. You can press ESC to exit.")
				Pause()
				FFPlay("media.d/ready.png", "", []string{"-fs"})
				fmt.Println("If you can saw that picture, the test is completely passed. If not, you should check your FFPlay.")
			} else {
				color.HiRed("The simple check returned \"false\", test FAILD.")
			}
			Pause()
		case "I", "i":
			Disp_init()
		case "F", "f":
			QuitAutoPlay = true
			time.Sleep(10 * time.Millisecond)
			Show_scene(SceneConf.SectionStrings()[1])
		case "P", "p":
			currentSceneID := Find_str(SceneConf.SectionStrings(), LatestScene)
			if currentSceneID >= 2 {
				QuitAutoPlay = true
				time.Sleep(10 * time.Millisecond)
				Show_scene(SceneConf.SectionStrings()[currentSceneID-1])
			}
		case "N", "n":
			if LatestScene == "(nil)" {
				LatestScene = SceneConf.SectionStrings()[1]
			}
			currentSceneID := Find_str(SceneConf.SectionStrings(), LatestScene)
			if currentSceneID <= len(SceneConf.SectionStrings())-1 {
				QuitAutoPlay = true
				time.Sleep(10 * time.Millisecond)
				Show_scene(SceneConf.SectionStrings()[currentSceneID+1])
			}
		case "J", "j":
			jump_to()
		case "M", "m":
			manually()
		case "A", "a":
			QuitAutoPlay = false
			go Auto_play()
		case "L", "l":
			list_scene()
		case "C", "c":
			Close_ffplay()
			RstDurController = true
			time.Sleep(10 * time.Millisecond)
		case "R", "r":
			RstDurController = true
			time.Sleep(10 * time.Millisecond)
		case "K", "k":
			RstDurController = true
			QuitAutoPlay = true
			time.Sleep(10 * time.Millisecond)
			Kill_all_ffplay()
			DispInited = false
		case "D", "d":
			debug_settings()
		case "U", "u":
			QuitAutoPlay = true
			RstDurController = true
			time.Sleep(100 * time.Millisecond)
			Kill_all_ffplay()
			FFConf = nil
			DispConf = nil
			MediaConf = nil
			SceneConf = nil
			Prepare_tmp(false)
			LiveDescFile = ""
			LatestScene = "(nil)"
			DispInited = false
			Clear()
		}
		cmd = ""
	}
}
