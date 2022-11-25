package v1

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/mr-chelyshkin/go-gridder/internal/sys"

	"github.com/bbengfort/x/pid"
	"golang.org/x/sync/errgroup"
)

const (
	exitCodeOk    = 0
	exitCodeErr   = 1
	exitCodeAbort = 134

	defaultShutdownTimeout = 30
	defaultPprofURL        = "0.0.0.0:6060"
	defaultPidPath         = "/var/run/"
)

type daemon struct {
	cfg config

	s service
	l logger

	ctx context.Context
	wg  *errgroup.Group
}

// NewDaemon return daemon object.
func NewDaemon(
	ctx context.Context,
	cfg config,
	s service,
	l logger,
) (*daemon, error) {
	if cfg == nil {
		return nil, fmt.Errorf("incoming config object is nil")
	}
	if s == nil {
		return nil, fmt.Errorf("incoming service object is nil")
	}
	if l == nil {
		return nil, fmt.Errorf("incoming logger object is nil")
	}

	if cfg.GetPidFilePath() == "" {
		cfg.SetPidFilePath(filepath.Join(defaultPidPath, fmt.Sprintf("%s.pid", s.GetName())))
	}
	if cfg.GetShutdownTimeoutSec() == 0 {
		cfg.SetShutdownTimeoutSec(defaultShutdownTimeout)
	}

	wg, ctx := errgroup.WithContext(ctx)
	return &daemon{
		cfg: cfg,

		s: s,
		l: l,

		ctx: ctx,
		wg:  wg,
	}, nil
}

// Run daemon.
func (d *daemon) Run() {
	runtime.GOMAXPROCS(d.cfg.GetMaxPoc())
	debug.SetGCPercent(d.cfg.GetGCPercent())

	p := pid.New(d.cfg.GetPidFilePath())
	if err := p.Save(); err != nil {
		d.l.Fatalf("[daemon] initPidFileErr", err.Error())
	}

	d.l.Infof(
		"starting daemon",
		map[string]string{
			"service": reflect.TypeOf(d.s).String(),
			"logger":  reflect.TypeOf(d.l).String(),
		},
	)
	d.l.Infof(
		"runtime info",
		map[string]string{
			"GOMAXPROCS":  strconv.Itoa(d.cfg.GetMaxPoc()),
			"GOGCPERCENT": strconv.Itoa(d.cfg.GetGCPercent()),
			"pid":         strconv.Itoa(p.PID),
			"path":        d.cfg.GetPidFilePath(),
		},
	)

	if d.cfg.GetPprof() {
		d.l.Infof("starting pprof", map[string]string{"url": "http://" + defaultPprofURL})
		d.wg.Go(func() error {
			_ = http.ListenAndServe(defaultPprofURL, nil)
			return nil
		})
	}

	exit := make(chan interface{})
	d.startService(exit)
	go func() {
		<-d.ctx.Done()
		_ = p.Free()

		d.l.Infof("waiting for application shutdown", map[string]string{"reason": d.ctx.Err().Error()})
		ticker := time.NewTicker(sys.ToTimeDuration(d.cfg.GetShutdownTimeoutSec()))
		for {
			select {
			case <-exit:
				d.l.Infof("daemon stopped", nil)
				os.Exit(exitCodeOk)
			case <-ticker.C:
				d.l.Warnf(
					"graceful shutdown watchdog triggered: forcing shutdown",
					map[string]string{
						"reason": "waiting too long, timeout: " + strconv.Itoa(d.cfg.GetShutdownTimeoutSec()) + "sec",
					},
				)
				os.Exit(exitCodeAbort)
			}
		}
	}()
	if err := d.wg.Wait(); err != nil {
		_ = p.Free()
		d.l.Errorf("application returned error", map[string]string{"err": err.Error()})
		os.Exit(exitCodeErr)
	}
}

func (d *daemon) startService(exit chan interface{}) {
	d.wg.Go(func() error {
		d.l.Infof("starting service", map[string]string{"name": d.s.GetName()})
		return d.s.Start(d.ctx)
	})
	d.wg.Go(func() error {
		<-d.ctx.Done()
		d.l.Infof("shutting down service", nil)
		defer func() { exit <- nil }()
		return d.s.Shutdown()
	})
}
