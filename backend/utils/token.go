package utils

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

// JWT Payloads
type Payloads struct {
	Exp      int64  `json:"exp"` // 过期时间(Unix时间戳)
	User_id  int64  `json:"user_id"`
	Username string `json:"username"`
}

// NewToken 生成新的JWT令牌
//
// 参数:
//   - payload interface{}: JWT负载信息,会被序列化为 JWT Claims
//   - key string: 用于签名的密钥
//
// 返回:
//   - string: 生成的token字符串，格式为: header.payload.signature
//   - error: 生成过程中的错误
//
// 注意:
//   - 使用 RFC 7519 标准的 NumericDate 格式(Unix时间戳)
//   - 签名算法使用标准的 HMAC-SHA256
func NewToken(payload interface{}, key string) (string, error) {
	// 将 payload 转换为 jwt.MapClaims
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var claims jwt.MapClaims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		return "", err
	}

	// 创建 token,使用 HS256 算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并生成完整的 token 字符串
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CheckToken 验证JWT令牌的有效性
//
// 参数:
//   - tokenString string: 待验证的JWT令牌字符串
//   - key string: 验证签名用的密钥
//   - payloadsDest interface{}: 用于存储解码后payload的结构体指针
//
// 返回:
//   - error: 验证过程中的错误，包括:
//   - jwt.ErrTokenExpired: 令牌已过期
//   - jwt.ErrSignatureInvalid: 签名验证失败
//   - jwt.ErrTokenMalformed: 令牌格式错误
//   - 其他错误: 解析或解码错误
//
// 注意:
//   - 使用标准的 JWT 验证流程
//   - 自动验证签名
//   - 如果 payload 包含 exp 字段,会自动验证过期时间
func CheckToken(tokenString string, key string, payloadsDest interface{}) error {
	// 解析并验证 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	if err != nil {
		return err
	}

	// 验证 token 是否有效
	if !token.Valid {
		return jwt.ErrTokenInvalidClaims
	}

	// 提取 claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.ErrTokenInvalidClaims
	}

	// 将 claims 转换为目标结构体
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return err
	}

	err = json.Unmarshal(claimsBytes, payloadsDest)
	if err != nil {
		return err
	}

	return nil
}
