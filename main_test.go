package main

//run go generate first
import (
	"os"
	"path/filepath"
	"testing"
)

func Test_setLogFile(t *testing.T) {
	type args struct {
		lf string
	}
	tests := []struct {
		name string
		args args
	}{
		{"log in another folder", args{"log/_test.log"}},
		{"log in this folder", args{"_test.log"}},
		{"log with empty string", args{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := setLogFile(tt.args.lf)
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				t.Fail()
			} else {
				os.Remove(fileName)
				os.Remove(filepath.Dir(fileName))
			}
		})
	}
}

/*func Test_getCurrentService(t *testing.T) {
	const providerURL = "http://localhost:4000"
	currentServices, err := getCurrentServices(providerURL)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", currentServices)
}*/
