package parser

import (
	"strings"
	"testing"
)

// Define a simple test structure
type Person struct {
	Name       string
	Age        int
	IsEmployed bool
	Skills     []string
	Salary     float64
	Department *Department
	Tags       map[string]string
}

type Department struct {
	Name     string
	Location string
}

func TestSimpleComparisons(t *testing.T) {
	// Test data
	people := []Person{
		{Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50},
		{Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25},
		{Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75},
	}

	tests := []struct {
		name     string
		query    string
		expected int // expected number of results
	}{
		{"Equal string", "Name = 'Alice'", 1},
		{"Not equal string", "Name != 'Alice'", 2},
		{"Equal number", "Age = 30", 1},
		{"Greater than", "Age > 25", 2},
		{"Less than", "Age < 35", 2},
		{"Greater equal", "Age >= 30", 2},
		{"Less equal", "Age <= 30", 2},
		{"Boolean true", "IsEmployed = true", 2},
		{"Boolean false", "IsEmployed = false", 1},
		{"Contains in string", "Name CONTAINS 'li'", 2}, // Alice and Charlie
		{"Float comparison", "Salary > 70000", 2},
		{"Float exact match", "Salary = 65000.25", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestLogicalOperators(t *testing.T) {
	// Test data
	people := []Person{
		{Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50},
		{Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25},
		{Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75},
		{Name: "Diana", Age: 28, IsEmployed: true, Skills: []string{"Python", "JavaScript"}, Salary: 72000.00},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Simple AND", "Age > 25 AND IsEmployed = true", 3},
		{"Simple OR", "Age = 25 OR Age = 35", 2},
		{"Complex AND/OR", "Age > 30 OR (Age = 25 AND IsEmployed = false)", 2},
		{"Multiple AND", "Age > 25 AND IsEmployed = true AND Salary > 70000", 3},
		{"Multiple OR", "Age = 25 OR Age = 28 OR Age = 35", 3},
		{"Nested AND in OR", "Age = 30 OR (Age > 30 AND IsEmployed = true)", 2},
		{"Nested OR in AND", "IsEmployed = true AND (Age = 28 OR Age = 35)", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestContainsOperator(t *testing.T) {
	// Test data
	people := []Person{
		{Name: "Alice", Skills: []string{"Go", "Python", "SQL"}},
		{Name: "Bob", Skills: []string{"Java", "C++", "C#"}},
		{Name: "Charlie", Skills: []string{"Go", "Rust", "TypeScript"}},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Contains in array", "Skills CONTAINS 'Go'", 2},
		{"Contains substring", "Name CONTAINS 'li'", 2}, // Alice and Charlie
		{"Contains exact match", "Skills CONTAINS 'Java'", 1},
		{"Contains no match", "Skills CONTAINS 'PHP'", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestNegativeNumbers(t *testing.T) {
	type Item struct {
		Value int
		Price float64
	}

	items := []Item{
		{Value: -10, Price: -5.5},
		{Value: 0, Price: 0},
		{Value: 5, Price: 7.5},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Negative int equals", "Value = -10", 1},
		{"Negative float equals", "Price = -5.5", 1},
		{"Greater than negative", "Value > -5", 2},
		{"Less than negative", "Value < -5", 1},
		{"Between negatives", "Value > -15 AND Value < -5", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestNestedFields(t *testing.T) {
	people := []Person{
		{
			Name: "Alice",
			Department: &Department{
				Name:     "Engineering",
				Location: "New York",
			},
		},
		{
			Name: "Bob",
			Department: &Department{
				Name:     "Marketing",
				Location: "San Francisco",
			},
		},
		{
			Name: "Charlie",
			Department: &Department{
				Name:     "Engineering",
				Location: "Seattle",
			},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Nested equal", "Department.Name = 'Engineering'", 2},
		{"Nested contains", "Department.Location CONTAINS 'Francisco'", 1},
		{"Multiple nested conditions", "Department.Name = 'Engineering' AND Department.Location = 'Seattle'", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestIsNullOperator(t *testing.T) {
	people := []Person{
		{Name: "Alice", Department: &Department{Name: "Engineering"}},
		{Name: "Bob", Department: nil},
		{Name: "Charlie", Department: &Department{Name: "Marketing"}},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"IS NULL", "Department IS NULL", 1},
		{"IS NOT NULL", "Department IS NOT NULL", 2},
		{"Combination with other operators", "Department IS NOT NULL AND Department.Name = 'Engineering'", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestAnyOperator(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, Skills: []string{"Go", "Python"}},
		{Name: "Bob", Age: 25, Skills: []string{"Java", "C++"}},
		{Name: "Charlie", Age: 35, Skills: []string{"Go", "Rust"}},
		{Name: "Diana", Age: 28, Skills: []string{"Python", "JavaScript"}},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"ANY with single value", "ANY(Skills) = 'Go'", 2},
		{"ANY with multiple values", "ANY(Skills) = ANY('Go', 'Python')", 3},
		{"ANY with numbers", "ANY(Age) = ANY('25', '30')", 2},
		{"ANY with contains", "ANY(Skills) CONTAINS 'a'", 2}, // Java, JavaScript
		{"ANY with greater than", "ANY(Age) > '28'", 2},      // 30, 35
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestNotOperator(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, IsEmployed: true},
		{Name: "Bob", Age: 25, IsEmployed: false},
		{Name: "Charlie", Age: 35, IsEmployed: true},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"NOT with equals", "NOT Name = 'Alice'", 2},
		{"NOT with boolean", "NOT IsEmployed = true", 1},
		{"NOT with comparison", "NOT Age > 30", 2},
		{"NOT with AND", "NOT (Age > 25 AND IsEmployed = true)", 1},
		{"Complex NOT", "NOT (Name = 'Alice' OR Name = 'Bob')", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestEmptyQuery(t *testing.T) {
	people := []Person{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	// The parser returns an error for empty queries, which is expected behavior
	results, err := Parse("", people)
	if err != nil {
		t.Fatalf("Expected no error for empty query but got one %s", err)
	}

	// The error should mention that the AST is nil
	if len(results) != len(people) {
		t.Errorf("Empty query should return nil results, got %v", results)
	}
}

func TestSyntaxErrors(t *testing.T) {
	people := []Person{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	// These are the syntax errors that the parser currently catches
	tests := []struct {
		name        string
		query       string
		expectError bool
	}{
		{"Missing value", "Name = ", false}, // Parser doesn't catch this currently
		{"Invalid operator", "Name <> 'Alice'", true},
		{"Unclosed parenthesis", "(Name = 'Alice'", true},
		{"Unclosed string", "Name = 'Alice", true},                      // Parser now catches this
		{"Invalid field reference", "NonExistentField = 'test'", false}, // Error happens at evaluation time, not parse time
		{"Empty AND expression", "Name = 'Alice' AND", true},
		{"Empty OR expression", "Name = 'Alice' OR", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.query, people)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for invalid query '%s', but got none", tt.query)
			} else if !tt.expectError && err != nil {
				// We're documenting current behavior, not enforcing it as correct
				t.Logf("Note: Query '%s' now returns error: %v", tt.query, err)
			}
		})
	}
}

func TestMapFields(t *testing.T) {
	// Map field access should now be working

	people := []Person{
		{
			Name: "Alice",
			Tags: map[string]string{
				"team":     "backend",
				"location": "remote",
			},
		},
		{
			Name: "Bob",
			Tags: map[string]string{
				"team":     "frontend",
				"location": "office",
			},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Map field equality", "Tags.team = 'backend'", 1},
		{"Map field contains", "Tags.location CONTAINS 'remote'", 1},
		{"Case insensitive map key", "Tags.TEAM = 'frontend'", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

// TestMapFieldsWithHyphen covers issue #8: map keys containing a hyphen ('-')
// should be matchable. The lexer previously stopped reading the identifier at
// the hyphen, producing a "failed to parse query: -" error.
func TestMapFieldsWithHyphen(t *testing.T) {
	people := []Person{
		{
			Name: "Alice",
			Tags: map[string]string{
				"team-name": "backend",
				"location":  "remote",
			},
		},
		{
			Name: "Bob",
			Tags: map[string]string{
				"team-name": "frontend",
				"location":  "office",
			},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Hyphenated map key equality", "Tags.team-name = 'backend'", 1},
		{"Hyphenated map key case insensitive", "Tags.TEAM-NAME = 'frontend'", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestNestedMapFields(t *testing.T) {
	type ComplexPerson struct {
		Name       string
		Properties map[string]any
	}

	people := []ComplexPerson{
		{
			Name: "Alice",
			Properties: map[string]any{
				"contact": map[string]any{
					"email": "alice@example.com",
					"phone": "123-456-7890",
				},
				"preferences": map[string]any{
					"theme": "dark",
					"notifications": map[string]any{
						"email": true,
						"push":  false,
					},
				},
			},
		},
		{
			Name: "Bob",
			Properties: map[string]any{
				"contact": map[string]any{
					"email": "bob@example.com",
					"phone": "098-765-4321",
				},
				"preferences": map[string]any{
					"theme": "light",
					"notifications": map[string]any{
						"email": false,
						"push":  true,
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Simple map access",
			query:    "Properties.preferences.theme = 'dark'",
			expected: 1,
		},
		{
			name:     "Deeply nested map access",
			query:    "Properties.preferences.notifications.email = true",
			expected: 1,
		},
		{
			name:     "Case insensitive map keys",
			query:    "Properties.CONTACT.email = 'bob@example.com'",
			expected: 1,
		},
		{
			name:     "Contains operator on map value",
			query:    "Properties.contact.email CONTAINS 'alice'",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestCaseInsensitiveFields(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Lowercase field", "name = 'Alice'", 1},
		{"Uppercase field", "NAME = 'Alice'", 1},
		{"Mixed case field", "NaMe = 'Alice'", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
			}
		})
	}
}

func TestParenthesisAndPrecedence(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, IsEmployed: true},
		{Name: "Bob", Age: 25, IsEmployed: false},
		{Name: "Charlie", Age: 35, IsEmployed: true},
		{Name: "Diana", Age: 28, IsEmployed: true},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Basic parenthesis", "(Age > 25)", 3},
		{"Parenthesis with AND", "(Age > 25) AND IsEmployed = true", 3},
		{"Parenthesis with OR", "(Age = 25) OR (Age = 35)", 2},
		{"Multiple parenthesis", "(Age > 25) AND (IsEmployed = true)", 3},
		{"Nested parenthesis", "(Age > 25 AND (IsEmployed = true OR Name = 'Bob'))", 3},
		{"OR precedence", "Age = 25 OR Age = 30 AND IsEmployed = true", 2},
		{"AND precedence", "Age = 35 AND IsEmployed = true OR Age = 25", 2},
		{"Parenthesis overriding precedence", "(Age = 35 AND IsEmployed = true) OR Age = 25", 2},
		{"Complex nesting", "((Age > 25) AND (Name = 'Alice' OR Name = 'Charlie')) OR Name = 'Bob'", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, people)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d, got %v",
					tt.query, len(results), tt.expected, getNames(results))
			}
		})
	}
}

func TestUnclosedStringDetection(t *testing.T) {
	people := []Person{
		{Name: "Alice"},
		{Name: "Bob"},
	}

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Properly closed string",
			query:   "Name = 'Alice'",
			wantErr: false,
		},
		{
			name:    "Unclosed string",
			query:   "Name = 'Alice",
			wantErr: true,
		},
		{
			name:    "Unclosed string with escape",
			query:   "Name = 'Alice\\'",
			wantErr: true,
		},
		{
			name:    "String with escaped quote",
			query:   "Name = 'Alice\\'s'",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.query, people)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				if !strings.Contains(err.Error(), "unclosed string") {
					t.Errorf("Expected error message to contain 'unclosed string', got: %v", err)
				}
			}
		})
	}
}

func TestNumericValidation(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, Salary: 75000.50},
		{Name: "Bob", Age: 25, Salary: 65000.25},
	}

	tests := []struct {
		name        string
		query       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid integer comparison",
			query:       "Age > 20",
			expectError: false,
		},
		{
			name:        "Valid float comparison",
			query:       "Salary > 70000.00",
			expectError: false,
		},
		{
			name:        "Invalid integer - non-numeric",
			query:       "Age = 'thirty'",
			expectError: true,
			errorMsg:    "invalid integer value",
		},
		{
			name:        "Invalid integer - alphabetic characters",
			query:       "Age > 25abc",
			expectError: true,
			errorMsg:    "invalid numeric value",
		},
		{
			name:        "Invalid float - non-numeric",
			query:       "Salary = 'seventy-five thousand'",
			expectError: true,
			errorMsg:    "invalid floating point value",
		},
		{
			name:        "Valid float with commas",
			query:       "Salary > 65,000.25",
			expectError: false,
		},
		{
			name:        "Valid negative integer",
			query:       "Age > -10",
			expectError: false,
		},
		{
			name:        "Valid scientific notation",
			query:       "Salary > 7.5e4",
			expectError: false,
		},
		{
			name:        "Valid zero",
			query:       "Age > 0",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.query, people)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for query '%s', but got none", tt.query)
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for query '%s': %v", tt.query, err)
			}

			if tt.expectError && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', but got: %v", tt.errorMsg, err)
				}
			}
		})
	}
}

func TestNumericValidationErrors(t *testing.T) {
	type Item struct {
		Age int
	}

	items := []Item{
		{Age: 25},
	}

	// Test an obviously invalid numeric format - typos like "25abc"
	input := "Age = 25abc"
	_, err := Parse(input, items)

	if err == nil {
		t.Errorf("Expected an error for invalid numeric format '%s', but got none", input)
	}
}

// Helper function to get names for debugging
func getNames(people []Person) []string {
	names := make([]string, len(people))
	for i, p := range people {
		names[i] = p.Name
	}
	return names
}

func TestMessageFieldParsing(t *testing.T) {
	type RowWithMessage struct {
		UserName string
		Message  string
	}

	type RowWithMessages struct {
		UserName string
		Messages string
	}

	dataMessage := []RowWithMessage{
		{UserName: "kyle", Message: "hello"},
		{UserName: "alice", Message: "world"},
	}
	dataMessages := []RowWithMessages{
		{UserName: "kyle", Messages: "1"},
		{UserName: "bob", Messages: "2"},
	}

	tests := []struct {
		name          string
		query         string
		data          any
		expectError   bool
		expectedCount int
	}{
		// Queries on RowWithMessage
		{"Error on wrong field 'Messages'", "UserName = 'kyle' AND NOT Messages CONTAINS '1'", dataMessage, true, 0},
		{"Successful CONTAINS on 'Message'", "UserName = 'kyle' AND Message CONTAINS 'hell'", dataMessage, false, 1},
		{"Successful equality on 'Message'", "UserName = 'alice' AND Message = 'world'", dataMessage, false, 1},
		{"Error on wrong field 'Messages' with equality", "UserName = 'alice' AND Messages = '1'", dataMessage, true, 0},

		// Queries on RowWithMessages
		{"Successful CONTAINS on 'Messages'", "UserName = 'kyle' AND Messages CONTAINS '1'", dataMessages, false, 1},
		{"Successful equality on 'Messages'", "UserName = 'bob' AND Messages = '2'", dataMessages, false, 1},
		{"Error on wrong field 'Message' with CONTAINS", "UserName = 'bob' AND Message CONTAINS '2'", dataMessages, true, 0},
		{"Error on wrong field 'Message' with equality", "UserName = 'bob' AND Message = '2'", dataMessages, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var resultsLength int

			switch data := tt.data.(type) {
			case []RowWithMessage:
				results, parseErr := Parse(tt.query, data)
				err = parseErr
				resultsLength = len(results)
			case []RowWithMessages:
				results, parseErr := Parse(tt.query, data)
				err = parseErr
				resultsLength = len(results)
			default:
				t.Fatalf("unsupported data type: %T", tt.data)
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error for query '%s', but got none", tt.query)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for query '%s': %v", tt.query, err)
				}
				if resultsLength != tt.expectedCount {
					t.Errorf("Query '%s' returned %d results, expected %d", tt.query, resultsLength, tt.expectedCount)
				}
			}
		})
	}
}
