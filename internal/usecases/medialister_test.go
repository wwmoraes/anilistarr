package usecases_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/memory"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestNewMediaLister(t *testing.T) {
	t.Parallel()

	type args struct {
		tracker usecases.Tracker
		mapper  usecases.Mapper
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "success",
			args: args{
				tracker: &testdata.Tracker{},
				mapper:  &adapters.Mapper{},
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "nil tracker",
			args: args{
				tracker: nil,
				mapper:  &adapters.Mapper{},
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "nil mapper",
			args: args{
				tracker: &testdata.Tracker{},
				mapper:  nil,
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := usecases.NewMediaLister(tt.args.tracker, tt.args.mapper)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMediaLister() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("NewMediaLister() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestSonarrMediaLister_Generate(t *testing.T) {
	t.Parallel()

	type fields struct {
		Tracker usecases.Tracker
		Mapper  usecases.Mapper
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.CustomList
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Tracker: testdata.SampleTracker,
				Mapper: &adapters.Mapper{
					Provider: testdata.Provider,
					Store:    testdata.SampleStore,
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    testdata.CustomListFromIDs(t, testdata.TargetIDs...),
			wantErr: false,
		},
		{
			name: "user not found",
			fields: fields{
				Tracker: &testdata.Tracker{},
				Mapper:  nil,
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "user id tracker error",
			fields: fields{
				Tracker: &testdata.Tracker{
					MediaLists: map[int][]string{},
				},
				Mapper: nil,
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "media list tracker error",
			fields: fields{
				Tracker: &testdata.Tracker{
					UserIds: map[string]int{
						testdata.Username: testdata.UserID,
					},
					MediaLists: nil,
				},
				Mapper: nil,
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "map IDs error",
			fields: fields{
				Tracker: testdata.SampleTracker,
				Mapper: &adapters.Mapper{
					Provider: nil,
					Store:    (memory.Memory)(nil),
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty target ID",
			fields: fields{
				Tracker: testdata.SampleTracker,
				Mapper: &adapters.Mapper{
					Provider: nil,
					Store: memory.Memory{
						"1": "",
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    entities.CustomList{},
			wantErr: false,
		},
		{
			name: "invalid target ID",
			fields: fields{
				Tracker: testdata.SampleTracker,
				Mapper: &adapters.Mapper{
					Provider: nil,
					Store: memory.Memory{
						"1": "a",
					},
				},
			},
			args: args{
				ctx:  context.TODO(),
				name: testdata.Username,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := &usecases.SonarrMediaLister{
				Tracker: tt.fields.Tracker,
				Mapper:  tt.fields.Mapper,
			}

			got, err := lister.Generate(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SonarrMediaLister.Generate() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SonarrMediaLister.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSonarrMediaLister_GetUserID(t *testing.T) {
	t.Parallel()

	type fields struct {
		Tracker usecases.Tracker
		Mapper  usecases.Mapper
	}

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := &usecases.SonarrMediaLister{
				Tracker: tt.fields.Tracker,
				Mapper:  tt.fields.Mapper,
			}

			got, err := lister.GetUserID(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SonarrMediaLister.GetUserID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("SonarrMediaLister.GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSonarrMediaLister_Close(t *testing.T) {
	t.Parallel()

	type fields struct {
		Tracker usecases.Tracker
		Mapper  usecases.Mapper
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Tracker: &testdata.Tracker{},
				Mapper: &adapters.Mapper{
					Store: memory.New(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := &usecases.SonarrMediaLister{
				Tracker: tt.fields.Tracker,
				Mapper:  tt.fields.Mapper,
			}

			if err := lister.Close(); (err != nil) != tt.wantErr {
				t.Errorf("SonarrMediaLister.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSonarrMediaLister_Refresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		Tracker usecases.Tracker
		Mapper  usecases.Mapper
	}

	type args struct {
		ctx    context.Context
		client usecases.Getter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Tracker: nil,
				Mapper: &adapters.Mapper{
					Provider: testdata.Provider,
					Store:    memory.New(),
				},
			},
			args: args{
				ctx:    context.TODO(),
				client: usecases.HTTPGetterAsGetter(&testdata.SampleClient),
			},
			wantErr: false,
		},
		{
			name: "no client",
			fields: fields{
				Tracker: nil,
				Mapper: &adapters.Mapper{
					Provider: testdata.Provider,
				},
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := &usecases.SonarrMediaLister{
				Tracker: tt.fields.Tracker,
				Mapper:  tt.fields.Mapper,
			}

			if err := lister.Refresh(tt.args.ctx, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("SonarrMediaLister.Refresh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
