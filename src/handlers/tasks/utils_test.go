package tasks_handler

import (
	"reflect"
	"testing"
)

func Test_taskProcessing(t *testing.T) {
	tests := []struct {
		name string
		task taskType
	}{
		{
			name: "без элементов",
			task: taskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				Elements:    []interface{}{},
			},
		},
		{
			name: "один элемент: текст",
			task: taskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				Elements: []interface{}{
					textElementType{
						Hash: "1",
						Body: "текст",
					},
				},
			},
		},
		{
			name: "один элемент: выбор варианта",
			task: taskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				Elements: []interface{}{
					choiceElementType{
						Hash:        "1",
						Description: "описание",
						Items: []choiceElementItemType{
							{
								Hash:     "1",
								Text:     "вариант 1",
								Selected: false,
							},
							{
								Hash:     "2",
								Text:     "вариант 2",
								Selected: true,
							},
						},
					},
				},
			},
		},
		{
			name: "один элемент: множественный выбор",
			task: taskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				Elements: []interface{}{
					multiChoiceElementType{
						Hash:        "1",
						Description: "описание",
						Items: []choiceElementItemType{
							{
								Hash:     "1",
								Text:     "вариант 1",
								Selected: true,
							},
							{
								Hash:     "2",
								Text:     "вариант 2",
								Selected: false,
							},
							{
								Hash:     "3",
								Text:     "вариант 3",
								Selected: true,
							},
						},
					},
				},
			},
		},
		{
			name: "один элемент: соответствие элементов",
			task: taskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				Elements: []interface{}{
					relationsElementType{
						Hash: "1",
						LeftItems: []RelationItemType{
							{
								Hash: "1",
								Text: "текст 1",
							},
							{
								Hash: "2",
								Text: "текст 2",
							},
						},
						RightItems: []RelationItemType{
							{
								Hash: "1",
								Text: "текст 1",
							},
						},
						Relations: []RelationType{
							{
								Left:  "1",
								Right: "1",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deserializeTask(serializeTask(tt.task)); !reflect.DeepEqual(got, tt.task) {
				t.Errorf("result task = %v, want %v", got, tt.task)
			}
		})
	}
}
