package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "", 0)
	writer.HidePartialLines()
	writer.Print("Hello ")
	writer.Print("Dan, ")
	writer.Print("how are you?\n")
	assert.Equal("Hello Dan, how are you?\n", buf.String(), "Print should not add newlines")
	buf.Reset()
	writer.Print("You\n feel\n  like\n   you're\n    falling...")
	assert.Equal("You\n feel\n  like\n   you're\n", buf.String(), "Print should only output full lines until Closed")
	buf.Reset()
	writer.Close()
	assert.Equal("    falling...\n", buf.String(), "Close should flush any unfinished line by appending a terminating newline")
}

func TestPrintNothing(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	writer := New(&buf, "", 0)
	writer.Print("")
	writer.Close()
	assert.Equal("", buf.String(), "Print and Close should output nothing if nothing was printed")
}

func TestTempLines(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "", 0)
	defer writer.Close()
	writer.Print("Hello ")
	assert.Equal("Hello ", buf.String())
	buf.Reset()
	writer.Print("Dan, ")
	assert.Equal("Dan, ", buf.String())
	buf.Reset()
	writer.Print("how are you?\n")
	assert.Equal("how are you?\n", buf.String())
}

func TestMultipleTempLines(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer1 = New(&buf, "", 0)
	var writer2 = New(&buf, "", 0)
	defer writer1.Close()
	defer writer2.Close()
	writer1.Print("Testing...")
	assert.Equal("Testing...", buf.String())
	buf.Reset()
	writer2.Print("Writing Code...")
	assert.Equal(" | Writing Code...", buf.String())
	buf.Reset()
	writer2.Print(" done.\n")
	assert.Equal("\rWriting Code... done.       \nTesting...", buf.String())
	buf.Reset()
	writer1.Print(" done.\n")
	assert.Equal(" done.\n", buf.String())
	buf.Reset()
	writer2.Print("Writing More Code...")
	assert.Equal("Writing More Code...", buf.String())
	buf.Reset()
	writer1.Print("Testing More...")
	assert.Equal(" | Testing More...", buf.String())
}

func TestMultipleTempLinesDiffWriters(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	var writer1 = New(&buf, "", 0)
	var writer2 = New(&buf2, "", 0)
	defer writer1.Close()
	defer writer2.Close()
	writer1.Print("Testing...")
	assert.Equal("Testing...", buf.String())
	buf.Reset()
	writer2.Print("Writing Code...")
	assert.Equal("Writing Code...", buf2.String())
	buf2.Reset()
	writer1.Print(" done.\n")
	assert.Equal(" done.\n", buf.String())
	assert.Equal("", buf2.String())
	buf.Reset()
	writer2.Print(" done.\n")
	assert.Equal("", buf.String())
	assert.Equal(" done.\n", buf2.String())
	buf2.Reset()
}

func TestAnsiColors(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "", 0)
	defer writer.Close()
	writer.Print("Here is @(red:some red text).\n")
	assert.Equal("Here is @(red:some red text).\n", buf.String())
	buf.Reset()
	writer.EnableColorTemplate()
	writer.Printf("Here is @(red:some red text).\n")
	assert.Equal("Here is \033[31msome red text\033[39m.\n", buf.String())
	buf.Reset()
	writer.Printf("Here is @(dim:some dim text).\n")
	assert.Equal("Here is \033[2msome dim text\033[0m.\n", buf.String())
	buf.Reset()
	writer.Printf("Here is some @(green)green text@(r) and @(garbage) and @(cyan:cyan text).\n")
	assert.Equal("Here is some \033[32mgreen text\033[0m and @(garbage) and \033[36mcyan text\033[39m.\n", buf.String())
	buf.Reset()
	writer.DisableColorTemplate()
	writer.Printf("Here is @(red:some red text).\n")
	assert.Equal("Here is @(red:some red text).\n", buf.String())
	buf.Reset()
}

func TestAnsiSpanningLines(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "\033[32m$$ ", 0)
	defer writer.Close()
	writer.Print("Hello, ")
	assert.Equal("\033[32m$$ \033[39mHello, ", buf.String(), "we auto-reset ansi colors after the prefix")
	buf.Reset()
	writer.Print("we're writing in red\033[31m")
	assert.Equal("we're writing in red\033[31m", buf.String())
	buf.Reset()
	writer.Print(" in red")
	assert.Equal(" in red", buf.String())
	buf.Reset()
	writer.Print("but\nnot when we're on a new line.\033[39m\n")
	assert.Equal("but\033[39m\n\033[32m$$ \033[39mnot when we're on a new line.\033[39m\n", buf.String())
	buf.Reset()
	writer.EnableColorTemplate()
	writer.Printf("@(blue:but not\ntemplated\nnewlines).\n")
	assert.Equal("\033[32m$$ \033[39m\033[34mbut not\033[39m\n\033[32m$$ \033[39mtemplated\n\033[32m$$ \033[39mnewlines\033[39m.\n", buf.String())
	buf.Reset()
}

func TestDisableColor(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "\033[32m$$ ", 0)
	defer writer.Close()
	var input = "\033[31mI \033[32mlike \033[33mcolors\033[39m\n"
	var withEscapes = "\033[32m$$ \033[39m\033[31mI \033[32mlike \033[33mcolors\033[39m\n"
	var withoutEscapes = "$$ I like colors\n"
	// Note: behavior is undefined when enabling/disabling color in the middle of partial lines
	writer.Print(input)
	assert.Equal(withEscapes, buf.String())
	buf.Reset()
	DisableColor()
	writer.Print(input)
	assert.Equal(withoutEscapes, buf.String())
	buf.Reset()
	writer.EnableColor()
	writer.Print(input)
	assert.Equal(withEscapes, buf.String(), "Enabling color on a specific Logger overrides global setting")
	buf.Reset()
	writer.DisableColor()
	writer.Print(input)
	assert.Equal(withoutEscapes, buf.String())
	buf.Reset()
	EnableColor()
	writer.Print(input)
	assert.Equal(withoutEscapes, buf.String(), "Disabling color on a specific Logger overrides global setting")
	buf.Reset()
}

func TestAddColorCode(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	AddAnsiColorCode("awesome", 1)
	AddAnsiColorCode("sauce", 36)
	var writer = New(&buf, "@(awesome,sauce:$$) ", 0)
	defer writer.Close()
	writer.EnableColorTemplate()
	writer.Printf("@(sauce,awesome)text@(r)\n")
	assert.Equal("\033[1m\033[36m$$\033[0m \033[36m\033[1mtext\033[0m\n", buf.String())
	buf.Reset()
	writer.Printf("@(awesome)this is all bright @(sauce:text), even this.@(r)\n")
	assert.Equal("\033[1m\033[36m$$\033[0m \033[1mthis is all bright \033[36mtext\033[39m, even this.\033[0m\n", buf.String())
	buf.Reset()
}

// non-english example text drawn mostly from http://www.columbia.edu/~fdc/utf8/

func TestTermWidthTruncation(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer1 = New(&buf, "@(green)$$ ", 0)
	defer writer1.Close()
	writer1.EnableColorTemplate()
	var writer2 = New(&buf, "@(red)$$ ", 0)
	defer writer2.Close()
	writer2.EnableColorTemplate()
	writer1.SetTerminalWidth(30) // Applies to both because they both write to buf
	writer1.Printf("@(yellow)𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸")
	assert.Equal("\033[32m$$ \033[39m\033[33m𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸", buf.String())
	buf.Reset()
	writer2.Printf("@(blue)ᚠᛇᚻᛒᛦᚦᚠᚱᚩᚠ")
	assert.Equal(" | \033[31m$$ \033[39m\033[34mᚠᛇᚻᛒᛦᚦᚠᚱᚩᚠ", buf.String())
	buf.Reset()
	writer1.Printf("1234567890σπασμένα1234567890")
	assert.Contains(buf.String(), "ᚠᛇᚻᛒᛦ", "We should try to show a little of each partial line if possible")
	buf.Reset()
}

func TestNonLatinRunes(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "我能吞下玻璃而不伤身体。", 0)
	defer writer.Close()
	writer.Print("أنا قادر على أكل الزجاج و هذا")
	assert.Equal("我能吞下玻璃而不伤身体。أنا قادر على أكل الزجاج و هذا", buf.String())
	buf.Reset()
	writer.Print(" لا يؤلمني.\n")
	assert.Equal(" لا يؤلمني.\n", buf.String())
	buf.Reset()
	writer.SetTerminalWidth(20)
	// This has a combining diacritic after/in the third character.
	writer.Print("ನನಗೆ ಹಾನಿ ಆಗದೆ, ನಾನು ಗಜನ್ನು ತಿನಬಹುದು")
	assert.Equal("我能吞下玻璃而不伤身体。ನನಗೆ...", buf.String())
}

func TestApplyTemplateEarly(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "", 0)
	defer writer.Close()
	writer.EnableColorTemplate()
	writer.Printf("@(red:%s)\n", "@(green:this is not green)")
	assert.Equal("\033[31m@(green:this is not green)\033[39m\n", buf.String())
	buf.Reset()
	SetOutput(&buf)
	SetFlags(0)
	EnableColorTemplate()
	defer SetOutput(os.Stderr)
	defer SetFlags(LstdFlags)
	defer DisableColorTemplate()
	Printf("@(red:%s)\n", "@(green:this is not green)")
	assert.Equal("\033[31m@(green:this is not green)\033[39m\n", buf.String())
	buf.Reset()
}

func TestCarriageReturns(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer1 = New(&buf, "", 0)
	defer writer1.Close()
	writer1.Print("Working...")
	assert.Equal("Working...", buf.String())
	buf.Reset()
	var writer2 = New(&buf, "", 0)
	defer writer2.Close()
	writer2.EnableColorTemplate()
	writer2.Printf("Progress: @(red:  0) percent.")
	assert.Equal(" | Progress: \033[31m  0\033[39m percent.", buf.String())
	buf.Reset()
	writer2.Printf("\rProgress: @(red:  1) percent.")
	assert.Equal("\rWorking... | Progress: \033[31m  1\033[39m percent.", buf.String())
	buf.Reset()
	writer2.Printf("\rProgress: @(red:  2) percent.\r")
	assert.Equal("\rWorking... | Progress: \033[31m  2\033[39m percent.", buf.String())
	buf.Reset()
	writer2.Printf("Progress: @(red: 33) percent.")
	assert.Equal("\rWorking... | Progress: \033[31m 33\033[39m percent.", buf.String())
	buf.Reset()
	writer2.Printf("\rProgress: @(blue) 6")
	assert.Equal("\rWorking... | Progress: \033[34m 6\033[39m\033[31m3\033[39m percent.", buf.String())
	buf.Reset()
}

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "", 0)
	defer writer.Close()
	writer.Replace("Hello Susan.")
	assert.Equal("Hello Susan.", buf.String())
	buf.Reset()
	writer.Replace("Hello Bob.")
	assert.Equal("\rHello Bob.  ", buf.String())
	buf.Reset()
	writer.Replacef("Hello %s.", "Al")
	assert.Equal("\rHello Al. ", buf.String())
	buf.Reset()
	writer.Replacef("Hello %s", "Ala")
	assert.Equal("\rHello Ala", buf.String())
	buf.Reset()
	writer.Replace("Hello Alan.")
	assert.Equal("n.", buf.String())
	buf.Reset()
}

func TestMultilineMode(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	writer1 := New(&buf, "", 0)
	lineUp := tput("cuu", "1")
	lineDown := tput("cud", "1")
	readBuf := func() string {
		s := buf.String()
		buf.Reset()
		s = strings.Replace(s, lineDown, "{DOWN}", -1)
		s = strings.Replace(s, lineUp, "{UP}", -1)
		return s
	}
	writer1.EnableMultilineMode()
	writer1.Print("writer1...")
	assert.Equal("writer1...", readBuf())
	writer2 := New(&buf, "", 0)
	writer2.Print("writer2...")
	assert.Equal("\nwriter2...", readBuf())
	writer1.Print(" working...")
	assert.Equal("{UP}\rwriter1... working...", readBuf())
	writer1.Print("  50 percent finished...")
	assert.Equal("  50 percent finished...", readBuf())
	writer2.Print(" working... ")
	assert.Equal("{DOWN}\rwriter2... working... ", readBuf())
	writer2.Print("done.\n")
	// Need to move up to the previous line, overwrite writer1's text, then only move down a line.
	// A newline is not necessary since we're only *completing* an existing line and not yet starting
	// a new line.
	assert.Equal("{UP}\rwriter2... working... done.                  {DOWN}\rwriter1... working...  50 percent finished...", readBuf())
	writer2.Print("working again...")
	assert.Equal("\nworking again...", readBuf())
	writer1.Print("\rwriter1... working... 100 percent. done.     \n")
	// Again, we don't append a newline, but this time, we move to the last line even though we're not
	// going to write anything else immediately.
	assert.Equal("{UP}\rwriter1... working... 100 percent. done.     {DOWN}\r", readBuf())
	writer1.Print("Hello")
	assert.Equal("\nHello", readBuf())
	writer1.Close()
	assert.Equal("{UP}\rHello           {DOWN}\rworking again...", readBuf())
	writer2.Close()
	assert.Equal("\n", readBuf())
}

func TestAutoNewlines(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer1 = New(&buf, "", 0)
	defer writer1.Close()
	var writer2 = New(&buf, "", 0)
	defer writer2.Close()
	writer1.Print("this is ")
	assert.Equal("this is ", buf.String())
	buf.Reset()
	EnableAutoNewlines()
	assert.Equal("", buf.String(), "AutoNewlines does not take effect immediately")
	buf.Reset()
	writer1.Print("this is a partial line.")
	assert.Equal("this is a partial line.\n", buf.String())
	buf.Reset()
	writer2.Print("this is a partial line.")
	assert.Equal("this is a partial line.\n", buf.String())
	buf.Reset()
	DisableAutoNewlines()
	writer1.Print("this is ")
	assert.Equal("this is ", buf.String())
	buf.Reset()
	writer1.EnableAutoNewlines()
	assert.Equal("", buf.String(), "AutoNewlines does not take effect immediately")
	buf.Reset()
	writer1.Print("this is a partial line.")
	assert.Equal("this is a partial line.\n", buf.String())
	buf.Reset()
	writer2.Print("this is a partial line.")
	assert.Equal("this is a partial line.", buf.String())
	buf.Reset()
}

func TestFlagElapsed(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer = New(&buf, "$$ ", Lelapsed)
	defer writer.Close()
	writer.Print("Testing... ")
	assert.Equal("$$ Testing... ", buf.String())
	buf.Reset()
	writer.Print("done.\n")
	assert.Equal("\r$$ (", buf.String()[:5])
	buf.Reset()
}

func TestFormatDuration(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("0.0ms", string(formatDuration(0*time.Microsecond)))
	assert.Equal("0.1ms", string(formatDuration(50*time.Microsecond)))
	assert.Equal("0.1ms", string(formatDuration(100*time.Microsecond)))
	assert.Equal("0.5ms", string(formatDuration(500*time.Microsecond)))
	assert.Equal("1.0ms", string(formatDuration(1000*time.Microsecond)))
	assert.Equal("9.9ms", string(formatDuration(9900*time.Microsecond)))
	assert.Equal(" 10ms", string(formatDuration(10000*time.Microsecond)))
	assert.Equal(" 99ms", string(formatDuration(99000*time.Microsecond)))
	assert.Equal("100ms", string(formatDuration(99500*time.Microsecond)))
	assert.Equal("100ms", string(formatDuration(100000*time.Microsecond)))
	assert.Equal("999ms", string(formatDuration(999000*time.Microsecond)))
	assert.Equal("1.00s", string(formatDuration(999500*time.Microsecond)))
	assert.Equal("1.00s", string(formatDuration(1000*time.Millisecond)))
	assert.Equal("9.99s", string(formatDuration(9990*time.Millisecond)))
	assert.Equal("10.0s", string(formatDuration(10000*time.Millisecond)))
	assert.Equal("99.9s", string(formatDuration(99900*time.Millisecond)))
	assert.Equal(" 100s", string(formatDuration(100*time.Second)))
	assert.Equal("10.0m", string(formatDuration(600*time.Second)))
	assert.Equal("99.9m", string(formatDuration(5994*time.Second)))
	assert.Equal(" 100m", string(formatDuration(5997*time.Second)))
	assert.Equal(" 100m", string(formatDuration(6000*time.Second)))
	assert.Equal("10.0h", string(formatDuration(600*time.Minute)))
	assert.Equal("99.9h", string(formatDuration(5994*time.Minute)))
	assert.Equal(" 100h", string(formatDuration(5997*time.Minute)))
	assert.Equal(" 100h", string(formatDuration(6000*time.Minute)))
	assert.Equal(" 100h", string(formatDuration(6000*time.Minute)))
	assert.Equal("9999h", string(formatDuration(9999*time.Hour)))
	assert.Equal("99999h", string(formatDuration(99999*time.Hour)))
	assert.Equal("999999h", string(formatDuration(999999*time.Hour)))
}

func TestLoggerInception(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	var writer1 = New(&buf, "", 0)
	defer writer1.Close()
	var writer2 = New(writer1, "prefix: ", 0)
	defer writer2.Close()
	writer2.Print("hello\n")
	assert.Equal("prefix: hello\n", buf.String())
	buf.Reset()
}

// XXX To make this really work, we'd need to stub out time.Now() in log.go.
// func TestPrefix(t *testing.T) {
// 	assert := assert.New(t)
// 	testEquivalence := func(template string, flags int) {
// 		var buf1 bytes.Buffer
// 		writer1 := New(&buf1, template, 0)
// 		defer writer1.Close()
// 		var buf2 bytes.Buffer
// 		writer2 := New(&buf2, "", flags)
// 		defer writer2.Close()
// 		readBuf1 := func() string {
// 			defer buf1.Reset()
// 			return buf1.String()
// 		}
// 		readBuf2 := func() string {
// 			defer buf2.Reset()
// 			return buf2.String()
// 		}
// 		writer1.Printf("Hi\nHello")
// 		writer2.Printf("Hi\nHello")
// 		assert.Equal(readBuf1(), readBuf2())
// 	}
// 	testEquivalence("{date} {time} ", Ldate|Ltime)
// 	testEquivalence("{date} {time micros} ", Ldate|Ltime|Lmicroseconds)
// 	testEquivalence("{isodate} ", Lisodate)
// 	testEquivalence("{isodate micros} ", Lisodate|Lmicroseconds)
// }

// TODO test &/or implement:
// - Set custom ANSI template regexp specifically or globally
// - Handle \b and \t characters intelligently
