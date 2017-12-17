// Test the linuxdir code

package linuxdir_test

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"davidb.org/x/gosure/linuxdir"
)

// Compare the output of the linuxdir call with our own manual loop
// doing the same thing less efficiently.  Do with a read-only
// directory to reduce the chance of a race between these.
func TestDirs(t *testing.T) {
	a, err := linuxdir.Readdir("/usr/bin")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d entries", len(a))

	// Read the same thing, with standard calls.
	b, err := ioutil.ReadDir("/usr/bin")
	if err != nil {
		t.Fatal(err)
	}
	sort.Sort((*nameSort)(&b))
	t.Logf("Got %d entries", len(b))

	if len(a) != len(b) {
		t.Fatal("Different number of entries in dir read")
	}

	for i := 0; i < len(a); i++ {
		if a[i].Name() != b[i].Name() {
			t.Fatalf("Directory read mismatch: %s, %s", a[i].Name(), b[i].Name())
		}
	}
}

type nameSort []os.FileInfo

func (p nameSort) Len() int           { return len(p) }
func (p nameSort) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p nameSort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
