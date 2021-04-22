package main

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cases := []struct {
		Text        string
		Cmd         string
		Args        []string
		ErrExpected bool
	}{
		{`!테스트`, "테스트", []string{}, false},
		{`!테스트 명령어1 명령어2`, "테스트", []string{"명령어1", "명령어2"}, false},
		{`!테스트 "명령어 1" '명령어 2'`, "테스트", []string{"명령어 1", "명령어 2"}, false},
		{`!테스트 명령어1 "명령어 2"`, "테스트", []string{"명령어1", "명령어 2"}, false},
		{`!테스트 명령어1 "명령어2`, "", []string{}, true},
		{`!테스트 명령어1 "명령어2'`, "", []string{}, true},
	}

	for _, tt := range cases {
		t.Run(tt.Cmd, func(t *testing.T) {
			cmd, args, err := ParseCommand(tt.Text)
			if (cmd != tt.Cmd) || !reflect.DeepEqual(tt.Args, args) || ((err == nil) == tt.ErrExpected) {
				t.Errorf("failed to test case: %s", tt.Text)
			}
		})
	}
}
