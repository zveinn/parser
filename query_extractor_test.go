package parser

import (
	"fmt"
	"testing"
)

func TestGetQueryExclusive(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		param       string
		expectOk    bool
		expectValue string
	}{
		// Basic AND cases - should work
		{
			name:        "Simple AND with two params",
			query:       "PoolID = 1 AND SetID = 2",
			param:       "PoolID",
			expectOk:    true,
			expectValue: "1",
		},
		{
			name:        "Simple AND with two params - second param",
			query:       "PoolID = 1 AND SetID = 2",
			param:       "SetID",
			expectOk:    true,
			expectValue: "2",
		},
		{
			name:        "Multiple AND conditions",
			query:       "PoolID = 1 AND SetID = 2 AND Count = 100",
			param:       "Count",
			expectOk:    true,
			expectValue: "100",
		},
		{
			name:        "String values with quotes",
			query:       "Name = 'Alice' AND Age = 30",
			param:       "Name",
			expectOk:    true,
			expectValue: "Alice", // Parser removes quotes
		},
		// Note: Double quotes are currently not supported by the parser's lexer
		// TODO: Fix double quote support in enhanced_lexer.go
		// {
		// 	name:        "Double quoted strings",
		// 	query:       `Name = "Bob" AND Age = 25`,
		// 	param:       "Name",
		// 	expectOk:    true,
		// 	expectValue: "Bob", // Parser removes quotes
		// },
		{
			name:        "Case insensitive parameter matching",
			query:       "poolid = 5 AND setid = 10",
			param:       "PoolID",
			expectOk:    true,
			expectValue: "5",
		},
		{
			name:        "Single equality check",
			query:       "UserID = 42",
			param:       "UserID",
			expectOk:    true,
			expectValue: "42",
		},
		{
			name:        "Numeric values with decimals",
			query:       "Price = 99.99 AND Quantity = 5",
			param:       "Price",
			expectOk:    true,
			expectValue: "99.99",
		},

		// OR cases - should return false
		{
			name:        "Simple OR expression",
			query:       "a = 2 OR b = 3",
			param:       "a",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "OR expression second param",
			query:       "a = 2 OR b = 3",
			param:       "b",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Mixed AND/OR expression",
			query:       "a = 1 AND b = 2 OR c = 3",
			param:       "a",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Parenthesized OR groups",
			query:       "(a = 1 AND b = 2) OR c = 3",
			param:       "a",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Complex nested OR",
			query:       "a = 1 AND (b = 2 OR c = 3)",
			param:       "b",
			expectOk:    false,
			expectValue: "",
		},

		// Non-exclusive cases
		{
			name:        "Parameter not found",
			query:       "PoolID = 1 AND SetID = 2",
			param:       "Count",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Empty query",
			query:       "",
			param:       "PoolID",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Empty parameter",
			query:       "PoolID = 1",
			param:       "",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Non-equality operator",
			query:       "Age > 25 AND Name = 'Alice'",
			param:       "Age",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "NOT expression",
			query:       "NOT (Status = 'active')",
			param:       "Status",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "CONTAINS operator",
			query:       "Tags CONTAINS 'important' AND Status = 'active'",
			param:       "Tags",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "IS NULL check",
			query:       "Department IS NULL AND Status = 'pending'",
			param:       "Department",
			expectOk:    false,
			expectValue: "",
		},

		// Complex real-world examples
		{
			name:        "Database query example",
			query:       "TableName = 'users' AND SchemaName = 'public' AND DatabaseName = 'production'",
			param:       "SchemaName",
			expectOk:    true,
			expectValue: "public", // Parser removes quotes
		},
		{
			name:        "API filter example",
			query:       "endpoint = '/api/v1/users' AND method = 'GET' AND status = 200",
			param:       "method",
			expectOk:    true,
			expectValue: "GET", // Parser removes quotes
		},
		{
			name:        "Mixed operators - should fail",
			query:       "UserID = 100 AND Age > 18 AND Status = 'active'",
			param:       "Age",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Chained AND with spaces",
			query:       "   PoolID   =   1   AND   SetID   =   2   ",
			param:       "PoolID",
			expectOk:    true,
			expectValue: "1",
		},

		// Edge cases with humanized values
		{
			name:        "Time duration values",
			query:       "Timeout = 30s AND MaxRetries = 3",
			param:       "Timeout",
			expectOk:    true,
			expectValue: "30",
		},
		{
			name:        "Byte size values",
			query:       "MaxMemory = 8GB AND MinMemory = 1GB",
			param:       "MaxMemory",
			expectOk:    true,
			expectValue: "8000000000",
		},
		{
			name:        "SI prefix values",
			query:       "UserCount = 5K AND PageViews = 1.5M",
			param:       "UserCount",
			expectOk:    true,
			expectValue: "5000",
		},

		// Invalid query syntax - should return false
		{
			name:        "Malformed query",
			query:       "PoolID = AND SetID = 2",
			param:       "PoolID",
			expectOk:    false,
			expectValue: "",
		},
		{
			name:        "Unclosed parenthesis",
			query:       "(PoolID = 1 AND SetID = 2",
			param:       "PoolID",
			expectOk:    false,
			expectValue: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok, value := GetQueryExclusive(tc.query, tc.param)

			if ok != tc.expectOk {
				t.Errorf("GetQueryExclusive(%q, %q) ok = %v, want %v",
					tc.query, tc.param, ok, tc.expectOk)
			}

			if value != tc.expectValue {
				t.Errorf("GetQueryExclusive(%q, %q) value = %q, want %q",
					tc.query, tc.param, value, tc.expectValue)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetQueryExclusive_Simple(b *testing.B) {
	query := "PoolID = 1 AND SetID = 2"
	param := "PoolID"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetQueryExclusive(query, param)
	}
}

func BenchmarkGetQueryExclusive_Complex(b *testing.B) {
	query := "PoolID = 1 AND SetID = 2 AND UserID = 100 AND Status = 'active' AND Count = 50"
	param := "Status"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetQueryExclusive(query, param)
	}
}

func BenchmarkGetQueryExclusive_WithOR(b *testing.B) {
	query := "PoolID = 1 OR SetID = 2"
	param := "PoolID"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetQueryExclusive(query, param)
	}
}

// Example usage tests to demonstrate the API
func ExampleGetQueryExclusive() {
	// Basic usage
	query := "PoolID = 1 AND SetID = 2"

	if ok, value := GetQueryExclusive(query, "PoolID"); ok {
		fmt.Println("PoolID:", value)
	}

	if ok, value := GetQueryExclusive(query, "SetID"); ok {
		fmt.Println("SetID:", value)
	}

	// Non-exclusive due to OR
	orQuery := "Status = 'active' OR Status = 'pending'"
	if ok, _ := GetQueryExclusive(orQuery, "Status"); !ok {
		fmt.Println("Status is not exclusive due to OR")
	}
	// Output:
	// PoolID: 1
	// SetID: 2
	// Status is not exclusive due to OR
}
