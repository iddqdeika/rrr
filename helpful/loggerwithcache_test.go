package helpful

import "testing"

func TestTrimRight(t *testing.T) {
	source := []string{"1", "2", "3", "4", "5"}
	res := pickRight(source, 3)
	if len(res) != 3 {
		t.Fatal("incorrect pickRight result len")
	}
}

func TestLoggerWithCache(t *testing.T) {
	source := []string{"1", "2", "3", "4", "5"}
	l, c := LoggerWithCache(3)
	for _, s := range source {
		l.Infof("some %v", s)
	}
	res := c.GetLastLogs()
	if len(res) != 3 {
		t.Fatal("incorrect cached log result len")
	}

}
