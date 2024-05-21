package core_database

import (
	"fmt"
	"testing"
)

func TestCreateColumnsQuery(t *testing.T) {
	type args struct {
		fieldsMap any
	}

	avatarId := "1234"

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "без параметров",
			args: args{
				fieldsMap: UpdateUserOptionsFields{},
			},
			want: "",
		},
		{
			name: "один строковый параметр (первый)",
			args: args{
				fieldsMap: UpdateUserOptionsFields{
					FirstName: CreateField("name"),
				},
			},
			want: "first_name = 'name'",
		},
		{
			name: "один строковый параметр (не первый)",
			args: args{
				fieldsMap: UpdateUserOptionsFields{
					AvatarId: CreateField(&avatarId),
				},
			},
			want: fmt.Sprintf("avatar_id = '%s'", avatarId),
		},
		{
			name: "несколько строковых параметров",
			args: args{
				fieldsMap: UpdateUserOptionsFields{
					FirstName: CreateField("name"),
					AvatarId:  CreateField(&avatarId),
				},
			},
			want: fmt.Sprintf("first_name = 'name', avatar_id = '%s'", avatarId),
		},
		{
			name: "строковый параметр + массив байтов",
			args: args{
				fieldsMap: UpdateTaskOptionsFields{
					Name:     CreateField("name"),
					Elements: CreateField([]byte{4, 0, 0, 0, 4, 0}),
				},
			},
			want: "name = 'name', elements = '\\x040000000400'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateColumnsQuery(tt.args.fieldsMap); got != tt.want {
				t.Errorf("CreateColumnsQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
