package main

import (
	// "bytes"
	"fmt"
	"log"
	"os"
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
func (s *Spawner) stop_tunnel(tun Tunnel) {
	p, exists := s.procs[tun.Id]
	if !exists {
		return
	}
	p.autoreboot_chan <- false
	delete(s.procs, tun.Id)
	if p.cmd.Process != nil {
		p.cmd.Process.Signal(syscall.SIGHUP)
	}
}

/* use this to properly start a tunnel */
func (s *Spawner) start_tunnel(tun Tunnel) {
    proc := kickstart(tun)
    s.procs[tun.Id] = &proc
    go track_exit(&proc)
    if tun.Autoreboot {
        go auto_reboot_on_sig(&proc)
    }
}

/* very basic id builder */
const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genId(n int) string {
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
            s.start_tunnel(t)
		}
	}

	return Spawner{tun_map, proc_map}
}

func track_exit(tun *Tunnel_Process) {
	if tun == nil {
		return
	}

	tun.cmd.Wait()
	log.Println("SSH session exited!")
	tun.status = Disconnected
	tun.autoreboot_chan <- true
}

/* attempts to start the SSH process if autoreboot_chan received a true value.  */
func auto_reboot_on_sig(proc *Tunnel_Process) {
	if proc == nil {
		return
	}

	s := <-proc.autoreboot_chan
	if !s {
		log.Println("Exiting autoreboot!")
		return
	}

	log.Println("Autorebooting!")

	tun := proc.tunnel
	cmd := exec.Command("ssh", "-o", "ExitOnForwardFailure=yes", "-N", "-L", fmt.Sprintf("%d:%s:%d", tun.Local_port, tun.Host, tun.Remote_port), tun.Conn_addr)
	cmd.Stderr = os.Stderr
	proc.cmd = cmd
	log.Println(cmd)
	// var stderrBuffer bytes.Buffer
	// cmd.Stderr = &stderrBuffer
	err := cmd.Start()
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
	cmd.Stderr = os.Stderr
	status := Online
	log.Println(cmd)
	// var stderrBuffer bytes.Buffer
	// cmd.Stderr = &stderrBuffer
	err := cmd.Start()
	if err != nil {
		status = Disconnected
	}
	/* 2 buffered channel is necessary to avoid deadblocks (i.e. removing a tunnel with no autoreboot) */
	return Tunnel_Process{cmd, tun, status, make(chan bool, 2)}
}
