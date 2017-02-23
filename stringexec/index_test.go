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
			expect: `echo "ok" && echo "ok"`,
		},
		testtable{
			isWindows: true,
			cmd: `echo "ok" \
    && echo "ok"`,
			expect: `echo "ok" && echo "ok"`,
		},
		testtable{
			isWindows: false,
			cmd: `666 gh-api-cli create-release -n release -o mh-cbon --guess \
       --ver !newversion! -c "changelog ghrelease --version !newversion!" \
      --draft !isprerelease!`,
			expect: `666 gh-api-cli create-release -n release -o mh-cbon --guess \
       --ver !newversion! -c "changelog ghrelease --version !newversion!" \
      --draft !isprerelease!`,
		},
		testtable{
			isWindows: true,
			cmd: `666 gh-api-cli create-release -n release -o mh-cbon --guess \
       --ver !newversion! -c "changelog ghrelease --version !newversion!" \
      --draft !isprerelease!`,
			expect: `666 gh-api-cli create-release -n release -o mh-cbon --guess --ver !newversion! -c "changelog ghrelease --version !newversion!" --draft !isprerelease!`,
		},
	}
	for i, test := range table {
		got := prepareCommand(test.cmd, test.isWindows)
		if got != test.expect {
			t.Errorf(
				"Test(%v): isWindows=%v Invalid command result\n%q\n%q\n",
				i,
				test.isWindows,
				test.expect,
				got,
			)
		}
	}
}
