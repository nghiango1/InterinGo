package handlers

import (
	"encoding/json"
	"fmt"
	"interingo/pkg/lexer"
	"interingo/pkg/parser"
	"interingo/pkg/share"
	"log/slog"
	"os"
	"strings"
	"testing"
)

// Diff returns a slice of indexes where a and b differ,
// and the runes from each string at those positions (useful when lengths differ).
type DiffEntry struct {
	Index int
	A     rune // 0 if index >= len(a)
	B     rune // 0 if index >= len(b)
}

func fmtDiffs(diffs []DiffEntry) string {
	if len(diffs) == 0 {
		return "<nil>"
	}
	var s strings.Builder
	for _, d := range diffs {
		a := fmt.Sprintf("%q", d.A)
		b := fmt.Sprintf("%q", d.B)
		if d.A == 0 {
			a = "<nil>"
		}
		if d.B == 0 {
			b = "<nil>"
		}
		fmt.Fprintf(&s, "(%d: %s vs %s) ", d.Index, a, b)
	}
	return s.String()
}

func Diff(a, b string) []DiffEntry {
	ar := []rune(a)
	br := []rune(b)
	n := max(len(br), len(ar))
	var diffs []DiffEntry
	for i := range n {
		var ra, rb rune
		if i < len(ar) {
			ra = ar[i]
		}
		if i < len(br) {
			rb = br[i]
		}
		if ra != rb {
			diffs = append(diffs, DiffEntry{Index: i, A: ra, B: rb})
		}
	}
	return diffs
}

func TestCommentLineFormat(t *testing.T) {
	input := `			//This is a coment line
		add(1, 2 * 3, 4 +		5);
// This is a second commented line`

	expectedOutput := `//This is a coment line
add(1, 2*3, 4+5);
// This is a second commented line
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	output := FormatedAST(program, FormattingOptions{}, 0)
	if output != expectedOutput {
		jsonData, err := json.Marshal(program)
		if err != nil {
			fmt.Printf("[ERROR] Got: %v\n", err)
		} else {
			fmt.Printf("[INFO] Data: (%v) %v\n", len(jsonData), string(jsonData))
		}

		t.Fatalf("Format failed, output doesn't match. Got: \n%v", fmtDiffs(Diff(output, expectedOutput)))
	}
}

func TestFunctionFormatHard(t *testing.T) {
	input := `// 5
let c = fn(x)  {
// hw
    return x;
    //asd'
}(12,2);
// return //as das/ahehehe  a/ a/ a a/a /a a/ a aa/ a 
        // as das ssdas dahuhuhuhuh u hu hu uh uh  h  uh h h h taht tso oos oso sso os os 
// return hehe(hhh, h, hh) // asasdasdasd
        // let c =
// if (a(c){a{c}/3 * 4 - a */} /)
`

	expectedOutput := `// 5
let c = fn(x) {
    // hw
    return x;
    //asd'
}(12, 2);
// return //as das/ahehehe  a/ a/ a a/a /a a/ a aa/ a 
// as das ssdas dahuhuhuhuh u hu hu uh uh  h  uh h h h taht tso oos oso sso os os 
// return hehe(hhh, h, hh) // asasdasdasd
// let c =
// if (a(c){a{c}/3 * 4 - a */} /)
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	tabSize := 4
	output := FormatedAST(program, FormattingOptions{
		TabSize: &tabSize,
	}, 0)
	if output != expectedOutput {
		jsonData, err := json.Marshal(program)
		if err != nil {
			slog.Error(fmt.Sprintf("Got: %v\n", err))
		} else {
			slog.Info(fmt.Sprintf("Data: (%v) %v\n", len(jsonData), string(jsonData)))
		}

		t.Fatalf("Format failed, output doesn't match. Got: \n%v", fmtDiffs(Diff(output, expectedOutput)))
	}
}

func TestMain(m *testing.M) {
	share.SetDefaultLog(slog.LevelDebug)
	os.Exit(m.Run())
}
