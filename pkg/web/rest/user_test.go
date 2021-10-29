package rest

import "testing"

func Test_validatePassword(t *testing.T) {
	type args struct {
		password []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "denies empty password",
			args: args{
				password: []byte(""),
			},
			want: false,
		},
		{
			name: "denies long but still unsafe password",
			args: args{
				password: []byte("aaaaaaaaaaaa"),
			},
			want: false,
		},
		{
			name: "denies too short password with required chars",
			args: args{
				password: []byte("Aa1232#"),
			},
			want: false,
		},
		{
			name: "accepts relatively secure password",
			args: args{
				password: []byte("Aa1kscnl232#"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePassword(tt.args.password); got != tt.want {
				t.Errorf("validatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "denies too short username",
			args: args{username: "aa"},
			want: false,
		},
		{
			name: "denies username with special characters",
			args: args{username: "usernamexyz+#"},
			want: false,
		},
		{
			name: "accepts acceptable username",
			args: args{username: "usernamexyz"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUsername(tt.args.username); got != tt.want {
				t.Errorf("validateUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
