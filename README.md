# Parser 🧠

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zveinn/parser)
![GitHub Issues](https://img.shields.io/github/issues/zveinn/parser)
![GitHub Stars](https://img.shields.io/github/stars/zveinn/parser?style=social)

**Parser** is a lightweight, efficient Go library for filtering slices of structs using a SQL-like query language. It enables in-memory data filtering without a database, ideal for applications like data processing, configuration management, or API response filtering. Built with Go generics, it offers type-safe queries and supports complex expressions, nested fields, and case-insensitive matching.

## 🚀 Features

- **SQL-Like Query Language**: Filter structs with intuitive queries (e.g., `Age > 25 AND Skills CONTAINS 'Go'`).
- **Type-Safe with Generics**: Works with any struct type using Go’s generics.
- **Nested Field Access**: Query nested structs and maps using dot notation (e.g., `Department.Name`).
- **Humanized Values Support**: Parse human-readable values like time units (`10m`, `2h30m`), byte units (`10GB`/`10GiB`, `2TB`/`2TiB`), SI prefixes (`1.5K`, `2.3M`), and comma-separated numbers (`1,000`) automatically.
- **Rich Operators**: Supports `=`, `!=`, `<`, `>`, `<=`, `>=`, `CONTAINS`, `IS NULL`, `ANY`, `NOT`, `AND`, `OR`.
- **Case-Insensitive Matching**: Field names, keywords (e.g., `AND`, `OR`), and string value comparisons are all case-insensitive.
- **Efficient Parsing**: Uses an enhanced lexer with support for negative numbers, scientific notation, and comma-separated numbers.
- **Robust Error Handling**: Detailed error messages for syntax and evaluation errors.
- **Zero Dependencies**: Pure Go implementation with built-in support for time, byte, and SI unit parsing.

## 📋 Requirements

- Go 1.24.1 or higher (for generics and latest features)
- No external dependencies

## 📦 Installation

Install the library using Go modules:

```bash
go get github.com/zveinn/parser
```

## 🔧 Usage

### Basic Example

Filter a slice of structs using a query:

```go
package main

import (
    "fmt"
    "log"
    "github.com/zveinn/parser"
)

type Person struct {
    Name       string
    Age        int
    IsEmployed bool
    Skills     []string
    Salary     float64
    Department *Department
}

type Department struct {
    Name     string
    Location string
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 30, IsEmployed: true, Skills: []string{"Go", "Python"}, Salary: 75000.50, Department: &Department{Name: "Engineering", Location: "New York"}},
        {Name: "Bob", Age: 25, IsEmployed: false, Skills: []string{"Java", "C++"}, Salary: 65000.25, Department: &Department{Name: "Engineering", Location: "Remote"}},
        {Name: "Charlie", Age: 35, IsEmployed: true, Skills: []string{"Go", "Rust"}, Salary: 85000.75, Department: nil},
    }

    results, err := parser.Parse("Age > 25 AND isemployed = true", people)
    if err != nil {
        log.Fatalf("Error parsing query: %v", err)
    }

    for _, p := range results {
        fmt.Printf("Match: %s (Age: %d)\n", p.Name, p.Age)
    }
}
```

**Output**:
```
Match: Alice (Age: 30)
Match: Charlie (Age: 35)
```

### Query Syntax

The query language supports a variety of operators and expressions:

#### Case-Insensitive String Matching

All string comparisons are **case-insensitive** by default. This applies to:
- Equality comparisons (`=`, `!=`)
- Ordering comparisons (`<`, `>`, `<=`, `>=`)
- Contains operations (`CONTAINS`)
- Array/slice element matching

**Examples:**
```sql
# These all match "Apple", "APPLE", "apple", etc.
Name = 'apple'
Name = 'APPLE'
Name = 'Apple'

# Case-insensitive NOT EQUAL
Name != 'samsung'  # Excludes "Samsung", "SAMSUNG", "samsung", etc.

# Case-insensitive CONTAINS
Description CONTAINS 'phone'  # Matches "iPhone", "PHONE", "Phone", etc.

# Case-insensitive in arrays
Tags CONTAINS 'premium'  # Matches array elements like "Premium", "PREMIUM", "premium"

# Case-insensitive ordering
Brand < 'b'  # "Apple", "APPLE", "apple" all evaluate as less than "b"
```

#### Comparison Operators
| Operator   | Description              | Example                |
|------------|--------------------------|------------------------|
| `=`        | Equal                    | `Name = 'Alice'`       |
| `!=`       | Not equal                | `Age != 30`            |
| `>`        | Greater than             | `Salary > 70,000`      |
| `<`        | Less than                | `Age < 35`             |
| `>=`       | Greater than or equal    | `Salary >= 75000.50`   |
| `<=`       | Less than or equal       | `Age <= 30`            |
| `CONTAINS` | String or slice contains | `Skills CONTAINS 'Go'` |

#### Logical Operators
| Operator | Description | Example                          |
|----------|-------------|----------------------------------|
| `AND`    | Logical AND | `Age > 25 AND IsEmployed = true` |
| `OR`     | Logical OR  | `Name = 'Alice' OR Name = 'Bob'` |
| `NOT`    | Logical NOT | `NOT (Age < 30)`                 |

#### Special Operators
| Operator      | Description               | Example                           |
|---------------|---------------------------|-----------------------------------|
| `IS NULL`     | Check for nil/zero value  | `Department IS NULL`              |
| `IS NOT NULL` | Check for non-nil value   | `Department IS NOT NULL`          |
| `ANY`         | Match any value in a list | `ANY(Skills) = ANY('Go', 'Rust')` |

#### Example Queries
```sql
# Basic filtering
Name = 'Alice'
Salary > 80,000
Skills CONTAINS 'Go'

# Time-based filtering (converted to seconds)
ResponseTime < 30s
Timeout > 5m
CacheExpiry < 2h
Uptime > 1d

# Byte size filtering
Memory > 8GB
Storage < 1TiB
BackupSize > 500MiB

# SI prefix filtering (uppercase only)
Population > 1.5M
Records < 10K
Distance >= 2.5G

# Nested fields and maps
Department.Location = 'Remote'
Tags.level = 'senior'

# Complex logic with mixed units
(Age > 30 AND Salary > 75,000) OR IsEmployed = false
ResponseTime < 1m AND Memory > 8GB AND Uptime > 1d
ANY(Skills) = 'Go' AND NOT (Department IS NULL)
```

### Advanced Usage

#### Nested Structs and Maps
Query nested fields or map values using dot notation:

```go
query := "Department.Name = 'Engineering' AND Tags.level = 'senior'"
results, err := parser.Parse(query, people)
```

#### Numeric Formats
The parser supports advanced numeric formats:
- Negative numbers: `Salary > -1000`
- Scientific notation: `Salary > 7.5e4`
- Comma-separated numbers: `Salary > 1,000,000.50`
- Time durations: `ResponseTime < 30s`, `Timeout > 2h30m`
- Byte sizes: `Memory > 8GB`, `Storage < 1TiB`
- SI prefixes: `Population > 1.5M`, `Count < 5K` (uppercase only)

#### Humanized Values
The parser automatically converts humanized values to their numeric equivalents with unambiguous parsing rules. Values are parsed in the following priority order:

1. **Time Duration Units** (parsed first to avoid conflicts)
2. **Byte Size Units** (decimal and binary)
3. **SI Prefixes** (case-sensitive, uppercase only)
4. **Comma-Separated Numbers**

**Time Duration Units:**
Time units are converted to total seconds and support compound durations:
```sql
# Single time units
ResponseTime < 30s           # 30 seconds
Timeout > 5m                 # 300 seconds (5 minutes)
CacheExpiry < 2h             # 7200 seconds (2 hours)
Retention > 7d               # 604800 seconds (7 days)

# Compound time units (multiple units combined)
Duration = 2h30m             # 9000 seconds (2 hours + 30 minutes)
Delay < 1m30s                # 90 seconds (1 minute + 30 seconds)
Uptime > 1d12h               # 129600 seconds (1 day + 12 hours)

# Supported time units:
# ns - nanoseconds, us/µs - microseconds, ms - milliseconds
# s - seconds, m - minutes, h - hours, d - days
```

**Byte Size Units (Decimal and Binary):**
```sql
# Decimal units (powers of 1000) - International System of Units
Drive.Size > 10GB            # 10,000,000,000 bytes
Memory > 1.5TB               # 1,500,000,000,000 bytes  
Storage < 500MB              # 500,000,000 bytes
Buffer < 100KB               # 100,000 bytes

# Binary units (powers of 1024) - Computer memory standards
Backup > 2.5GiB              # 2,684,354,560 bytes
Cache > 512MiB               # 536,870,912 bytes
Temp < 100KiB                # 102,400 bytes
Archive > 1TiB               # 1,099,511,627,776 bytes

# Supported byte units:
# Decimal: B, KB, MB, GB, TB, PB, EB, ZB, YB
# Binary: B, KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB
```

**SI Prefixes (Case-Sensitive, Uppercase Only):**
SI prefixes are now case-sensitive and only recognize uppercase letters to avoid conflicts with time units:
```sql
Population > 1.5M            # 1,500,000 (mega = 10^6)
Count < 5K                   # 5,000 (kilo = 10^3)
Records >= 2.3G              # 2,300,000,000 (giga = 10^9)
Distance < 500K              # 500,000 (kilo = 10^3)

# Supported SI prefixes (uppercase only):
# K (kilo = 10^3), M (mega = 10^6), G (giga = 10^9)
# T (tera = 10^12), P (peta = 10^15), E (exa = 10^18)
# Z (zettabyte = 10^21), Y (yottabyte = 10^24)

# Note: Lowercase prefixes (k, m, g, etc.) are NOT supported
# to avoid conflicts with time units (m = minutes, s = seconds)
```

**Comma-Separated Numbers:**
```sql
Price > 1,000,000            # 1000000
Users >= 50,000              # 50000
Transactions < 2,500         # 2500
```

**Example with Real Data:**
```go
type Server struct {
    Name         string
    Memory       int64  // in bytes
    Storage      int64  // in bytes
    ResponseTime int64  // in seconds
    Uptime       int64  // in seconds
}

servers := []Server{
    {Name: "web1", Memory: 8589934592, Storage: 536870912000, ResponseTime: 30, Uptime: 86400},    // 8GB, 500GB, 30s, 1 day
    {Name: "db1", Memory: 34359738368, Storage: 2199023255552, ResponseTime: 120, Uptime: 604800}, // 32GB, 2TB, 2m, 7 days
}

// Query using time units (converted to seconds)
results, _ := parser.Parse("ResponseTime < 1m AND Uptime > 1d", servers)

// Query using decimal byte units (powers of 1000)
results, _ := parser.Parse("Memory > 16GB AND Storage < 1TB", servers)

// Query using binary byte units (powers of 1024) 
results, _ := parser.Parse("Memory > 16GiB AND Storage < 1TiB", servers)

// Query using SI prefixes (case-sensitive, uppercase only)
results, _ := parser.Parse("Memory > 8G AND Storage > 500M", servers) // Treating as generic numbers

// Mixed units work correctly due to unambiguous parsing
results, _ := parser.Parse("ResponseTime < 2m AND Memory > 8GB AND Uptime > 1d", servers)
```

**Important Notes:**
- Time units (`m`, `s`, `h`, `d`) take precedence over SI prefixes
- SI prefixes are case-sensitive and only recognize uppercase (`K`, `M`, `G`, etc.)
- Byte units support both decimal (GB, MB) and binary (GiB, MiB) standards
- The parser automatically resolves conflicts by checking units in priority order

### Parsing Rules and Conflict Resolution

The parser uses a priority-based system to handle potential conflicts between different unit types:

1. **Time Units First**: `10m` is always parsed as 10 minutes (600 seconds), never as 10 milli-units
2. **Byte Units Second**: `10GB` is parsed as 10 gigabytes (10,000,000,000 bytes)  
3. **SI Prefixes Third**: `10K` is parsed as 10,000 using the kilo prefix
4. **Comma-Separated Last**: `10,000` is parsed as ten thousand

**Case Sensitivity Rules:**
- Time units are case-insensitive: `10M` = `10m` = 10 minutes
- Byte units are case-sensitive: `10GB` ≠ `10gb` (only `10GB` is valid)
- SI prefixes are case-sensitive: `10K` is valid, `10k` is not supported
- This prevents conflicts like `m` (minutes) vs `m` (milli-prefix)

**Examples of Conflict Resolution:**
```sql
# These are unambiguous and work as expected:
Duration < 5m              # 5 minutes = 300 seconds (time unit)
Size > 5MB                 # 5 megabytes = 5,000,000 bytes (byte unit)  
Count > 5K                 # 5 thousand = 5,000 (SI prefix)

# These demonstrate the priority system:
Value > 10m                # Always 10 minutes (600 seconds), never 10 milli-units
Storage > 10M              # 10 megabytes if comparing to bytes, otherwise 10 million
Population > 10M           # 10 million (SI prefix) when comparing to numbers
```

### Performance Considerations
Based on benchmark results:
- **Efficient for Small to Medium Datasets**: Queries on datasets of 10–1000 structs are fast, with simple queries (e.g., `Age > 30`) taking microseconds.
- **Unit Parsing Overhead**: Time, byte, and SI unit parsing adds minimal overhead and is optimized for common cases.
- **Reflection Overhead**: Minimal reflection is used during evaluation, with no reflection during query compilation.
- **Scalability**: Performance scales linearly with dataset size. For very large datasets (>10,000 items), consider batching.
- **Query Complexity**: Complex queries with nested logic or `ANY` operators are slightly slower but optimized with short-circuit evaluation.
- **Memory Usage**: Low memory footprint, with minimal allocations for simple queries (benchmarks show 1–2 allocations per query).

### Error Handling
The parser provides detailed error messages:

```go
_, err := parser.Parse("Age >", people)
if err != nil {
    fmt.Println(err) // Output: "failed to parse query: unexpected EOF"
}

_, err = parser.Parse("InvalidField = 10", people)
if err != nil {
    fmt.Println(err) // Output: "evaluation error: field 'InvalidField' not found"
}
```

## 🛠️ Building and Testing

Clone the repository and build:

```bash
git clone https://github.com/zveinn/parser.git
cd parser
```

Run tests to verify functionality:

```bash
go test -v ./...
```

Run benchmarks to measure performance:

```bash
go test -bench=. ./...
```

## 📚 Documentation

- **API Reference**: Available via [GoDoc](https://pkg.go.dev/github.com/zveinn/parser).
- **Examples**: See the [examples/](examples/) directory for sample queries (create this directory if needed).
- **Source Code Insights**:
  - `parser.go`: Core parsing logic with AST evaluation.
  - `enhanced_lexer.go`: Tokenization with support for advanced numeric formats.
  - `parser_test.go`: Comprehensive test suite for all operators and edge cases.

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/my-feature`).
3. Commit your changes (`git commit -m 'Add my feature'`).
4. Push to the branch (`git push origin feature/my-feature`).
5. Open a Pull Request.

## 📜 License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

## 🌟 Acknowledgements

- Built by [zveinn](https://github.com/zveinn).
- Inspired by SQL query engines and libraries like [rql](https://github.com/dvaldivia/rql).
- Thanks to the Go community for feedback and inspiration.

---

⭐ **Star this project** if you find it useful!  
💬 Report issues or suggest features in [Issues](https://github.com/zveinn/parser/issues).
