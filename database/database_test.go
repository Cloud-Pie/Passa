package database

import (
	"reflect"
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
			if got := searchQuery(tt.args.currentStates, tt.args.searchName); got != tt.want {
				t.Errorf("SearchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitializeDB(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create db file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeDB(&ymlparser.Config{})
		})
	}
}

func Test_insertDB(t *testing.T) {
	InitializeDB(&ymlparser.Config{})
	type args struct {
		newState ymlparser.State
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "add to db",
			args: args{newState: ymlparser.State{
				Name: "my new state",
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			InsertState(tt.args.newState)
		})
	}
}

func Test_readAll(t *testing.T) {
	dropDB()
	InitializeDB(&ymlparser.Config{})
	myConfig := ymlparser.Config{
		States: []ymlparser.State{
			{
				Name: "mystate1",
			},
			{
				Name: "mystate2",
			},
		},
	}
	for s := range myConfig.States {

		InsertState(myConfig.States[s])

	}
	tests := []struct {
		name string
		want []ymlparser.State
	}{
		{
			name: "Check is read correctly",
			want: myConfig.States,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadAllStates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dropDB(t *testing.T) {
	InitializeDB(&ymlparser.Config{})
	tests := []struct {
		name string
	}{
		{
			name: "DROP!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dropDB()
		})
	}
}
