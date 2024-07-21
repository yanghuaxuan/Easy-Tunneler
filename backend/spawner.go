package main

import (
	// "bytes"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"syscall"
	"time"

	"math/rand"
	// "time"
)

const AUTOREBOOT_TIMEOUT = time.Second * 30

type TunnelStatus int

const (
	Disconnected TunnelStatus = iota
	Loading
	Online
)

type Tunnel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Local_port  int    `json:"local_port"`
	Host        string `json:"host"`
	Remote_port int    `json:"remote_port"`
	Conn_addr   string `json:"conn_addr"`
	Autoreboot  bool   `json:"autoreboot"`
}

type Tunnel_Process struct {
	cmd             *exec.Cmd
	tunnel          Tunnel
	status          TunnelStatus
	autoreboot_chan chan bool
	// history string
}

type Spawner struct {
	tunnels map[string]Tunnel
	procs   map[string]*Tunnel_Process
}

/* stops a tunnel, if it exists in proc */
func (s *Spawner) stop_tunnel(tunId string) {
	_, exists := s.tunnels[tunId]
	if (!exists) {
		return
	}
	p, exists := s.procs[tunId]
	if !exists {
		return
	}
	p.autoreboot_chan <- false
	delete(s.procs, tunId)
	if p.cmd.Process != nil {
		p.cmd.Process.Signal(syscall.SIGHUP)
	}
}

/* use this to properly start a tunnel */
func (s *Spawner) start_tunnel(tunId string) {
	tun, exists := s.tunnels[tunId]
	if (!exists) {
		return
	}

    proc := kickstart(tun)
    s.procs[tunId] = &proc
    go track_exit(&proc)
    if tun.Autoreboot {
        go auto_reboot_on_sig(&proc)
    }
}

/* very basic id builder */
func genId(n int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range n {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func init_spawner(tun []Tunnel) Spawner {
	tun_map := make(map[string]Tunnel)
	proc_map := make(map[string]*Tunnel_Process)

    s := Spawner{tun_map, proc_map}

	for i := range tun {
		t := tun[i]
		s.tunnels[t.Id] = t
		if t.Enabled {
			// proc := kickstart(t)
			// log.Println("init_spawner: ", proc.tunnel)
			// proc_map[t.Id] = &proc
			// go track_exit(&proc)
			// if t.Autoreboot {
			// 	go auto_reboot_on_sig(&proc)
			// }
            s.start_tunnel(t.Id)
		}
	}

	return Spawner{tun_map, proc_map}
}

func track_exit(tun *Tunnel_Process) {
	if tun == nil {
		return
	}

	tun.cmd.Wait()
	slog.Debug("SSH session exited!")
	tun.status = Disconnected
	tun.autoreboot_chan <- true
}

func log_tunnel(tun Tunnel, rc io.ReadCloser) {
	b := make([]byte, 1024)
	for {
		n, err := rc.Read(b);
		if (n == 0) {
			continue
		}
		slog.Warn(fmt.Sprintf("Message from %s ->\n%s", tun.Name, b[:n]))
		if err == io.EOF {
			break
		}
	}
}

/* attempts to start the SSH process if autoreboot_chan received a true value.  */
func auto_reboot_on_sig(proc *Tunnel_Process) {
	if proc == nil {
		return
	}

	s := <-proc.autoreboot_chan
	if !s {
		slog.Debug("Exiting autoreboot!")
		return
	}

	slog.Debug("Autorebooting!")

	tun := proc.tunnel
	cmd := exec.Command("ssh", "-o", "ExitOnForwardFailure=yes", "-N", "-L", fmt.Sprintf("%d:%s:%d", tun.Local_port, tun.Host, tun.Remote_port), tun.Conn_addr)
	// cmd.Stderr = os.Stderr
	stderr, err := cmd.StderrPipe()
	if (err == nil) {
		go log_tunnel(tun, stderr)
	} else {
		slog.Warn("Cannot log a SSH session!")
	}
	proc.cmd = cmd
	slog.Debug(cmd.String())
	// var stderrBuffer bytes.Buffer
	// cmd.Stderr = &stderrBuffer
	err = cmd.Start()
	proc.status = Online
	if err != nil {
		proc.status = Disconnected
	} else {
		go track_exit(proc)
	}

	time.Sleep(AUTOREBOOT_TIMEOUT)
	go auto_reboot_on_sig(proc)
}

/* start SSH session for tunnel and return its process */
func kickstart(tun Tunnel) Tunnel_Process {
	cmd := exec.Command("/usr/bin/ssh", "-o", "ExitOnForwardFailure yes", "-N", "-L", fmt.Sprintf("%d:%s:%d", tun.Local_port, tun.Host, tun.Remote_port), tun.Conn_addr)
	// cmd.Stderr = os.Stderr
	stderr, err := cmd.StderrPipe()
	if (err == nil) {
		go log_tunnel(tun, stderr)
	} else {
		slog.Warn("Cannot log a SSH session!")
	}
	status := Online
	slog.Debug(cmd.String())
	// var stderrBuffer bytes.Buffer
	// cmd.Stderr = &stderrBuffer
	err = cmd.Start()
	if err != nil {
		status = Disconnected
	}
	/* 2 buffered channel is necessary to avoid deadblocks (i.e. removing a tunnel with no autoreboot) */
	return Tunnel_Process{cmd, tun, status, make(chan bool, 2)}
}
