package language

import (
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		expectedLang  string
		expectedError error
	}{
		{
			name:          "simple english",
			text:          "Therefore, it can be presented as 2d space with threshold functions, that splits it into Reliable and Not reliable areas.",
			expectedLang:  "en",
			expectedError: nil,
		},
		{
			name:          "Turkce bir ornek",
			text:          "Bu bir Turkce ornek olup, Turkce bulmasi beklenmektedir",
			expectedLang:  "tr",
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualLang, _, actualErr := detectLanguage(test.text)
			if actualErr != test.expectedError {
				t.Fatalf("expected err %v, got %v", test.expectedError, actualErr)
			}

			if actualLang != test.expectedLang {
				t.Errorf("expected language %s, got %s", test.expectedLang, actualLang)
			}
		})
	}
}

func BenchmarkDetectLanguage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		detectLanguage("Therefore, it can be presented as 2d space with threshold functions, that splits it into Reliable and Not reliable areas.")
	}
}
