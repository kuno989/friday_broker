package connect

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

func (s *Server) VmResolution(){
	args := []string {
		"controlvm",
		s.Config.VBoxName,
		"setvideomodehint",
		"1152",
		"864",
		"64",
	}
	run, err := s.RunCommand(s.Config.VBoxManage, args...)
	if err != nil {
		logrus.Error("Vbox 해상도 변경 중 에러", err)
	}
	logrus.Infof("Vbox 해상도 변경 \n%s",run)
}

func (s *Server) VmDump(sha256 string) string {
	tempMemoryDump := fmt.Sprintf("%s/%s.dmp",s.Config.VBoxDumpPath, sha256)
	args := []string {
		"debugvm",
		s.Config.VBoxName,
		"dumpvmcore",
		"--filename",
		tempMemoryDump,
	}
	run, err := s.RunCommand(s.Config.VBoxManage, args...)
	if err != nil {
		logrus.Error("메모리 덤프 중 에러", err)
	}
	logrus.Infof("메모리 덤프 완료 \n%s",run)
	return tempMemoryDump
}

func (s *Server) VmStart(){
	args := []string {
		"startvm",
		s.Config.VBoxName,
	}
	run, err := s.RunCommand(s.Config.VBoxManage, args...)
	if err != nil {
		logrus.Error("Vbox 시작 중 에러")
	}
	logrus.Infof("Vbox 시작 \n%s",run)
}

func (s *Server) VmStop(){
	args := []string {
		"controlvm",
		s.Config.VBoxName,
		"poweroff",
	}
	run, err := s.RunCommand(s.Config.VBoxManage, args...)
	if err != nil {
		logrus.Error("Vbox 중지 중 에러")
	}
	logrus.Infof("Vbox 중지 \n%s",run)
}

func (s *Server) VmRestore(){
	args := []string {
		"snapshot",
		s.Config.VBoxName,
		"restore",
		s.Config.VBoxSnapshot,
	}
	run, err := s.RunCommand(s.Config.VBoxManage, args...)
	if err != nil {
		logrus.Error("Snapshot 복구 중 에러")
	}
	logrus.Infof("SnapShot 복구중 \n%s",run)
}

func (s *Server) RunCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}