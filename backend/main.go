package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

const AUTOSAVE_INTERVAL = time.Second * 60
const TUN_SAVE_FILE = ".tunnels.json"

type Tunnels_header struct {
	Tunnels []Tunnel `json:"tunnels"`
}

func save_tunnels(tun []Tunnel) {
	i := Tunnels_header{
		tun,
	}
	dat, err := json.Marshal(i)
	if err != nil {
		slog.Error("Cannot save to .tunnels.json!", slog.Any("error", err))
		return
	}
	os.WriteFile(TUN_SAVE_FILE, dat, 0644)
}

func main() {
	var router *gin.Engine

	if os.Getenv("EASY_TUNNELER_PROD") == "1" {
		router = gin.New()
		router.Use(gin.Recovery())
		router.Use(static.Serve("/", static.LocalFile("./public", false)))
	} else {
		router = gin.Default()
		/* relax CORS for development */
		slog.SetLogLoggerLevel(slog.LevelDebug)
		router.Use(cors.Default())
		slog.Info("You are running Easy-Tunneler in non-production mode. The frontend side is not served by the in this mode.To switch to production mode, set EASY_TUNNELER_PROD=1 in your environment.")
	}

	var tunnels []Tunnel

	dat, err := os.ReadFile(".tunnels.json")
	if err == nil {
		slog.Info("./tunnels.json found! Loading saved configuration")
		var f Tunnels_header
		err = json.Unmarshal(dat, &f)
		if err != nil {
			slog.Error("Error occured while processing tunnels.json: ", slog.Any("error", err))
			return
		}
		tunnels = f.Tunnels
	} else {
		slog.Info("tunnels.json not found. Starting from a fresh instance.")
		tunnels = make([]Tunnel, 0)
	}

	ssh_path, err := try_ssh()
	if (err != nil) {
		slog.Error("Cannot find SSH! Do you have OpenSSH installed?")
		return 
	}


	spawner := init_spawner(tunnels, ssh_path)

	/* tunnel.json autosaver*/
	stop_autosave := make(chan bool)
	autosave_ticker := time.NewTicker(AUTOSAVE_INTERVAL)
	go func() {
		for {
			select {
			case <-stop_autosave:
				return
			case <-autosave_ticker.C:
				slog.Debug("Autosaving!")
				a := make([]Tunnel, 0)
				for i := range spawner.tunnels {
					a = append(a, spawner.tunnels[i])
				}
				save_tunnels(a)
			}
		}
	}()

	const apiv1 = "/api/v1"

	router.GET(apiv1+"/tunnel_status", func(c *gin.Context) {
		t := make([]interface{}, 0)
		for i := range spawner.tunnels {
			if spawner.tunnels[i].Enabled {
				id := spawner.tunnels[i].Id
				t = append(t, struct {
					Tunnel Tunnel       `json:"tunnel"`
					Status TunnelStatus `json:"status"`
				}{
					spawner.tunnels[i],
					spawner.procs[id].status,
				})
			} else {
				t = append(t, struct {
					Tunnel Tunnel		`json:"tunnel"`
				}{
					spawner.tunnels[i],
				})
			}
		}
		c.JSON(200, gin.H{
			"tunnel_status": t,
		})
	})

	router.POST(apiv1+"/remove_tunnel", func(c *gin.Context) {
		var req struct {
			Id string `json:"id"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"status": "Invalid JSON!",
			})
			return
		}

		_, exists := spawner.tunnels[req.Id]
		if !exists {
			c.JSON(400, gin.H{
				"status": "Tunnel not found!",
			})
			return
		}

		spawner.stop_tunnel(req.Id)

		delete(spawner.tunnels, req.Id)

		c.JSON(200, gin.H{
			"status": "Tunnel deleted",
		})
	})

	router.POST(apiv1+"/add_tunnel", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Enabled     bool   `json:"enabled"`
			Local_port  int    `json:"local_port"`
			Host        string `json:"host"`
			Remote_port int    `json:"remote_port"`
			Conn_addr   string `json:"conn_addr"`
			Autoreboot  bool   `json:"autoreboot"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"status": "Invalid JSON!",
			})
			return
		}
		tid := genId(16)
		for {
			_, exists := spawner.tunnels[tid]
			if (!exists) {
				break
			}
		}
		t := Tunnel{
			tid,
			req.Name,
			req.Enabled,
			req.Local_port,
			req.Host,
			req.Remote_port,
			req.Conn_addr,
			req.Autoreboot,
		}
		spawner.tunnels[t.Id] = t
		if req.Enabled {
            spawner.start_tunnel(t.Id)
		}
	})

	router.PATCH(apiv1+"/update_tunnel", func(c *gin.Context) {
		var newT Tunnel
		if err := c.BindJSON(&newT); err != nil {
			c.JSON(400, gin.H{
				"status": "Invalid JSON!",
			})
			return
		}
		oldT, exists := spawner.tunnels[newT.Id]
		if !exists {
			c.JSON(400, gin.H{
				"status": "Tunnel provided to update does not exist!",
			})
			return
		}
		if !newT.Enabled && oldT.Enabled {
			spawner.stop_tunnel(oldT.Id)
			spawner.tunnels[newT.Id] = newT
		} else if newT.Enabled && !oldT.Enabled {
			spawner.tunnels[newT.Id] = newT
            spawner.start_tunnel(oldT.Id)
		} else {
			spawner.stop_tunnel(oldT.Id)
			spawner.tunnels[newT.Id] = newT
            spawner.start_tunnel(oldT.Id)
		}

		c.JSON(200, gin.H{
			"status": "Updated tunnel settings",
		})
	})

	fmt.Println()
	fmt.Println("===================================================")
	fmt.Println("Easy-Tunneler running at http://localhost:4140.")
	fmt.Println("===================================================")
	fmt.Println()

	err = router.Run("localhost:4140")
	if (err != nil) {
		slog.Error(fmt.Sprintf("%s", err))
	}
}
