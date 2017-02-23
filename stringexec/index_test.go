package stringexec

import "testing"

type testtable struct {
	isWindows bool
	cmd       string
	expect    string
}

func TestPrepareCommand(t *testing.T) {

	table := []testtable{
		testtable{
			isWindows: false,
			cmd: `echo "ok"
echo "ok"`,
			expect: `echo "ok" \
 && echo "ok"`,
		},
		testtable{
			isWindows: false,
			cmd: `echo "ok" \
&& echo "ok"`,
			expect: `echo "ok" \
&& echo "ok"`,
		},
		testtable{
			isWindows: false,
			cmd: `echo "ok" \
    && echo "ok"`,
			expect: `echo "ok" \
    && echo "ok"`,
		},
		testtable{
			isWindows: true,
			cmd: `echo "ok"
echo "ok"`,
			expect: `echo "ok" && echo "ok"`,
		},
		testtable{
			isWindows: true,
			cmd: `echo "ok" \
&& echo "ok"`,
			expect: `echo "ok" \
&& echo "ok"`,
		},
		testtable{
			isWindows: true,
			cmd: `echo "ok" \
    && echo "ok"`,
			expect: `echo "ok" \
    && echo "ok"`,
		},
	}
	for i, test := range table {
		got := prepareCommand(test.cmd, test.isWindows)
		if got != test.expect {
			t.Errorf(
				"Test(%v): isWindows=%v Invalid command result\nexpected=\n%q\ngot=\n%q\n",
				i,
				test.isWindows,
				test.expect,
				got,
			)
		}
	}
}
