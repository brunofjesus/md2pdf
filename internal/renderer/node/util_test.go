package node

import (
	"testing"
)

func TestJoinBase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		base string
		dest string
		want string
	}{
		{
			name: "local base with relative file",
			base: "/home/user/docs",
			dest: "img.png",
			want: "/home/user/docs/img.png",
		},
		{
			name: "local base strips leading ./",
			base: "/home/user/docs",
			dest: "./img.png",
			want: "/home/user/docs/img.png",
		},
		{
			name: "local base with nested relative path",
			base: "/home/user/docs",
			dest: "sub/img.png",
			want: "/home/user/docs/sub/img.png",
		},
		{
			name: "http base with relative file",
			base: "https://example.com",
			dest: "img.png",
			want: "https://example.com/img.png",
		},
		{
			name: "http base with trailing slash avoids double slash",
			base: "https://example.com/",
			dest: "img.png",
			want: "https://example.com/img.png",
		},
		{
			name: "http base strips leading ./",
			base: "https://example.com",
			dest: "./img.png",
			want: "https://example.com/img.png",
		},
		{
			name: "http base with nested relative path",
			base: "https://example.com/assets",
			dest: "sub/img.png",
			want: "https://example.com/assets/sub/img.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := joinBase(tt.base, tt.dest); got != tt.want {
				t.Errorf("joinBase(%q, %q) = %q, want %q", tt.base, tt.dest, got, tt.want)
			}
		})
	}
}
