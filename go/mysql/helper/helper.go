//
// helper.go
//
package helper

import (
    "context"
    "database/sql"
    "database/sql/driver"
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "strconv"
    "strings"
    "sync"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"
)


var (
    ErrDuplicateKey = errors.New("duplicate key")
)


type IdCounter struct {
    mu sync.Mutex
    id uint64
}


type Condition struct {
    Column                   string      `json:"column"`
    Value                    interface{} `json:"value"`
    LikeFlag                 bool        `json:"likeFlag"`
    InFlag                   bool        `json:"inFlag"`
    LessThanFlag             bool        `json:"lessThanFlag"`
    LessThanOrEqualToFlag    bool        `json:"lessThanOrEqualToFlag"`
    GreaterThanFlag          bool        `json:"greaterThanFlag"`
    GreaterThanOrEqualToFlag bool        `json:"greaterThanOrEqualToFlag"`
    EqualToFlag              bool        `json:"equalToFlag"`
    NotEqualToFlag           bool        `json:"notEqualToFlag"`
}


type UpdateField[T any] struct {
    SetNull bool `json:"setNull"`
    Value   T    `json:"value"`
}


type Executor interface {
    Exec(query string, args ...any) (sql.Result, error)
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
    Query(query string, args ...any) (*sql.Rows, error)
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRow(query string, args ...any) *sql.Row
    QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}


//
// Append assignment.
//
func AppendUpdateAssignment[T any](assignments *[]string, args *[]any, col string, f *UpdateField[T]) {
    if f == nil {
        return
    }
    *assignments = append(*assignments, col + " = ?")
    if f.SetNull {
        *args = append(*args, nil)
    } else {
        *args = append(*args, f.Value)
    }
}


//
// Generate query for where clause.
//
func GenerateQueryWhere(conditions []*Condition) (string, []interface{}, error) {
    if len(conditions) == 0 {
        return "", nil, nil
    }

    var queries []string
    var args []interface{}

    for _, condition := range conditions {
        if condition == nil || condition.Column == "" {
            continue
        }

        switch {
        case condition.EqualToFlag:
            queries = append(queries, fmt.Sprintf("%s = ?", condition.Column))
            args = append(args, condition.Value)
        case condition.InFlag:
            v, ok := condition.Value.(string)
            if !ok || strings.TrimSpace(v) == "" {
                return "", nil, fmt.Errorf("Invalid IN condition: column=%s\n", condition.Column)
            }
            slice := strings.Split(v, ",")
            var tmpItems []string
            for _, item := range slice {
                trimmed := strings.TrimSpace(item)
                if trimmed != "" {
                    tmpItems = append(tmpItems, trimmed)
                }
            }
            if len(tmpItems) == 0 {
                return "", nil, fmt.Errorf("Missing IN list: column=%s\n", condition.Column)
            }
            placeholders := strings.TrimRight(strings.Repeat("?,", len(tmpItems)), ",")
            queries = append(queries, fmt.Sprintf("%s IN (%s)", condition.Column, placeholders))
            for _, item := range tmpItems {
                args = append(args, item)
            }
        case condition.NotEqualToFlag:
            queries = append(queries, fmt.Sprintf("%s <> ?", condition.Column))
            args = append(args, condition.Value)
        case condition.LessThanFlag:
            queries = append(queries, fmt.Sprintf("%s < ?", condition.Column))
            args = append(args, condition.Value)
        case condition.LessThanOrEqualToFlag:
            queries = append(queries, fmt.Sprintf("%s <= ?", condition.Column))
            args = append(args, condition.Value)
        case condition.GreaterThanFlag:
            queries = append(queries, fmt.Sprintf("%s > ?", condition.Column))
            args = append(args, condition.Value)
        case condition.GreaterThanOrEqualToFlag:
            queries = append(queries, fmt.Sprintf("%s >= ?", condition.Column))
            args = append(args, condition.Value)
        default:
            return "", nil, fmt.Errorf("Missing operator flag: column=%s\n", condition.Column)
        }
    }

    var query strings.Builder
    if len(queries) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(queries, " AND "))
    }

    return query.String(), args, nil
}


//
// Truncate runes.
//
func TruncateRunes(s string, max int) string {
    if max <= 0 {
        return ""
    }

    if utf8.RuneCountInString(s) <= max {
        return s
    }

    r := []rune(s)

    if max <= 3 {
        return string(r[:max])
    }

    return string(r[:max]) + "..."
}


//
// Generate an ID.
//
// Version:
//   - 2026-04-28: Added.
//
func (idc *IdCounter) GenerateID() uint64 {
    idc.mu.Lock()
    defer idc.mu.Unlock()
    id := uint64(time.Now().UnixNano())
    if id <= idc.id {
        id = idc.id + 1
    }
    idc.id = id
    return id
}


//
// Scan uint8 value.
//
// Version:
//   - 2026-05-03: Added.
//
func ScanUint8(name string, value any) (uint8, error) {
    if value == nil {
        return 0, fmt.Errorf("missing required parameter: %s=null", name)
    }

    switch v := value.(type) {
    case uint8:
        return v, nil
    case uint16:
        if v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case uint32:
        if v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case uint64:
        if v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case uint:
        if v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case int8:
        if v < 0 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case int16:
        if v < 0 || v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case int32:
        if v < 0 || v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case int64:
        if v < 0 || v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case int:
        if v < 0 || v > math.MaxUint8 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint8(v), nil
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return 0, fmt.Errorf("invalid parameter: %s=%s", name, string(v))
        }
        return uint8(n), nil
    case string:
        n, err := strconv.ParseUint(v, 10, 8)
        if err != nil {
            return 0, fmt.Errorf("invalid parameter: %s=%s", name, v)
        }
        return uint8(n), nil
    default:
        return 0, fmt.Errorf("unsupported parameter: %s_type=%T", name, value)
    }
}


//
// Scan uint16 value.
//
// Version:
//   - 2026-07-08: Added.
//
func ScanUint16(name string, value any) (uint16, error) {
    if value == nil {
        return 0, fmt.Errorf("missing required parameter: %s=null", name)
    }

    switch v := value.(type) {
    case uint8:
        return uint16(v), nil
    case uint16:
        return v, nil
    case uint32:
        if v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case uint64:
        if v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case uint:
        if v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case int8:
        if v < 0 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case int16:
        if v < 0 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case int32:
        if v < 0 || v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case int64:
        if v < 0 || v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case int:
        if v < 0 || v > math.MaxUint16 {
            return 0, fmt.Errorf("invalid parameter: %s=%d", name, v)
        }
        return uint16(v), nil
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 16)
        if err != nil {
            return 0, fmt.Errorf("invalid parameter: %s=%s", name, string(v))
        }
        return uint16(n), nil
    case string:
        n, err := strconv.ParseUint(v, 10, 16)
        if err != nil {
            return 0, fmt.Errorf("invalid parameter: %s=%s", name, v)
        }
        return uint16(n), nil
    default:
        return 0, fmt.Errorf("unsupported parameter: %s_type=%T", name, value)
    }
}


type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
    if len(j) == 0 {
        return nil, nil
    }

    if err := j.Validate(); err != nil {
        return nil, err
    }

    return []byte(j), nil
}

func (j *JSON) Scan(src any) error {
    if j == nil {
        return fmt.Errorf("missing required parameter: json=null")
    }

    if src == nil {
        *j = nil
        return nil
    }

    var b []byte

    switch v := src.(type) {
    case []byte:
        b = append([]byte(nil), v...)
    case string:
        b = []byte(v)
    default:
        return fmt.Errorf("unsupported parameter: json: type=%T", src)
    }

    *j = JSON(b)

    if err := j.Validate(); err != nil {
        return err
    }

    return nil
}


func (j *JSON) Validate() error {
    if j == nil {
        return fmt.Errorf("missing required parameter: json=null")
    }
    if len(*j) == 0 {
        return fmt.Errorf("invalid parameter: json=%q", "empty")
    }
    if len(*j) > 4096 {
        return fmt.Errorf("invalid parameter: json=%q max_bytes=4096", "too long")
    }
    if !json.Valid([]byte(*j)) {
        return fmt.Errorf("invalid parameter: json=%q", string(*j))
    }

    return nil
}

func (j JSON) RawMessage() json.RawMessage {
    return json.RawMessage(j)
}
