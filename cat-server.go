package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"net/http"
	"regexp"
)

type Cat struct {
	Message  string `json:"message"`
	Position string `json:"position"`
	Picture  string `json:"picture"`
}

func sendResponse(ctx echo.Context) error {
	cat := new(Cat)
	cat.Message = ""
	cat.Position = ""
	cat.Picture = ""

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
		}
	}
	return ctx.JSON(http.StatusOK, cat)
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

func main() {
	srv := echo.New()

	srv.Pre(unfuckPath())
	srv.Pre(middleware.RemoveTrailingSlash())

	srv.GET("/", sendResponse)
	srv.GET("/:message", sendResponse)
	srv.GET("/:message/:position", sendResponse)
	srv.GET("/:message/:position/:picture", sendResponse)

	srv.Run(standard.New(":8080"))
}
