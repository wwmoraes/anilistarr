//go:build !pure

package sqlite_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/wwmoraes/anilistarr/internal/drivers/sqlite"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/finalizers"
)

type functor = testdata.Functor[*sqlite.SQLite]

func newSQLite(tb testing.TB) *sqlite.SQLite {
	tb.Helper()

	db, err := sqlite.New("file::memory:")
	if err != nil {
		tb.Fatal(err)
	}

	return db
}

func putStrings(entries ...[2]string) functor {
	return functor(func(tb testing.TB, db *sqlite.SQLite) *sqlite.SQLite {
		tb.Helper()

		var err error

		for _, pair := range entries {
			err = db.SetString(context.TODO(), pair[0], pair[1])
			if err != nil {
				tb.Fatal(err)
			}
		}

		return db
	})
}

func putMedias(medias ...*entities.Media) functor {
	return functor(func(tb testing.TB, db *sqlite.SQLite) *sqlite.SQLite {
		tb.Helper()

		var err error

		for _, media := range medias {
			err = db.PutMedia(context.TODO(), media)
			if err != nil {
				tb.Fatal(err)
			}
		}

		return db
	})
}

func closeSQLite(tb testing.TB, db *sqlite.SQLite) *sqlite.SQLite {
	tb.Helper()

	err := db.Close()
	if err != nil {
		tb.Fatal(err)
	}

	return db
}

func sortMedia(a, b *entities.Media) int {
	return strings.Compare(a.SourceID, b.SourceID)
}

//nolint:gocognit // needed to test cancelled context side-effects
func newInterruptCtx(tb testing.TB, successes uint, caller string) *testdata.MockContext {
	tb.Helper()

	doneChan := make(chan struct{})

	var outDoneChan <-chan struct{} = doneChan

	ctx := testdata.MockContext{}

	ctx.On("Value", mock.Anything).Return(nil)
	errCall := ctx.On("Err").Return(nil)

	ctx.On("Done").Return(outDoneChan).Run(func(_ mock.Arguments) {
		if !testdata.CallerMatches(tb, caller) {
			return
		}

		count := ctx.TestData().Get("count").Uint(0)
		count++

		if count <= successes {
			ctx.TestData().Set("count", count)

			return
		}

		closed := ctx.TestData().Get("closed").Bool(false)
		if !closed {
			// make sure the channel is empty
			select {
			default:
			case <-doneChan:
			}

			close(doneChan)
			ctx.TestData().Set("closed", true)
			errCall.Return(context.Canceled)
		}
	})

	tb.Cleanup(func() {
		defer func() {
			if err := recover(); err != nil {
				tb.Log(err)
			}
		}()

		close(doneChan)
	})

	return &ctx
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		dataSourceName string
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		assertValue assert.ValueAssertionFunc
		name        string
		args        args
	}{
		{
			name: "memory",
			args: args{
				dataSourceName: "file:?mode=memory",
			},
			assertError: require.NoError,
			assertValue: assert.NotNil,
		},
		{
			name: "invalid source name",
			args: args{
				dataSourceName: "file:?;",
			},
			assertError: require.Error,
			assertValue: assert.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sqlite.New(tt.args.dataSourceName)
			tt.assertError(t, err)

			tt.assertValue(t, got)
		})
	}
}

func TestSQLite_Close(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		fields      fields
		name        string
	}{
		{
			name: "memory",
			fields: fields{
				db: newSQLite(t),
			},
			assertError: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.fields.db.Close()
			tt.assertError(t, err)
		})
	}
}

func TestSQLite_GetString(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx context.Context
		key string
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		name        string
		fields      fields
		args        args
		want        string
	}{
		{
			name: "empty db",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:        "",
			assertError: require.Error,
		},
		{
			name: "empty key",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "",
			},
			want:        "",
			assertError: require.Error,
		},
		{
			name: "success",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), putStrings(
					[2]string{"foo", "bar"},
				)),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:        "bar",
			assertError: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer finalizers.Close(tt.fields.db)

			got, err := tt.fields.db.GetString(tt.args.ctx, tt.args.key)
			tt.assertError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLite_SetString(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx     context.Context
		key     string
		value   string
		options []usecases.CacheOption
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		name        string
		fields      fields
		args        args
	}{
		{
			name: "success",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "bar",
			},
			assertError: require.NoError,
		},
		{
			name: "empty key",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "",
				value: "bar",
			},
			assertError: require.Error,
		},
		{
			name: "empty value",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "",
			},
			assertError: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.fields.db.SetString(
				tt.args.ctx,
				tt.args.key,
				tt.args.value,
				tt.args.options...,
			)
			tt.assertError(t, err)
		})
	}
}

func TestSQLite_GetMedia(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		fields      fields
		want        *entities.Media
		args        args
		name        string
	}{
		{
			name: "empty",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				id:  "foo",
			},
			want:        nil,
			assertError: require.Error,
		},
		{
			name: "db error",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), closeSQLite),
			},
			args: args{
				ctx: context.TODO(),
				id:  "",
			},
			want:        nil,
			assertError: require.Error,
		},
		{
			name: "success",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), putMedias(
					&entities.Media{
						SourceID: "foo",
						TargetID: "bar",
					},
				)),
			},
			args: args{
				ctx: context.TODO(),
				id:  "foo",
			},
			want: &entities.Media{
				SourceID: "foo",
				TargetID: "bar",
			},
			assertError: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.fields.db.GetMedia(tt.args.ctx, tt.args.id)
			tt.assertError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLite_GetMediaBulk(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx context.Context
		ids []string
	}

	tests := []struct {
		assertError require.ErrorAssertionFunc
		name        string
		fields      fields
		args        args
		want        []*entities.Media
	}{
		{
			name: "nil",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				ids: nil,
			},
			want:        []*entities.Media{},
			assertError: require.NoError,
		},
		{
			name: "null",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{},
			},
			want:        []*entities.Media{},
			assertError: require.NoError,
		},
		{
			name: "db error",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), closeSQLite),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{},
			},
			want:        nil,
			assertError: require.Error,
		},
		{
			name: "no matches",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo"},
			},
			want:        []*entities.Media{},
			assertError: require.NoError,
		},
		{
			name: "partial match",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), putMedias(
					&entities.Media{
						SourceID: "foo",
						TargetID: "bar",
					},
				)),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo", "baz"},
			},
			want: []*entities.Media{
				{
					SourceID: "foo",
					TargetID: "bar",
				},
			},
			assertError: require.NoError,
		},
		{
			name: "total match",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), putMedias(
					&entities.Media{
						SourceID: "foo",
						TargetID: "bar",
					},
					&entities.Media{
						SourceID: "baz",
						TargetID: "qux",
					},
				)),
			},
			args: args{
				ctx: context.TODO(),
				ids: []string{"foo", "baz"},
			},
			want: []*entities.Media{
				{
					SourceID: "foo",
					TargetID: "bar",
				},
				{
					SourceID: "baz",
					TargetID: "qux",
				},
			},
			assertError: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.fields.db.GetMediaBulk(tt.args.ctx, tt.args.ids)
			tt.assertError(t, err)

			slices.SortFunc(got, sortMedia)
			slices.SortFunc(tt.want, sortMedia)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLite_PutMedia(t *testing.T) {
	t.Parallel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx   context.Context
		media *entities.Media
	}

	tests := []struct {
		args      args
		wantError error
		fields    fields
		name      string
	}{
		{
			name: "success",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				media: &entities.Media{
					SourceID: "1",
					TargetID: "91",
				},
			},
			wantError: nil,
		},
		{
			name: "empty source ID error",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				media: &entities.Media{
					SourceID: "",
					TargetID: "91",
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "empty target ID error",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				media: &entities.Media{
					SourceID: "1",
					TargetID: "",
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "closed db error",
			fields: fields{
				db: testdata.Compose(t, newSQLite(t), closeSQLite),
			},
			args: args{
				ctx: context.TODO(),
				media: &entities.Media{
					SourceID: "1",
					TargetID: "91",
				},
			},
			wantError: usecases.ErrStatusFailedPrecondition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.fields.db.PutMedia(tt.args.ctx, tt.args.media)
			require.ErrorIs(t, err, tt.wantError)
		})
	}
}

func TestSQLite_PutMediaBulk(t *testing.T) {
	t.Parallel()

	cancelledCtx, cancel := context.WithCancel(context.TODO())
	cancel()

	type fields struct {
		db *sqlite.SQLite
	}

	type args struct {
		ctx    context.Context
		medias []*entities.Media
	}

	tests := []struct {
		wantError error
		fields    fields
		name      string
		args      args
	}{
		{
			name: "single new",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "91",
					},
				},
			},
			wantError: nil,
		},
		{
			name: "multiple new",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "91",
					},
					{
						SourceID: "2",
						TargetID: "92",
					},
				},
			},
			wantError: nil,
		},
		{
			name: "invalid source id",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "",
						TargetID: "91",
					},
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "invalid target id",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "",
					},
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "empty medias",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx:    context.TODO(),
				medias: []*entities.Media{},
			},
			wantError: nil,
		},
		{
			name: "nil medias",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx:    context.TODO(),
				medias: nil,
			},
			wantError: nil,
		},
		{
			name: "done context",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: cancelledCtx,
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "91",
					},
				},
			},
			wantError: usecases.ErrStatusFailedPrecondition,
		},
		{
			name: "cancelled on put",
			fields: fields{
				db: newSQLite(t),
			},
			args: args{
				ctx: newInterruptCtx(t, 0, "modernc.org/sqlite.(*stmt).exec"),
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "91",
					},
				},
			},
			wantError: usecases.ErrStatusAborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.fields.db.PutMediaBulk(tt.args.ctx, tt.args.medias)
			require.ErrorIs(t, err, tt.wantError)
		})
	}
}

func ExampleNew() {
	db, err := sqlite.New("file:///dev/null?;;")
	fmt.Println(db)
	fmt.Println(err)
	// Output:
	// <nil>
	// failed to connect to database: invalid semicolon separator in query
}

func ExampleSQLite_GetMediaBulk() {
	db, err := sqlite.New("file::memory:?cache=shared&_pragma=journal_mode=wal&_txlock=exclusive")
	fmt.Println("New err:", err)

	if err != nil {
		return
	}

	err = db.PutMediaBulk(context.Background(), []*entities.Media{
		{
			SourceID: "foo",
			TargetID: "bar",
		},
		{
			SourceID: "baz",
			TargetID: "qux",
		},
	})
	fmt.Println("PutMediaBulk err:", err)

	medias, err := db.GetMediaBulk(context.Background(), []string{"foo", "baz"})
	fmt.Println("GetMediaBulk err:", err)

	for _, media := range medias {
		fmt.Println(media.SourceID, "=", media.TargetID)
	}

	// Unordered output:
	// New err: <nil>
	// PutMediaBulk err: <nil>
	// GetMediaBulk err: <nil>
	// foo = bar
	// baz = qux
}

func ExampleNew_schema_failure() {
	db, err := sqlite.New("file::memory:?_pragma=query_only=true")
	fmt.Println("DB:", db)
	fmt.Println("Error:", err)

	// Output:
	// DB: <nil>
	// Error: failed to execute schema queries: attempt to write a readonly database (8)
}
