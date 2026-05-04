//
// otp_challenge.go
//
package email

import (
    "database/sql"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultOTPChallengeTableName = "account_email_otp_challenges"
)

var (
    otpChallengeIDCounter = &helper.IdCounter{}
)


//
// OTPChallenge.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPChallenge struct {
    ID                 uint64     `json:"id,string"`
    Email              string     `json:"email"`
    Purpose            OTPPurpose `json:"purpose"`
    RequestedIP        *string    `json:"requestedIp,omitempty"`
    RequestedUserAgent *string    `json:"requestedUserAgent,omitempty"`
    CreatedAt          time.Time  `json:"createdAt"`
}


//
// OTPChallengeStore.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPChallengeStore struct {
    db        *sql.DB
    tableName string
}


//
// OTPChallengeSelectOption.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPChallengeSelectOption struct {
    ID           *uint64     `json:"id,string,omitempty"`
    Email        *string     `json:"email,omitempty"`
    EmailLike    *string     `json:"emailLike,omitempty"`
    Purpose      *OTPPurpose `json:"status,omitempty"`
    CreatedAtGTE *time.Time  `json:"createdAtGte,omitempty"`
    CreatedAtLTE *time.Time  `json:"createdAtLte,omitempty"`
    OrderBy      string      `json:"orderBy"`
    OrderByDesc  bool        `json:"orderByDesc"`
    Limit        int         `json:"limit"`
    Offset       int         `json:"offset"`
}


//
// Generate OTP challenge ID.
//
// Version:
//   - 2026-05-03: Added.
//
func GenerateOTPChallengeID() uint64 {
    return otpChallengeIDCounter.GenerateID()
}


//
// Create new OTP challenge store.
//
// Version:
//   - 2026-05-03: Added.
//
func NewOTPChallengeStore(db *sql.DB, tableName string) (*OTPChallengeStore, error) {
    if db == nil {
        return nil, fmt.Errorf("failed to create account email otp challenge store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account email otp challenge store: missing required parameter: table_name=empty")
    }

    return &OTPChallengeStore{
        db:        db,
        tableName: tableName,
    }, nil
}


//
// Validate ID.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *OTPChallenge) ValidateID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: otp_challenge=null")
    }
    if c.ID == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate email.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *OTPChallenge) ValidateEmail() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: otp_challenge=null")
    }
    if c.Email == "" {
        return fmt.Errorf("invalid parameter: email=empty")
    }
    if utf8.RuneCountInString(c.Email) > 255 {
        return fmt.Errorf("invalid parameter: email=%q", helper.TruncateRunes(c.Email, 255))
    }
    return nil
}


//
// Validate OTP purpose.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *OTPChallenge) ValidatePurpose() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: otp_challenge=null")
    }
    if !c.Purpose.IsValid() {
        return fmt.Errorf("invalid parameter: otp_purpose=%d", c.Purpose)
    }
    return nil
}


//
// Validate requested IP.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *OTPChallenge) ValidateRequestedIP() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: otp_challenge=null")
    }
    if c.RequestedIP == nil {
        return nil
    }
    if *c.RequestedIP == "" {
        return fmt.Errorf("invalid parameter: requested_ip=empty")
    }
    if utf8.RuneCountInString(*c.RequestedIP) > 45 {
        return fmt.Errorf("invalid parameter: requested_ip=%q", helper.TruncateRunes(*c.RequestedIP, 45))
    }
    return nil
}


//
// Validate requested user agent.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *OTPChallenge) ValidateRequestedUserAgent() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: otp_challenge=null")
    }
    if c.RequestedUserAgent == nil {
        return nil
    }
    if *c.RequestedUserAgent == "" {
        return fmt.Errorf("invalid parameter: requested_user_agent=empty")
    }
    if utf8.RuneCountInString(*c.RequestedUserAgent) > 255 {
        return fmt.Errorf("invalid parameter: requested_user_agent=%q", helper.TruncateRunes(*c.RequestedUserAgent, 255))
    }
    return nil
}


//
// Create account email OTP challenges table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *OTPChallengeStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account email otp challenges table: missing required parameter: otp_challenge_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account email otp challenges table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account email otp challenges table: missing required parameter: table_name=empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s VARCHAR(255) NOT NULL COMMENT 'Email',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Purpose',
            %s VARCHAR(45) NULL COMMENT 'RequestedIP',
            %s VARCHAR(255) NULL COMMENT 'RequestedUserAgent',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            PRIMARY KEY (%s),
            KEY idx_account_app_email_otp_challenge_email (%s),
            KEY idx_account_app_email_otp_challenge_purpose (%s)
        );`,
        s.tableName,
        ColID,
        ColEmail,
        ColPurpose,
        ColRequestedIP,
        ColRequestedUserAgent,
        ColCreatedAt,
        ColID,
        ColEmail,
        ColPurpose,
    )

    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create account email otp challenges table: %w", err)
    }

    return nil
}


//
// Insert account email OTP challenge.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPChallengeStore) Insert(row *OTPChallenge) error {
    if s == nil {
        return fmt.Errorf("failed to insert account email otp challenge: missing required parameter: otp_challenge_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account email otp challenge: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account email otp challenge: missing required parameter: table_name=empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account email otp challenge: missing required parameter: otp_challenge=null")
    }
    if err := row.ValidateEmail(); err != nil {
        return fmt.Errorf("failed to insert account email otp challenge: %w", err)
    }
    if err := row.ValidatePurpose(); err != nil {
        return fmt.Errorf("failed to insert account email otp challenge: %w", err)
    }
    if err := row.ValidateRequestedIP(); err != nil {
        return fmt.Errorf("failed to insert account email otp challenge: %w", err)
    }
    if err := row.ValidateRequestedUserAgent(); err != nil {
        return fmt.Errorf("failed to insert account email otp challenge: %w", err)
    }

    if row.ID == 0 {
        row.ID = GenerateOTPChallengeID()
    }

    now := time.Now()
    if row.CreatedAt.IsZero() {
        row.CreatedAt = now
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColEmail,
        ColPurpose,
        ColRequestedIP,
        ColRequestedUserAgent,
        ColCreatedAt,
    )

    if _, err := s.db.Exec(
        query,
        row.ID,
        row.Email,
        row.Purpose,
        row.RequestedIP,
        row.RequestedUserAgent,
        row.CreatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account email otp challenge: %w", err)
    }

    return nil
}


//
// Select account email otp challenge by ID.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPChallengeStore) SelectByID(id uint64) (*OTPChallenge, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email otp challenge by id: missing required parameter: otp_challenge_store=null")
    }   
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account email otp challenge by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email otp challenge by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account email otp challenge by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &OTPChallenge{}
    err := s.db.QueryRow(query, id).Scan(
        &result.ID,
        &result.Email,
        &result.Purpose,
        &result.RequestedIP,
        &result.RequestedUserAgent,
        &result.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account email otp challenge by id: %w", err)
    }

    return result, nil
}


//
// Select account email otp challenges.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPChallengeStore) Select(option *OTPChallengeSelectOption) ([]*OTPChallenge, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email otp challenges: missing required parameter: otp_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account email otp challenges: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email otp challenges: missing required parameter: table_name=empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account email otp challenges: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    rows, err := s.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account email otp challenges: %w", err)
    }
    defer rows.Close()

    var result []*OTPChallenge
    for rows.Next() {
        row := &OTPChallenge{}
        if err := rows.Scan(
            &row.ID,
            &row.Email,
            &row.Purpose,
            &row.RequestedIP,
            &row.RequestedUserAgent,
            &row.CreatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select account email otp challenges: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account email otp challenges: %w", err)
    }

    return result, nil
}


//
// Count account email otp challenges.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPChallengeStore) Count(option *OTPChallengeSelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account email otp challenges: missing required parameter: otp_challenge_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account email otp challenges: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account email otp challenges: missing required parameter: table_name=empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account email otp challenges: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.db.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account email otp challenges: %w", err)
    }

    return result, nil
}


// 
// Delete account email otp challenge by ID.
// 
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPChallengeStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account email otp challenge by id: missing required parameter: otp_challenge_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account email otp challenge by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account email otp challenge by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account email otp challenge by id: invalid parameter: id=0")
    }
    
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.db.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account email otp challenge by id: %w", err)
    }   
    
    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-03: Added.
//
func (o *OTPChallengeSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

    if o.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *o.ID)
    }
    if o.Email != nil {
        conditions = append(conditions, ColEmail + " = ?")
        args = append(args, *o.Email)
    }
    if o.EmailLike != nil {
        conditions = append(conditions, ColEmail + " LIKE ?")
        args = append(args, "%" + *o.EmailLike + "%")
    }
    if o.Purpose != nil {
        conditions = append(conditions, ColPurpose + " = ?")
        args = append(args, *o.Purpose)
    }
    if o.CreatedAtGTE != nil {
        conditions = append(conditions, ColCreatedAt + " >= ?")
        args = append(args, *o.CreatedAtGTE)
    }
    if o.CreatedAtLTE != nil {
        conditions = append(conditions, ColCreatedAt + " <= ?")
        args = append(args, *o.CreatedAtLTE)
    }

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }
    
    if o.OrderBy != "" {
        query.WriteString(" ORDER BY ")
        query.WriteString(o.OrderBy)
        if o.OrderByDesc {
            query.WriteString(" DESC")
        }
    }

    if o.Limit > 0 {
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, o.Limit, o.Offset)
    }

    return query.String(), args
}


//
// Validate account email otp challenge select option.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *OTPChallengeSelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: otp_challenge_select_option=null")
    }

    if o.ID != nil {
        if err := (&OTPChallenge{ID: *o.ID}).ValidateID(); err != nil {
            return err
        }
    }

    if o.Email != nil {
        if err := (&OTPChallenge{Email: *o.Email}).ValidateEmail(); err != nil {
            return err
        }
    }

    if o.EmailLike != nil {
        if err := (&OTPChallenge{Email: *o.EmailLike}).ValidateEmail(); err != nil {
            return err
        }
    }

    if o.Purpose != nil {
        if err := (&OTPChallenge{Purpose: *o.Purpose}).ValidatePurpose(); err != nil {
            return err
        }
    }

    if o.CreatedAtGTE != nil && o.CreatedAtLTE != nil && o.CreatedAtGTE.After(*o.CreatedAtLTE) {
        return fmt.Errorf("invalid parameter: created_at_gte=%s created_at_lte=%s", *o.CreatedAtGTE, *o.CreatedAtLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColID,
            ColEmail,
            ColPurpose,
            ColCreatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%s", o.OrderBy)
        }
    }

    if o.Limit < 0 {
        return fmt.Errorf("invalid parameter: limit=%d", o.Limit)
    }
    if o.Offset < 0 {
        return fmt.Errorf("invalid parameter: offset=%d", o.Offset)
    }

    return nil
}
