package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailBodyWithFooter(t *testing.T) {
	cases := []struct {
		key           MailKey
		data          map[string]any
		expectSubject string
		expectBody    string
	}{
		{
			key: MailKeyNotifyScholarshipDeadline,
			data: map[string]any{
				"Scholarships": "奨学金詳細",
			},
			expectSubject: "【クレクレ奨学金】まもなく締切(1週間後)の奨学金のお知らせ",
			expectBody:    "もうすぐ締め切りの奨学金があります。\nご確認の上、お早めにご応募ください。",
		},
	}

	for _, tt := range cases {
		t.Run(string(tt.key), func(t *testing.T) {
			subject, body, err := tt.key.MailBodyWithFooter(tt.data)
			require.NoError(t, err)
			assert.Contains(t, subject, tt.expectSubject)
			assert.Contains(t, body, tt.expectBody)
		})
	}
}
