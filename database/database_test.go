package database

import (
	"testing"

	"github.com/Cloud-Pie/Passa/ymlparser"
)

func TestSearchQuery(t *testing.T) {
	type args struct {
		currentStates []ymlparser.State
		searchName    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Find in Config",
			args: args{
				currentStates: []ymlparser.State{{Name: "state-with-7"}},
				searchName:    "state-with-7",
			},
			want: 0,
		},
		{
			name: "Don't find in Config",
			args: args{
				currentStates: ymlparser.ParseStatesfile("../test/passa-states-test.yml").States,
				searchName:    "dummy-State",
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SearchQuery(tt.args.currentStates, tt.args.searchName); got != tt.want {
				t.Errorf("SearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
