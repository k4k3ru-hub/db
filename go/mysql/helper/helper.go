//
// helper.go
//
package helper

import (
    "fmt"
    "strings"
    "sync"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"
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


