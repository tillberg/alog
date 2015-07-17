package log

import (
    "bytes"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
    var buf bytes.Buffer
    var me = New(&buf, "", 0)
    me.Print("Hello ")
    me.Print("Dan, ")
    me.Print("how are you?\n")
    assert.Equal(t, "Hello Dan, how are you?\n", buf.String(), "Print should not add newlines")
    buf.Reset()
    me.Print("You\n feel\n  like\n   you're\n    falling...")
    assert.Equal(t, "You\n feel\n  like\n   you're\n", buf.String(), "Print should only output full lines until Closed")
    buf.Reset()
    me.Close()
    assert.Equal(t, "    falling...\n", buf.String(), "Close should flush any unfinished line by appending a terminating newline")
    buf.Reset()
    me.Print("")
    me.Close()
    assert.Equal(t, "", buf.String(), "Print and Close should output nothing if nothing was printed")
}
