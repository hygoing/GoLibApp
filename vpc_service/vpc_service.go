package main

import (
	"GoLibApp/net/ssh_server"
	"fmt"
	"jd.com/cc/jstack-cc-ops/api/cfg"
	"jd.com/cc/jstack-cc-ops/api/log"
	"jd.com/cc/jstack-cc-ops/api/service"
	"jd.com/cc/jstack-cc-ops/api/utils"
	"jd.com/cc/jstack-cc-ops/model"
	"strings"
	"time"
)

type Manager struct {
	workers map[string]*Worker
}

func (m *Manager) Setup() {
	m.workers = make(map[string]*Worker)
}

func (m *Manager) process() {
	oldKeyMap := make(map[string]map[string]*model.RuntimeDetect)
	doDbTimer := time.NewTicker(time.Millisecond * time.Duration(cfg.ServiceConfig.SyncDbIntervalMs))
	for {
		select {
		case <-doDbTimer.C:
			// select data from db
			filterAdminStatus := "admin_status"
			descDetectsOpts := &model.DescribesParams{
				Filter: []*model.FilterItem{
					&model.FilterItem{
						Field: &filterAdminStatus,
						Value: utils.RUNTIME_DETECT_ADMIN_STATUS_UP,
					},
				},
			}
			detectPoints, _, err := service.RtDetectService.Describes(descDetectsOpts)
			if err != nil {
				return
			}
			detectors := make(map[string][]*model.RuntimeDetect)
			for _, detectPoint := range detectPoints {
				detectors[detectPoint.IpAddress] = append(detectors[detectPoint.IpAddress], detectPoint)
			}
			newKeyMap := m.genkeyMap(detectors)
			m.diffAndUpdate(oldKeyMap, newKeyMap)
			oldKeyMap = newKeyMap
		}
	}
}

func (m *Manager) diffAndUpdate(oldKeyMap, newKeyMap map[string]map[string]*model.RuntimeDetect) {
	for ip, newKeys := range newKeyMap {
		if _, ok := oldKeyMap[ip]; !ok {
			worker := &Worker{
				sshIp:        ip,
				detectPoints: []*model.RuntimeDetect{},
				exit:         make(chan interface{}),
			}
			for _, newDetectPoint := range newKeys {
				worker.detectPoints = append(worker.detectPoints, newDetectPoint)
			}
			m.workers[ip] = worker
		} else {
			updated := false
			worker := &Worker{
				sshIp:        ip,
				detectPoints: []*model.RuntimeDetect{},
				exit:         make(chan interface{}),
			}
			for newKey, newDetectPoint := range newKeys {
				if _, ok := oldKeyMap[ip][newKey]; !ok {
					updated = true
				}
				worker.detectPoints = append(worker.detectPoints, newDetectPoint)
			}
			for oldKey, _ := range oldKeyMap[ip] {
				if _, ok := newKeys[oldKey]; !ok {
					updated = true
					break
				}
			}
			if updated {
				m.workers[ip] = worker
			}
		}
	}
	for ip, _ := range oldKeyMap {
		if _, ok := newKeyMap[ip]; !ok {
			m.workers[ip].close()
			delete(m.workers, ip)
		}
	}

	for _, worker := range m.workers {
		go worker.process()
	}
}

func (m *Manager) genkeyMap(cache map[string][]*model.RuntimeDetect) map[string]map[string]*model.RuntimeDetect {
	keyMap := make(map[string]map[string]*model.RuntimeDetect)
	for ip, detectPoints := range cache {
		for _, detectPoint := range detectPoints {
			if _, ok := keyMap[ip]; !ok {
				keyMap[ip] = make(map[string]*model.RuntimeDetect)
			}
			keyMap[ip][m.key(detectPoint)] = detectPoint
		}
	}
	return keyMap
}

func (m *Manager) key(v *model.RuntimeDetect) string {
	key := fmt.Sprintf("%s@%s@%d@%d", v.DstIpAddress, v.Protocol, &v.Port, v.Detection)
	return key
}

type Worker struct {
	sshIp        string
	detectPoints []*model.RuntimeDetect
	exit         chan interface{}
}

func (w *Worker) process() {
	icmpResMap := map[int]string{utils.DETECTION_ACTIVE: "alive", utils.DETECTION_NEGTIVE: "unreachable"}
	tcpResMap := map[int]string{utils.DETECTION_ACTIVE: "0", utils.DETECTION_NEGTIVE: "100"}

	timer := time.NewTicker(time.Millisecond * time.Duration(cfg.ServiceConfig.PingIntervalMs))
	select {
	case <-timer.C:
		sshClient := ssh_server.Init(w.sshIp, "root", "siqin*****")
		icmpList := make([]*model.RuntimeDetect, 0)
		tcpList := make([]*model.RuntimeDetect, 0)
		detectPointMap := make(map[string]*model.RuntimeDetect)

		for _, detectPoint := range w.detectPoints {
			switch detectPoint.Protocol {
			case utils.RUNTIME_DETECT_PING:
				icmpList = append(icmpList, detectPoint)
				detectPointMap[w.key(detectPoint.DstIpAddress, utils.RUNTIME_DETECT_PING, 0)] = detectPoint
			case utils.RUNTIME_DETECT_CURL:
				tcpList = append(tcpList, detectPoint)
			}
		}

		// icmp
		icmpResult := make(map[string]string)
		icmpCmd := fmt.Sprintf("fping -r 0 -t %d", cfg.ServiceConfig.IcmpTimeout)
		for _, detector := range icmpList {
			icmpCmd = fmt.Sprintf("%s %s", detector.DstIpAddress)
		}
		icmpCmd = fmt.Sprintf("%s %s", icmpCmd, "| awk '{print $1\":\"$3}")
		result, err := sshClient.Run(icmpCmd)
		if err != nil {
			log.Error("[ssh %s run %s error: %s]", w.sshIp, icmpCmd, err.Error())
		}
		lines := strings.Split(strings.Trim(string(result), "\n"), "\n")
		for _, line := range lines {
			ipAndRes := strings.Split(line, ":")
			icmpResult[w.key(ipAndRes[0], utils.RUNTIME_DETECT_PING, 0)] = ipAndRes[1]
		}

		// valid icmp result
		for key, res := range icmpResult {
			detectPoint := detectPointMap[key]
			if res == icmpResMap[detectPoint.Detection] {
				log.Info("[ssh %s fping %s, detection: %d success]", w.sshIp, detectPoint.DstIpAddress, detectPoint.Detection)
			} else {
				log.Error("[ssh %s fping %s, detection: %d failed]", w.sshIp, detectPoint.DstIpAddress, detectPoint.Detection)
			}
		}

		// tcp
		for _, detectPoint := range tcpList {
			tcpCmd := fmt.Sprintf("hping %s -c 1 -p %d | grep -v Process", detectPoint.Description, detectPoint.Port)
			result, err := sshClient.Run(tcpCmd)
			if err != nil {
				continue
			}

			// valid tcp result
			lines = strings.Split(strings.Trim(string(result), "\n"), "\n")
			res := w.resolveResult(lines[1])
			if res == tcpResMap[detectPoint.Detection] {
				log.Info("[ssh %s hping %s:%d, detection: %d success]", w.sshIp, detectPoint.DstIpAddress, detectPoint.Port, detectPoint.Detection)
			} else {
				log.Error("[ssh %s hping %s:%d, detection: %d failed]", w.sshIp, detectPoint.DstIpAddress, detectPoint.Port, detectPoint.Detection)
			}
		}
	case <-w.exit:
		log.Info("[worker exit, sshIp: %s]", w.sshIp)
		return
	}

}

func (w *Worker) close() {
	close(w.exit)
}

func (w *Worker) key(ip string, protocol, port int) string {
	key := fmt.Sprintf("%s@%d@%d", ip, protocol, port)
	return key
}

func (w *Worker) resolveResult(str string) string {
	strArray := strings.Split(str, ", ")
	end := strings.Index(strArray[2], "%")
	return strArray[2][:end]
}
