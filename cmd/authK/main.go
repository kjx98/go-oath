package main

import (
	"flag"
	"fmt"
	"github.com/kjx98/go-ncurses"
	"github.com/kjx98/go-oath"
	"github.com/op/go-logging"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	MaxLines = 10
)

var logg = logging.MustGetLogger("authK")

func printOtp(acct *account) []byte {
	var ss string
	// display name ... OTP
	ss = fmt.Sprintf("%-20.20s %-10.10s %06d", acct.Name, acct.Issuer,
		acct.otp.Now())
	return []byte(ss)
}

func main() {
	var acctJson string
	var bOverWrite bool
	// $DJI, $COMPX
	flag.StringVar(&acctJson, "acct", "", "Accounts in JSON format")
	flag.BoolVar(&bOverWrite, "force", false, "force overwrite acct db")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: authK [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	accts := readAcct(acctJson, bOverWrite)
	nSymbols := len(accts)
	if nSymbols > MaxLines {
		nSymbols = MaxLines
	}
	logg.Infof("Start init ncurses for %d accts", nSymbols)
	for _, acct := range accts {
		logg.Infof("%v %v %06d\n", acct.Name, acct.Issuer,
			acct.otp.Now())
	}
	//os.Exit(0)
	running := true
	w, _ := ncurses.Initscr()
	defer ncurses.Endwin()
	defer func() {
		if r := recover(); r != nil {
			ncurses.Endwin()
			fmt.Printf("panic:\n%s\n", r)
			os.Exit(-1)
		}
	}()
	logg.Info("Enable color ncurses")
	// Enable color mode
	ncurses.StartColor()

	// Define color pairs
	ncurses.AddColorPair("bw", ncurses.ColorGreen, ncurses.ColorBlack)
	ncurses.AddColorPair("wb", ncurses.ColorBlue, ncurses.ColorBlack)
	ncurses.AddColorPair("rb", ncurses.ColorRed, ncurses.ColorWhite)
	ncurses.AddColorPair("yb", ncurses.ColorYellow, ncurses.ColorBlack)

	// Set cursor visiblity to hidden
	ncurses.SetCursor(ncurses.CURSOR_HIDDEN)
	// Automatically refresh after each command
	w.AutoRefresh = true
	// Set color for stdscr-window to system defaults.
	w.Wbkgd("std")
	w.SetColor("wb")
	{
		sTitle := "OATH-OTP authenticator"
		w.Move(0, 10)
		w.Write([]byte(sTitle))
	}
	nRows := uint16(nSymbols)
	var w1 *ncurses.Window
	if ww, err := ncurses.NewWindow("otp", ncurses.Position{0, 1},
		ncurses.Size{80, nRows}); err == nil {
		w1 = ww
		w1.AutoRefresh = true
		w1.SetScrolling(true)
		nRows++

		if w2, err := ncurses.NewWindow("log", ncurses.Position{0, nRows},
			ncurses.Size{80, 20}); err == nil {
			w2.AutoRefresh = true
			w2.SetScrolling(true)
			logInit(w2)
		}
	} else {
		w1 = w
	}
	// catch  SIGTERM, SIGINT, SIGUP
	{
		c := make(chan os.Signal, 10)
		signal.Notify(c)
		go func() {
			for s := range c {
				switch s {
				case os.Kill, os.Interrupt, syscall.SIGTERM:
					logg.Info("退出", s)
					running = false
					ExitFunc()
				case syscall.SIGQUIT:
					logg.Info("Quit", s)
					running = false
					ExitFunc()
				default:
					logg.Info("Got signal", s)
				}
			}
		}()
	}
	ts := time.Now()
	bFirst := true
	for running {
		tNow := time.Now()
		tNowS := tNow.Format("01-02 15:04:05")
		tRun := tNow.Sub(ts)
		w.SetColor("yb")
		w.Move(0, 40)
		w.Write([]byte(tNowS))
		w.Move(0, 56)
		fmt.Fprintf(w, "UpTime: %02d:%02d:%02d", int(tRun.Hours()),
			int(tRun.Minutes())%60, int(tRun.Seconds())%60)
		if cha := w.Getch(); cha == 'q' || cha == 'Q' || cha == 'x' {
			break
		}
		if nowS := tNow.Unix(); !bFirst && nowS%oath.Interval != 0 {
			time.Sleep(time.Millisecond * 200)
			runtime.Gosched()
			continue
		}
		for i := 0; i < nSymbols; i++ {
			w1.Move(uint16(i), 0)
			//w1.SetColor("rb")
			w1.SetColor("bw")
			// Blue
			//w1.SetColor("wb")
			w1.Write(printOtp(&accts[i]))
		}
		if bFirst {
			bFirst = false
		} else {
			time.Sleep(time.Millisecond * 200)
		}
		runtime.Gosched()
	}
	logInit(os.Stderr)
}

func ExitFunc() {
	//beeep.Alert("orSync Quit!", "try to exit", "")
	logg.Warning("开始退出...")
	logg.Warning("执行退出...")
	logg.Warning("结束退出...")
	// wait 3 seconds
	//time.Sleep(time.Second * 3)
	//os.Exit(1)
}

// `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
func logInit(w io.Writer) {
	var format = logging.MustStringFormatter(
		`%{color}%{time:01-02 15:04:05}  ▶ %{level:.4s} %{color:reset} %{message}`,
	)

	if w != os.Stderr {
		format = logging.MustStringFormatter(
			`%{time:01-02 15:04:05}  ▶ %{level:.4s} %{message}`,
		)
	}
	logback := logging.NewLogBackend(w, "", 0)
	logfmt := logging.NewBackendFormatter(logback, format)
	b := logging.SetBackend(logfmt)
	b.SetLevel(logging.INFO, "")
}

func init() {
	logInit(os.Stderr)
}
