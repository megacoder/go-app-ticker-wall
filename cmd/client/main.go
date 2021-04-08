package main

import (
	"context"
	"time"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	tickerManager "github.com/polygon-io/go-app-ticker-wall/ticker_manager"
	"github.com/polygon-io/nanovgo"
	"github.com/polygon-io/nanovgo/perfgraph"
	"github.com/sirupsen/logrus"

	"github.com/kelseyhightower/envconfig"
	tombv2 "gopkg.in/tomb.v2"
)

const maxMessageSize = 1024 * 1024 * 1 // 1MB

var (
	// AnimationDuration is the length of animations.
	AnimationDuration = 750 // ms
)

var cfg ServiceConfig

type ServiceConfig struct {
	// Service details
	LogLevel string `split_words:"true" default:"DEBUG"`
	Leader   string `split_words:"true" default:"localhost:6886"`

	// Local Presentation Settings:
	ScreenWidth  int `split_words:"true" default:"1200"`
	ScreenHeight int `split_words:"true" default:"300"`
	ScreenIndex  int `split_words:"true" default:"2"`
}

func run() error {
	// Global top level context.
	tomb, ctx := tombv2.WithContext(context.Background())

	// Parse Env Vars:
	err := envconfig.Process("", &cfg)
	if err != nil {
		return err
	}

	// Set Log Levels
	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.WithField("err", err).Warn("parse log level")
	} else {
		logrus.SetLevel(l)
	}

	if err := glfw.Init(gl.ContextWatcher); err != nil {
		return err
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(cfg.ScreenWidth, cfg.ScreenHeight, "Polygon Ticker Wall", nil, nil)
	if err != nil {
		return err
	}
	window.MakeContextCurrent()

	// ctx, err := nanovgo.NewContext(0)
	nanoCtx, err := nanovgo.NewContext(0)
	defer nanoCtx.Delete()
	if err != nil {
		panic(err)
	}

	glfw.SwapInterval(0)
	createFonts(nanoCtx)

	// Ticker Manager. By default we believe we are the only one. Once we connect to leader we will get updated info.
	mgr := tickerManager.NewDefaultManager(&tickerManager.PresentationData{
		ScreenWidth:        cfg.ScreenWidth,
		ScreenHeight:       cfg.ScreenHeight,
		ScreenGlobalOffset: 0,
		TickerBoxWidth:     cfg.ScreenWidth,
		ScreenIndex:        cfg.ScreenIndex,
		NumberOfScreens:    1,
		GlobalViewportSize: int64(cfg.ScreenWidth),
	})

	// Ticker wall client.
	tickerWallClient := NewTickerWallClient(&cfg, mgr)
	defer tickerWallClient.Close()

	// Create GRPC connection for the client.
	if err := tickerWallClient.CreateGRPCClient(); err != nil {
		return err
	}

	// Load initial tickers
	if err := tickerWallClient.LoadTickers(ctx); err != nil {
		return err
	}

	// tomb will context the context
	tomb.Go(func() error {
		return tickerWallClient.Run(ctx)
	})

	return createRenderingLoop(ctx, nanoCtx, window, mgr)
}

func createRenderingLoop(ctx context.Context, nanoCtx *nanovgo.Context, window *glfw.Window, mgr tickerManager.TickerManager) error {
	fps := perfgraph.NewPerfGraph("Frame Time", "sans")
	fbWidth, fbHeight := window.GetFramebufferSize()
	winWidth, winHeight := window.GetSize()
	pixelRatio := float32(fbWidth) / float32(winWidth)
	gl.Viewport(0, 0, fbWidth, fbHeight)

	nanoCtx.SetFontFace("sans")
	nanoCtx.SetTextAlign(nanovgo.AlignLeft | nanovgo.AlignTop)
	nanoCtx.SetTextLineHeight(1.2)

	specialMessage := true
	startTimer := time.Now().Add(1 * time.Minute)
	startTimer = startTimer.Truncate(time.Minute)
	specialMessageTimeActivate := startTimer.UnixNano() / int64(time.Millisecond)
	logrus.Info("activation time: ", specialMessageTimeActivate)

	for !window.ShouldClose() {
		fps.UpdateGraph()
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.Enable(gl.CULL_FACE)
		gl.Disable(gl.DEPTH_TEST)
		nanoCtx.BeginFrame(winWidth, winHeight, pixelRatio)
		// nanoCtx.Save()

		t := time.Now().UnixNano() / int64(time.Millisecond*10)
		// println(t)
		// Actual application drawing.
		renderTickers(nanoCtx, mgr, t)

		if specialMessage {
			renderSpecialMessage(nanoCtx, mgr, t, "Very Important Special Message... Read it!", int(specialMessageTimeActivate), 5000)
		}

		// nanoCtx.Restore()
		fps.RenderGraph(nanoCtx, -50, -50)
		nanoCtx.EndFrame()
		gl.Enable(gl.DEPTH_TEST)
		window.SwapBuffers()
		glfw.PollEvents()
		// time.Sleep(time.Millisecond * 16)
	}

	return ctx.Err()
}

func createFonts(ctx *nanovgo.Context) {
	ctx.CreateFont("sans", "fonts/Roboto-Regular.ttf")
	ctx.CreateFont("sans-light", "fonts/Roboto-Light.ttf")
	ctx.CreateFont("sans-bold", "fonts/Roboto-Bold.ttf")
}

func main() {
	if err := run(); err != nil {
		logrus.WithError(err).Error("Program exiting")
	}
}
