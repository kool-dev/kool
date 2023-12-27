package utils

import (
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func MakeFastLoading(loadingMsg, loadedMsg string, w io.Writer) *spinner.Spinner {
	return makeLoading(loadingMsg, loadedMsg, 14, w)
}

func MakeSlowLoading(loadingMsg, loadedMsg string, w io.Writer) *spinner.Spinner {
	return makeLoading(loadingMsg, loadedMsg, 21, w)
}

func makeLoading(loadingMsg, loadedMsg string, charset int, w io.Writer) (s *spinner.Spinner) {
	s = spinner.New(
		spinner.CharSets[charset],
		100*time.Millisecond,
		// spinner.WithColor("green"),
		spinner.WithFinalMSG(loadedMsg+"\n"),
		spinner.WithSuffix(" "+loadingMsg),
	)
	s.Prefix = " "
	s.HideCursor = true
	if fh, isFh := w.(*os.File); isFh {
		s.WriterFile = fh
	} else {
		s.Writer = w
	}
	s.Start()
	return
}
