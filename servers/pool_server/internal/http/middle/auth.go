package middle

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data/Repo/User"
	"github.com/tanjl855/tan_go_im/data/entity"
	"github.com/tanjl855/tan_go_im/pkg/im_auth"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/db"
	"net/http"
	"time"
)

var (
	userRepo = User.NewIUserRepo(db.DB.MongoDB)
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		userClaims, err := im_auth.GetClaimFromToken(token)
		if userClaims == nil || err != nil {
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
			} else {
				ctx.JSON(http.StatusUnauthorized, err)
			}

			log.Error("用户鉴权, 提取token信息错误")
			ctx.Abort()
			return
		}
		userJson, err := db.DB.Rdb.Get(ctx, "auth:"+userClaims.UID).Result()
		if err != nil && err != redis.Nil {
			ctx.JSON(http.StatusInternalServerError, err)
			log.Error("用户鉴权, redis取数据错误")
			ctx.Abort()
			return
		}
		user := &entity.User{}
		if userJson != "" {
			err = json.Unmarshal([]byte(userJson), user)
			if err != nil || user.Version != userClaims.Version {
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, err)
					log.Error("用户鉴权, redis数据json解析错误")
				} else {
					ctx.JSON(http.StatusUnauthorized, err)
				}
			}
		} else {
			userInfo, err := userRepo.GetUserInfoByUID(ctx, userClaims.UID)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				ctx.JSON(http.StatusInternalServerError, err)
				log.Error("用户鉴权，数据库查询错误")
				ctx.Abort()
				return
			}
			if userInfo == nil || userInfo.Id.IsZero() {
				ctx.JSON(http.StatusUnauthorized, err)
				log.Error("用户鉴权，数据库查询错误")
				ctx.Abort()
				return
			}
			userJsonBytes, err := json.Marshal(entity.TransformFromModel(userInfo))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				log.Error("用户鉴权，用户数据解析到json错误")
				ctx.Abort()
				return
			}

			err = db.DB.Rdb.Set(ctx, "auth:"+userInfo.UID, string(userJsonBytes), 3600*time.Second).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				log.Error("用户鉴权，用户信息存入redis缓存错误")
				ctx.Abort()
				return
			}
			user = entity.TransformFromModel(userInfo)
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

// GetUserFromCtx 从上下文中获取用户信息
func GetUserFromCtx(ctx *gin.Context) *entity.User {
	user, isexist := ctx.Get("user")
	if isexist {
		return user.(*entity.User)
	}
	return nil
}
