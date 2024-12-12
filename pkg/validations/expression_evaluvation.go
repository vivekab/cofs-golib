package golibvalidations

import (
	"fmt"
	"regexp"
	"strings"

	golibarray "github.com/vivekab/golib/pkg/array"
)

const (
	OPERATOR_AND = "&&"
	OPERATOR_OR  = "||"

	NULL = "null"
)

// expression: field_name in value
// operators supported: in, not_in
// values: comma separated values
// examples: expression: "status in 0,2"
// "status not_in 3"
// "status in Pending Activation,Active"
// "type not_in 4 && status in 0,2"
// "type in 4 && status in 0,2 || status in 3"
// (status in 2 || (status in 4 && type in physical)) && (status in 3 && type in virtual)
func EvaluvateGolibExp(expression string, valueMap map[string]string) bool {
	stackOfOperations := []string{}

	push := func(operator string) {
		trimmedOperator := strings.TrimSpace(operator)
		if trimmedOperator != "" {
			stackOfOperations = append(stackOfOperations, trimmedOperator)
		}
	}

	pop := func() string {
		length := len(stackOfOperations)
		if length == 0 {
			return ""
		}

		lastElement := stackOfOperations[length-1]
		stackOfOperations = stackOfOperations[:length-1]
		return lastElement
	}

	tempStr := ""
	for i, char := range expression {
		switch char {
		case '(':
			push(tempStr)
			tempStr = ""
		case ')':
			if len(stackOfOperations) >= 2 {
				// Extract & evaluate the last 3 elements from the stack
				exp2 := tempStr
				operator := pop()
				exp1 := pop()

				rslt := evaluvateTuple(exp1, operator, exp2, valueMap)

				// Push the result back to the stack
				tempStr = fmt.Sprintf("%t", rslt)
				push(tempStr)
				tempStr = ""
			}
		case '&':
			if expression[i+1] == '&' {
				push(tempStr)
				push(OPERATOR_AND)
				tempStr = ""
			}
		case '|':
			if expression[i+1] == '|' {
				push(tempStr)
				push(OPERATOR_OR)
				tempStr = ""
			}

		default:
			tempStr += string(char)

		}
	}

	push(tempStr)

	for {
		if len(stackOfOperations) == 1 {
			break
		}
		if len(stackOfOperations) == 2 {
			// Invalid expression
			break
		}

		exp2 := pop()
		operator := pop()
		exp1 := pop()

		rslt := evaluvateTuple(exp1, operator, exp2, valueMap)

		push(fmt.Sprintf("%t", rslt))
	}

	return evaluvateSingleExpression(golibarray.SafeIndex(stackOfOperations, 0, ""), valueMap)
}

func evaluvateTuple(exp1, condition, exp2 string, valueMap map[string]string) bool {
	return evaluvateOperator(condition,
		evaluvateSingleExpression(exp1, valueMap),
		evaluvateSingleExpression(exp2, valueMap),
	)

}

func evaluvateOperator(operator string, left bool, right bool) bool {

	if operator == OPERATOR_AND {
		return left && right
	}

	if operator == OPERATOR_OR {
		return left || right
	}

	return false
}

func evaluvateSingleExpression(expression string, valueMap map[string]string) bool {
	expression = strings.TrimSpace(expression)

	if expression == "true" {
		return true
	}
	if expression == "false" {
		return false
	}

	expressionParts := strings.SplitN(expression, " ", 3)
	if len(expressionParts) != 3 {
		return false
	}

	fieldName := strings.TrimSpace(expressionParts[0])
	operator := strings.TrimSpace(expressionParts[1])
	ref := strings.TrimSpace(expressionParts[2])
	if ref == "empty" {
		ref = ""
	}

	refValues := strings.Split(ref, ",")

	actValue, ok := valueMap[fieldName]
	if !ok {
		return false
	}

	return operator == "in" && golibarray.Contains(refValues, actValue) ||
		operator == "not_in" && !golibarray.Contains(refValues, actValue) ||
		operator == "is" && actValue == refValues[0] ||
		operator == "is_not" && actValue != refValues[0] ||
		operator == "matches" && isRegexMatching(ref, actValue) ||
		operator == "is_valid" && isValidOfType(refValues[0], actValue)
}

func isRegexMatching(pattern, value string) bool {

	pattern = strings.TrimSpace(pattern)
	value = strings.TrimSpace(value)

	match, _ := regexp.MatchString(pattern, value)

	return match
}

func isValidOfType(typeName, value string) bool {

	predefined_regexs := map[string]string{
		"number":            `^\d+$`,
		"mcc_list":          `^([0-9]{4}(,[0-9]{4})*)?$`,
		"country_code_list": `^([A-Z]{2}(,[A-Z]{2})*)?$`,
		"ssn":               `^[0-9]{9}$`,
		"card_spend_types":  `^(atm_withdrawal|atm_deposit|pos|ecom|wallet|debit_push|debit_pull)(,(atm_withdrawal|atm_deposit|pos|ecom|wallet|debit_push|debit_pull))*$`,
		// Add more predefined regexs here
	}

	if pattern, ok := predefined_regexs[typeName]; ok {
		return isRegexMatching(pattern, value)
	}

	return false
}
