//
// otp.go
//
package email

import (
    "crypto/rand"
	"crypto/sha256"
    "database/sql"
	"encoding/hex"
    "fmt"
    "math/big"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultOTPTableName = "account_email_otps"

    DefaultCodeLength = 6

    DefaultMaxAttemptCount uint8 = 3

    DefaultExpiresIn = 10 * time.Minute

    DefaultLockedUntilIn = 15 * time.Minute
)


//
// OTP.
//
// Version:
//   - 2026-05-03: Added.
//
type OTP struct {
    Email        string     `json:"email"`
    Purpose      OTPPurpose `json:"purpose"`
    Status       OTPStatus  `json:"status"`
    CodeHash     string     `json:"codeHash"`
    ExpiresAt    time.Time  `json:"expiresAt"`
    AttemptCount uint8      `json:"attemptCount"`
    LastSentAt   time.Time  `json:"lastSentAt"`
    LockedUntil  *time.Time `json:"lockedUntil,omitempty"`
}


//
// OTPStore.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPStore struct {
    executor  helper.Executor
    tableName string
}


//
// OTPSelectOption.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPSelectOption struct {
    Email        *string     `json:"email,omitempty"`
    EmailLike    *string     `json:"emailLike,omitempty"`
    Purpose      *OTPPurpose `json:"purpose,omitempty"`
    Status       *OTPStatus  `json:"status,omitempty"`
    ExpiresAtGTE *time.Time  `json:"expiresAtGte,omitempty"`
    ExpiresAtLTE *time.Time  `json:"expiresAtLte,omitempty"`
    OrderBy      string      `json:"orderBy"`
    OrderByDesc  bool        `json:"orderByDesc"`
    Limit        int         `json:"limit"`
    Offset       int         `json:"offset"`
}


//
// OTPUpdateOption.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPUpdateOption struct {
    Email              string     `json:"email"`
    Purpose            OTPPurpose `json:"purpose"`
    Status             *OTPStatus `json:"status,omitempty"`
    CodeHash           *string    `json:"codeHash,omitempty"`
    ExpiresAt          *time.Time `json:"expiresAt,omitempty"`
    AttemptCount       *uint8     `json:"attemptCount,omitempty"`
    LastSentAt         *time.Time `json:"lastSentAt,omitempty"`
    LockedUntil        *time.Time `json:"lockedUntil,omitempty"`
    LockedUntilSetNull bool       `json:"lockedUntilSetNull"`
}


//
// Get default expires at.
//
// Version:
//   - 2026-05-04: Added.
//
func DefaultExpiresAt() time.Time {
    return time.Now().UTC().Add(DefaultExpiresIn)
}


//
// Generate code.
//
// Version:
//   - 2026-05-04: Added.
//
func GenerateCode(codeLength int) (string, error) {
	// Calculate max value (10^codeLength).
	max := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(int64(codeLength)),
		nil,
	)

	// Generate secure random number.
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	// Zero padding (e.g. 000123).
	return fmt.Sprintf("%0*d", codeLength, n.Uint64()), nil
}


//
// Hash code.
//
// Version:
//   - 2026-05-04: Added.
//
func HashCode(code string, maxCodeLength int) (string, error) {
	// Normalize.
	code = strings.TrimSpace(code)
	if code == "" {
		return "", fmt.Errorf("failed to hash code: missing required parameter: code=empty")
	}
    if len(code) > maxCodeLength {
        return "", fmt.Errorf("failed to hash code: invalid parameter: code=%q", helper.TruncateRunes(code, maxCodeLength))
    }

	// Hash.
	sum := sha256.Sum256([]byte(code))

	return hex.EncodeToString(sum[:]), nil
}


//
// Create new OTP store.
//
// Version:
//   - 2026-05-03: Added.
//
func NewOTPStore(executor helper.Executor, tableName string) (*OTPStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create account email otp store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account email otp store: missing required parameter: table_name=%q", "empty")
    }

    return &OTPStore{
        executor:  executor,
        tableName: tableName,
    }, nil
}


//
// Validate account email OTP email.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateOTPEmail(email string) error {
    if email == "" {
        return fmt.Errorf("email=%q", "empty")
    }
    if utf8.RuneCountInString(email) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 email=%q", "too long")
    }
    return nil
}


//
// Validate account email OTP email.
//
// Version:
//   - 2026-05-12: Added.
//
func (o *OTP) ValidateEmail() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPEmail(o.Email)
}


//
// Validate account email OTP purpose.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateOTPPurpose(purpose OTPPurpose) error {
    if !purpose.IsValid() {
        return fmt.Errorf("invalid parameter: purpose=%q", purpose)
    }
    return nil
}


//
// Validate account email OTP purpose.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidatePurpose() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPPurpose(o.Purpose)
}


//
// Validate account email OTP status.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateOTPStatus(status OTPStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate account email OTP status.
//
// Version:
//   - 2026-05-09: Added.
//
func (o *OTP) ValidateStatus() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPStatus(o.Status)
}


//
// Validate account email OTP code hash.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateOTPCodeHash(codeHash string) error {
    if codeHash == "" {
        return fmt.Errorf("invalid parameter: code_hash=%q", "empty")
    }
    if utf8.RuneCountInString(codeHash) > 128 {
        return fmt.Errorf("invalid parameter: max_length=128 code_hash=%q", "too long")
    }
    return nil
}


//
// Validate account email OTP code hash.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateCodeHash() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPCodeHash(o.CodeHash)
}


//
// Validate account email OTP expires at.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateOTPExpiresAt(expiresAt time.Time) error {
    if expiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=%q", "empty")
    }
    if !expiresAt.UTC().After(time.Now().UTC()) {
        return fmt.Errorf("invalid parameter: expires_at=%q", expiresAt.Format(time.RFC3339))
    }
    return nil
}


//
// Validate account email OTP expires at.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateExpiresAt() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPExpiresAt(o.ExpiresAt)
}


//
// Validate account email OTP last sent at.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateOTPLastSentAt(lastSentAt time.Time) error {
    if lastSentAt.IsZero() {
        return fmt.Errorf("invalid parameter: last_sent_at=%q", "empty")
    }
    if lastSentAt.UTC().After(time.Now().UTC()) {
        return fmt.Errorf("invalid parameter: last_sent_at=%q", lastSentAt.Format(time.RFC3339))
    }
    return nil
}


//
// Validate account email OTP last sent at.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateLastSentAt() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPLastSentAt(o.LastSentAt)
}


//
// Validate account email OTP locked until.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateOTPLockedUntil(lockedUntil *time.Time) error {
    if lockedUntil == nil {
        return nil
    }
    if lockedUntil.IsZero() {
        return fmt.Errorf("invalid parameter: locked_until=%q", "empty")
    }
    if !lockedUntil.UTC().After(time.Now().UTC()) {
        return fmt.Errorf("invalid parameter: locked_until=%s", lockedUntil.Format(time.RFC3339))
    }
    return nil
}


//
// Validate account email OTP locked until.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateLockedUntil() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_email_otp=null")
    }
    return ValidateOTPLockedUntil(o.LockedUntil)
}


//
// Create account email OTPs table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *OTPStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account email otps table: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create account email otps table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account email otps table: missing required parameter: table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s VARCHAR(255) NOT NULL COMMENT 'Email',
            %s VARCHAR(32) NOT NULL COMMENT 'Purpose',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(255) NOT NULL COMMENT 'Code hash',
            %s DATETIME NOT NULL COMMENT 'Expires at',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Attempt count',
            %s DATETIME NOT NULL COMMENT 'Last sent at',
            %s DATETIME NULL COMMENT 'Locked until',
            PRIMARY KEY (%s, %s),
            KEY idx_account_app_email_otp_status (%s)
        );`,
        s.tableName,
        ColEmail,
        ColPurpose,
        ColStatus,
        ColCodeHash,
        ColExpiresAt,
        ColAttemptCount,
        ColLastSentAt,
        ColLockedUntil,
        ColEmail, ColPurpose,
        ColStatus,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account email otps table: %w", err)
    }

    return nil
}


//
// Insert account email OTP.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Insert(row *OTP) error {
    if s == nil {
        return fmt.Errorf("failed to insert account email otp: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert account email otp: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account email otp: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account email otp: missing required parameter: otp=null")
    }
    if err := row.ValidateEmail(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidatePurpose(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidateCodeHash(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidateExpiresAt(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidateLastSentAt(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }
    if err := row.ValidateLockedUntil(); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColEmail,
        ColPurpose,
        ColStatus,
        ColCodeHash,
        ColExpiresAt,
        ColAttemptCount,
        ColLastSentAt,
        ColLockedUntil,
    )

    if _, err := s.executor.Exec(
        query,
        row.Email,
        row.Purpose,
        row.Status,
        row.CodeHash,
        row.ExpiresAt,
        row.AttemptCount,
        row.LastSentAt,
        row.LockedUntil,
    ); err != nil {
        return fmt.Errorf("failed to insert account email otp: %w", err)
    }

    return nil
}


//
// Upsert account email OTP.
//
// Version:
//   - 2026-05-05: Added.
//
func (s *OTPStore) Upsert(row *OTP) error {
    if s == nil {
        return fmt.Errorf("failed to upsert account email otp: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to upsert account email otp: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to upsert account email otp: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to upsert account email otp: missing required parameter: otp=null")
    }
    if err := row.ValidateEmail(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidatePurpose(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidateCodeHash(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidateExpiresAt(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidateLastSentAt(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }
    if err := row.ValidateLockedUntil(); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE %s = VALUES(%s), %s = VALUES(%s), %s = VALUES(%s), %s = VALUES(%s), %s = VALUES(%s), %s = VALUES(%s);",
        s.tableName,
        ColEmail,
        ColPurpose,
        ColStatus,
        ColCodeHash,
        ColExpiresAt,
        ColAttemptCount,
        ColLastSentAt,
        ColLockedUntil,
        ColStatus, ColStatus,
        ColCodeHash, ColCodeHash,
        ColExpiresAt, ColExpiresAt,
        ColAttemptCount, ColAttemptCount,
        ColLastSentAt, ColLastSentAt,
        ColLockedUntil, ColLockedUntil,
    )

    if _, err := s.executor.Exec(
        query,
        row.Email,
        row.Purpose,
        row.Status,
        row.CodeHash,
        row.ExpiresAt,
        row.AttemptCount,
        row.LastSentAt,
        row.LockedUntil,
    ); err != nil {
        return fmt.Errorf("failed to upsert account email otp: %w", err)
    }

    return nil
}


//
// Select account email otp by email and purpose.
//
// Version:
//   - 2026-05-07: Added.
//
func (s *OTPStore) SelectByEmailAndPurpose(email string, purpose OTPPurpose) (*OTP, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: missing required parameter: otp_store=null")
    }   
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: missing required parameter: table_name=%q", "empty")
    }
    if err := (&OTP{Email: email}).ValidateEmail(); err != nil {
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: %w", err)
    }
    if err := (&OTP{Purpose: purpose}).ValidatePurpose(); err != nil {
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: %w", err)
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? and %s = ? LIMIT 1;", s.tableName, ColEmail, ColPurpose)

    result := &OTP{}
    err := s.executor.QueryRow(query, email, purpose).Scan(
        &result.Email,
        &result.Purpose,
        &result.Status,
        &result.CodeHash,
        &result.ExpiresAt,
        &result.AttemptCount,
        &result.LastSentAt,
        &result.LockedUntil,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account email otp by email and purpose: %w", err)
    }

    return result, nil
}


//
// Select account email otps.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Select(option *OTPSelectOption) ([]*OTP, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email otps: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account email otps: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email otps: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account email otps: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account email otps: %w", err)
    }
    defer rows.Close()

    var result []*OTP
    for rows.Next() {
        row := &OTP{}
        if err := rows.Scan(
            &row.Email,
            &row.Purpose,
            &row.Status,
            &row.CodeHash,
            &row.ExpiresAt,
            &row.AttemptCount,
            &row.LastSentAt,
            &row.LockedUntil,
        ); err != nil {
            return nil, fmt.Errorf("failed to select account email otps: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account email otps: %w", err)
    }

    return result, nil
}


//
// Count account email otps.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Count(option *OTPSelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account email otps: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count account email otps: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account email otps: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account email otps: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account email otps: %w", err)
    }

    return result, nil
}


//
// Update account email OTP.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Update(option *OTPUpdateOption) error {
    if s == nil {
        return fmt.Errorf("failed to update account email otp: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update account email otp: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account email otp: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return fmt.Errorf("failed to update account email otp: %w", err)
    }

    assignments := make([]string, 0, 7)
    args := make([]any, 0, 9)

    if option.Status != nil {
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.CodeHash != nil {
        assignments = append(assignments, ColCodeHash + " = ?")
        args = append(args, *option.CodeHash)
    }

    if option.ExpiresAt != nil {
        assignments = append(assignments, ColExpiresAt + " = ?")
        args = append(args, *option.ExpiresAt)
    }

    if option.AttemptCount != nil {
        assignments = append(assignments, ColAttemptCount + " = ?")
        args = append(args, *option.AttemptCount)
    }

    if option.LastSentAt != nil {
        assignments = append(assignments, ColLastSentAt + " = ?")
        args = append(args, *option.LastSentAt)
    }

    if option.LockedUntilSetNull {
        assignments = append(assignments, ColLockedUntil + " = NULL")
    } else if option.LockedUntil != nil {
        assignments = append(assignments, ColLockedUntil + " = ?")
        args = append(args, *option.LockedUntil)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account email otp: invalid parameter: assignments=empty")
    }

    args = append(args, option.Email, option.Purpose)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ? AND %s = ?;", s.tableName, strings.Join(assignments, ", "), ColEmail, ColPurpose)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account email otp: %w", err)
    }

    return nil
}


// 
// Delete account email otp by email.
// 
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) DeleteByEmail(email string) error {
    if s == nil {
        return fmt.Errorf("failed to delete account email otp by email: missing required parameter: otp_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete account email otp by email: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account email otp by email: missing required parameter: table_name=%q", "empty")
    }
    if err := (&OTP{Email: email}).ValidateEmail(); err != nil {
        return fmt.Errorf("failed to delete account email otp by email: %w", err)
    }
    
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColEmail)

    if _, err := s.executor.Exec(query, email); err != nil {
        return fmt.Errorf("failed to delete account email otp by email: %w", err)
    }   
    
    return nil
} 


//
// Select usable account email otp by email and purpose.
//
func (s *OTPStore) SelectUsableByEmailAndPurpose(email string, p OTPPurpose, now time.Time) (*OTP, error) {
    otp, err := s.SelectByEmailAndPurpose(email, p)
    if err != nil {
        return nil, err
    }

    if otp == nil {
        return nil, nil
    }

    // Check whether OTP has been locked.
    if otp.LockedUntil != nil && otp.LockedUntil.After(now) {
        return otp, fmt.Errorf("forbidden: locked_until=%q", otp.LockedUntil)
    }

    // Check whether OTP is expired.
    if !otp.ExpiresAt.After(now) {
        return otp, fmt.Errorf("expired: expired_at=%q", otp.ExpiresAt)
    }

    // Check whether OTP is already verified.
    if otp.Status == OTPStatusVerified {
        return otp, fmt.Errorf("forbidden: otp=%q", "already verified")
    }

    // Check whether OTP is not active.
    if otp.Status != OTPStatusActive {
        return otp, fmt.Errorf("forbidden: otp=%q", "not active")
    }

    return otp, nil
}


//
// Build query.
//
// Version:
//   - 2025-05-03: Added.
//
func (o *OTPSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

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
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.ExpiresAtGTE != nil {
        conditions = append(conditions, ColExpiresAt + " >= ?")
        args = append(args, *o.ExpiresAtGTE)
    }
    if o.ExpiresAtLTE != nil {
        conditions = append(conditions, ColExpiresAt + " <= ?")
        args = append(args, *o.ExpiresAtLTE)
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
// Validate account email otp select option.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *OTPSelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: otp_select_option=null")
    }

    if o.Email != nil {
        if err := (&OTP{Email: *o.Email}).ValidateEmail(); err != nil {
            return err
        }
    }

    if o.EmailLike != nil {
        if err := (&OTP{Email: *o.EmailLike}).ValidateEmail(); err != nil {
            return err
        }
    }

    if o.Purpose != nil {
        if err := (&OTP{Purpose: *o.Purpose}).ValidatePurpose(); err != nil {
            return err
        }
    }

    if o.Status != nil {
        if err := (&OTP{Status: *o.Status}).ValidateStatus(); err != nil {
            return err
        }
    }

    if o.ExpiresAtGTE != nil {
        if err := (&OTP{ExpiresAt: *o.ExpiresAtGTE}).ValidateExpiresAt(); err != nil {
            return err
        }
    }

    if o.ExpiresAtLTE != nil {
        if err := (&OTP{ExpiresAt: *o.ExpiresAtLTE}).ValidateExpiresAt(); err != nil {
            return err
        }
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColEmail,
            ColPurpose,
            ColStatus,
            ColExpiresAt,
            ColAttemptCount,
            ColLastSentAt,
            ColLockedUntil:
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


//
// Validate account email otp update option.
//
// Version:
//   - 2025-05-03: Added.
//
func (o *OTPUpdateOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: otp_update_option=null")
    }

    if o.Email == "" {
        return fmt.Errorf("invalid parameter: email=empty")
    }

    if err := (&OTP{Purpose: o.Purpose}).ValidatePurpose(); err != nil {
        return err
    }

    if o.Status != nil {
        otp := OTP{
            Status: *o.Status,
        }
        if err := otp.ValidateStatus(); err != nil {
            return err
        }
    }

    if o.CodeHash != nil {
        otp := OTP{
            CodeHash: *o.CodeHash,
        }
        if err := otp.ValidateCodeHash(); err != nil {
            return err
        }
    }

    if o.ExpiresAt != nil {
        otp := OTP{
            ExpiresAt: *o.ExpiresAt,
        }
        if err := otp.ValidateExpiresAt(); err != nil {
            return err
        }
    }

    if o.LastSentAt != nil {
        otp := OTP{
            LastSentAt: *o.LastSentAt,
        }
        if err := otp.ValidateLastSentAt(); err != nil {
            return err
        }
    }

    if o.LockedUntil != nil {
        otp := OTP{
            LockedUntil: o.LockedUntil,
        }
        if err := otp.ValidateLockedUntil(); err != nil {
            return err
        }
    }

    if o.LockedUntil != nil && o.LockedUntilSetNull {
        return fmt.Errorf("invalid parameter: locked_until and locked_until_set_null are mutually exclusive")
    }

    return nil
}
