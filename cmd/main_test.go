package main

import "testing"

func TestHasAnyPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		file     string
		prefixes []string
		want     bool
	}{
		{
			name:     "exact match",
			file:     "Marktakteure.xml",
			prefixes: []string{"Marktakteure"},
			want:     true,
		},
		{
			name:     "underscore suffix",
			file:     "Marktakteure_20241016.xml",
			prefixes: []string{"Marktakteure"},
			want:     true,
		},
		{
			name:     "no separator between prefix and text",
			file:     "MarktakteureUndRollen.xml",
			prefixes: []string{"Marktakteure"},
			want:     false,
		},
		{
			name:     "longer prefix takes precedence",
			file:     "MarktakteureUndRollen.xml",
			prefixes: []string{"MarktakteureUndRollen"},
			want:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := hasAnyPrefix(tt.file, tt.prefixes)
			if got != tt.want {
				t.Fatalf("hasAnyPrefix(%q, %v) = %t, want %t", tt.file, tt.prefixes, got, tt.want)
			}
		})
	}
}
