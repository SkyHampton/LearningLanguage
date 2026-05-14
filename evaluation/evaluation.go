package evaluation

import (
	"bytes"
	"fmt"
	"learningLanguage/ast"
	"strings"
)

const (
	INTTYPE = iota
	BOOLTYPE
	STRINGTYPE
	FLOATTYPE
)

type Data struct {
	dataType    int
	intValue    int32
	boolValue   bool
	stringValue string
	floatValue  float32
}

type List struct {
	dataType int
	arr      []Data
}

type Evaluator struct {
	errors      []string
	variableMap map[string]Data
	listMap     map[string]List
}

var variableTypes = map[string]int{
	"int":    INTTYPE,
	"bool":   BOOLTYPE,
	"string": STRINGTYPE,
	"float":  FLOATTYPE,
}

var dataTypeToString = map[int]string{
	INTTYPE:    "INT",
	BOOLTYPE:   "BOOL",
	STRINGTYPE: "STRING",
	FLOATTYPE:  "FLOAT",
}

func New() Evaluator {
	eval := Evaluator{}
	eval.errors = []string{}
	eval.variableMap = make(map[string]Data)
	eval.listMap = make(map[string]List)
	return eval
}

func (e *Evaluator) Errors() []string {
	return e.errors
}

func (e *Evaluator) ResetErrors() {
	e.errors = []string{}
}

func (e *Evaluator) EvaluateProgram(program *ast.Program) string {
	var output bytes.Buffer
	for _, statement := range program.Statements {
		output.WriteString(e.evaluateStatement(statement))
	}

	return output.String()
}

func (e *Evaluator) evaluateStatement(statement ast.Statement) string {
	var output string
	createStmt, ok := statement.(*ast.CreateStatement)
	if ok {
		output = e.evaluateCreateStatement(createStmt)
	}

	setStmt, ok := statement.(*ast.SetStatement)
	if ok {
		output = e.evaluateSetStatement(setStmt)
	}

	ifStmt, ok := statement.(*ast.IfStatement)
	if ok {
		output = e.evaluateIfStatement(ifStmt)
	}

	whileStmt, ok := statement.(*ast.WhileStatement)
	if ok {
		output = e.evaluateWhileStatement(whileStmt)
	}

	countStmt, ok := statement.(*ast.CountStatement)
	if ok {
		output = e.evaluateCountStatement(countStmt)
	}

	structStmt, ok := statement.(*ast.StructStatement)
	if ok {
		output = e.evaluateStructStatement(structStmt)
	}

	appendStmt, ok := statement.(*ast.AppendStatement)
	if ok {
		output = e.evaluateAppendStatement(appendStmt)
	}

	printStmt, ok := statement.(*ast.PrintStatement)
	if ok {
		value, list := e.evaluateExpression(printStmt.Value)
		if len(list.arr) == 0 {
			switch value.dataType {
			case INTTYPE:
				output = fmt.Sprintf("%d", value.intValue)
			case BOOLTYPE:
				output = fmt.Sprintf("%t", value.boolValue)
			case FLOATTYPE:
				output = fmt.Sprintf("%.5g", value.floatValue)
			case STRINGTYPE:
				output = strings.Trim(value.stringValue, "\"")
			}
		} else {
			switch list.dataType {
			case INTTYPE:
				output = "["
				for _, item := range list.arr {
					output += fmt.Sprintf("%d, ", item.intValue)
				}
			case BOOLTYPE:
				output = "["
				for _, item := range list.arr {
					output += fmt.Sprintf("%t, ", item.boolValue)
				}
			case FLOATTYPE:
				output = "["
				for _, item := range list.arr {
					output += fmt.Sprintf("%f, ", item.floatValue)
				}
			case STRINGTYPE:
				output = "["
				for _, item := range list.arr {
					output += fmt.Sprintf("%s, ", item.stringValue)
				}
			}
			output = strings.TrimSuffix(output, ", ")
			output += "]"
		}
		if printStmt.NewLine {
			output += "\n"
		}
	}

	return output
}

func (e *Evaluator) evaluateCreateStatement(statement *ast.CreateStatement) string {
	if statement.Name.IsList {
		e.listMap[statement.Name.Value] = List{dataType: variableTypes[statement.Name.DataType]}
	} else {
		e.variableMap[statement.Name.Value] = Data{dataType: variableTypes[statement.Name.DataType]}
	}
	return ""
}

func (e *Evaluator) evaluateSetStatement(statement *ast.SetStatement) string {
	var name string
	if statement.Name.Attribute != "" {
		name = fmt.Sprintf("%s.%s", statement.Name.Value, statement.Name.Attribute)
	} else {
		name = statement.Name.Value
	}
	_, okVar := e.variableMap[name]
	_, okList := e.listMap[name]
	if !okVar && !okList {
		e.errors = append(e.errors, fmt.Sprintf("Variable %s has not been created.", name))
		return ""
	}
	value, list := e.evaluateExpression(statement.Value)
	if statement.Name.Index != nil && okList && len(list.arr) == 0 {
		indexValue, _ := e.evaluateExpression(statement.Name.Index)
		if indexValue.dataType != INTTYPE {
			e.errors = append(e.errors, "Cannot use non-integer values in list index")
			return ""
		}
		if int(indexValue.intValue) > len(e.listMap[name].arr) || indexValue.intValue < 0 {
			e.errors = append(e.errors, "Index out of range")
			return ""
		}
		if value.dataType != e.listMap[name].dataType {
			e.errors = append(e.errors, "Mismatching list data types")
			return ""
		}
		e.listMap[name].arr[indexValue.intValue-1] = value
		return ""
	}
	if value != (Data{}) && value.dataType != e.variableMap[name].dataType {
		e.errors = append(e.errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
			dataTypeToString[e.variableMap[name].dataType], dataTypeToString[value.dataType]))
		return ""
	}
	if len(list.arr) != 0 && list.dataType != e.listMap[name].dataType {
		e.errors = append(e.errors, "Mismatching list data types")
		return ""
	}
	if len(list.arr) != 0 {
		e.listMap[name] = list
		return ""
	}
	if len(list.arr) == 0 {
		e.variableMap[name] = value
	}
	return ""
}

func (e *Evaluator) evaluateStructStatement(statement *ast.StructStatement) string {
	for _, attribute := range statement.Attributes {
		attributeName := fmt.Sprintf("%s.%s", statement.StructIdent.Value, attribute.Value)
		if statement.Values[attribute.Value] == nil {
			e.variableMap[attributeName] = Data{dataType: variableTypes[attribute.DataType]}
		} else {
			value, _ := e.evaluateExpression(statement.Values[attribute.Value])
			if variableTypes[attribute.DataType] != value.dataType {
				e.errors = append(e.errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
					attribute.DataType, dataTypeToString[value.dataType]))
				return ""
			}
		}
	}
	return ""
}

func (e *Evaluator) evaluateIfStatement(statement *ast.IfStatement) string {
	ret := ""
	conditionData, _ := e.evaluateExpression(statement.Condition)
	if conditionData.dataType != BOOLTYPE {
		e.errors = append(e.errors, "Cannot use non-boolean expression in if statement condition.")
		return ""
	} else {
		if conditionData.boolValue {
			for _, stmt := range statement.IfTrue {
				ret += e.evaluateStatement(stmt)
			}
		} else {
			for _, stmt := range statement.Else {
				ret += e.evaluateStatement(stmt)
			}
		}
	}

	return ret
}

func (e *Evaluator) evaluateWhileStatement(statement *ast.WhileStatement) string {
	ret := ""
	conditionData, _ := e.evaluateExpression(statement.Condition)
	if conditionData.dataType != BOOLTYPE {
		e.errors = append(e.errors, "Cannot use non-boolean expression in while statement condition.")
		return ""
	} else {
		for conditionData.boolValue {
			for _, stmt := range statement.LoopStatements {
				ret += e.evaluateStatement(stmt)
			}
			conditionData, _ = e.evaluateExpression(statement.Condition)
		}
	}

	return ret
}

func (e *Evaluator) evaluateCountStatement(statement *ast.CountStatement) string {
	ret := ""
	fromValue, _ := e.evaluateExpression(statement.From)
	toValue, _ := e.evaluateExpression(statement.To)
	byValue, _ := e.evaluateExpression(statement.By)
	if fromValue.dataType != INTTYPE || toValue.dataType != INTTYPE || byValue.dataType != INTTYPE {
		e.errors = append(e.errors, "Cannot use non-integer values in counting statement from, to, or by.")
		return ret
	}

	counterName := statement.Counter.Value
	e.variableMap[counterName] = fromValue
	for e.variableMap[counterName].intValue <= toValue.intValue {
		for _, stmt := range statement.LoopStatements {
			ret += e.evaluateStatement(stmt)
		}
		tempData := e.variableMap[counterName]
		tempData.intValue += byValue.intValue
		e.variableMap[counterName] = tempData
	}

	return ret
}

func (e *Evaluator) evaluateAppendStatement(statement *ast.AppendStatement) string {
	listName := statement.List.Value
	list, ok := e.listMap[listName]
	if !ok {
		e.errors = append(e.errors, fmt.Sprintf("List %s does not exist.", listName))
		return ""
	}

	appendValue, _ := e.evaluateExpression(statement.Value)
	if appendValue == (Data{}) {
		e.errors = append(e.errors, "Invalid append data")
		return ""
	}

	if appendValue.dataType != list.dataType {
		e.errors = append(e.errors, "Mismatching append data type")
		return ""
	}

	list.arr = append(list.arr, appendValue)
	e.listMap[listName] = list

	return ""
}

func (e *Evaluator) evaluateLenExpression(statement *ast.LengthExpression) Data {
	list, ok := e.listMap[statement.List.Value]
	if !ok {
		e.errors = append(e.errors, fmt.Sprintf("List %s does not exist.", statement.List.Value))
	}

	return Data{dataType: INTTYPE, intValue: int32(len(list.arr))}
}

func (e *Evaluator) evaluateIntLit(expression *ast.IntegerLiteral) Data {
	return Data{dataType: INTTYPE, intValue: expression.Value}
}

func (e *Evaluator) evaluateBoolLit(expression *ast.BooleanLiteral) Data {
	return Data{dataType: BOOLTYPE, boolValue: expression.Value}
}

func (e *Evaluator) evaluateFloatLit(expression *ast.FloatLiteral) Data {
	return Data{dataType: FLOATTYPE, floatValue: expression.Value}
}

func (e *Evaluator) evaluateStringLit(expression *ast.StringLiteral) Data {
	return Data{dataType: STRINGTYPE, stringValue: expression.Value}
}

func (e *Evaluator) evaluateListLit(expression *ast.ListLiteral) List {
	expList := expression.List
	var dataList List
	for _, exp := range expList {
		data, _ := e.evaluateExpression(exp)
		if data.dataType != dataList.dataType {
			e.errors = append(e.errors, "Mismatching list elements")
		}
		dataList.arr = append(dataList.arr, data)
	}

	return dataList
}

func (e *Evaluator) evaluateExpression(expression ast.Expression) (Data, List) {
	var value Data = Data{}
	var list List = List{}

	intLit, ok := expression.(*ast.IntegerLiteral)
	if ok {
		value = e.evaluateIntLit(intLit)
	}

	boolLit, ok := expression.(*ast.BooleanLiteral)
	if ok {
		value = e.evaluateBoolLit(boolLit)
	}

	floatLiteral, ok := expression.(*ast.FloatLiteral)
	if ok {
		value = e.evaluateFloatLit(floatLiteral)
	}

	strLit, ok := expression.(*ast.StringLiteral)
	if ok {
		value = e.evaluateStringLit(strLit)
	}

	listLit, ok := expression.(*ast.ListLiteral)
	if ok {
		list = e.evaluateListLit(listLit)
	}

	lenExp, ok := expression.(*ast.LengthExpression)
	if ok {
		value = e.evaluateLenExpression(lenExp)
	}

	identifierExp, ok := expression.(*ast.Identifier)
	if ok {
		value, list = e.evaluateIdentifier(identifierExp)
	}

	prefixExp, ok := expression.(*ast.PrefixExpression)
	if ok {
		value = e.evaluatePrefixExp(prefixExp)
	}

	infixExp, ok := expression.(*ast.InfixExpression)
	if ok {
		value = e.evaluateInfixExp(infixExp)
	}

	return value, list
}

func (e *Evaluator) evaluateIdentifier(identifier *ast.Identifier) (Data, List) {
	var name string
	var list List
	if identifier.Attribute != "" {
		name = fmt.Sprintf("%s.%s", identifier.Value, identifier.Attribute)
	} else {
		name = identifier.Value
	}
	value, ok := e.variableMap[name]
	if !ok {
		list, ok = e.listMap[name]
		if !ok {
			e.errors = append(e.errors, fmt.Sprintf("Variable %s does not exist.", identifier.Value))
			return Data{}, List{}
		} else {
			if identifier.Index != nil {
				indexValue, _ := e.evaluateExpression(identifier.Index)
				if indexValue.dataType != INTTYPE {
					e.errors = append(e.errors, "Cannot use non-integer values in list index")
					return Data{}, List{}
				}
				if int(indexValue.intValue) > len(e.listMap[name].arr) || indexValue.intValue < 0 {
					e.errors = append(e.errors, "Index out of range")
					return Data{}, List{}
				}
				return list.arr[indexValue.intValue-1], List{}
			}
			return Data{}, list
		}
	}

	return value, List{}
}

func (e *Evaluator) evaluatePrefixExp(expression *ast.PrefixExpression) Data {
	value, _ := e.evaluateExpression(expression.Right)
	switch expression.Operator {
	case "-":
		value.intValue = -1 * value.intValue
	case "!":
		value.boolValue = !value.boolValue
	}

	return value
}

func (e *Evaluator) evaluateInfixExp(expression *ast.InfixExpression) Data {
	leftValue, _ := e.evaluateExpression(expression.Left)
	rightValue, _ := e.evaluateExpression(expression.Right)
	var retValue Data

	switch expression.Operator {
	case "+":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left + right)
		} else {
			retValue.floatValue = left + right
		}

	case "-":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left - right)
		} else {
			retValue.floatValue = left - right
		}

	case "*":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left * right)
		} else {
			retValue.floatValue = left * right
		}

	case "/":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left / right)
		} else {
			retValue.floatValue = left / right
		}
	case "==":
		retValue.dataType = BOOLTYPE
		switch leftValue.dataType {
		case BOOLTYPE:
			retValue.boolValue = leftValue.boolValue == rightValue.boolValue
		case INTTYPE:
			retValue.boolValue = leftValue.intValue == rightValue.intValue
		case FLOATTYPE:
			retValue.boolValue = leftValue.floatValue == rightValue.floatValue
		case STRINGTYPE:
			retValue.boolValue = leftValue.stringValue == rightValue.stringValue
		}
	case "!=":
		retValue.dataType = BOOLTYPE
		switch leftValue.dataType {
		case BOOLTYPE:
			retValue.boolValue = leftValue.boolValue != rightValue.boolValue
		case INTTYPE:
			retValue.boolValue = leftValue.intValue != rightValue.intValue
		case FLOATTYPE:
			retValue.boolValue = leftValue.floatValue != rightValue.floatValue
		case STRINGTYPE:
			retValue.boolValue = leftValue.stringValue != rightValue.stringValue
		}
	case ">":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			e.errors = append(e.errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue > rightValue.intValue
	case ">=":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			e.errors = append(e.errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue >= rightValue.intValue
	case "<":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			e.errors = append(e.errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue < rightValue.intValue
	case "<=":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			e.errors = append(e.errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue <= rightValue.intValue
	case "or":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != BOOLTYPE && rightValue.dataType != BOOLTYPE {
			e.errors = append(e.errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue || rightValue.boolValue
	case "and":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != BOOLTYPE && rightValue.dataType != BOOLTYPE {
			e.errors = append(e.errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue && rightValue.boolValue
	}
	return retValue
}
