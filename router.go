package main

import (
	"fmt"
	"net/http"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupRoutes(db datastore) {
	app := iris.New()

	cur_story := story{Id: 0, DBPtr: db}
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	verl := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		verl.Post("/add", func(ctx iris.Context) {
			var msg incomingMsg
			err := ctx.ReadJSON(&msg)
			if err != nil {
				log.WithFields(log.Fields{
					"invalidReq_error": err,
				}).Error("Unable to parse request from client!")

				ctx.Values().Set("message", "Malformed packet or empty terminal id")
				ctx.StatusCode(http.StatusUnprocessableEntity)
				return
			}
			status, resp := cur_story.addWord(msg.Word)
			ctx.StatusCode(status)
			ctx.JSON(resp)
		})
	}
	port := fmt.Sprintf(":%s", viper.GetString("Verloop.Port"))
	app.Run(iris.Addr(port), iris.WithCharset(viper.GetString("CharSet")))
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("<b>404 not found</b>")
}
