package code

const (
	SUCCESS                       = 200
	INVALID_PARAMS                = 400
	ERROR                         = 500
	ID_NOT_EMPTY                  = 4001
	ERROR_TOKEN_EMPTY             = 4002
	ERROR_TOKEN_INVALID           = 4003
	ERROR_TOKEN_EXPIRE            = 4004
	ERROR_USER_NOT_EXIST          = 4005
	ERROR_USER_EXIST              = 4006
	ERROR_USER_PWD                = 4007
	ERROR_MERCHANT_NOT_EXIST      = 4008
	ERROR_MERCHANT_EXIST          = 4009
	ERROR_SHOP_BUSINESS_EXIST     = 4010
	ERROR_SHOP_BUSINESS_NOT_EXIST = 4011
	ERROR_SKU_CODE_EXIST          = 4012
	ERROR_SKU_CODE_NOT_EXIST      = 4013
	ERROR_SHOP_ID_NOT_EXIST       = 4014
	ERROR_SHOP_ID_EXIST           = 4015

	DB_DUPLICATE_ENTRY = 50000

	ERROR_EMAIL_SEND          = 50001
	ERROR_VERIFY_CODE_EMPTY   = 50002
	ERROR_VERIFY_CODE_INVALID = 50003
	ERROR_VERIFY_CODE_EXPIRE  = 50004
)
