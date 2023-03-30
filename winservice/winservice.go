package winservice

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
	"github.com/yxcloud1/go-comm/logger"
)

type ServiceProgram interface {
	Run() error
	Kill() error
}

type WinService struct {
	prog   ServiceProgram
	config *service.Config
}

type ServiceInstance struct {
	kill func() error
	run  func() error
}

func (s *ServiceInstance) Run() error {
	if s.run != nil {
		return s.run()
	}
	return nil
}

func (s *ServiceInstance) Kill() error {
	if s.kill != nil {
		return s.kill()
	}
	return nil
}

func (p *WinService) Start(s service.Service) error {
	go p.prog.Run()
	logger.TxtLog("服务已经启动")
	return nil
}

func (p *WinService) Stop(s service.Service) error {
	err := p.prog.Kill()
	if err != nil {
		logger.TxtErr(p.config.Name, err)
	}
	logger.TxtLog("服务已经停止")
	return err
}

func runAsNoService(run func() error, kill func() error) error {
	sigs := make(chan os.Signal, 1)
	defer func(){
		close(sigs)
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	run()
	<-sigs
	return kill()
}

func RunAsService(serviceName string, displayName string, description string, run func() error, kill func() error) {

	svrCfg := &service.Config{
		Name:        serviceName,
		DisplayName: displayName,
		Description: description,
	}

	//	logFile, _ := os.OpenFile(filepath.Dir(os.Args[0])+"/log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	svr := &WinService{
		config: svrCfg,
		prog: &ServiceInstance{
			kill: kill,
			run:  run,
		},
	}
	s, err := service.New(svr, svrCfg)
	if err != nil {
		log.Fatal(err)
	}
	/*	log.SetOutput(logFile)


	 */
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err := s.Install()
			if err == nil {
				logger.TxtLog(serviceName, "服务安装成功！")
			} else {
				logger.TxtErr(serviceName, err)
			}
			return
		}
		if (os.Args[1]) == "remove" {
			err = s.Uninstall()
			if err == nil {
				logger.TxtLog(serviceName, "服务卸载成功！")
			} else {
				logger.TxtErr(serviceName, err)
			}
			return
		}
		if os.Args[1] == "run" {
			if run != nil {
				runAsNoService(run, kill)
			}
			return
		}
	}
	err = s.Run()
	if err != nil {
		logger.TxtErr(serviceName, serviceName, err)
	}

}
