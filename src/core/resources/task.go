package core_resources

import (
	flatbuffers "github.com/google/flatbuffers/go"
	api_tasks "github.com/xantinium/project-a-backend/api/tasks"
)

type textElementType struct {
	Hash string
	Body string
}

type choiceElementItemType struct {
	Hash     string
	Text     string
	Selected bool
}

type choiceElementType struct {
	Hash        string
	Description string
	Items       []choiceElementItemType
}

type multiChoiceElementType = struct {
	Hash        string
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
	Hash       string
	LeftItems  []RelationItemType
	RightItems []RelationItemType
	Relations  []RelationType
}

type taskElementType = interface{}

type taskType struct {
	Id          int
	Name        string
	Description string
	Elements    []taskElementType
}

type readElementFunc func(*api_tasks.Element, int) bool

func getElementBytes(element api_tasks.Element) []byte {
	data := make([]byte, 0)

	for i := 0; i < element.DataLength(); i++ {
		data = append(data, byte(element.Data(i)))
	}

	return data
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

func deserializeTextElement(element api_tasks.Element) textElementType {
	textElement := api_tasks.GetRootAsTextElement(getElementBytes(element), 0)

	return textElementType{
		Hash: string(element.Hash()),
		Body: string(textElement.Body()),
	}
}

func deserializeChoiceElement(element api_tasks.Element) choiceElementType {
	choiceElement := api_tasks.GetRootAsChoiceElement(getElementBytes(element), 0)

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

	return choiceElementType{
		Hash:        string(element.Hash()),
		Description: string(choiceElement.Description()),
		Items:       items,
	}
}

func deserializeMultiChoiceElement(element api_tasks.Element) multiChoiceElementType {
	return deserializeChoiceElement(element)
}

func deserializeRelationsElement(element api_tasks.Element) relationsElementType {
	relationsElement := api_tasks.GetRootAsRelationsElement(getElementBytes(element), 0)

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

	return relationsElementType{
		Hash:       string(element.Hash()),
		LeftItems:  leftItems,
		RightItems: rightItems,
		Relations:  relations,
	}
}

func DeserializeTask(data []byte) taskType {
	task := api_tasks.GetRootAsTask(data, 0)

	elements := DeserializeElements(task.Elements, task.ElementsLength())

	return taskType{
		Id:          int(task.Id()),
		Name:        string(task.Name()),
		Description: string(task.Description()),
		Elements:    elements,
	}
}

func serializeTextElement(element textElementType) []byte {
	b := &flatbuffers.Builder{}

	body := b.CreateString(element.Body)

	api_tasks.TextElementStart(b)
	api_tasks.TextElementAddBody(b, body)
	offset := api_tasks.TextElementEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}

func serializeChoiceElement(element choiceElementType) []byte {
	b := &flatbuffers.Builder{}

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
	offset = api_tasks.ChoiceElementEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}

func serializeMultiChoiceElement(element choiceElementType) []byte {
	return serializeChoiceElement(element)
}

func serializeRelationsElement(element relationsElementType) []byte {
	b := &flatbuffers.Builder{}

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
	api_tasks.RelationsElementAddLeftItems(b, leftItems)
	api_tasks.RelationsElementAddRightItems(b, rightItems)
	api_tasks.RelationsElementAddRelations(b, relations)
	offset := api_tasks.RelationsElementEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}

func SerializeElements(elements []taskElementType) []byte {
	b := &flatbuffers.Builder{}

	offsets := make([]flatbuffers.UOffsetT, 0, len(elements))

	for _, el := range elements {
		var data []byte

		switch element := el.(type) {
		case textElementType:
			data = serializeTextElement(element)
		case choiceElementType:
			data = serializeChoiceElement(element)
		case multiChoiceElementType:
			data = serializeMultiChoiceElement(element)
		case relationsElementType:
			data = serializeRelationsElement(element)
		}

		offset := b.CreateByteVector(data)

		offsets = append(offsets, offset)
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks.TaskElementsStart(b)
	api_tasks.TaskElementsAddElements(b, offset)
	offset = api_tasks.TaskElementsEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}

func SerializeTask(task taskType) []byte {
	b := &flatbuffers.Builder{}

	name := b.CreateString(task.Name)
	description := b.CreateString(task.Description)

	offsets := make([]flatbuffers.UOffsetT, 0, len(task.Elements))

	for _, el := range task.Elements {
		var hash flatbuffers.UOffsetT
		var elType api_tasks.ElementsTypes
		var data []byte

		switch element := el.(type) {
		case textElementType:
			hash = b.CreateString(element.Hash)
			elType = api_tasks.ElementsTypesTEXT
			data = serializeTextElement(element)
		case choiceElementType:
			hash = b.CreateString(element.Hash)
			elType = api_tasks.ElementsTypesCHOICE
			data = serializeChoiceElement(element)
		case multiChoiceElementType:
			hash = b.CreateString(element.Hash)
			elType = api_tasks.ElementsTypesMULTI_CHOICE
			data = serializeMultiChoiceElement(element)
		case relationsElementType:
			hash = b.CreateString(element.Hash)
			elType = api_tasks.ElementsTypesRELATIONS
			data = serializeRelationsElement(element)
		}

		offset := b.CreateByteVector(data)

		api_tasks.ElementStart(b)
		api_tasks.ElementAddHash(b, hash)
		api_tasks.ElementAddType(b, elType)
		api_tasks.ElementAddData(b, offset)
		offset = api_tasks.ElementEnd(b)
		offsets = append(offsets, offset)
	}

	offset := b.CreateVectorOfTables(offsets)

	api_tasks.TaskStart(b)
	api_tasks.TaskAddId(b, uint32(task.Id))
	api_tasks.TaskAddName(b, name)
	api_tasks.TaskAddDescription(b, description)
	api_tasks.TaskAddElements(b, offset)
	offset = api_tasks.TaskEnd(b)
	b.Finish(offset)

	return b.FinishedBytes()
}
