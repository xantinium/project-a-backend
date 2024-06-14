package core_resources

import (
	"reflect"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/xantinium/project-a-backend/api/tasks"
)

func Test_taskProcessing(t *testing.T) {
	tests := []struct {
		name string
		task TaskType
	}{
		{
			name: "без элементов",
			task: TaskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				IsPrivate:   true,
				Elements:    []taskElementType{},
			},
		},
		{
			name: "один элемент: текст",
			task: TaskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				IsPrivate:   true,
				Elements: []taskElementType{
					{
						Hash: "1",
						Type: tasks.ElementsTypesTEXT,
						textElementData: textElementType{
							Body: "текст",
						},
					},
				},
			},
		},
		{
			name: "один элемент: выбор варианта",
			task: TaskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				IsPrivate:   true,
				Elements: []taskElementType{
					{
						Hash: "1",
						Type: tasks.ElementsTypesCHOICE,
						choiceElementData: choiceElementType{
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
		},
		{
			name: "один элемент: множественный выбор",
			task: TaskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				IsPrivate:   true,
				Elements: []taskElementType{
					{
						Hash: "1",
						Type: tasks.ElementsTypesMULTI_CHOICE,
						multiChoiceElementData: multiChoiceElementType{
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
		},
		{
			name: "один элемент: соответствие элементов",
			task: TaskType{
				Id:          4,
				Name:        "название",
				Description: "описание",
				IsPrivate:   true,
				Elements: []taskElementType{
					{
						Hash: "1",
						Type: tasks.ElementsTypesRELATIONS,
						relationsElementData: relationsElementType{
							Description: "описание",
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
		},
	}
	for _, tt := range tests {
		b := &flatbuffers.Builder{}
		offset := SerializeTask(b, tt.task)
		b.Finish(offset)

		t.Run(tt.name, func(t *testing.T) {
			if got := DeserializeTask(b.FinishedBytes()); !reflect.DeepEqual(got, tt.task) {
				t.Errorf("result task = %v, want %v", got, tt.task)
			}
		})
	}
}

func Test_elementsProcessing(t *testing.T) {
	tests := []struct {
		name     string
		elements []taskElementType
	}{
		{
			name: "один элемент: выбор варианта",
			elements: []taskElementType{
				{
					Hash: "1",
					Type: tasks.ElementsTypesCHOICE,
					choiceElementData: choiceElementType{
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
			name: "несколько элементов: текст + выбор варианта",
			elements: []taskElementType{
				{
					Hash: "1",
					Type: tasks.ElementsTypesTEXT,
					textElementData: textElementType{
						Body: "Текст",
					},
				},
				{
					Hash: "2",
					Type: tasks.ElementsTypesCHOICE,
					choiceElementData: choiceElementType{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeserializeElementsFromBytes(SerializeElements(tt.elements)); !reflect.DeepEqual(got, tt.elements) {
				t.Errorf("result task = %v, want %v", got, tt.elements)
			}
		})
	}
}
