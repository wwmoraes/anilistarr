package animelists_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wwmoraes/anilistarr/internal/drivers/animelists"
)

func TestAnilist2TVDBMetadata_GetTargetID(t *testing.T) {
	t.Parallel()

	entry := animelists.Anilist2TVDBMetadata{
		AnilistID: 1,
		TvdbID:    91,
	}

	assert.Equal(t, "91", entry.GetTargetID())
}

func TestAnilist2TVDBMetadata_GetSourceID(t *testing.T) {
	t.Parallel()

	entry := animelists.Anilist2TVDBMetadata{
		AnilistID: 1,
		TvdbID:    91,
	}
	assert.Equal(t, "1", entry.GetSourceID())
}

func TestAnilist2TVDBMetadata_Valid(t *testing.T) {
	t.Parallel()

	type fields struct {
		AnilistID uint64
		TvdbID    uint64
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				AnilistID: 1,
				TvdbID:    91,
			},
			want: true,
		},
		{
			name: "empty source",
			fields: fields{
				AnilistID: 0,
				TvdbID:    91,
			},
			want: false,
		},
		{
			name: "empty target",
			fields: fields{
				AnilistID: 1,
				TvdbID:    0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry := animelists.Anilist2TVDBMetadata{
				AnilistID: tt.fields.AnilistID,
				TvdbID:    tt.fields.TvdbID,
			}
			assert.Equal(t, tt.want, entry.Valid())
		})
	}
}
