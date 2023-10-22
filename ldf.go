/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-10-20 13:01:43
 * @LastEditTime: 2023-10-22 14:28:05
 * @LastEditors: FunctionSir
 * @Description:
 * @FilePath: /MikuMikuLive/ldf.go
 */
package main

import (
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-ini/ini"
)

func LDF_expand(name string) {
	ldfLines := Read_lines(name)
	var tmp = []string{}
	var fileToWrite = "/dev/null"
	for i := range ldfLines {
		switch ldfLines[i] {
		case "<BEGIN-FFMPEG-CONF>":
			tmp = []string{}
			fileToWrite = Unify_path(TmpDir, true) + "ffmpeg.conf"
		case "<BEGIN-DISPLAY-CONF>":
			tmp = []string{}
			fileToWrite = Unify_path(TmpDir, true) + "display.conf"
		case "<BEGIN-MEDIA-CONF>":
			tmp = []string{}
			fileToWrite = Unify_path(TmpDir, true) + "media.conf"
		case "<BEGIN-SCENE-CONF>":
			tmp = []string{}
			fileToWrite = Unify_path(TmpDir, true) + "scene.conf"
		case "<END-FFMPEG-CONF>", "<END-DISPLAY-CONF>", "<END-MEDIA-CONF>", "<END-SCENE-CONF>":
			Create_file(fileToWrite, tmp)
		default:
			tmp = append(tmp, ldfLines[i])
		}
	}
}

func Check_conf_files_gened() bool {
	if File_exist(Unify_path(TmpDir, true)+"ffmpeg.conf") && (len(Read_lines(Unify_path(TmpDir, true)+"ffmpeg.conf")) > 0) {
		color.HiGreen("The LDF file contains a FFConf, will use it instead of the internal one.")
		UseExtFFConf = true
		Pause()
	}
	if !File_exist(Unify_path(TmpDir, true)+"display.conf") || !(len(Read_lines(Unify_path(TmpDir, true)+"display.conf")) > 0) {
		color.HiYellow("The LDF file doesn't contain display conf, will use the internal one.")
		color.HiYellow("Warning: will use internal display conf!")
		UseIntDispConf = true
		Pause()
	}
	if !File_exist(Unify_path(TmpDir, true)+"scene.conf") || !(len(Read_lines(Unify_path(TmpDir, true)+"scene.conf")) > 0) {
		color.HiRed("Error: no scene conf found in the LDF file, check FAILED!")
		return false
	}
	if !File_exist(Unify_path(TmpDir, true)+"media.conf") || !(len(Read_lines(Unify_path(TmpDir, true)+"media.conf")) > 0) {
		color.HiRed("Error: no media conf found in the LDF file, check FAILED!")
		return false
	}
	return true
}

func Conf_loader() {
	var err error
	if UseExtFFConf {
		FFConf, err = ini.Load(TmpDir + "ffmpeg.conf")
		FFConf.Reload()
		Err_handle(err)
	} else {
		FFConf, err = ini.Load("conf.d/ffmpeg.conf")
		FFConf.Reload()
		Err_handle(err)
	}
	if UseIntDispConf {
		DispConf, err = ini.Load("conf.d/display.conf")
		DispConf.Reload()
		Err_handle(err)
	} else {
		DispConf, err = ini.Load(TmpDir + "display.conf")
		DispConf.Reload()
		Err_handle(err)
	}
	MediaConf, err = ini.Load(TmpDir + "media.conf")
	MediaConf.Reload()
	Err_handle(err)
	SceneConf, err = ini.Load(TmpDir + "scene.conf")
	SceneConf.Reload()
	Err_handle(err)
	return
}

func Disp_init() {
	QuitAutoPlay = true
	RstDurController = true
	time.Sleep(10 * time.Millisecond)
	Kill_all_ffplay()
	disps := DispConf.SectionStrings()
	for i := 1; i < len(disps); i++ {
		defaultMedia := DispConf.Section(disps[i]).Key("Default").String()
		mediaFile := MediaConf.Section(defaultMedia).Key("File").String()
		FFPlay(mediaFile, disps[i], []string{"$PROTECTED$"})
	}
	DispInited = true
	return
}

func Show_scene(scene string) {
	LatestScene = scene
	RstDurController = true
	time.Sleep(10 * time.Millisecond)
	Close_ffplay()
	toPlay := SceneConf.Section(scene).Key("Play").String()
	playList := strings.Split(toPlay, "+")
	for i := range playList {
		extArgs := []string{}
		playEntry := strings.Split(playList[i], "@")
		if MediaConf.Section(playEntry[0]).HasKey("Loop") {
			extArgs = append(extArgs, "-loop", MediaConf.Section(playEntry[0]).Key("Loop").String())
		}
		if MediaConf.Section(playEntry[0]).HasKey("NoVideo") {
			if MediaConf.Section(playEntry[0]).Key("NoVideo").String() == "true" {
				extArgs = append(extArgs, "-vn", "-nodisp")
			}
		}
		if MediaConf.Section(playEntry[0]).HasKey("NoAudio") {
			if MediaConf.Section(playEntry[0]).Key("NoAudio").String() == "true" {
				extArgs = append(extArgs, "-an")
			}
		}
		FFPlay(MediaConf.Section(playEntry[0]).Key("File").String(), playEntry[1], extArgs)
	}
	if SceneConf.Section(scene).HasKey("Dur") && SceneConf.Section(scene).Key("Dur").String() != "inf" {
		RstDurController = false
		go dur_controller(scene)
	}
	return
}

func dur_controller(scene string) {
	splitedKeyDur := strings.Split(SceneConf.Section(scene).Key("Dur").String(), ":")
	if splitedKeyDur[0] == "T" {
		dur, err := time.ParseDuration(splitedKeyDur[1])
		Err_handle(err)
		for i := time.Duration(0); i < dur; i += 1 * time.Millisecond {
			switch RstDurController {
			case true:
				time.Sleep(1 * time.Millisecond)
			default:
				return
			}
		}
		Close_ffplay()
	} else if splitedKeyDur[0] == "M" {
		for i := time.Duration(0); i < Get_video_dur(MediaConf.Section(splitedKeyDur[1]).Key("File").String()); i += 1 * time.Millisecond {
			switch RstDurController {
			case true:
				time.Sleep(1 * time.Millisecond)
			default:
				return
			}
		}
		Close_ffplay()
	}
	return
}

func Auto_play() {
	if !DispInited {
		Disp_init()
		time.Sleep(4 * time.Second)
	}
	if QuitAutoPlay {
		return
	}
	if LatestScene == "(nil)" {
		LatestScene = SceneConf.SectionStrings()[1]
	}
	sceneSections := []string{}
	if SceneConf != nil {
		sceneSections = SceneConf.SectionStrings()
	}
	for i := Find_str(sceneSections, LatestScene); SceneConf != nil && len(sceneSections) >= 2 && i != -1 && i < len(sceneSections); i++ {
		Show_scene(sceneSections[i])
		splitedKeyDur := strings.Split(SceneConf.Section(sceneSections[i]).Key("Dur").String(), ":")
		if splitedKeyDur[0] == "T" {
			dur, err := time.ParseDuration(splitedKeyDur[1])
			Err_handle(err)
			for i := time.Duration(0); i < dur; i += 1 * time.Millisecond {
				switch QuitAutoPlay {
				case false:
					time.Sleep(1 * time.Millisecond)
				default:
					return
				}
			}
		} else if splitedKeyDur[0] == "M" {
			for i := time.Duration(0); MediaConf != nil && i < Get_video_dur(MediaConf.Section(splitedKeyDur[1]).Key("File").String()); i += 1 * time.Millisecond {
				switch QuitAutoPlay {
				case true:
					time.Sleep(1 * time.Millisecond)
				default:
					return
				}
			}
		}
	}
	return
}
