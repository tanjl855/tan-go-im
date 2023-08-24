package im_rsp

var msgMap = map[int]string{
	SUCCESS: "ok",
	ERROR:   "fail",

	ParameterValidationError: "参数解析错误",
	UserEmailExist:           "邮箱已被注册",
	UserEmailRegistering:     "邮箱在被注册ing",
	ProcessError:             "流程错误",
	VerifyCodeError:          "验证错误",
	UserNotExist:             "用户不存在",
	UserPasswordError:        "密码错误",
	VerifyPassword:           "密码校验错误",
	AuthError:                "权限验证错误",
	GroupNotExist:            "群不存在",
	PasswordError:            "密码错误",
}

func getMsg(code int) string {
	msg, ok := msgMap[code]
	if ok {
		return msg
	}
	return msgMap[ERROR]
}
