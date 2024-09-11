package fonbnk

import (
	"os"
	"testing"
)

func Test_generateSignature(t *testing.T) {
	clientSecret := os.Getenv("FONBNK_CLIENT_SECRET")

	type args struct {
		clientSecret string
		timestamp    string
		endpoint     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Generate correct signature",
			args: args{
				clientSecret: clientSecret,
				timestamp:    "1726048611284",
				endpoint:     "/api/offramp/limits?type=mobile_money&country=KE",
			},
			// Generated with the official TS example + sandbox client secret
			// We want to ensure correctness in signature generation
			want: "pesqxOs0AVLg4CJV/zfA0nsNF4TkfCUAo8yQ+T8eRLg=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSignature(tt.args.clientSecret, tt.args.timestamp, tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
