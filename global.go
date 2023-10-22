/*
 * @Author: Guo Yaoze, Qu Zhixuan, Sun Longyu, Pu Zhengyu, Zhang Yuntao, Peng Huanran, Yang Qihao
 * @License: AGPLv3
 * @Date: 2023-10-19 23:19:31
 * @LastEditTime: 2023-10-22 14:36:50
 * @Description: Global consts, vars and funcs.
 * @FilePath: /MikuMikuLive/global.go
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-ini/ini"
)

const (
	VER           string = "0.1-alpha" // Version.
	VER_CODENAME  string = "Miku"      // Version codename.
	FFMPEG_CHK_TS int    = 1024        //FFMpeg check output threshold.
)

var (
	NoWelcome        bool      = false       // Controls whether the welcome info will appear when you start it.
	ProgramName      string    = ""          // This is == os.Args[0].
	HttpPort         int       = 8230        // This is http srv port.
	HttpAddr         string    = "127.0.0.1" // Addr to listen on for http.
	TmpDir           string    = "tmp/"      // Dir to store the tmp files.
	LiveDescFile     string    = ""          // LDL File.
	UseExtFFConf     bool      = false       // Use external FFmpeg conf.
	UseIntDispConf   bool      = false       // Use internal display conf.
	TmpDirPerm       int       = 0700        // Tmp dir perm.
	PlayWithBorder   bool      = false       // Use FFPlay with border, for testing use.
	NoAlwaysOnTop    bool      = false       // Don't let FFPlay always on top, for testing use.
	FFPlayPIDs                 = []int{}     // FFPlay PIDs
	FFConf           *ini.File = nil         // FFMpeg conf file.
	DispConf         *ini.File = nil         // Display conf file.
	MediaConf        *ini.File = nil         // Media conf file.
	SceneConf        *ini.File = nil         // Scene conf file.
	RstDurController bool      = true        // As its name.
	QuitAutoPlay     bool      = true        // As its name.
	DispInited       bool      = false       // Is your display initialiezd?
	LatestScene      string    = ""          // LatestScene.
)

func Is_windows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func Err_handle(err error) bool {
	if err != nil {
		color.HiRed("Error: " + err.Error() + ".")
		return true
	}
	return false
}

func File_exist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	} else if os.IsExist(err) {
		return true
	} else {
		return false
	}
}

func Read_lines(name string) []string {
	var r = []string{}
	f, e := os.Open(name)
	Err_handle(e)
	defer func() {
		e := f.Close()
		c := 0
		for Err_handle(e) && c <= 8 {
			e = f.Close()
			c++
		}
	}()
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		r = append(r, fileScanner.Text())
	}
	return r
}

func Unify_path(p string, f bool) string {
	if p == "" {
		return p
	}
	check_EOP := func(path string, flag bool) bool {
		condition := ((!strings.HasSuffix(path, "/")) && flag) || (strings.HasSuffix(path, "/") && (!flag)) || ((!strings.HasSuffix(path, "\\")) && flag) || (strings.HasSuffix(path, "\\") && (!flag))
		if condition {
			return false
		}
		return true
	}
	if Is_windows() {
		p = strings.ReplaceAll(p, "/", "\\")
	}
	if check_EOP(p, false) && f {
		p = p + "/"
	}
	for check_EOP(p, true) && Is_windows() && (!f) {
		p = p[:len(p)-2]
	}
	return p
}

func Append_lines(name string, lines []string) []error {
	var errs = []error{}
	s, e := os.Stat(name)
	if Err_handle(e) {
		errs = append(errs, e)
	}
	f, e := os.OpenFile(name, os.O_RDWR|os.O_APPEND, s.Mode())
	if Err_handle(e) {
		errs = append(errs, e)
	}
	defer func() {
		e := f.Close()
		c := 0
		for Err_handle(e) && c <= 8 {
			e = f.Close()
			c++
		}
	}()
	tmp := []byte{0}
	if s.Size() > 0 {
		_, e = f.Seek(-1, io.SeekEnd)
		if Err_handle(e) {
			errs = append(errs, e)
		}
		_, _ = f.Read(tmp)
	}
	if s.Size() != 0 && string(tmp) != "\n" {
		f.Seek(0, io.SeekEnd)
		if Err_handle(e) {
			errs = append(errs, e)
		}
		_, e := f.WriteString("\n")
		if Err_handle(e) {
			errs = append(errs, e)
		}
	} else {
		f.Seek(0, io.SeekEnd)
	}
	for i := 0; i < len(lines); i++ {
		if !strings.HasSuffix(lines[i], "\n") {
			lines[i] = lines[i] + "\n"
		}
		_, e := f.WriteString(lines[i])
		if Err_handle(e) {
			errs = append(errs, e)
		}
	}
	return errs
}

func Create_file(name string, lines []string) error {
	var s string = ""
	f, e := os.Create(name)
	Err_handle(e)
	if lines != nil {
		for i := 0; i < len(lines); i++ {
			if strings.HasSuffix(lines[i], "\n") {
				s = s + lines[i]
			} else {
				s = s + lines[i] + "\n"
			}
		}
		_, e = f.WriteString(s)
		Err_handle(e)
	}
	return e
}

// Clear screen.
func Clear() {
	if Is_windows() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("bash", "-c", "clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Clean the tmp dir.
func Prepare_tmp(mkNewTmp bool) {
	if mkNewTmp {
		err := os.RemoveAll(Unify_path(TmpDir, false))
		os.Mkdir(Unify_path(TmpDir, false), fs.FileMode(TmpDirPerm))
		Err_handle(err)
	} else {
		Create_file(Unify_path(TmpDir, true)+"ffmpeg.conf", nil)
		Create_file(Unify_path(TmpDir, true)+"display.conf", nil)
		Create_file(Unify_path(TmpDir, true)+"media.conf", nil)
		Create_file(Unify_path(TmpDir, true)+"scene.conf", nil)
	}
}

func Pause() {
	var tmp string
	fmt.Print("Press Enter to continue...")
	fmt.Scanln(&tmp)
}

func Kill_all_ffplay() {
	for i := range FFPlayPIDs {
		if !Is_windows() {
			kill := exec.Command("kill", "-9", strconv.Itoa(int(math.Abs(float64(FFPlayPIDs[i])))))
			kill.Run()
		} else {
			kill := exec.Command("taskkill", "/f", "/pid", strconv.Itoa(int(math.Abs(float64(FFPlayPIDs[i])))))
			kill.Run()
		}
	}
	FFPlayPIDs = []int{}
}

func FFPlay(media string, disp string, extArgs []string) int {
	var ff *exec.Cmd
	if disp == "" {
		args := []string{media}
		if !PlayWithBorder {
			args = append(args, "-noborder")
		}
		if !NoAlwaysOnTop {
			args = append(args, "-alwaysontop")
		}
		if extArgs != nil {
			args = append(args, extArgs...)
		}
		ff = exec.Command(FFConf.Section("FFPLAY").Key("Exec").String(), args...)
	} else {
		left := DispConf.Section(disp).Key("Left").String()
		top := DispConf.Section(disp).Key("Top").String()
		width := DispConf.Section(disp).Key("Width").String()
		height := DispConf.Section(disp).Key("Height").String()
		args := []string{media, "-left", left, "-top", top, "-x", width, "-y", height}
		if !PlayWithBorder {
			args = append(args, "-noborder")
		}
		if !NoAlwaysOnTop {
			args = append(args, "-alwaysontop")
		}
		if len(extArgs) >= 1 {
			if extArgs[0] != "$PROTECTED$" {
				args = append(args, extArgs...)
			} else {
				args = append(args, extArgs[1:]...)
			}
		}
		ff = exec.Command(FFConf.Section("FFPLAY").Key("Exec").String(), args...)
	}
	ff.Start()
	if (len(extArgs) == 0) || (len(extArgs) >= 1 && (extArgs[0] != "$PROTECTED$")) {
		FFPlayPIDs = append(FFPlayPIDs, ff.Process.Pid)
	} else {
		FFPlayPIDs = append(FFPlayPIDs, -ff.Process.Pid)
	}
	return ff.Process.Pid
}

func Get_video_dur(name string) time.Duration {
	ff := exec.Command(FFConf.Section("FFPROBE").Key("Exec").String(), name)
	ffOut, err := ff.CombinedOutput()
	Err_handle(err)
	durStr := strings.Split(strings.Split(string(ffOut), "Duration: ")[1], ", ")[0]
	dur, err := time.ParseDuration(durStr)
	Err_handle(err)
	return dur
}

func Scan_str() string {
	var s string
	var err error
	r := bufio.NewReader(os.Stdin)
	s, err = r.ReadString('\n')
	s = strings.TrimSpace(s)
	Err_handle(err)
	return s
}

func Close_ffplay() {
	for i := 0; i < len(FFPlayPIDs); i++ {
		if FFPlayPIDs[i] > 0 {
			if !Is_windows() {
				kill := exec.Command("kill", "-9", strconv.Itoa(FFPlayPIDs[i]))
				kill.Run()
			} else {
				kill := exec.Command("taskkill", "/f", strconv.Itoa(FFPlayPIDs[i]))
				kill.Run()
			}
			if i == len(FFPlayPIDs)-1 {
				FFPlayPIDs = FFPlayPIDs[:i]
			} else {
				FFPlayPIDs = append(FFPlayPIDs[:i], FFPlayPIDs[i+1:]...)
			}
			i--
		}
	}
}

// Find a str in a []string. If not found, it will return -1.
func Find_str(source []string, target string) int {
	for i := 0; i < len(source); i++ {
		if source[i] == target {
			return i
		}
	}
	return -1
}
