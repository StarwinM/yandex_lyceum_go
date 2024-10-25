package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if expression == "" {
		return 0, errors.New("empty expression")
	}

	numbers := make([]float64, 0)
	operators := make([]rune, 0)

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	numStr := ""
	for i, char := range expression {
		switch {
		case unicode.IsDigit(char) || char == '.':
			numStr += string(char)
			if i == len(expression)-1 || !(unicode.IsDigit(rune(expression[i+1])) || rune(expression[i+1]) == '.') {
				num, err := strconv.ParseFloat(numStr, 64)
				if err != nil {
					return 0, fmt.Errorf("invalid number: %s", numStr)
				}
				numbers = append(numbers, num)
				numStr = ""
			}
		case char == '(':
			operators = append(operators, char)
		case char == ')':
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if err := executeOperation(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		case char == '+' || char == '-' || char == '*' || char == '/':
			if char == '-' && (i == 0 || expression[i-1] == '(') {
				numStr = "-"
				continue
			}
			for len(operators) > 0 && operators[len(operators)-1] != '(' &&
				precedence[operators[len(operators)-1]] >= precedence[char] {
				if err := executeOperation(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			operators = append(operators, char)
		default:
			return 0, fmt.Errorf("invalid character: %c", char)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return 0, errors.New("mismatched parentheses")
		}
		if err := executeOperation(&numbers, &operators); err != nil {
			return 0, err
		}
	}

	if len(numbers) != 1 {
		return 0, errors.New("invalid expression")
	}

	return numbers[0], nil
}

func executeOperation(numbers *[]float64, operators *[]rune) error {
	if len(*numbers) < 2 {
		return errors.New("not enough operands")
	}

	b := (*numbers)[len(*numbers)-1]
	a := (*numbers)[len(*numbers)-2]
	op := (*operators)[len(*operators)-1]

	*numbers = (*numbers)[:len(*numbers)-2]
	*operators = (*operators)[:len(*operators)-1]

	var result float64
	switch op {
	case '+':
		result = a + b
	case '-':
		result = a - b
	case '*':
		result = a * b
	case '/':
		if b == 0 {
			return errors.New("division by zero")
		}
		result = a / b
	}

	*numbers = append(*numbers, result)
	return nil
}
