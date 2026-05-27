package evaluation

import (
	"bytes"
	"fmt"
	"learningLanguage/ast"
	"math"
	"slices"
	"strings"
)

// enumeration of data types
const (
	INTTYPE = iota
	BOOLTYPE
	STRINGTYPE
	FLOATTYPE
)

// data structure, used to return multiple data types in one function return
type Data struct {
	dataType    int
	intValue    int32
	boolValue   bool
	stringValue string
	floatValue  float32
}

// list of data
type List struct {
	dataType int
	arr      []Data
}

// evaluator object, contains errors, a map of variables, and a map of lists
type Evaluator struct {
	errors      []string
	variableMap map[string]Data
	listMap     map[string]List
}

var numericDataTypes = []int{INTTYPE, FLOATTYPE}

// variable type mapping to convert data type in AST node to enumerated data type
var variableTypes = map[string]int{
	"int":    INTTYPE,
	"bool":   BOOLTYPE,
	"string": STRINGTYPE,
	"float":  FLOATTYPE,
}

// mapping of enumerated data type to string for printing purposes
var dataTypeToString = map[int]string{
	INTTYPE:    "INT",
	BOOLTYPE:   "BOOL",
	STRINGTYPE: "STRING",
	FLOATTYPE:  "FLOAT",
}

// evaluator constructor
func New() Evaluator {
	eval := Evaluator{}
	eval.errors = []string{}
	eval.variableMap = make(map[string]Data)
	eval.listMap = make(map[string]List)
	return eval
}

// get method for evaluator errors
func (e *Evaluator) Errors() []string {
	return e.errors
}

// helper function to reset errors, mainly used by the REPL to ensure that errors don't softlock it
func (e *Evaluator) ResetErrors() {
	e.errors = []string{}
}

// to evaluate a program, evaluate each statement in it
func (e *Evaluator) EvaluateProgram(program *ast.Program) string {
	var output bytes.Buffer
	for _, statement := range program.Statements {
		output.WriteString(e.evaluateStatement(statement))
	}

	if len(e.errors) > 0 {
		return "Errors were found, no output."
	}
	return output.String()
}

// evaluate each statement type
func (e *Evaluator) evaluateStatement(statement ast.Statement) string {
	var output string
	// evaluate create statement
	createStmt, ok := statement.(*ast.CreateStatement)
	if ok {
		output = e.evaluateCreateStatement(createStmt)
	}

	// evaluate set statement
	setStmt, ok := statement.(*ast.SetStatement)
	if ok {
		output = e.evaluateSetStatement(setStmt)
	}

	// evaluate if statement
	ifStmt, ok := statement.(*ast.IfStatement)
	if ok {
		output = e.evaluateIfStatement(ifStmt)
	}

	// evaluate while statement
	whileStmt, ok := statement.(*ast.WhileStatement)
	if ok {
		output = e.evaluateWhileStatement(whileStmt)
	}

	// evaluate count statement
	countStmt, ok := statement.(*ast.CountStatement)
	if ok {
		output = e.evaluateCountStatement(countStmt)
	}

	// evaluate struct statement
	structStmt, ok := statement.(*ast.StructStatement)
	if ok {
		output = e.evaluateStructStatement(structStmt)
	}

	// evaluate append statement
	appendStmt, ok := statement.(*ast.AppendStatement)
	if ok {
		output = e.evaluateAppendStatement(appendStmt)
	}

	// evaluate print statement
	printStmt, ok := statement.(*ast.PrintStatement)
	if ok {
		// evaluate the expression to be printed
		value, list := e.evaluateExpression(printStmt.Value)
		// if it is not a list, output single value
		if len(list.arr) == 0 {
			switch value.dataType {
			// if expression is an int
			case INTTYPE:
				output = fmt.Sprintf("%d", value.intValue)
			// if expression is a bool
			case BOOLTYPE:
				output = fmt.Sprintf("%t", value.boolValue)
			// if expression is a float
			case FLOATTYPE:
				output = fmt.Sprintf("%.5g", value.floatValue)
			// if expression is a string
			case STRINGTYPE:
				output = strings.Trim(value.stringValue, "\"")
			}
		} else {
			// print lists
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
	// if the identifier in the create statement is a list, create an empty list in the list map
	if statement.Name.IsList {
		e.listMap[statement.Name.Value] = List{dataType: variableTypes[statement.Name.DataType]}
	} else {
		// otherwise, create an unassigned variable in the variable map
		e.variableMap[statement.Name.Value] = Data{dataType: variableTypes[statement.Name.DataType]}
	}
	return ""
}

func (e *Evaluator) evaluateSetStatement(statement *ast.SetStatement) string {
	var name string
	// if the identifier has an attribute, reformat the name of the identifier
	if statement.Name.Attribute != "" {
		name = fmt.Sprintf("%s.%s", statement.Name.Value, statement.Name.Attribute)
	} else {
		name = statement.Name.Value
	}

	// check if the variable or list is created
	_, okVar := e.variableMap[name]
	_, okList := e.listMap[name]
	// if it is not created, send an error
	if !okVar && !okList {
		e.errors = append(e.errors, fmt.Sprintf("Variable %s has not been created.", name))
		return ""
	}
	// evaluate expression being assigned
	value, list := e.evaluateExpression(statement.Value)
	// if there is an index value and the list is created, set the list at the index value (list[1] = <value>)
	// len(list.arr) == 0 is to determine that we are not trying to assign a list element to a list
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
	// if it is not a list and the data type does not match the created variable's type, send an error
	if value != (Data{}) && value.dataType != e.variableMap[name].dataType {
		e.errors = append(e.errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
			dataTypeToString[e.variableMap[name].dataType], dataTypeToString[value.dataType]))
		return ""
	}
	// if we are assigning a list literal and the data does not match the created lists data type, send an error
	if len(list.arr) != 0 && list.dataType != e.listMap[name].dataType {
		e.errors = append(e.errors, "Mismatching list data types")
		return ""
	}
	// if it is a list and no errors were found, assign the list in the list map
	if len(list.arr) != 0 {
		e.listMap[name] = list
		return ""
	}
	// if it is a variable and no errors were found, assign the variable in the variable map
	if len(list.arr) == 0 {
		e.variableMap[name] = value
	}
	return ""
}

func (e *Evaluator) evaluateStructStatement(statement *ast.StructStatement) string {
	// iterate over every attribute
	for _, attribute := range statement.Attributes {
		attributeName := fmt.Sprintf("%s.%s", statement.StructIdent.Value, attribute.Value)
		// if there is no assigned value to the attribute, assign an empty data struct
		if statement.Values[attribute.Value] == nil {
			e.variableMap[attributeName] = Data{dataType: variableTypes[attribute.DataType]}
		} else {
			// otherwise, evaluate the assigned value
			value, _ := e.evaluateExpression(statement.Values[attribute.Value])
			// if there is a data mismatch, send an error
			if variableTypes[attribute.DataType] != value.dataType {
				e.errors = append(e.errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
					attribute.DataType, dataTypeToString[value.dataType]))
				return ""
			} else {
				// if no data mismatch, assign the attribute combined name (myStruct.attribute) to the evaluated expression
				e.variableMap[attributeName] = value
			}
		}
	}
	return ""
}

func (e *Evaluator) evaluateIfStatement(statement *ast.IfStatement) string {
	ret := ""
	conditionData, _ := e.evaluateExpression(statement.Condition)
	// if the condition does not evaluate into a boolean, send an error
	if conditionData.dataType != BOOLTYPE {
		e.errors = append(e.errors, "Cannot use non-boolean expression in if statement condition.")
		return ""
	} else {
		// if the boolean is true, execute the iftrue statements
		if conditionData.boolValue {
			for _, stmt := range statement.IfTrue {
				ret += e.evaluateStatement(stmt)
			}
			// otherwise, execute the else statements
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
	// if condition expression is not evaluated as a boolean, send an error
	if conditionData.dataType != BOOLTYPE {
		e.errors = append(e.errors, "Cannot use non-boolean expression in while statement condition.")
		return ""
	} else {
		// execute the loop statements so long as the condition value is true
		for conditionData.boolValue {
			for _, stmt := range statement.LoopStatements {
				ret += e.evaluateStatement(stmt)
			}
			// must update condition on each loop
			conditionData, _ = e.evaluateExpression(statement.Condition)
		}
	}

	return ret
}

func (e *Evaluator) evaluateCountStatement(statement *ast.CountStatement) string {
	ret := ""
	// evaluate from, to, and by expressions
	fromValue, _ := e.evaluateExpression(statement.From)
	toValue, _ := e.evaluateExpression(statement.To)
	byValue, _ := e.evaluateExpression(statement.By)
	// if from, to, and by are not integers, send an error
	if fromValue.dataType != INTTYPE || toValue.dataType != INTTYPE || byValue.dataType != INTTYPE {
		e.errors = append(e.errors, "Cannot use non-integer values in counting statement from, to, or by.")
		return ret
	}

	// get counter identifier
	counterName := statement.Counter.Value
	// set counter variable to the initial from value
	e.variableMap[counterName] = fromValue
	// iterate from <from> to <to>
	for e.variableMap[counterName].intValue <= toValue.intValue {
		// execute loop statements
		for _, stmt := range statement.LoopStatements {
			ret += e.evaluateStatement(stmt)
		}
		tempData := e.variableMap[counterName]
		// increment the counter variable by the <by> value
		tempData.intValue += byValue.intValue
		e.variableMap[counterName] = tempData
	}

	return ret
}

func (e *Evaluator) evaluateAppendStatement(statement *ast.AppendStatement) string {
	listName := statement.List.Value
	// check if the list exists, if not send an error
	list, ok := e.listMap[listName]
	if !ok {
		e.errors = append(e.errors, fmt.Sprintf("List %s does not exist.", listName))
		return ""
	}

	// evaluate appended expression
	appendValue, _ := e.evaluateExpression(statement.Value)
	// if the append expression is not a single value (not a list), then send an error
	if appendValue == (Data{}) {
		e.errors = append(e.errors, "Invalid append data")
		return ""
	}

	// if the appended data does not match the lists data type, send an error
	if appendValue.dataType != list.dataType {
		e.errors = append(e.errors, "Mismatching append data type")
		return ""
	}

	// if no errors, append value to the list
	list.arr = append(list.arr, appendValue)
	e.listMap[listName] = list

	return ""
}

func (e *Evaluator) evaluateLenExpression(statement *ast.LengthExpression) Data {
	list, ok := e.listMap[statement.List.Value]
	// check if list exists, if not send an error
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
		// if any element does not match the list data type, send an error
		if data.dataType != dataList.dataType {
			e.errors = append(e.errors, "Mismatching list elements")
		}
		dataList.arr = append(dataList.arr, data)
	}

	return dataList
}

// considering an expression can either be a single data value or a list, we return both with one being an empty struct
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
	// if the identifier has an attribute, reformat the name to <identifier>.<attribute>
	if identifier.Attribute != "" {
		name = fmt.Sprintf("%s.%s", identifier.Value, identifier.Attribute)
	} else {
		name = identifier.Value
	}
	value, ok := e.variableMap[name]
	// if the variable does not exist, check for a list
	if !ok {
		list, ok = e.listMap[name]
		// if no list, then the identifier does not exist, send an error
		if !ok {
			e.errors = append(e.errors, fmt.Sprintf("Variable %s does not exist.", identifier.Value))
			return Data{}, List{}
		} else {
			// if it is a list, check if we are looking for a particular element of the list
			if identifier.Index != nil {
				indexValue, _ := e.evaluateExpression(identifier.Index)
				// if the index is not an int send an error
				if indexValue.dataType != INTTYPE {
					e.errors = append(e.errors, "Cannot use non-integer values in list index")
					return Data{}, List{}
				}
				// if the index is out of bounds, send an error
				if int(indexValue.intValue) > len(e.listMap[name].arr) || indexValue.intValue < 0 {
					e.errors = append(e.errors, "Index out of range")
					return Data{}, List{}
				}
				// otherwise, get the element at the specified index (-1 because of 1-based indexing)
				return list.arr[indexValue.intValue-1], List{}
			}
			return Data{}, list
		}
	}

	return value, List{}
}

func (e *Evaluator) evaluatePrefixExp(expression *ast.PrefixExpression) Data {
	value, _ := e.evaluateExpression(expression.Right)
	// evaluate - and ! prefix operators
	switch expression.Operator {
	case "-":
		if value.dataType != INTTYPE && value.dataType != FLOATTYPE {
			e.errors = append(e.errors, "Cannot perform negation on a non-numeric value")
			return Data{}
		}
		value.intValue = -1 * value.intValue
	case "!":
		if value.dataType != BOOLTYPE {
			e.errors = append(e.errors, "Cannot perform NOT on a non-boolean value")
			return Data{}
		}
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
		// check if the values are non-numeric, if so send an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform arithmetic operations on non-numeric values.")
			return Data{}
		}
		// if either side is a float, the result will be a float
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		// cast left side to a float
		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		// cast right side to a float
		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		// if the ret value is not a float, cast the addition of left and right to an int and return
		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left + right)
		} else {
			retValue.floatValue = left + right
		}

	case "-":
		// check if the values are non-numeric, if so send an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform arithmetic operations on non-numeric values.")
			return Data{}
		}
		// if either side is a float, the result will be a float
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		// cast left side to a float
		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		// cast right side to a float
		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		// if the ret value is not a float, cast the addition of left and right to an int and return
		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left - right)
		} else {
			retValue.floatValue = left - right
		}

	case "*":
		// check if the values are non-numeric, if so send an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform arithmetic operations on non-numeric values.")
			return Data{}
		}
		// if either side is a float, the result will be a float
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		// cast left side to a float
		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		// cast right side to a float
		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		// if the ret value is not a float, cast the addition of left and right to an int and return
		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left * right)
		} else {
			retValue.floatValue = left * right
		}

	case "/":
		// check if the values are non-numeric, if so send an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform arithmetic operations on non-numeric values.")
			return Data{}
		}
		// cast left side to a float
		var left float32
		switch leftValue.dataType {
		case INTTYPE:
			left += float32(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		// cast right side to a float
		var right float32
		switch rightValue.dataType {
		case INTTYPE:
			right += float32(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		// send error if dividing by 0
		if right == 0 {
			e.errors = append(e.errors, "Cannot divide by 0.")
			return Data{}
		}

		// if the result has no numbers after the decimal place, the result is an int, otherwise, it is a float
		if math.Mod(float64(left/right), 1) == 0 {
			retValue.dataType = INTTYPE
		} else {
			retValue.dataType = FLOATTYPE
		}

		// if the ret value is not a float, cast the addition of left and right to an int and return
		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int32(left / right)
		} else {
			retValue.floatValue = left / right
		}
	case "==":
		retValue.dataType = BOOLTYPE
		switch leftValue.dataType {
		// compare values depending on data type
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
		// compare values depending on data type
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
		// if left or right is not numerical, return an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform comparison operations on non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue > rightValue.intValue
	case ">=":
		retValue.dataType = BOOLTYPE
		// if left or right is not numerical, return an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform comparison operations on non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue >= rightValue.intValue
	case "<":
		retValue.dataType = BOOLTYPE
		// if left or right is not numerical, return an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform comparison operations on non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue < rightValue.intValue
	case "<=":
		retValue.dataType = BOOLTYPE
		// if left or right is not numerical, return an error
		if !slices.Contains(numericDataTypes, leftValue.dataType) || !slices.Contains(numericDataTypes, rightValue.dataType) {
			e.errors = append(e.errors, "Cannot perform comparison operations on non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue <= rightValue.intValue
	case "or":
		retValue.dataType = BOOLTYPE
		// if left or right are non-booleans, return an error
		if leftValue.dataType != BOOLTYPE || rightValue.dataType != BOOLTYPE {
			e.errors = append(e.errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue || rightValue.boolValue
	case "and":
		retValue.dataType = BOOLTYPE
		// if left or right are non-booleans, return an error
		if leftValue.dataType != BOOLTYPE || rightValue.dataType != BOOLTYPE {
			e.errors = append(e.errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue && rightValue.boolValue
	}
	return retValue
}
