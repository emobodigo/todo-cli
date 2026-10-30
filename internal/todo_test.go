package internal

import (
	"os"
	"strings"
	"testing"
)

func TestDeletePersistsRemainingTasks(t *testing.T) {
	const header = "No,Task,CreatedAt,Done\n"
	rows := []string{
		"1,first,2026-10-28 20:16:50,false\n",
		"2,middle,2026-10-28 20:16:51,true\n",
		"3,last,2026-10-28 20:16:52,false\n",
	}
	for _, tt := range []struct {
		id      string
		removed int
	}{
		{"1", 0},
		{"2", 1},
		{"3", 2},
		{"99", -1},
	} {
		t.Run(tt.id, func(t *testing.T) {
			t.Chdir(t.TempDir())
			if err := os.WriteFile("todo.csv", []byte(header+strings.Join(rows, "")), 0o600); err != nil {
				t.Fatal(err)
			}
			Delete(tt.id)
			got, err := os.ReadFile("todo.csv")
			if err != nil {
				t.Fatal(err)
			}
			want := header
			for i, row := range rows {
				if i != tt.removed {
					want += row
				}
			}
			if string(got) != want {
				t.Errorf("persisted CSV = %q, want %q", got, want)
			}
		})
	}
}
