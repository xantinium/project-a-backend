package core_resources

import (
	flatbuffers "github.com/google/flatbuffers/go"
	api_tasks "github.com/xantinium/project-a-backend/api/tasks"
)

type textElementType struct {
	Body string
}

type choiceElementItemType struct {
	Hash     string
	Text     string
	Selected bool
}

type choiceElementType struct {
	Description string
	Items       []choiceElementItemType
}

type multiChoiceElementType = struct {
	Description string
	Items       []choiceElementItemType
}

type RelationItemType struct {
	Hash string
	Text string
}

type RelationType struct {
	Left  string
	Right string
}

type relationsElementType struct {
	Description string
	LeftItems   []RelationItemType
	RightItems  []RelationItemType
	Relations   []RelationType
}

type taskElementType struct {
	Hash                   string
	Type                   api_tasks.ElementsTypes
	textElementData        textElementType
	choiceElementData      choiceElementType
	multiChoiceElementData multiChoiceElementType
	relationsElementData   relationsElementType
}

type TaskType struct {
	Id          int
	Name        string
	Description string
	IsPrivate   bool
	Elements    []taskElementType
}

type readElementFunc func(*api_tasks.Element, int) bool

func deserializeTextElement(element api_tasks.Element) taskElementType {
	textElement := element.TextElementData(nil)

	return taskElementType{
		Hash: string(element.Hash()),
		Type: api_tasks.ElementsTypesTEXT,
		textElementData: textElementType{
			Body: string(textElement.Body()),
		},
	}
}

func deserializeChoiceElement(element api_tasks.Element) taskElementType {
	choiceElement := element.ChoiceElementData(nil)

	items := make([]choiceElementItemType, 0, choiceElement.ItemsLength())

	for i := 0; i < choiceElement.ItemsLength(); i++ {
		var item api_tasks.ChoiceElementItem
		choiceElement.Items(&item, i)

		items = append(items, choiceElementItemType{
			Hash:     string(item.Hash()),
			Text:     string(item.Text()),
			Selected: item.Selected(),
		})
	}

	return taskElementType{
		Hash: string(element.Hash()),
		Type: api_tasks.ElementsTypesCHOICE,
		choiceElementData: choiceElementType{
			Description: string(choiceElement.Description()),
			Items:       items,
		},
	}
}

func deserializeMultiChoiceElement(element api_tasks.Element) taskElementType {
	choiceElement := element.MultiChoiceElementData(nil)

	items := make([]choiceElementItemType, 0, choiceElement.ItemsLength())

	for i := 0; i < choiceElement.ItemsLength(); i++ {
		var item api_tasks.ChoiceElementItem
		choiceElement.Items(&item, i)

		items = append(items, choiceElementItemType{
			Hash:     string(item.Hash()),
			Text:     string(item.Text()),
			Selected: item.Selected(),
		})
	}

	return taskElementType{
		Hash: string(element.Hash()),
		Type: api_tasks.ElementsTypesMULTI_CHOICE,
		multiChoiceElementData: multiChoiceElementType{
			Description: string(choiceElement.Description()),
			Items:       items,
		},
	}
}

func deserializeRelationsElement(element api_tasks.Element) taskElementType {
	relationsElement := element.RelationsElementData(nil)

	leftItems := make([]RelationItemType, 0, relationsElement.LeftItemsLength())
	rightItems := make([]RelationItemType, 0, relationsElement.RightItemsLength())
	relations := make([]RelationType, 0, relationsElement.RightItemsLength())

	for i := 0; i < relationsElement.LeftItemsLength(); i++ {
		var item api_tasks.RelationItem
		relationsElement.LeftItems(&item, i)

		leftItems = append(leftItems, RelationItemType{
			Hash: string(item.Hash()),
			Text: string(item.Text()),
		})
	}

	for i := 0; i < relationsElement.RightItemsLength(); i++ {
		var item api_tasks.RelationItem
		relationsElement.RightItems(&item, i)

		rightItems = append(rightItems, RelationItemType{
			Hash: string(item.Hash()),
			Text: string(item.Text()),
		})
	}

	for i := 0; i < relationsElement.RelationsLength(); i++ {
		var relation api_tasks.Relation
		relationsElement.Relations(&relation, i)

		relations = append(relations, RelationType{
			Left:  string(relation.Left()),
			Right: string(relation.Right()),
		})
	}

	return taskElementType{
		Hash: string(element.Hash()),
		Type: api_tasks.ElementsTypesRELATIONS,
		relationsElementData: relationsElementType{

			Description: string(relationsElement.Description()),
			LeftItems:   leftItems,
			RightItems:  rightItems,
			Relations:   relations,
		},
	}
}

func DeserializeElements(readElement readElementFunc, elementsNum int) []taskElementType {
	elements := make([]taskElementType, 0, elementsNum)

	for i := 0; i < elementsNum; i++ {
		var element api_tasks.Element

		readElement(&element, i)

		switch element.Type() {
		case api_tasks.ElementsTypesTEXT:
			elements = append(elements, deserializeTextElement(element))
		case api_tasks.ElementsTypesCHOICE:
			elements = append(elements, deserializeChoiceElement(element))
		case api_tasks.ElementsTypesMULTI_CHOICE:
			elements = append(elements, deserializeMultiChoiceElement(element))
		case api_tasks.ElementsTypesRELATIONS:
			elements = append(elements, deserializeRelationsElement(element))
		}
	}

	return elements
}

func DeserializeElementsFromBytes(data []byte) []taskElementType {
	taskElements := api_tasks.GetRootAsTaskElements(data, 0)

	return DeserializeElements(taskElements.Elements, taskElements.ElementsLength())
}

func DeserializeTask(data []byte) TaskType {
	task := api_tasks.GetRootAsTask(data, 0)

	elements := DeserializeElements(task.Elements, task.ElementsLength())

	return TaskType{
		Id:          int(task.Id()),
		Name:        string(task.Name()),
		Description: string(task.Description()),
		IsPrivate:   task.IsPrivate(),
		Elements:    elements,
	}
}

func serializeTextElement(b *flatbuffers.Builder, element textElementType) flatbuffers.UOffsetT {
	body := b.CreateString(element.Body)

	api_tasks.TextElementStart(b)
	api_tasks.TextElementAddBody(b, body)
	return api_tasks.TextElementEnd(b)
}

func serializeChoiceElement(b *flatbuffers.Builder, element choiceElementType) flatbuffers.UOffsetT {
	description := b.CreateString(element.Description)

	offsets := make([]flatbuffers.UOffsetT, 0, len(element.Items))

	for _, item := range element.Items {
		hash := b.CreateString(item.Hash)
		text := b.CreateString(item.Text)

		api_tasks.ChoiceElementItemStart(b)
		api_tasks.ChoiceElementItemAddHash(b, hash)
		api_tasks.ChoiceElementItemAddText(b, text)
		api_tasks.ChoiceElementItemAddSelected(b, item.Selected)
		offset := api_tasks.ChoiceElementItemEnd(b)
		offsets = append(offsets, offset)
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks.ChoiceElementStart(b)
	api_tasks.ChoiceElementAddDescription(b, description)
	api_tasks.ChoiceElementAddItems(b, offset)
	return api_tasks.ChoiceElementEnd(b)
}

func serializeMultiChoiceElement(b *flatbuffers.Builder, element choiceElementType) flatbuffers.UOffsetT {
	return serializeChoiceElement(b, element)
}

func serializeRelationsElement(b *flatbuffers.Builder, element relationsElementType) flatbuffers.UOffsetT {
	description := b.CreateString(element.Description)

	leftItemsOffsets := make([]flatbuffers.UOffsetT, 0, len(element.LeftItems))
	rightItemsOffsets := make([]flatbuffers.UOffsetT, 0, len(element.RightItems))
	relationsOffsets := make([]flatbuffers.UOffsetT, 0, len(element.Relations))

	for _, item := range element.LeftItems {
		hash := b.CreateString(item.Hash)
		text := b.CreateString(item.Text)

		api_tasks.RelationItemStart(b)
		api_tasks.RelationItemAddHash(b, hash)
		api_tasks.RelationItemAddText(b, text)
		offset := api_tasks.RelationItemEnd(b)
		leftItemsOffsets = append(leftItemsOffsets, offset)
	}

	for _, item := range element.RightItems {
		hash := b.CreateString(item.Hash)
		text := b.CreateString(item.Text)

		api_tasks.RelationItemStart(b)
		api_tasks.RelationItemAddHash(b, hash)
		api_tasks.RelationItemAddText(b, text)
		offset := api_tasks.RelationItemEnd(b)
		rightItemsOffsets = append(rightItemsOffsets, offset)
	}

	for _, item := range element.Relations {
		left := b.CreateString(item.Left)
		right := b.CreateString(item.Right)

		api_tasks.RelationStart(b)
		api_tasks.RelationAddLeft(b, left)
		api_tasks.RelationAddRight(b, right)
		offset := api_tasks.RelationEnd(b)
		relationsOffsets = append(relationsOffsets, offset)
	}

	leftItems := b.CreateVectorOfTables(leftItemsOffsets)
	rightItems := b.CreateVectorOfTables(rightItemsOffsets)
	relations := b.CreateVectorOfTables(relationsOffsets)

	api_tasks.RelationsElementStart(b)
	api_tasks.RelationsElementAddDescription(b, description)
	api_tasks.RelationsElementAddLeftItems(b, leftItems)
	api_tasks.RelationsElementAddRightItems(b, rightItems)
	api_tasks.RelationsElementAddRelations(b, relations)
	return api_tasks.RelationsElementEnd(b)
}

func serializeElement(b *flatbuffers.Builder, element *taskElementType) flatbuffers.UOffsetT {
	hash := b.CreateString(element.Hash)

	var addElementData func()

	switch element.Type {
	case api_tasks.ElementsTypesTEXT:
		textElement := serializeTextElement(b, element.textElementData)
		addElementData = func() {
			api_tasks.ElementAddTextElementData(b, textElement)
		}
	case api_tasks.ElementsTypesCHOICE:
		choiceElement := serializeChoiceElement(b, element.choiceElementData)
		addElementData = func() {
			api_tasks.ElementAddChoiceElementData(b, choiceElement)
		}
	case api_tasks.ElementsTypesMULTI_CHOICE:
		multiChoiceElement := serializeMultiChoiceElement(b, element.multiChoiceElementData)
		addElementData = func() {
			api_tasks.ElementAddMultiChoiceElementData(b, multiChoiceElement)
		}
	case api_tasks.ElementsTypesRELATIONS:
		relationsElement := serializeRelationsElement(b, element.relationsElementData)
		addElementData = func() {
			api_tasks.ElementAddRelationsElementData(b, relationsElement)
		}
	}

	api_tasks.ElementStart(b)
	api_tasks.ElementAddHash(b, hash)
	api_tasks.ElementAddType(b, element.Type)
	addElementData()
	return api_tasks.ElementEnd(b)
}

func SerializeElements(elements []taskElementType) []byte {
	b := &flatbuffers.Builder{}

	offsets := make([]flatbuffers.UOffsetT, 0, len(elements))

	for _, el := range elements {
		offsets = append(offsets, serializeElement(b, &el))
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks.TaskElementsStart(b)
	api_tasks.TaskElementsAddElements(b, offset)
	offset = api_tasks.TaskElementsEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}

func SerializeTask(b *flatbuffers.Builder, task TaskType) flatbuffers.UOffsetT {
	name := b.CreateString(task.Name)
	description := b.CreateString(task.Description)

	offsets := make([]flatbuffers.UOffsetT, 0, len(task.Elements))

	for _, el := range task.Elements {
		offsets = append(offsets, serializeElement(b, &el))
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks.TaskStart(b)
	api_tasks.TaskAddId(b, uint32(task.Id))
	api_tasks.TaskAddName(b, name)
	api_tasks.TaskAddDescription(b, description)
	api_tasks.TaskAddIsPrivate(b, task.IsPrivate)
	api_tasks.TaskAddElements(b, offset)

	return api_tasks.TaskEnd(b)
}
