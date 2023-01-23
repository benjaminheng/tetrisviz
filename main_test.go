package main

import "testing"

func TestGetOutputFilename(t *testing.T) {
	type args struct {
		inputFilename string
		format        OutputFormat
	}
	var tests = []struct {
		name     string
		args     args
		expected string
	}{
		{
			"tetrisviz extension, pikchr output",
			args{
				inputFilename: "abc.tetrisviz",
				format:        OutputFormatPikchr,
			},
			"abc.pikchr",
		},
		{
			"tetrisviz extension, svg output",
			args{
				inputFilename: "abc.tetrisviz",
				format:        OutputFormatSVG,
			},
			"abc.svg",
		},
		{
			"non-tetrisviz extension, svg output",
			args{
				inputFilename: "abc.something",
				format:        OutputFormatSVG,
			},
			"abc.something.svg",
		},
		{
			"non-tetrisviz extension, pikchr output",
			args{
				inputFilename: "abc.something",
				format:        OutputFormatPikchr,
			},
			"abc.something.pikchr",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := getOutputFilename(tt.args.inputFilename, tt.args.format)
			if result != tt.expected {
				t.Errorf("(%+v): expected %+v, got %+v", tt.args, tt.expected, result)
			}

		})
	}
}
