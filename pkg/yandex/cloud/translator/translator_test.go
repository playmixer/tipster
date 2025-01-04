package translator

import "testing"

func TestTranslator_Translate(t *testing.T) {
	type args struct {
		sourceLang string
		targetLang string
		text       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "hi",
			args: args{
				sourceLang: "ru",
				targetLang: "en",
				text:       "привет",
			},
			want:    "hi",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := New("")
			got, err := tr.Translate(tt.args.sourceLang, tt.args.targetLang, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translator.Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Translator.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}
