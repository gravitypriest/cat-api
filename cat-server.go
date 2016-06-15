package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"net/http"
	"regexp"
	"strings"
)

type Cat struct {
	Message  string `json:"message"`
	Position string `json:"position"`
	Picture  string `json:"picture"`
}

func generateHtml(cat *Cat) string {
	html := `<html>
				<header>
			 		<title>%s</title>
			 		<link rel="stylesheet" href="__assets/cat.css"/>
			 	</header>
			    <body style="background:black">
			        <div class="bgstyle" style="background-image:url(__assets/%s.jpg);">
			        	<span class="caption caption%s">%s</span>
			        </div>
			        <audio autoplay loop>
			        	<source src="__assets/music.mp3" type="audio/mpeg"/>
			        	<source src="__assets/music.ogg" type="audio/ogg"/>
		        	</audio>
    		   	</body>
			 </html>`
	return fmt.Sprintf(html, cat.Message, cat.Picture, cat.Position, cat.Message)
}

func sendResponse(ctx echo.Context) error {
	cat := new(Cat)
	cat.Message = "come on now"
	cat.Position = "2"
	cat.Picture = "cat"

	params := ctx.ParamNames()
	for _, p := range params {
		switch p {
		case "message":
			cat.Message = ctx.Param(p)
		case "position":
			cat.Position = ctx.Param(p)
		case "picture":
			cat.Picture = ctx.Param(p)
		default:
			// [shrugging intensifies]
		}
	}
	html := generateHtml(cat)
	return ctx.HTML(http.StatusOK, html)
}

func unfuckPath() echo.MiddlewareFunc {
	// remove duplicate slashes
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()
			url := req.URL()
			path := url.Path()

			reg, _ := regexp.Compile("(/+)")
			path = reg.ReplaceAllString(path, "/")

			req.SetURI(path)
			url.SetPath(path)
			return next(ctx)
		}
	}
}

func fuckRouting() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()
			url := req.URL()
			path := url.Path()

			if strings.Contains(path, "/__assets/") {
				pathEnd := strings.SplitN(path, "/__assets/", 2)[1]
				path = "/__assets/" + pathEnd
				req.SetURI(path)
				url.SetPath(path)
			}

			return next(ctx)
		}
	}
}

func main() {
	srv := echo.New()

	srv.Pre(unfuckPath())
	srv.Pre(middleware.RemoveTrailingSlash())
	srv.Pre(fuckRouting())

	srv.GET("/", sendResponse)
	srv.GET("/__assets", func(ctx echo.Context) error {
		return ctx.HTML(http.StatusForbidden, "nah dude")
	})
	srv.GET("/:message", sendResponse)
	srv.GET("/:message/:picture", sendResponse)
	srv.GET("/:message/:picture/:position", sendResponse)
	srv.GET("/:message/:picture/:position/*", sendResponse)
	srv.Static("/__assets", "assets")
	srv.Run(standard.New(":8080"))
}
