package log

import (
    "bytes"
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
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
    writer.Print("Here is @[red:some red text].\n")
    assert.Equal("Here is @[red:some red text].\n", buf.String())
    buf.Reset()
    writer.EnableColorTemplate()
    writer.Print("Here is @[red:some red text].\n")
    assert.Equal("Here is \033[31msome red text\033[39m.\n", buf.String())
    buf.Reset()
    writer.Print("Here is @[dim:some dim text].\n")
    assert.Equal("Here is \033[2msome dim text\033[0m.\n", buf.String())
    buf.Reset()
    writer.Print("Here is some @[green]green text@[r] and @[garbage] and @[cyan:cyan text].\n")
    assert.Equal("Here is some \033[32mgreen text\033[0m and @[garbage] and \033[36mcyan text\033[39m.\n", buf.String())
    buf.Reset()
    writer.DisableColorTemplate()
    writer.Print("Here is @[red:some red text].\n")
    assert.Equal("Here is @[red:some red text].\n", buf.String())
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
    writer.Print("we're writing\033[31m")
    assert.Equal("we're writing\033[31m", buf.String())
    buf.Reset()
    writer.Print(" in red")
    assert.Equal(" in red", buf.String())
    buf.Reset()
    writer.Print("even\nwhen we're on a new line.\033[39m\n")
    assert.Equal("even\033[39m\n\033[32m$$ \033[39m\033[31mwhen we're on a new line.\033[39m\n", buf.String())
    buf.Reset()
    writer.EnableColorTemplate()
    writer.Print("@[blue:templated\nnewlines\ntoo].\n")
    assert.Equal("\033[32m$$ \033[39m\033[34mtemplated\033[39m\n\033[32m$$ \033[39m\033[34mnewlines\033[39m\n\033[32m$$ \033[39m\033[34mtoo\033[39m.\n", buf.String())
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
    AddAnsiCode("awesome", 1)
    AddAnsiCode("sauce", 36)
    var writer = New(&buf, "@[awesome,sauce:$$] ", 0)
    defer writer.Close()
    writer.EnableColorTemplate()
    writer.Print("@[sauce,awesome]text@[r]\n")
    assert.Equal("\033[1m\033[36m$$\033[0m \033[36m\033[1mtext\033[0m\n", buf.String())
    buf.Reset()
    writer.Print("@[awesome]this is all bright @[sauce:text], even this.@[r]\n")
    assert.Equal("\033[1m\033[36m$$\033[0m \033[1mthis is all bright \033[36mtext\033[39m, even this.\033[0m\n", buf.String())
    buf.Reset()
}

// non-english example text drawn mostly from http://www.columbia.edu/~fdc/utf8/

func TestTermWidthTruncation(t *testing.T) {
    assert := assert.New(t)
    var buf bytes.Buffer
    var writer1 = New(&buf, "@[green]$$ ", 0)
    defer writer1.Close()
    writer1.EnableColorTemplate()
    var writer2 = New(&buf, "@[red]$$ ", 0)
    defer writer2.Close()
    writer2.EnableColorTemplate()
    writer1.SetTerminalWidth(30) // Applies to both because they both write to buf
    writer1.Print("@[yellow]𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸")
    assert.Equal("\033[32m$$ \033[39m\033[33m𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸𐌸", buf.String())
    buf.Reset()
    writer2.Print("@[blue]ᚠᛇᚻᛒᛦᚦᚠᚱᚩᚠ")
    assert.Equal(" | \033[31m$$ \033[39m\033[34mᚠᛇᚻᛒᛦᚦᚠᚱᚩᚠ", buf.String())
    buf.Reset()
    writer1.Print("1234567890σπασμένα1234567890")
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
    writer.Printf("@[red:%s]\n", "@[green:this is not green]")
    assert.Equal("\033[31m@[green:this is not green]\033[39m\n", buf.String())
    buf.Reset()
    SetOutput(&buf)
    SetFlags(0)
    EnableColorTemplate()
    defer SetOutput(os.Stderr)
    defer SetFlags(LstdFlags)
    defer DisableColorTemplate()
    Printf("@[red:%s]\n", "@[green:this is not green]")
    assert.Equal("\033[31m@[green:this is not green]\033[39m\n", buf.String())
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
    writer2.Print("Progress: @[red:  0] percent.")
    assert.Equal(" | Progress: \033[31m  0\033[39m percent.", buf.String())
    buf.Reset()
    writer2.Print("\rProgress: @[red:  1] percent.")
    assert.Equal("\rWorking... | Progress: \033[31m  1\033[39m percent.", buf.String())
    buf.Reset()
    writer2.Print("\rProgress: @[red:  2] percent.\r")
    assert.Equal("\rWorking... | Progress: \033[31m  2\033[39m percent.", buf.String())
    buf.Reset()
    writer2.Print("Progress: @[red: 33] percent.")
    assert.Equal("\rWorking... | Progress: \033[31m 33\033[39m percent.", buf.String())
    buf.Reset()
    writer2.Print("\rProgress: @[blue] 6")
    assert.Equal("\rWorking... | Progress: \033[34m 6\033[39m\033[31m3\033[39m percent.", buf.String())
    buf.Reset()
}

func TestMultilineMode(t *testing.T) {
    assert := assert.New(t)
    var buf bytes.Buffer
    var writer1 = New(&buf, "", 0)
    writer1.EnableMultilineMode()
    defer writer1.Close()
    writer1.Print("writer1...")
    assert.Equal("writer1...", buf.String())
    buf.Reset()
    var writer2 = New(&buf, "", 0)
    defer writer2.Close()
    writer2.Print("writer2...")
    assert.Equal("\nwriter2...", buf.String())
    buf.Reset()
    writer1.Print(" working...")
    assert.Equal("\033[1Fwriter1... working...", buf.String())
    buf.Reset()
    writer1.Print("  50 percent finished...")
    assert.Equal("  50 percent finished...", buf.String())
    buf.Reset()
    writer2.Print("working... ")
    assert.Equal("\033[1Ewriter2... working... ", buf.String())
    buf.Reset()
    writer2.Print("done.\n")
    // Need to move up to the previous line, overwrite writer1's text, then only move down a line.
    // A newline is not necessary since we're only *completing* an existing line and not yet starting
    // a new line.
    assert.Equal("\033[1Fwriter2... working... done.                   \033[1Ewriter1... working...  50 percent finished...", buf.String())
    buf.Reset()
    writer2.Print("working again...")
    assert.Equal("\nworking again...", buf.String())
    buf.Reset()
    writer1.Print("\rwriter1... working... 100 percent. done.     \n")
    // Again, we don't append a newline, but this time, we move to the last line even though we're not
    // going to write anything else immediately.
    assert.Equal("\033[1Fwriter1... working... 100 percent. done.     \033[1E", buf.String())
    buf.Reset()
}

// TODO test &/or implement:
// - Set custom ANSI template regexp specifically or globally
// - Add option to auto-append newlines with each Print/Printf for stock `log` compatibility
// - Add duration output flag? "(37.2 secs) Downloading stuff... done."
// - Experiment with multiple lines of temp output? Probably doesn't work.
// - Test & implement proper bookkeeping around SetOutput and tempLoggers
