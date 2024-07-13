package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/gorilla/websocket"
)

// type tunnel struct {
// 	Name string

// }

// type get_tunnels_resp struct {
// 	Tunnels []
// }

type Tunnels_header struct {
    Tunnels []Tunnel        `json:"tunnels"`
}


func main() {
    // upgrader := websocket.Upgrader{
    // 	ReadBufferSize:  1024,
    // 	WriteBufferSize: 1024,
    // 	// CheckOrigin: func (r *http.Request) bool {
    // 	// 	fmt.Printf("[xdddd] Origin is %s\n", r.Header.Get("origin"))
    // 	// 	fmt.Printf("[xdddd] Host %s\n", r.Header.Get("Host"))
    // 	// 	return true
    // 	// },
    // }
    router := gin.Default()
    // hub := Hub{
    // 	clients:    make(map[*Client]bool),
    // 	register:   make(chan *Client),
    // 	unregister: make(chan *Client),
    // }
    router.Use(cors.Default())

    var tunnels []Tunnel

    dat, err := os.ReadFile(".tunnels.json")
    if (err == nil) {
        fmt.Println("./tunnels.json found! Loading saved configuration")
        var f Tunnels_header
        err = json.Unmarshal(dat, &f)
        // fmt.Println("=====", f.Tunnels)
        if (err != nil) {
            fmt.Println("Error occured while processing tunnels.json: ", err)
            return
        }
        tunnels = f.Tunnels
    } else {
        fmt.Println("tunnels.json not found.")
        tunnels = make([]Tunnel, 0)
    }

    spawner := init_spawner(tunnels)
    
    const apiv1 = "/api/v1"

    router.GET("/", func(c *gin.Context) {
        router.LoadHTMLFiles("index.html")
        c.HTML(200, "index.html", gin.H{})
    })

    router.GET(apiv1 + "/tunnel_status", func(c *gin.Context) {
        t := make([]interface{}, 0)
        for i := range spawner.tunnels {
            if (spawner.tunnels[i].Enabled) {
                id := spawner.tunnels[i].Id
                fmt.Println(*(spawner.procs[id].tunnel))
                t = append(t, struct {
                    Tunnel Tunnel       `json:"tunnel"`
                    Status TunnelStatus `json:"status"`
                }{
                    spawner.tunnels[i],
                    spawner.procs[id].status,
                })
            } else {
                t = append(t, struct {
                    tunnel Tunnel
                }{
                    spawner.tunnels[i],
                })
            }
        }
        c.JSON(200, gin.H{
            "tunnel_status": t,
        })
    })

    router.POST(apiv1 + "/remove_tunnel", func(c *gin.Context) {
        var req struct {
            Id string       `json:"id"`
        }

        if err := c.BindJSON(&req); err != nil {
            c.JSON(400, gin.H {
                "status": "Invalid JSON!",
            })
            return
        }


        proc, exists := spawner.procs[req.Id]
        if (!exists) {
            c.JSON(400, gin.H {
                "status": "Tunnel not found!",
            })
            return
        }
        proc.autoreboot_chan <- false
        log.Println("xd")
        // close(proc.autoreboot_chan)
        delete(spawner.procs, req.Id)
        delete(spawner.tunnels, req.Id) 

        /* this should never be nil, but if it does, it's probably safe to be quiet about it */
        if (proc.cmd.Process != nil) {
            proc.cmd.Process.Signal(syscall.SIGHUP)
        }

        c.JSON(200, gin.H {
            "status": "Tunnel deleted",
        })
    })

    // router.GET("/api/v1/get_tunnels")

    // router.GET("/ws", func(ctx *gin.Context) {
    // 	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    // 	if err != nil {
    // 		fmt.Println(err)
    // 		return
    // 	}
    // 	client := Client{&hub, conn}
    // 	hub.register <- &client
    // 	go client.readPump()
    // })

    // go hub.handle_events()
    router.Run("0.0.0.0:4140")
}
