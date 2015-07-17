package log

import (
    "bytes"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
    var buf bytes.Buffer
    var writer = New(&buf, "", 0)
    writer.SetCurrLineVisible(false)
    writer.Print("Hello ")
    writer.Print("Dan, ")
    writer.Print("how are you?\n")
    assert.Equal(t, "Hello Dan, how are you?\n", buf.String(), "Print should not add newlines")
    buf.Reset()
    writer.Print("You\n feel\n  like\n   you're\n    falling...")
    assert.Equal(t, "You\n feel\n  like\n   you're\n", buf.String(), "Print should only output full lines until Closed")
    buf.Reset()
    writer.Close()
    assert.Equal(t, "    falling...\n", buf.String(), "Close should flush any unfinished line by appending a terminating newline")
    buf.Reset()
    writer.Print("")
    writer.Close()
    assert.Equal(t, "", buf.String(), "Print and Close should output nothing if nothing was printed")
}

func TestTempLines(t *testing.T) {
    var buf bytes.Buffer
    var writer = New(&buf, "", 0)
    writer.Print("Hello ")
    assert.Equal(t, "Hello ", buf.String())
    buf.Reset()
    writer.Print(" Dan, ")
    assert.Equal(t, "Dan, ", buf.String())
    buf.Reset()
    writer.Print("how are you?\n")
    assert.Equal(t, "how are you?\n", buf.String())
}

func TestMultipleTempLines(t *testing.T) {
    var buf bytes.Buffer
    var writer1 = New(&buf, "", 0)
    var writer2 = New(&buf, "", 0)
    writer1.Print("Testing...")
    assert.Equal(t, "Testing... ", buf.String())
    buf.Reset()
    writer2.Print("Writing Code...")
    assert.Equal(t, "\rTesting... | Writing Code...", buf.String())
    buf.Reset()
    writer2.Print(" done.\n")
    assert.Equal(t, "\rWriting Code... done        \nTesting...", buf.String())
    buf.Reset()
    writer1.Print(" done.\n")
    assert.Equal(t, " done.\n", buf.String())
}
