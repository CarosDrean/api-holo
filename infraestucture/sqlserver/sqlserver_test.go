package sqlserver

import (
	"database/sql"
	"testing"
	"time"

	"api-holo/model"

	"github.com/stretchr/testify/assert"
)

func TestBuildSQLInsert(t *testing.T) {
	tableTest := []struct {
		table   string
		fields  []string
		fieldID string
		want    string
	}{
		{
			table:   "cashboxes",
			fields:  []string{"responsable", "country", "user_id", "account"},
			fieldID: "id",
			want:    "INSERT INTO cashboxes (responsable,country,user_id,account) VALUES (@p1,@p2,@p3,@p4) RETURNING id, d_InsertDate",
		},
		{
			table:   "nothing",
			fields:  []string{},
			fieldID: "id",
			want:    "INSERT INTO nothing () VALUES () RETURNING id, d_InsertDate",
		},
		{
			table:   "one",
			fields:  []string{"one_field"},
			fieldID: "id",
			want:    "INSERT INTO one (one_field) VALUES (@p1) RETURNING id, d_InsertDate",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLInsert(tt.table, tt.fields, tt.fieldID))
	}
}

func TestBuildSQLUpdateByID(t *testing.T) {
	tableTest := []struct {
		table   string
		fields  []string
		fieldID string
		want    string
	}{
		{
			table:   "cashboxes",
			fields:  []string{"responsable", "country", "user_id", "account"},
			fieldID: "id",
			want:    "UPDATE cashboxes SET responsable = @p1, country = @p2, user_id = @p3, account = @p4, d_UpdateDate = GETDATE() WHERE id = @p5",
		},
		{
			table:   "nothing",
			fields:  []string{},
			fieldID: "",
			want:    "",
		},
		{
			table:   "one",
			fields:  []string{"one_field"},
			fieldID: "id",
			want:    "UPDATE one SET one_field = @p1, d_UpdateDate = GETDATE() WHERE id = @p2",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLUpdateByID(tt.table, tt.fields, tt.fieldID))
	}
}

func TestBuildSQLSelect(t *testing.T) {
	tableTest := []struct {
		table   string
		fields  []string
		fieldID string
		want    string
	}{
		{
			table:   "cashboxes",
			fields:  []string{"responsable", "country", "user_id", "account"},
			fieldID: "id",
			want:    "SELECT id, responsable, country, user_id, account, d_InsertDate, d_UpdateDate FROM cashboxes",
		},
		{
			table:   "nothing",
			fields:  []string{},
			fieldID: "id",
			want:    "",
		},
		{
			table:   "one",
			fields:  []string{"one_field"},
			fieldID: "id",
			want:    "SELECT id, one_field, d_InsertDate, d_UpdateDate FROM one",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLSelect(tt.table, tt.fields, tt.fieldID))
	}
}

func TestBuildSQLSelectFields(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			want:   "SELECT responsable, country, user_id, account FROM cashboxes",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			want:   "SELECT one_field FROM one",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLSelectFields(tt.table, tt.fields))
	}
}

func TestBuildSQLWhere(t *testing.T) {
	fakeDate := time.Date(2021, 4, 28, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	tableTest := []struct {
		name      string
		fields    model.Fields
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "where with ILIKE",
			fields: model.Fields{
				{Name: "id", Value: []uint{1, 2, 3}, Operator: model.In},
			},
			wantQuery: "WHERE id IN (1,2,3)",
			wantArgs:  []interface{}{},
		},
		{
			name: "where with all operators",
			fields: model.Fields{
				{Name: "name", Value: "Alejandro"},
				{Name: "age", Value: 30, ChainingKey: model.Or},
				{Name: "course", Value: "Go"},
				{Name: "id", Value: []uint{1, 4, 9}, Operator: model.In},
				{Name: "DESCRIPTION", Value: "%golang%", Operator: model.Ilike},
				{Name: "certificates", Value: 3, Operator: model.GreaterThan},
				{Name: "is_active", Value: true},
			},
			wantQuery: "WHERE name = @p1 AND age = @p2 OR course = @p3 AND id IN (1,4,9) AND description ILIKE @p4 AND certificates > @p5 AND is_active = @p6",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: "Alejandro"}, sql.NamedArg{Name: "p2", Value: 30}, sql.NamedArg{Name: "p3", Value: "Go"}, sql.NamedArg{Name: "p4", Value: "%golang%"}, sql.NamedArg{Name: "p5", Value: 3}, sql.NamedArg{Name: "p6", Value: true}},
		},
		{
			name: "where with operators and string ILIKE",
			fields: model.Fields{
				{Name: "country", Value: "COLOMBIA"},
				{Name: "currency_id", Value: 3, ChainingKey: model.Or},
				{Name: "enable", Value: true},
				{Name: "code", Value: []string{"COL", "COP"}, Operator: model.In},
			},
			wantQuery: "WHERE country = @p1 AND currency_id = @p2 OR enable = @p3 AND code IN ('COL','COP')",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: "COLOMBIA"}, sql.NamedArg{Name: "p2", Value: 3}, sql.NamedArg{Name: "p3", Value: true}},
		},
		{
			name: "where with operators and NOT NULL",
			fields: model.Fields{
				{Name: "country", Value: "COLOMBIA"},
				{Name: "currency_id", Value: 3, ChainingKey: model.Or},
				{Name: "begins_at", Value: "fake", Operator: model.IsNull},
				{Name: "enable", Value: true},
				{Name: "code", Value: []string{"COL", "COP"}, Operator: model.In},
			},
			wantQuery: "WHERE country = @p1 AND currency_id = @p2 OR begins_at IS NULL AND enable = @p3 AND code IN ('COL','COP')",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: "COLOMBIA"}, sql.NamedArg{Name: "p2", Value: 3}, sql.NamedArg{Name: "p3", Value: true}},
		},
		{
			name: "where with aliased",
			fields: model.Fields{
				{Source: "contracts", Name: "employer_id", Value: 777},
				{Source: "contracts", Name: "pay_frequency_id", Value: 2, ChainingKey: model.Or},
				{Source: "contracts", Name: "is_active", Value: true},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE contracts.employer_id = @p1 AND contracts.pay_frequency_id = @p2 OR contracts.is_active = @p3 AND contract_statuses.description ILIKE @p4",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: 777}, sql.NamedArg{Name: "p2", Value: 2}, sql.NamedArg{Name: "p3", Value: true}, sql.NamedArg{Name: "p4", Value: "ACTIVE"}},
		},
		{
			name: "where with aliased where some fields have missing source",
			fields: model.Fields{
				{Name: "employer_id", Value: 19},
				{Name: "pay_frequency_id", Value: 1, ChainingKey: model.Or},
				{Name: "is_active", Value: false},
				{Source: "contract_statuses", Name: "description", Value: "CREATED", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = @p1 AND pay_frequency_id = @p2 OR is_active = @p3 AND contract_statuses.description ILIKE @p4",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: 19}, sql.NamedArg{Name: "p2", Value: 1}, sql.NamedArg{Name: "p3", Value: false}, sql.NamedArg{Name: "p4", Value: "CREATED"}},
		},
		{
			name: "where with group conditions",
			fields: model.Fields{
				{Name: "employer_id", Value: 1},
				{Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Name: "is_active", Value: true, ChainingKey: model.Or},
				{GroupClose: true, Name: "is_staff", Value: false},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = @p1 AND pay_frequency_id = @p2 AND (is_active = @p3 OR is_staff = @p4) AND contract_statuses.description ILIKE @p5",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: 1}, sql.NamedArg{Name: "p2", Value: 2}, sql.NamedArg{Name: "p3", Value: true}, sql.NamedArg{Name: "p4", Value: false}, sql.NamedArg{Name: "p5", Value: "ACTIVE"}},
		},
		{
			name: "where with group conditions and with missing GroupClose key",
			fields: model.Fields{
				{Name: "employer_id", Value: 1},
				{Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Name: "is_active", Value: true, ChainingKey: model.Or},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = @p1 AND pay_frequency_id = @p2 AND (is_active = @p3 OR contract_statuses.description ILIKE @p4)",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: 1}, sql.NamedArg{Name: "p2", Value: 2}, sql.NamedArg{Name: "p3", Value: true}, sql.NamedArg{Name: "p4", Value: "ACTIVE"}},
		},
		{
			name: "where with group conditions and aliases - complex",
			fields: model.Fields{
				{Source: "c", Name: "employer_id", Value: 1},
				{Source: "c", Name: "termination_date", Operator: model.IsNotNull},
				{Source: "c", Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Source: "cs", Name: "description", Operator: model.Ilike, Value: "activo", ChainingKey: model.Or},
				{GroupOpen: true, Source: "cs", Name: "description", Operator: model.Ilike, Value: "creado"},
				{GroupClose: true, Source: "c", Name: "hire_date", Operator: model.LessThanOrEqualTo, Value: fakeDate},
			},
			wantQuery: "WHERE c.employer_id = @p1 AND c.termination_date IS NOT NULL AND c.pay_frequency_id = @p2 AND (cs.description ILIKE @p3 OR (cs.description ILIKE @p4 AND c.hire_date <= @p5))",
			wantArgs:  []interface{}{sql.NamedArg{Name: "p1", Value: 1}, sql.NamedArg{Name: "p2", Value: 2}, sql.NamedArg{Name: "p3", Value: "activo"}, sql.NamedArg{Name: "p4", Value: "creado"}, sql.NamedArg{Name: "p5", Value: "2021-04-28"}},
		},
	}

	for _, tt := range tableTest {
		gotQuery, gotArgs := BuildSQLWhere(tt.fields)
		assert.Equal(t, tt.wantQuery, gotQuery)
		assert.Equal(t, tt.wantArgs, gotArgs)
	}
}

func TestBuildSQLUpdateBy(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		by     string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			by:     "user_id",
			want:   "UPDATE cashboxes SET responsable = @p1, country = @p2, user_id = @p3, account = @p4, d_UpdateDate = now() WHERE user_id = @p5",
		},
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			by:     "responsable",
			want:   "UPDATE cashboxes SET responsable = @p1, country = @p2, user_id = @p3, account = @p4, d_UpdateDate = now() WHERE responsable = @p5",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			by:     "user_id",
			want:   "UPDATE one SET one_field = @p1, d_UpdateDate = now() WHERE user_id = @p2",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLUpdateBy(tt.table, tt.fields, tt.by))
	}
}

func TestColumnsAliased(t *testing.T) {
	tableTest := []struct {
		aliased string
		fields  []string
		want    string
	}{
		{
			aliased: "b",
			fields:  []string{"title", "slug", "content", "poster"},
			want:    "b.id, b.title, b.slug, b.content, b.poster, b.d_InsertDate, b.d_UpdateDate",
		},
		{
			aliased: "nothing",
			fields:  []string{},
			want:    "",
		},
		{
			aliased: "one",
			fields:  []string{"one_field"},
			want:    "one.id, one.one_field, one.d_InsertDate, one.d_UpdateDate",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, ColumnsAliased(tt.fields, tt.aliased))
	}
}

func TestBuildSQLOrderBy(t *testing.T) {
	tests := []struct {
		name  string
		sorts model.SortFields
		want  string
	}{
		{
			name:  "Without sort order specification",
			sorts: model.SortFields{{Name: "id"}, {Name: "begins_at"}},
			want:  "ORDER BY id ASC, begins_at ASC",
		},
		{
			name:  "With sort order specification",
			sorts: model.SortFields{{Name: "id", Order: model.Desc}, {Name: "begins_at", Order: model.Asc}},
			want:  "ORDER BY id DESC, begins_at ASC",
		},
		{
			name:  "One field sort",
			sorts: model.SortFields{{Name: "id"}},
			want:  "ORDER BY id ASC",
		},
		{
			name:  "Without field sorts",
			sorts: model.SortFields{},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSQLOrderBy(tt.sorts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_buildIN(t *testing.T) {
	tableTest := []struct {
		field     model.Field
		wantQuery string
	}{
		{
			field: model.Field{
				Name: "id", Value: []uint{1, 2, 3}, Operator: model.In,
			},
			wantQuery: "id IN (1,2,3)",
		},
		{
			field: model.Field{
				Name: "employee_id", Value: []int{5, 6, 7}, Operator: model.In,
			},
			wantQuery: "employee_id IN (5,6,7)",
		},
		{
			field: model.Field{
				Name: "marital_status", Value: []string{"SINGLE"}, Operator: model.In,
			},
			wantQuery: "marital_status IN ('SINGLE')",
		},
		{
			field: model.Field{
				Name: "employee_id", Value: "fake", Operator: model.In,
			},
			wantQuery: "employee_id = 0",
		},
		{
			field: model.Field{
				Name: "contract_id", Value: []uint{}, Operator: model.In,
			},
			wantQuery: "contract_id = 0",
		},
	}

	for _, tt := range tableTest {
		gotQuery := BuildIN(tt.field)
		assert.Equal(t, tt.wantQuery, gotQuery)
	}
}
