package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

func TestMedia_Valid(t *testing.T) {
	t.Parallel()

	type fields struct {
		SourceID entities.SourceID
		TargetID entities.TargetID
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid",
			fields: fields{
				SourceID: "1",
				TargetID: "91",
			},
			want: true,
		},
		{
			name: "empty source ID",
			fields: fields{
				SourceID: "",
				TargetID: "91",
			},
			want: false,
		},
		{
			name: "empty target ID",
			fields: fields{
				SourceID: "1",
				TargetID: "",
			},
			want: false,
		},
		{
			name: "zero source ID",
			fields: fields{
				SourceID: "0",
				TargetID: "91",
			},
			want: false,
		},
		{
			name: "zero target ID",
			fields: fields{
				SourceID: "1",
				TargetID: "0",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			media := &entities.Media{
				SourceID: tt.fields.SourceID,
				TargetID: tt.fields.TargetID,
			}
			assert.Equal(t, tt.want, media.Valid())
		})
	}
}
