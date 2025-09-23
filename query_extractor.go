package parser

import (
	"strings"
)

// GetQueryExclusive extracts the value of a parameter from a query string
// if and only if that parameter appears exclusively in an AND context.
// Returns (true, value) if the parameter is found with an exclusive AND condition.
// Returns (false, "") if the parameter is not found, appears in an OR context,
// or if the query structure makes it non-exclusive.
//
// Examples:
//   - GetQueryExclusive("PoolID = 1 AND SetID = 2", "PoolID") returns (true, "1")
//   - GetQueryExclusive("a = 2 OR b = 3", "a") returns (false, "")
//   - GetQueryExclusive("(a = 1 AND b = 2) OR c = 3", "a") returns (false, "")
func GetQueryExclusive(query, param string) (ok bool, value string) {
	if query == "" || param == "" {
		return false, ""
	}

	// Parse the query to build the AST
	query = normalizeHumanizedValues(query)
	l := NewEnhancedLexer(query)
	p := NewParser(l)

	ast, err := p.ParseQuery()
	if err != nil || ast == nil {
		return false, ""
	}

	// Check if the entire expression is exclusive AND context
	return extractExclusiveValue(ast, param)
}

// extractExclusiveValue traverses the AST to find exclusive parameter values
func extractExclusiveValue(expr Expression, param string) (bool, string) {
	switch e := expr.(type) {
	case *ComparisonExpression:
		// Direct comparison - check if it matches our parameter
		if strings.EqualFold(e.Field, param) && e.Operator == EQ {
			return true, e.Value
		}
		return false, ""

	case *ConjunctionExpression:
		// AND expression - all parts must be AND-ed together
		// Check each sub-expression
		for _, subExpr := range e.Expressions {
			if ok, val := extractExclusiveValue(subExpr, param); ok {
				// Found the parameter in an AND context
				// Verify no OR expressions exist at this level
				if !containsOrExpression(e) {
					return true, val
				}
			}
		}
		return false, ""

	case *OrExpression:
		// OR expression - not exclusive, return false
		return false, ""

	case *NotExpression:
		// NOT expression - not an exclusive positive match
		return false, ""

	case *AnyExpression:
		// ANY expression - check if it matches our parameter
		if strings.EqualFold(e.Field, param) && e.Operator == EQ && len(e.Values) == 1 {
			return true, e.Values[0]
		}
		return false, ""

	default:
		return false, ""
	}
}

// containsOrExpression checks if an expression tree contains any OR expressions
func containsOrExpression(expr Expression) bool {
	switch e := expr.(type) {
	case *OrExpression:
		return true

	case *ConjunctionExpression:
		for _, subExpr := range e.Expressions {
			if containsOrExpression(subExpr) {
				return true
			}
		}
		return false

	case *NotExpression:
		return containsOrExpression(e.Expression)

	default:
		return false
	}
}

// Helper function to clean and extract values from quoted strings
func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	// Remove quotes if present
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
		   (value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}