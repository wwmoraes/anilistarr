package adapters_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/goccy/go-json"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/memory"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/functional"
)

func TestMapper_MapIDs(t *testing.T) {
	t.Parallel()

	type fields struct {
		Provider adapters.Provider[adapters.Metadata]
		Store    adapters.Store
	}

	type args struct {
		ctx context.Context
		ids []entities.SourceID
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entities.TargetID
		wantErr bool
	}{
		{
			name: "full hit",
			fields: fields{
				Provider: newJSONLocalProvider(t),
				Store: memory.Memory{
					"foo": "bar",
					"baz": "qux",
				},
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo", "baz"},
			},
			want:    []string{"bar", "qux"},
			wantErr: false,
		},
		{
			name: "full miss",
			fields: fields{
				Provider: newJSONLocalProvider(t),
				Store:    memory.Memory{},
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo", "baz"},
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "partial hit",
			fields: fields{
				Provider: newJSONLocalProvider(t),
				Store: memory.Memory{
					"foo": "bar",
				},
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo", "baz"},
			},
			wantErr: false,
			want:    []string{"bar"},
		},
		{
			name: "store failure",
			fields: fields{
				Provider: newJSONLocalProvider(t),
				Store:    (memory.Memory)(nil),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapper := &adapters.Mapper{
				Provider: tt.fields.Provider,
				Store:    tt.fields.Store,
			}
			defer mapper.Close()

			got, err := mapper.MapIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mapper.MapIDs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mapper.MapIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapper_Refresh(t *testing.T) {
	t.Parallel()

	type fields struct {
		Provider adapters.Provider[adapters.Metadata]
		Store    adapters.Store
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
				Provider: newJSONLocalProvider(t),
				Store:    memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: false,
		},
		{
			name: "provider error",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs:   &testdata.MemoryFS{},
					Name: "",
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: true,
		},
		{
			name: "empty source ID",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs: &testdata.MemoryFS{
						testLocalName: functional.Unwrap(json.Marshal([]adapters.Metadata{
							memoryRawMetadata{
								AnilistID: "",
								TheTvdbID: "90",
							},
						})),
					},
					Name: testLocalName,
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: true,
		},
		{
			name: "empty target ID",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs: &testdata.MemoryFS{
						testLocalName: functional.Unwrap(json.Marshal([]adapters.Metadata{
							memoryRawMetadata{
								AnilistID: "1",
								TheTvdbID: "",
							},
						})),
					},
					Name: testLocalName,
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: true,
		},
		{
			name: "zero source ID",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs: &testdata.MemoryFS{
						testLocalName: functional.Unwrap(json.Marshal([]adapters.Metadata{
							memoryMetadata{
								AnilistID: 0,
								TheTvdbID: 90,
							},
						})),
					},
					Name: testLocalName,
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: false,
		},
		{
			name: "zero target ID",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs: &testdata.MemoryFS{
						testLocalName: functional.Unwrap(json.Marshal([]adapters.Metadata{
							memoryMetadata{
								AnilistID: 90,
								TheTvdbID: 0,
							},
						})),
					},
					Name: testLocalName,
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: false,
		},
		{
			name: "getter failure",
			fields: fields{
				Provider: &adapters.JSONLocalProvider[memoryMetadata]{
					Fs:   &testdata.MemoryFS{},
					Name: testLocalName,
				},
				Store: memory.New(),
			},
			args: args{
				ctx:    context.TODO(),
				client: nil,
			},
			wantErr: true,
		},
		{
			name: "store failure",
			fields: fields{
				Provider: newJSONLocalProvider(t),
				Store:    (memory.Memory)(nil),
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

			mapper := &adapters.Mapper{
				Provider: tt.fields.Provider,
				Store:    tt.fields.Store,
			}
			defer mapper.Close()

			if err := mapper.Refresh(tt.args.ctx, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("Mapper.Refresh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
