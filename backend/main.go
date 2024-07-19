package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	// "github.com/gorilla/websocket"
)

const tun_save_file = ".tunnels.json"

type Tunnels_header struct {
	Tunnels []Tunnel `json:"tunnels"`
}

func save_tunnels(tun []Tunnel) {
	i := Tunnels_header{
		tun,
	}
	dat, err := json.Marshal(i)
	if err != nil {
		log.Println("Cannot save to .tunnels.json!", err)
		return
	}
	os.WriteFile(tun_save_file, dat, 0644)
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery())

	if os.Getenv("EASY_TUNNELER_PROD") == "1" {
		// var serv embed.FS
		gin.DefaultWriter = nil
		router.Use(static.Serve("/", static.LocalFile("./public", false)))
	} else {
		/* relax CORS for development */
		router.Use(cors.Default())
		log.Println("You are running Easy-Tunneler in non-production mode. The frontend side is not served by the in this mode.To switch to production mode, set EASY_TUNNELER_PROD=1 in your environment.")
	}

	var tunnels []Tunnel

	dat, err := os.ReadFile(".tunnels.json")
	if err == nil {
		log.Println("./tunnels.json found! Loading saved configuration")
		var f Tunnels_header
		err = json.Unmarshal(dat, &f)
		// fmt.Println("=====", f.Tunnels)
		if err != nil {
			log.Println("Error occured while processing tunnels.json: ", err)
			return
		}
		tunnels = f.Tunnels
	} else {
		log.Println("tunnels.json not found.")
		tunnels = make([]Tunnel, 0)
	}

	spawner := init_spawner(tunnels)

	/* tunnel.json autosaver*/
	stop_autosave := make(chan bool)
	autosave_ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-stop_autosave:
				return
			case <-autosave_ticker.C:
				log.Println("Autosaving!")
				a := make([]Tunnel, 0)
				for i := range spawner.tunnels {
					a = append(a, spawner.tunnels[i])
				}
				save_tunnels(a)
			}
		}
	}()

	const apiv1 = "/api/v1"

	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLFiles("index.html")
		c.HTML(200, "index.html", gin.H{})
	})

	router.GET(apiv1+"/tunnel_status", func(c *gin.Context) {
		t := make([]interface{}, 0)
		for i := range spawner.tunnels {
			if spawner.tunnels[i].Enabled {
				id := spawner.tunnels[i].Id
				// fmt.Println(*(spawner.procs[id].tunnel))
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

		spawner.stop_tunnel(spawner.tunnels[req.Id])

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
            spawner.start_tunnel(t)
		}
	})

	router.PATCH(apiv1+"/update_tunnel", func(c *gin.Context) {
		var req Tunnel
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(400, gin.H{
				"status": "Invalid JSON!",
			})
			return
		}
		t, exists := spawner.tunnels[req.Id]
		if !exists {
			c.JSON(400, gin.H{
				"status": "Tunnel provided to update does not exist!",
			})
			return
		}
        /* TODO */
        if req.Autoreboot != t.Autoreboot {
            c.JSON(400, gin.H{
                "status": "Not implemented (yet)",
            })
			return
        }
		if !req.Enabled && t.Enabled {
			spawner.stop_tunnel(t)
		}
		if req.Enabled && !t.Enabled {
            spawner.start_tunnel(t)
		}
		spawner.tunnels[req.Id] = req
		c.JSON(200, gin.H{
			"status": "Updated tunnel settings",
		})
	})

	fmt.Println()
	fmt.Println("===================================================")
	fmt.Println("Easy-Tunneler running at http://localhost:4140.")
	fmt.Println("===================================================")
	fmt.Println()

	router.Run("localhost:4140")
}
