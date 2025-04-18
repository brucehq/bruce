package system

import (
	"bruce/exe"
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	sys     *SystemInfo
	sysLock = new(sync.RWMutex)
)

type SystemInfo struct {
	OSType                string
	OSID                  string
	OSVersionID           string
	OSArch                string
	OsName                string
	PackageHandler        string
	PackageHandlerPath    string
	CurrentUser           *user.User
	CurrentGroup          *user.Group
	UseSudo               bool
	CanUpdateServices     bool
	ServiceControllerPath string
	ServiceController     string
	ModifiedTemplates     []string
}

func InitializeSysInfo() error {
	s := Get()
	// defaults
	s.OSArch = runtime.GOARCH
	if runtime.GOOS == "linux" {
		s.OSType = "linux"
		u, err := user.Current()
		if err != nil {
			log.Error().Err(err).Msg("user should exist to operate")
			return err
		}
		s.CurrentUser = u
		p := GetLinuxPackageHandler()
		s.PackageHandlerPath = p
		s.PackageHandler = path.Base(p)
		svcInfo, err := GetLinuxServiceController()
		if err != nil {
			s.CanUpdateServices = false
		} else {
			s.ServiceControllerPath = svcInfo
			s.ServiceController = path.Base(svcInfo)
		}
		hasOsData := ReadLinuxOsData(s)
		if !hasOsData {
			log.Error().Msgf("could not read os data must gather it another way?")
		}
		if s.OSArch == "amd64" {
			s.OSArch = "x86_64"
		}
	}
	s.Save()
	return nil
}

func ReadLinuxOsData(s *SystemInfo) bool {
	if _, err := os.Stat("/etc/os-release"); !os.IsNotExist(err) {
		readFile, err := os.Open("/etc/os-release")

		if err != nil {
			log.Error().Err(err).Msg("could not read /etc/os-release")
			return false
		}
		fs := bufio.NewScanner(readFile)
		fs.Split(bufio.ScanLines)

		for fs.Scan() {
			fields := strings.Split(fs.Text(), "=")
			if fields[0] == "ID" {
				s.OSID = s.wash(fields[1])
			}
			if fields[0] == "VERSION_ID" {
				s.OSVersionID = s.wash(fields[1])
			}
			if fields[0] == "VERSION_CODENAME" {
				s.OsName = s.wash(fields[1])
			}
			if s.OSVersionID != "" {
				return true
			}
			//fmt.Println(fs.Text())
		}

		readFile.Close()
	}
	s.OSArch = s.wash(exe.Run("uname -m", "").Get())
	return false
}

func (s *SystemInfo) wash(input string) string {
	bs := strings.ToLower(strings.Trim(input, " "))
	us, err := strconv.Unquote(bs)
	if err != nil {
		//log.Debug().Err(err).Msg("could not unquote string")
		return bs
	}
	return us
}

func GetLinuxServiceController() (string, error) {
	// We only support systemctl for now
	sysPath := exe.HasExecInPath("systemctl")
	if sysPath == "" {
		return "", fmt.Errorf("systemctl not found on this system")
	}
	return sysPath, nil
}

// Get function returns the currently set global system information to be used.
func Get() *SystemInfo {
	sysLock.RLock()
	defer sysLock.RUnlock()
	if sys == nil {
		sys = &SystemInfo{}
	}
	return sys
}

// Save saves.
func (s *SystemInfo) Save() {
	sysLock.Lock()
	defer sysLock.Unlock()
	sys = s
}

func (s *SystemInfo) AddModifiedTemplate(local string) {
	s.ModifiedTemplates = append(s.ModifiedTemplates, local)
	s.Save()
}
