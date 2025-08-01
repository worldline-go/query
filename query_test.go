package query

import (
	"net/url"
	"reflect"
	"testing"
)

func ptr(i uint64) *uint64 {
	return &i
}

func TestParseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    *Query
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				query: "name=foo,a%2Fb&age=1&sort=-age&limit=10&offset=5&fields=id,name",
			},
			want: &Query{
				Select: []string{"id", "name"},
				Where: []Expression{
					newExpressionCmp(OperatorIn, "name", []string{"foo", "a/b"}),
					newExpressionCmp(OperatorEq, "age", "1"),
				},
				Sort: []ExpressionSort{
					{
						Field: "age",
						Desc:  true,
					},
				},
				Offset: ptr(5),
				Limit:  ptr(10),
			},
			wantErr: false,
		},
		{
			name: "test 2",
			args: args{
				query: "name=foo|nick=bar&age=1&sort=age&limit=10",
			},
			want: &Query{
				Where: []Expression{
					ExpressionLogic{
						Operator: OperatorOr,
						List: []Expression{
							newExpressionCmp(OperatorEq, "name", "foo"),
							newExpressionCmp(OperatorEq, "nick", "bar"),
						},
					},
					newExpressionCmp(OperatorEq, "age", "1"),
				},
				Sort: []ExpressionSort{
					{
						Field: "age",
						Desc:  false,
					},
				},
				Limit: ptr(10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseQuery() error = %s, wantErr %#v", err, tt.wantErr)
			}

			if got.Limit != nil && tt.want.Limit != nil {
				if *got.Limit != *tt.want.Limit {
					t.Fatalf("ParseQuery() Limit = %d, want %d", *got.Limit, *tt.want.Limit)
				}
			}
			if got.Offset != nil && tt.want.Offset != nil {
				if *got.Offset != *tt.want.Offset {
					t.Fatalf("ParseQuery() Offset = %d, want %d", *got.Offset, *tt.want.Offset)
				}
			}

			if !reflect.DeepEqual(got.Where, tt.want.Where) {
				t.Fatalf("ParseQuery() = \n%#v\n, want \n%#v\n", got.Where, tt.want.Where)
			}

			if !reflect.DeepEqual(got.Sort, tt.want.Sort) {
				t.Fatalf("ParseQuery() = \n%#v\n, want \n%#v\n", got.Sort, tt.want.Sort)
			}

			if !reflect.DeepEqual(got.Select, tt.want.Select) {
				t.Fatalf("ParseQuery() = \n%#v\n, want \n%#v\n", got.Select, tt.want.Select)
			}
		})
	}
}

func Test_URLQuery(t *testing.T) {
	testURL := "http://example.com?name=foo|nick=foo&age=1&sort=age&limit=10&offset=5&fields=id,name#test"
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	if parsedURL.RawQuery != "name=foo|nick=foo&age=1&sort=age&limit=10&offset=5&fields=id,name" {
		t.Fatalf("parsed URL query does not match expected value: %s", parsedURL.RawQuery)
	}
}
