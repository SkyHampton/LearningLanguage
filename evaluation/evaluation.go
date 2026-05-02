package evaluation

import (
	"bytes"
	"fmt"
	"learningLanguage/ast"
	"strconv"
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
	intValue    int64
	boolValue   bool
	stringValue string
	floatValue  float64
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

var variableMap = make(map[string]Data)
var errors []string

func EvaluateProgram(program *ast.Program) (string, []string) {
	var output bytes.Buffer
	errors = []string{}
	for _, statement := range program.Statements {
		output.WriteString(evaluateStatement(statement))
	}

	return output.String(), errors
}

func evaluateStatement(statement ast.Statement) string {
	var output string
	createStmt, ok := statement.(*ast.CreateStatement)
	if ok {
		output = evaluateCreateStatement(createStmt)
	}

	setStmt, ok := statement.(*ast.SetStatement)
	if ok {
		output = evaluateSetStatement(setStmt)
	}

	ifStmt, ok := statement.(*ast.IfStatement)
	if ok {
		output = evaluateIfStatement(ifStmt)
	}

	whileStmt, ok := statement.(*ast.WhileStatement)
	if ok {
		output = evaluateWhileStatement(whileStmt)
	}

	structStmt, ok := statement.(*ast.StructStatement)
	if ok {
		output = evaluateStructStatement(structStmt)
	}

	printStmt, ok := statement.(*ast.PrintStatement)
	if ok {
		value := evaluateExpression(printStmt.Value)
		switch value.dataType {
		case INTTYPE:
			output = strconv.FormatInt(value.intValue, 10) + "\n"
		case BOOLTYPE:
			output = strconv.FormatBool(value.boolValue) + "\n"
		case FLOATTYPE:
			output = strconv.FormatFloat(value.floatValue, 'f', -1, 64) + "\n"
		case STRINGTYPE:
			output = strings.Trim(value.stringValue, "\"") + "\n"
		}
	}

	return output
}

func evaluateCreateStatement(statement *ast.CreateStatement) string {
	variableMap[statement.Name.Value] = Data{dataType: variableTypes[statement.Name.DataType]}
	return ""
}

func evaluateSetStatement(statement *ast.SetStatement) string {
	var name string
	if statement.Name.Attribute != "" {
		name = fmt.Sprintf("%s.%s", statement.Name.Value, statement.Name.Attribute)
	} else {
		name = statement.Name.Value
	}
	_, ok := variableMap[name]
	if !ok {
		errors = append(errors, fmt.Sprintf("Variable %s has not been created.", name))
		return ""
	}
	value := evaluateExpression(statement.Value)
	if value.dataType != variableMap[name].dataType {
		errors = append(errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
			dataTypeToString[variableMap[name].dataType], dataTypeToString[value.dataType]))
		return ""
	}
	variableMap[name] = value
	return ""
}

func evaluateStructStatement(statement *ast.StructStatement) string {
	for _, attribute := range statement.Attributes {
		attributeName := fmt.Sprintf("%s.%s", statement.StructIdent.Value, attribute.Value)
		if statement.Values[attribute.Value] == nil {
			variableMap[attributeName] = Data{dataType: variableTypes[attribute.DataType]}
		} else {
			value := evaluateExpression(statement.Values[attribute.Value])
			if variableTypes[attribute.DataType] != value.dataType {
				errors = append(errors, fmt.Sprintf("Expected data type: %s. Got data type: %s",
					attribute.DataType, dataTypeToString[value.dataType]))
				return ""
			}
		}
	}
	return ""
}

func evaluateIfStatement(statement *ast.IfStatement) string {
	ret := ""
	conditionData := evaluateExpression(statement.Condition)
	if conditionData.dataType != BOOLTYPE {
		errors = append(errors, "Cannot use non-boolean expression in if statement condition.")
		return ""
	} else {
		if conditionData.boolValue {
			for _, stmt := range statement.IfTrue {
				ret += evaluateStatement(stmt)
			}
		} else {
			for _, stmt := range statement.Else {
				ret += evaluateStatement(stmt)
			}
		}
	}

	return ret
}

func evaluateWhileStatement(statement *ast.WhileStatement) string {
	ret := ""
	conditionData := evaluateExpression(statement.Condition)
	if conditionData.dataType != BOOLTYPE {
		errors = append(errors, "Cannot use non-boolean expression in while statement condition.")
		return ""
	} else {
		for conditionData.boolValue {
			for _, stmt := range statement.LoopStatements {
				ret += evaluateStatement(stmt)
			}
			conditionData = evaluateExpression(statement.Condition)
		}
	}

	return ret
}

func evaluateExpression(expression ast.Expression) Data {
	var value Data

	intLit, ok := expression.(*ast.IntegerLiteral)
	if ok {
		value = evaluateIntLit(intLit)
	}

	boolLit, ok := expression.(*ast.BooleanLiteral)
	if ok {
		value = evaluateBoolLit(boolLit)
	}

	floatLiteral, ok := expression.(*ast.FloatLiteral)
	if ok {
		value = evaluateFloatLit(floatLiteral)
	}

	strLit, ok := expression.(*ast.StringLiteral)
	if ok {
		value = evaluateStringLit(strLit)
	}

	identifierExp, ok := expression.(*ast.Identifier)
	if ok {
		value = evaluateIdentifier(identifierExp)
	}

	prefixExp, ok := expression.(*ast.PrefixExpression)
	if ok {
		value = evaluatePrefixExp(prefixExp)
	}

	infixExp, ok := expression.(*ast.InfixExpression)
	if ok {
		value = evaluateInfixExp(infixExp)
	}

	return value
}

func evaluateIntLit(expression *ast.IntegerLiteral) Data {
	return Data{dataType: INTTYPE, intValue: expression.Value}
}

func evaluateBoolLit(expression *ast.BooleanLiteral) Data {
	return Data{dataType: BOOLTYPE, boolValue: expression.Value}
}

func evaluateFloatLit(expression *ast.FloatLiteral) Data {
	return Data{dataType: FLOATTYPE, floatValue: expression.Value}
}

func evaluateStringLit(expression *ast.StringLiteral) Data {
	return Data{dataType: STRINGTYPE, stringValue: expression.Value}
}

func evaluateIdentifier(identifier *ast.Identifier) Data {
	var name string
	if identifier.Attribute != "" {
		name = fmt.Sprintf("%s.%s", identifier.Value, identifier.Attribute)
	} else {
		name = identifier.Value
	}
	value, ok := variableMap[name]
	if !ok {
		errors = append(errors, fmt.Sprintf("Variable %s does not exist.", identifier.Value))
		return Data{}
	}

	return value
}

func evaluatePrefixExp(expression *ast.PrefixExpression) Data {
	value := evaluateExpression(expression.Right)
	switch expression.Operator {
	case "-":
		value.intValue = -1 * value.intValue
	case "!":
		value.boolValue = !value.boolValue
	}

	return value
}

func evaluateInfixExp(expression *ast.InfixExpression) Data {
	leftValue := evaluateExpression(expression.Left)
	rightValue := evaluateExpression(expression.Right)
	var retValue Data

	switch expression.Operator {
	case "+":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float64
		switch leftValue.dataType {
		case INTTYPE:
			left += float64(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float64
		switch rightValue.dataType {
		case INTTYPE:
			right += float64(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int64(left + right)
		} else {
			retValue.floatValue = left + right
		}

	case "-":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float64
		switch leftValue.dataType {
		case INTTYPE:
			left += float64(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float64
		switch rightValue.dataType {
		case INTTYPE:
			right += float64(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int64(left - right)
		} else {
			retValue.floatValue = left - right
		}

	case "*":
		if leftValue.dataType == FLOATTYPE || rightValue.dataType == FLOATTYPE {
			retValue.dataType = FLOATTYPE
		}

		var left float64
		switch leftValue.dataType {
		case INTTYPE:
			left += float64(leftValue.intValue)
		case FLOATTYPE:
			left += leftValue.floatValue
		}

		var right float64
		switch rightValue.dataType {
		case INTTYPE:
			right += float64(rightValue.intValue)
		case FLOATTYPE:
			right += rightValue.floatValue
		}

		if retValue.dataType != FLOATTYPE {
			retValue.intValue = int64(left * right)
		} else {
			retValue.floatValue = left * right
		}

	case "/":
		retValue.dataType = INTTYPE
		retValue.intValue = leftValue.intValue / rightValue.intValue
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
			errors = append(errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue > rightValue.intValue
	case ">=":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			errors = append(errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue >= rightValue.intValue
	case "<":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			errors = append(errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue < rightValue.intValue
	case "<=":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != INTTYPE && leftValue.dataType != FLOATTYPE && rightValue.dataType != INTTYPE && rightValue.dataType != FLOATTYPE {
			errors = append(errors, "Cannot perform perform quanitative comparisons with non-numeric values.")
			return Data{}
		}
		retValue.boolValue = leftValue.intValue <= rightValue.intValue
	case "or":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != BOOLTYPE && rightValue.dataType != BOOLTYPE {
			errors = append(errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue || rightValue.boolValue
	case "and":
		retValue.dataType = BOOLTYPE
		if leftValue.dataType != BOOLTYPE && rightValue.dataType != BOOLTYPE {
			errors = append(errors, "Cannot perform perform logical operations with non-booleans.")
			return Data{}
		}
		retValue.boolValue = leftValue.boolValue && rightValue.boolValue
	}
	return retValue
}
