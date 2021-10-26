package controllers

import (
	"fmt"
	"go-revel-crud/app"
	"go-revel-crud/app/db"
	"go-revel-crud/app/models"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"status":  http.StatusOK,
		"success": true,
	})
}

func (c App) bytesToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (c App) Health() revel.Result {
	ctx := c.Request.Context()

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	hostName, err := os.Hostname()
	if err != nil {
		c.Log.Errorf("could not get hostname: %v", err)
	}

	currentTime := time.Now()
	tZone, offset := currentTime.Zone()

	dbUp := true
	err = db.DB().Ping()
	if err != nil {
		dbUp = false
		c.Log.Errorf("db ping failed because err=[%v]", err)
	}

	hello, err := (&models.Checker{}).One(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not db select 1: %v", err)
	}

	version, err := (&models.Checker{}).Version(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not db version: %v", err)
	}

	return c.RenderJSON(map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
		"time": map[string]interface{}{
			"now":      currentTime,
			"timezone": tZone,
			"offset":   offset,
		},
		"version":    app.AppVersion,
		"build_time": app.BuildTime,
		"db": map[string]interface{}{
			"type":    "postgres",
			"up":      dbUp,
			"hello":   hello,
			"version": version,
		},
		"server": map[string]interface{}{
			"hostname":   hostName,
			"cpu":        runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
			"memory": map[string]interface{}{
				"alloc":       fmt.Sprintf("%v MB", c.bytesToMb(memStats.Alloc)),
				"total_alloc": fmt.Sprintf("%v MB", c.bytesToMb(memStats.TotalAlloc)),
				"sys":         fmt.Sprintf("%v MB", c.bytesToMb(memStats.Sys)),
				"num_gc":      memStats.NumGC,
			},
			"goarch":   runtime.GOARCH,
			"goos":     runtime.GOOS,
			"compiler": runtime.Compiler,
		},
	})
}
