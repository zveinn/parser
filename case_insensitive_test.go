package parser

import (
	"testing"
)

// TestCaseInsensitiveStringComparisons tests case-insensitive string comparisons
func TestCaseInsensitiveStringComparisons(t *testing.T) {
	type Product struct {
		Name     string
		Category string
		Brand    string
		Tags     []string
		Details  map[string]string
	}

	products := []Product{
		{
			Name:     "iPhone 15 Pro",
			Category: "Electronics",
			Brand:    "Apple",
			Tags:     []string{"Smartphone", "Premium", "5G"},
			Details:  map[string]string{"color": "Black", "storage": "256GB"},
		},
		{
			Name:     "SAMSUNG Galaxy S24",
			Category: "ELECTRONICS",
			Brand:    "Samsung",
			Tags:     []string{"SMARTPHONE", "Android", "premium"},
			Details:  map[string]string{"COLOR": "White", "Storage": "512gb"},
		},
		{
			Name:     "Google Pixel 8",
			Category: "electronics",
			Brand:    "GOOGLE",
			Tags:     []string{"smartphone", "ANDROID", "ai"},
			Details:  map[string]string{"Color": "GREEN", "STORAGE": "128GB"},
		},
		{
			Name:     "MacBook Pro",
			Category: "Computers",
			Brand:    "apple",
			Tags:     []string{"Laptop", "Professional", "M3"},
			Details:  map[string]string{"color": "Silver", "storage": "1TB"},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
		validate func([]Product) bool
	}{
		{
			name:     "Case insensitive EQ - lowercase query, mixed case data",
			query:    "Category = 'electronics'",
			expected: 3,
			validate: func(results []Product) bool {
				for _, p := range results {
					if p.Name == "MacBook Pro" {
						return false
					}
				}
				return true
			},
		},
		{
			name:     "Case insensitive EQ - uppercase query, mixed case data",
			query:    "Brand = 'APPLE'",
			expected: 2,
			validate: func(results []Product) bool {
				foundIPhone := false
				foundMacBook := false
				for _, p := range results {
					if p.Name == "iPhone 15 Pro" {
						foundIPhone = true
					}
					if p.Name == "MacBook Pro" {
						foundMacBook = true
					}
				}
				return foundIPhone && foundMacBook
			},
		},
		{
			name:     "Case insensitive NE - should exclude case-insensitive matches",
			query:    "Brand != 'samsung'",
			expected: 3,
			validate: func(results []Product) bool {
				for _, p := range results {
					if p.Name == "SAMSUNG Galaxy S24" {
						return false
					}
				}
				return true
			},
		},
		{
			name:     "Case insensitive CONTAINS in string field",
			query:    "Name CONTAINS 'phone'",
			expected: 1,
		},
		{
			name:     "Case insensitive CONTAINS with uppercase query",
			query:    "Name CONTAINS 'PHONE'",
			expected: 1,
		},
		{
			name:     "Case insensitive CONTAINS with mixed case",
			query:    "Name CONTAINS 'GaLaXy'",
			expected: 1,
		},
		{
			name:     "Case insensitive array CONTAINS",
			query:    "Tags CONTAINS 'smartphone'",
			expected: 3,
		},
		{
			name:     "Case insensitive array CONTAINS - uppercase query",
			query:    "Tags CONTAINS 'PREMIUM'",
			expected: 2,
		},
		{
			name:     "Case insensitive with AND operator",
			query:    "Category = 'ELECTRONICS' AND Brand = 'google'",
			expected: 1,
			validate: func(results []Product) bool {
				return len(results) == 1 && results[0].Name == "Google Pixel 8"
			},
		},
		{
			name:     "Case insensitive with OR operator",
			query:    "Brand = 'APPLE' OR Brand = 'google'",
			expected: 3,
		},
		{
			name:     "Case insensitive map field access",
			query:    "Details.color = 'black'",
			expected: 1,
		},
		{
			name:     "Case insensitive map field with mixed case",
			query:    "Details.storage = '256gb'",
			expected: 1,
		},
		{
			name:     "Case insensitive LT comparison",
			query:    "Brand < 'b'",
			expected: 2, // "Apple" and "apple" both < "b" when lowercased
		},
		{
			name:     "Case insensitive GT comparison",
			query:    "Brand > 'f'",
			expected: 2, // "GOOGLE" and "Samsung" both > "f" when lowercased
		},
		{
			name:     "Case insensitive LE comparison",
			query:    "Brand <= 'apple'",
			expected: 2,
		},
		{
			name:     "Case insensitive GE comparison",
			query:    "Brand >= 'SAMSUNG'",
			expected: 1,
		},
		{
			name:     "Complex query with multiple case insensitive operations",
			query:    "(Category = 'ELECTRONICS' AND Tags CONTAINS 'android') OR Brand = 'apple'",
			expected: 4,
		},
		{
			name:     "Case insensitive with NOT operator",
			query:    "NOT (Brand = 'SAMSUNG')",
			expected: 3,
		},
		{
			name:     "Case insensitive ANY operator with single value",
			query:    "ANY(Tags) = 'SMARTPHONE'",
			expected: 3,
		},
		{
			name:     "Case insensitive ANY operator with multiple values",
			query:    "ANY(Tags) = ANY('laptop', 'SMARTPHONE', 'Ai')",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, products)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
				t.Logf("Results:")
				for _, r := range results {
					t.Logf("  - %s (Brand: %s, Category: %s)", r.Name, r.Brand, r.Category)
				}
			}

			if tt.validate != nil && !tt.validate(results) {
				t.Errorf("Query '%s' validation failed", tt.query)
			}
		})
	}
}

// TestCaseInsensitiveWithSpecialCharacters tests case insensitivity with special characters
func TestCaseInsensitiveWithSpecialCharacters(t *testing.T) {
	type Document struct {
		Title   string
		Content string
		Tags    []string
	}

	docs := []Document{
		{
			Title:   "COVID-19 Report",
			Content: "Analysis of COVID-19 impact",
			Tags:    []string{"COVID-19", "Health", "2020"},
		},
		{
			Title:   "Post-COVID-19 Analysis",
			Content: "Recovery after covid-19",
			Tags:    []string{"covid-19", "recovery", "2021"},
		},
		{
			Title:   "E-Commerce Trends",
			Content: "Growth in e-commerce sector",
			Tags:    []string{"E-Commerce", "Business", "Digital"},
		},
		{
			Title:   "AI/ML Overview",
			Content: "Introduction to AI/ML technologies",
			Tags:    []string{"AI/ML", "Technology", "Future"},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Case insensitive with hyphen",
			query:    "Title CONTAINS 'covid-19'",
			expected: 2,
		},
		{
			name:     "Case insensitive with uppercase hyphenated term",
			query:    "Tags CONTAINS 'COVID-19'",
			expected: 2,
		},
		{
			name:     "Case insensitive with mixed case and hyphen",
			query:    "Content CONTAINS 'Covid-19'",
			expected: 2,
		},
		{
			name:     "Case insensitive with slash",
			query:    "Title = 'ai/ml overview'",
			expected: 1,
		},
		{
			name:     "Case insensitive E-Commerce variations",
			query:    "Tags CONTAINS 'e-commerce'",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, docs)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
				t.Logf("Results:")
				for _, r := range results {
					t.Logf("  - %s", r.Title)
				}
			}
		})
	}
}

// TestCaseInsensitiveEdgeCases tests edge cases for case insensitive comparisons
func TestCaseInsensitiveEdgeCases(t *testing.T) {
	type Item struct {
		ID          string
		Description string
		Attributes  []string
	}

	items := []Item{
		{ID: "1", Description: "", Attributes: []string{"", "Empty", "NULL"}},
		{ID: "2", Description: "EMPTY", Attributes: []string{"empty", "EMPTY", "Empty"}},
		{ID: "3", Description: "null", Attributes: []string{"null", "NULL", "Null"}},
		{ID: "4", Description: "True", Attributes: []string{"TRUE", "true", "True"}},
		{ID: "5", Description: "FALSE", Attributes: []string{"false", "FALSE", "False"}},
		{ID: "6", Description: " SPACES ", Attributes: []string{" spaces ", "SPACES", " Spaces"}},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Empty string comparison",
			query:    "Description = ''",
			expected: 1,
		},
		{
			name:     "Case insensitive 'empty' match",
			query:    "Description = 'empty'",
			expected: 1,
		},
		{
			name:     "Case insensitive 'null' as string",
			query:    "Description = 'NULL'",
			expected: 1,
		},
		{
			name:     "Case insensitive boolean-like string 'true'",
			query:    "Description = 'true'",
			expected: 1,
		},
		{
			name:     "Case insensitive boolean-like string 'FALSE'",
			query:    "Description = 'false'",
			expected: 1,
		},
		{
			name:     "String with spaces - case insensitive",
			query:    "Description = ' spaces '",
			expected: 1,
		},
		{
			name:     "Array contains empty variations",
			query:    "Attributes CONTAINS 'EMPTY'",
			expected: 2,
		},
		{
			name:     "Array contains null variations",
			query:    "Attributes CONTAINS 'null'",
			expected: 2,
		},
		{
			name:     "Array contains true variations",
			query:    "Attributes CONTAINS 'TRUE'",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, items)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
				t.Logf("Results:")
				for _, r := range results {
					t.Logf("  - ID: %s, Description: '%s'", r.ID, r.Description)
				}
			}
		})
	}
}

// TestCaseInsensitiveUnicode tests case insensitive comparisons with Unicode characters
func TestCaseInsensitiveUnicode(t *testing.T) {
	type Book struct {
		Title  string
		Author string
		Tags   []string
	}

	books := []Book{
		{
			Title:  "Über den Wolken",
			Author: "Müller",
			Tags:   []string{"Über", "Wolken", "German"},
		},
		{
			Title:  "ÜBER DEN WOLKEN",
			Author: "MÜLLER",
			Tags:   []string{"ÜBER", "WOLKEN", "GERMAN"},
		},
		{
			Title:  "über den wolken",
			Author: "müller",
			Tags:   []string{"über", "wolken", "german"},
		},
		{
			Title:  "Café Life",
			Author: "José",
			Tags:   []string{"Café", "Life", "Spanish"},
		},
		{
			Title:  "CAFÉ LIFE",
			Author: "JOSÉ",
			Tags:   []string{"CAFÉ", "LIFE", "SPANISH"},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Unicode German umlaut - case insensitive",
			query:    "Author = 'müller'",
			expected: 3,
		},
		{
			name:     "Unicode German umlaut uppercase query",
			query:    "Author = 'MÜLLER'",
			expected: 3,
		},
		{
			name:     "Unicode Spanish accent - case insensitive",
			query:    "Author = 'josé'",
			expected: 2,
		},
		{
			name:     "Unicode in CONTAINS",
			query:    "Title CONTAINS 'über'",
			expected: 3,
		},
		{
			name:     "Unicode in array CONTAINS",
			query:    "Tags CONTAINS 'café'",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(tt.query, books)
			if err != nil {
				t.Fatalf("Error parsing query '%s': %v", tt.query, err)
			}

			if len(results) != tt.expected {
				t.Errorf("Query '%s' returned %d results, expected %d", tt.query, len(results), tt.expected)
				t.Logf("Results:")
				for _, r := range results {
					t.Logf("  - Title: %s, Author: %s", r.Title, r.Author)
				}
			}
		})
	}
}