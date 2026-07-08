//
// otp.go
//
package account

import (
    "crypto/hmac"
    "crypto/rand"
	"crypto/sha256"
    "database/sql"
	"encoding/hex"
    "errors"
    "fmt"
    "math/big"
    "strings"
    "time"
    "unicode"

    "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultOTPTableName = "account_otps"

    DefaultCodeLength = 6

    DefaultMaxAttemptCount uint8 = 3

    DefaultExpiresIn = 10 * time.Minute

    DefaultLockedUntilIn = 15 * time.Minute
)

var (
    otpIDCounter = &helper.IdCounter{}
)


//
// OTP.
//
// Version:
//   - 2026-05-03: Added.
//
type OTP struct {
    ID              uint64    
    Status          OTPStatus 
    Channel         OTPChannel
    Purpose         OTPPurpose
    DestinationHash string    
    CodeHash        string    
    ConsumedAt      *time.Time
    ExpiresAt       time.Time 
    AttemptCount    uint8     
    LastSentAt      time.Time 
    LockedUntil     *time.Time
    CreatedAt       time.Time 
    UpdatedAt       time.Time 
}


//
// OTPStore.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPStore struct {
    tableName        string
    accountTableName string
}


type OTPInsertParams struct {
    ID              uint64
    Status          OTPStatus
    Channel         OTPChannel
    Purpose         OTPPurpose
    DestinationHash string
    CodeHash        string
    ConsumedAt      *time.Time
    ExpiresAt       time.Time
    AttemptCount    uint8 
    LastSentAt      time.Time
    LockedUntil     *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
    InsertIgnore    bool
}


type OTPSelectParams struct {
    Status       *OTPStatus 
    Channel      *OTPChannel
    Purpose      *OTPPurpose
    ExpiresAtGTE *time.Time 
    ExpiresAtLTE *time.Time 
    OrderBy      string     
    OrderByDesc  bool      
    Limit        int        
    Offset       int        
}


//
// OTPUpdateParams.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPUpdateParams struct {
    Status             *OTPStatus
    Purpose            *OTPPurpose
    CodeHash           *string   
    ExpiresAt          *time.Time
    AttemptCount       *uint8    
    LastSentAt         *time.Time
    LockedUntil        *time.Time
    LockedUntilSetNull bool      
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
// Generate OTP ID.
//
// Version:
//   - 2026-06-25: Added.
//
func GenerateOTPID() uint64 {
    return otpIDCounter.GenerateID()
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
func NewOTPStore(tableName string, accountTableName string) (*OTPStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account otp store: missing required parameter: table_name=%q", "empty")
    }
    accountTableName = strings.TrimSpace(accountTableName)
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account otp store: missing required parameter: account_table_name=%q", "empty")
    }

    return &OTPStore{
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Normalize destination raw.
//
// Version:
//   - 2026-06-26: Added.
//
func NormalizeDestinationRaw(channel OTPChannel, destinationRaw string) (string, error) {
    // Guard.
    if err := channel.Validate(); err != nil {
        return "", err
    }
    if destinationRaw == "" {
        return "", fmt.Errorf("invalid parameter: destination_raw=%q", "empty")
    }

    switch channel {
    case OTPChannelEmail:
        return NormalizeEmailDestinationRaw(destinationRaw)
    case OTPChannelSMS:
        return NormalizePhoneDestinationRaw(destinationRaw)
    default:
        return "", fmt.Errorf("invalid parameter: otp_channel=%q", channel)
    }
}


//
// Normalize email destination raw.
//
// Version:
//   - 2026-06-26: Added.
//
func NormalizeEmailDestinationRaw(emailRaw string) (string, error) {
    email := strings.ToLower(strings.TrimSpace(emailRaw))

    if email == "" {
        return "", fmt.Errorf("invalid parameter: email=%q", "empty")
    }
    if strings.Count(email, "@") != 1 {
        return "", fmt.Errorf("invalid parameter: email=%q", email)
    }

    parts := strings.Split(email, "@")
    if parts[0] == "" || parts[1] == "" {
        return "", fmt.Errorf("invalid parameter: email=%q", email)
    }

    return email, nil
}


//
// Normalize phone destination raw.
//
// Version:
//   - 2026-06-26: Added.
//
func NormalizePhoneDestinationRaw(phoneRaw string) (string, error) {
    phone := strings.TrimSpace(phoneRaw)

    var b strings.Builder

    for _, r := range phone {
        switch {
        case r >= '0' && r <= '9':
            b.WriteRune(r)

        case r >= '０' && r <= '９':
            b.WriteRune('0' + (r - '０'))

        case r == '+':
            if b.Len() != 0 {
                return "", fmt.Errorf("invalid parameter: phone=%q", phoneRaw)
            }
            b.WriteRune(r)

        case unicode.IsSpace(r), r == '-', r == 'ー', r == '−', r == '(', r == ')':
            continue
    
        default:
            return "", fmt.Errorf("invalid parameter: phone_contains_invalid_character")
        }
    }   
            
    normalizedPhone := b.String()
    if normalizedPhone == "" { 
        return "", fmt.Errorf("invalid parameter: phone=%q", "empty")
    }

    return normalizedPhone, nil
}


//
// Hash destination raw.
//
// Version:
//   - 2026-06-29: Added.
//
func HashDestinationRaw(secret []byte, destinationRaw string) (string, error) {
    // Guard.
    if len(secret) == 0 {
        return "", fmt.Errorf("failed to hash destination raw: missing required parameter: secret=%q", "empty")
    }
    if destinationRaw == "" {
        return "", fmt.Errorf("failed to hash destination raw: missing required parameter: destination_raw=%q", "empty")
    }

    mac := hmac.New(sha256.New, secret)
    if _, err := mac.Write([]byte(destinationRaw)); err != nil {
        return "", fmt.Errorf("failed to hash destination raw: %w", err)
    }

    destinationHash := hex.EncodeToString(mac.Sum(nil))

    if err := ValidateOTPDestinationHash(destinationHash); err != nil {
        return "", fmt.Errorf("failed to hash destination raw: %w", err)
    }

    return destinationHash, nil
}


//
// Validate account OTP ID.
//
// Version:
//   - 2026-06-25: Added.
//
func ValidateOTPID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account OTP ID.
//
// Version:
//   - 2026-06-25: Added.
//
func (o *OTP) ValidateID() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPID(o.ID)
}


//
// Validate account OTP channel.
//
// Version:
//   - 2026-06-25: Added.
//
func ValidateOTPChannel(channel OTPChannel) error {
    if err := channel.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account OTP channel.
//
// Version:
//   - 2026-06-25: Added.
//
func (o *OTP) ValidateChannel() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPChannel(o.Channel)
}


//
// Validate account OTP purpose.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateOTPPurpose(purpose OTPPurpose) error {
    if err := purpose.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account OTP purpose.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidatePurpose() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPPurpose(o.Purpose)
}


//
// Validate account OTP status.
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
// Validate account OTP status.
//
// Version:
//   - 2026-05-09: Added.
//
func (o *OTP) ValidateStatus() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPStatus(o.Status)
}


//
// Validate account OTP destination hash.
//
// Version:
//   - 2026-06-25: Added.
//
func ValidateOTPDestinationHash(destinationHash string) error {
    if destinationHash == "" {
        return fmt.Errorf("invalid parameter: destination_hash=%q", "empty")
    }
    if len(destinationHash) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 destination_hash=%q", "too long")
    }
    return nil
}


//
// Validate account OTP destination hash.
//
// Version:
//   - 2026-06-25: Added.
//
func (o *OTP) ValidateDestinationHash() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPDestinationHash(o.DestinationHash)
}


//
// Validate account OTP code hash.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateOTPCodeHash(codeHash string) error {
    if codeHash == "" {
        return fmt.Errorf("invalid parameter: code_hash=%q", "empty")
    }
    if len(codeHash) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 code_hash=%q", "too long")
    }
    return nil
}


//
// Validate account OTP code hash.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateCodeHash() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPCodeHash(o.CodeHash)
}


//
// Validate account OTP consumed at.
//
// Version:
//   - 2026-06-25: Added.
//
func ValidateOTPConsumedAt(consumedAt *time.Time) error {
    if consumedAt == nil {
        return nil
    }
    if consumedAt.IsZero() {
        return fmt.Errorf("invalid parameter: consumed_at=%q", "empty")
    }
    if consumedAt.UTC().After(time.Now().UTC()) {
        return fmt.Errorf("invalid parameter: consumed_at=%q", consumedAt.Format(time.RFC3339))
    }
    return nil
}


//
// Validate account OTP consumed at.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateConsumedAt() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPConsumedAt(o.ConsumedAt)
}


//
// Validate account OTP expires at.
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
// Validate account OTP expires at.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateExpiresAt() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPExpiresAt(o.ExpiresAt)
}


//
// Validate account OTP last sent at.
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
// Validate account OTP last sent at.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateLastSentAt() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPLastSentAt(o.LastSentAt)
}


//
// Validate account OTP locked until.
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
// Validate account OTP locked until.
//
// Version:
//   - 2026-05-03: Added.
//
func (o *OTP) ValidateLockedUntil() error {
    if o == nil {
        return fmt.Errorf("missing required parameter: account_otp=null")
    }
    return ValidateOTPLockedUntil(o.LockedUntil)
}


//
// Create account OTPs table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *OTPStore) CreateTable(executor helper.Executor) error {
    if s == nil {
        return fmt.Errorf("failed to create account otps table: missing required parameter: otp_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account otps table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create account otps table: missing required parameter: executor=null")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Channel',
            %s SMALLINT UNSIGNED NOT NULL COMMENT 'Purpose',
            %s VARCHAR(255) NOT NULL COMMENT 'Destination hash',
            %s VARCHAR(255) NOT NULL COMMENT 'Code hash',
            %s DATETIME NULL COMMENT 'Consumed at',
            %s DATETIME NOT NULL COMMENT 'Expires at',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Attempt count',
            %s DATETIME NOT NULL COMMENT 'Last sent at',
            %s DATETIME NULL COMMENT 'Locked until',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_%s_sta_cha_pur_des_has (%s, %s, %s, %s),
            KEY idx_%s_status_expires_at (%s, %s)
        );`,
        s.tableName,
        ColID,
        ColStatus,
        ColChannel,
        ColPurpose,
        ColDestinationHash,
        ColCodeHash,
        ColConsumedAt,
        ColExpiresAt,
        ColAttemptCount,
        ColLastSentAt,
        ColLockedUntil,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        s.tableName, ColStatus, ColChannel, ColPurpose, ColDestinationHash,
        s.tableName, ColStatus, ColExpiresAt,
    )

    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account otps table: %w", err)
    }

    return nil
}


//
// Insert account OTP.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Insert(executor helper.Executor, params *OTPInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert account otp: missing required parameter: otp_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account otp: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert account otp: missing required parameter: executor=null")
    }
    if params == nil {
        return fmt.Errorf("failed to insert account otp: missing required parameter: otp_insert_params=null")
    }
    if err := ValidateOTPStatus(params.Status); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPChannel(params.Channel); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPPurpose(params.Purpose); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPDestinationHash(params.DestinationHash); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPCodeHash(params.CodeHash); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPConsumedAt(params.ConsumedAt); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPExpiresAt(params.ExpiresAt); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPLastSentAt(params.LastSentAt); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }
    if err := ValidateOTPLockedUntil(params.LockedUntil); err != nil {
        return fmt.Errorf("failed to insert account otp: %w", err)
    }

    if params.ID == 0 {
        params.ID = GenerateOTPID()
    }

    now := time.Now().UTC()
    if params.CreatedAt.IsZero() {
        params.CreatedAt = now
    }
    if params.UpdatedAt.IsZero() {
        params.UpdatedAt = now
    }

    queryPrefix := "INSERT"
    if params.InsertIgnore {
        queryPrefix = "INSERT IGNORE"
    }

    query := fmt.Sprintf(
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColStatus,
        ColChannel,
        ColPurpose,
        ColDestinationHash,
        ColCodeHash,
        ColConsumedAt,
        ColExpiresAt,
        ColAttemptCount,
        ColLastSentAt,
        ColLockedUntil,
    )

    if _, err := executor.Exec(
        query,
        params.ID,
        params.Status,
        params.Channel,
        params.Purpose,
        params.DestinationHash,
        params.CodeHash,
        params.ConsumedAt,
        params.ExpiresAt,
        params.AttemptCount,
        params.LastSentAt,
        params.LockedUntil,
    ); err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            return fmt.Errorf("failed to insert account otp: %w", helper.ErrDuplicateKey)
        }
        return fmt.Errorf("failed to insert account otp: %w", err)
    }

    return nil
}


//
// Select latest usable active OTP and normalize its state if needed.
//
// Version:
//   - 2026-05-07: Added.
//
func (s *OTPStore) SelectLatestUsableAndNormalizeIfNeeded(executor helper.Executor, channel OTPChannel, purpose OTPPurpose, destinationHash string) (*OTP, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select latest usable account otp: missing required parameter: otp_store=null")
    }   
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select latest usable account otp: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select latest usable account otp: missing required parameter: executor=null")
    }
    if err := ValidateOTPChannel(channel); err != nil {
        return nil, fmt.Errorf("failed to select latest usable account otp: %w", err)
    }
    if err := ValidateOTPPurpose(purpose); err != nil {
        return nil, fmt.Errorf("failed to select latest usable account otp: %w", err)
    }

    // Generate query.
    query := fmt.Sprintf(
        "SELECT * FROM %s WHERE %s = ? AND %s = ? AND %s = ? AND %s = ? ORDER BY %s DESC LIMIT 1;",
        s.tableName, ColStatus, ColChannel, ColPurpose, ColDestinationHash, ColCreatedAt,
    )

    // Execute.
    otp := &OTP{}
    err := executor.QueryRow(query, OTPStatusActive, channel, purpose, destinationHash).Scan(
        &otp.ID,
        &otp.Status,
        &otp.Channel,
        &otp.Purpose,
        &otp.DestinationHash,
        &otp.CodeHash,
        &otp.ConsumedAt,
        &otp.ExpiresAt,
        &otp.AttemptCount,
        &otp.LastSentAt,
        &otp.LockedUntil,
        &otp.CreatedAt,
        &otp.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select latest usable account otp: %w", err)
    }

    now := time.Now().UTC()

    // Check whether OTP has been locked.
    if otp.LockedUntil != nil && otp.LockedUntil.After(now) {
        return nil, fmt.Errorf("failed to select latest usable account otp: %w: locked_until=%q", helper.ErrForbidden, otp.LockedUntil)
    }

    // Check whether OTP is expired.
    if !otp.ExpiresAt.After(now) {
        expiredStatus := OTPStatusExpired
        updateParams := &OTPUpdateParams{
            Status: &expiredStatus,
        }
        if err := s.UpdateByID(executor, otp.ID, updateParams); err != nil {
            return nil, fmt.Errorf("failed to select latest usable account otp: %w", err)
        }
        return nil, fmt.Errorf("failed to select latest usable account otp: %w: expired_at=%q", helper.ErrExpired, otp.ExpiresAt)
    }

    // Check whether OTP is already verified.
    if otp.Status == OTPStatusVerified {
        return nil, fmt.Errorf("failed to select latest usable account otp: %w: otp=%q", helper.ErrForbidden, "already verified")
    }

    // Check whether OTP is not active.
    if otp.Status != OTPStatusActive {
        return nil, fmt.Errorf("failed to select latest usable account otp: %w: otp=%q", helper.ErrForbidden, "not active")
    }

    return otp, nil
}


//
// Select account otps.
//
// Version:
//   - 2026-05-04: Added.
//
//func (s *OTPStore) Select(executor helper.Executor, params *OTPSelectParams) ([]*OTP, error) {
//    if s == nil {
//        return nil, fmt.Errorf("failed to select account otps: missing required parameter: otp_store=null")
//    }
//    if s.tableName == "" {
//        return nil, fmt.Errorf("failed to select account otps: missing required parameter: table_name=%q", "empty")
//    }
//    if executor == nil {
//        return nil, fmt.Errorf("failed to select account otps: missing required parameter: executor=null")
//    }
//    if err := params.Validate(); err != nil {
//        return nil, fmt.Errorf("failed to select account otps: %w", err)
//    }
// 
//    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)
// 
//    rows, err := executor.Query(query, args...)
//    if err != nil {
//        return nil, fmt.Errorf("failed to select account otps: %w", err)
//    }
//    defer rows.Close()
// 
//    var result []*OTP
//    for rows.Next() {
//        row := &OTP{}
//        if err := rows.Scan(
//            &row.Email,
//            &row.Purpose,
//            &row.Status,
//            &row.CodeHash,
//            &row.ExpiresAt,
//            &row.AttemptCount,
//            &row.LastSentAt,
//            &row.LockedUntil,
//        ); err != nil {
//            return nil, fmt.Errorf("failed to select account otps: %w", err)
//        }
// 
//        result = append(result, row)
//    }
// 
//    if err := rows.Err(); err != nil {
//        return nil, fmt.Errorf("failed to select account otps: %w", err)
//    }
// 
//    return result, nil
//}


//
// Count account otps.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) Count(executor helper.Executor, params *OTPSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account otps: missing required parameter: otp_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account otps: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count account otps: missing required parameter: executor=null")
    }
    if err := params.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account otps: %w", err)
    }

    query, args := params.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account otps: %w", err)
    }

    return result, nil
}


//
// Update account OTP.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) UpdateByID(executor helper.Executor, id uint64, params *OTPUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update account otp by id: missing required parameter: otp_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account otp by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update account otp by id: missing required parameter: executor=null")
    }
    if err := params.Validate(); err != nil {
        return fmt.Errorf("failed to update account otp by id: %w", err)
    }

    assignments := make([]string, 0, 7)
    args := make([]any, 0, 9)

    if params.Status != nil {
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *params.Status)
    }
    if params.CodeHash != nil {
        assignments = append(assignments, ColCodeHash + " = ?")
        args = append(args, *params.CodeHash)
    }
    if params.ExpiresAt != nil {
        assignments = append(assignments, ColExpiresAt + " = ?")
        args = append(args, *params.ExpiresAt)
    }
    if params.AttemptCount != nil {
        assignments = append(assignments, ColAttemptCount + " = ?")
        args = append(args, *params.AttemptCount)
    }
    if params.LastSentAt != nil {
        assignments = append(assignments, ColLastSentAt + " = ?")
        args = append(args, *params.LastSentAt)
    }
    if params.LockedUntilSetNull {
        assignments = append(assignments, ColLockedUntil + " = NULL")
    } else if params.LockedUntil != nil {
        assignments = append(assignments, ColLockedUntil + " = ?")
        args = append(args, *params.LockedUntil)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account otp by id: invalid parameter: assignments=%q", "empty")
    }

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    args = append(args, id)

    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account otp: %w", err)
    }

    return nil
}


// 
// Delete account otp by ID.
// 
// Version:
//   - 2026-05-04: Added.
//
func (s *OTPStore) DeleteByID(executor helper.Executor, id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account otp by id: missing required parameter: otp_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account otp by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete account otp by id: missing required parameter: executor=null")
    }
    if err := ValidateOTPID(id); err != nil {
        return fmt.Errorf("failed to delete account otp by id: %w", err)
    }
    
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account otp by id: %w", err)
    }   
    
    return nil
} 


//
// Select usable account otp by email and purpose.
//
//func (s *OTPStore) SelectUsableByEmailAndPurpose(executor helper.Executor, email string, p OTPPurpose, now time.Time) (*OTP, error) {
//    otp, err := s.SelectByEmailAndPurpose(executor, email, p)
//    if err != nil {
//        return nil, err
//    }
// 
//    if otp == nil {
//        return nil, nil
//    }
// 
//    // Check whether OTP has been locked.
//    if otp.LockedUntil != nil && otp.LockedUntil.After(now) {
//        return otp, fmt.Errorf("forbidden: locked_until=%q", otp.LockedUntil)
//    }
// 
//    // Check whether OTP is expired.
//    if !otp.ExpiresAt.After(now) {
//        return otp, fmt.Errorf("expired: expired_at=%q", otp.ExpiresAt)
//    }
// 
//    // Check whether OTP is already verified.
//    if otp.Status == OTPStatusVerified {
//        return otp, fmt.Errorf("forbidden: otp=%q", "already verified")
//    }
// 
//    // Check whether OTP is not active.
//    if otp.Status != OTPStatusActive {
//        return otp, fmt.Errorf("forbidden: otp=%q", "not active")
//    }
// 
//    return otp, nil
//}


//
// Build query.
//
// Version:
//   - 2025-05-03: Added.
//
func (o *OTPSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.Channel != nil {
        conditions = append(conditions, ColChannel + " = ?")
        args = append(args, *o.Channel)
    }
    if o.Purpose != nil {
        conditions = append(conditions, ColPurpose + " = ?")
        args = append(args, *o.Purpose)
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
// Validate account otp select params.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *OTPSelectParams) Validate() error {
    // Guard.
    if o == nil {
        return nil
    }

    if o.Status != nil {
        if err := ValidateOTPStatus(*o.Status); err != nil {
            return err
        }
    }
    if o.Channel != nil {
        if err := ValidateOTPChannel(*o.Channel); err != nil {
            return err
        }
    }
    if o.Purpose != nil {
        if err := ValidateOTPPurpose(*o.Purpose); err != nil {
            return err
        }
    }
    if o.ExpiresAtGTE != nil {
        if err := ValidateOTPExpiresAt(*o.ExpiresAtGTE); err != nil {
            return err
        }
    }
    if o.ExpiresAtLTE != nil {
        if err := ValidateOTPExpiresAt(*o.ExpiresAtLTE); err != nil {
            return err
        }
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColID,
             ColStatus,
             ColChannel,
             ColPurpose,
             ColConsumedAt,
             ColExpiresAt,
             ColAttemptCount,
             ColLastSentAt,
             ColLockedUntil,
             ColCreatedAt,
             ColUpdatedAt:
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
// Validate account otp update params.
//
// Version:
//   - 2025-05-03: Added.
//
func (p *OTPUpdateParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: otp_update_params=null")
    }

    if p.Status != nil {
        if err := ValidateOTPStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Purpose != nil {
        if err := ValidateOTPPurpose(*p.Purpose); err != nil {
            return err
        }
    }
    if p.CodeHash != nil {
        if err := ValidateOTPCodeHash(*p.CodeHash); err != nil {
            return err
        }
    }
    if p.ExpiresAt != nil {
        if err := ValidateOTPExpiresAt(*p.ExpiresAt); err != nil {
            return err
        }
    }
    if p.LastSentAt != nil {
        if err := ValidateOTPLastSentAt(*p.LastSentAt); err != nil {
            return err
        }
    }
    if p.LockedUntil != nil {
        if err := ValidateOTPLockedUntil(p.LockedUntil); err != nil {
            return err
        }
    }
    if p.LockedUntil != nil && p.LockedUntilSetNull {
        return fmt.Errorf("invalid parameter: locked_until and locked_until_set_null are mutually exclusive")
    }

    return nil
}
