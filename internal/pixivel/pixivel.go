package pixivel

import (
	"Pixivel/internal/config"
	"Pixivel/internal/database"

	"Pixivel/internal/database/redis"
	"Pixivel/internal/pixiv"
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
)

type Pixivel struct {
	Database *database.Database

	PixivAPI  *pixiv.AppPixiv
	RedisPool *redis.RedisPool
	Redis     *redis.RedisClient
}

func GetHandler() *Pixivel {
	conf := config.Read()
	db := database.GetDB(conf)
	db.Migrate()

	redisPool := redis.NewRedisPool(conf)
	api := pixiv.AppPixivAPI()
	api.BaseAPI.HookAccessToken(func(token string) {
		log.Println(token)
	})
	api.BaseAPI.SetAuth(conf.Pixiv.AccessToken, conf.Pixiv.RefreshToken)
	return &Pixivel{
		Database:  db,
		PixivAPI:  api,
		RedisPool: redisPool,
		Redis:     redisPool.NewRedisClient(),
	}
}

//Cache is a middleware of gin that caches all HTTP responses.
func (px *Pixivel) Cache(h gin.HandlerFunc, expire int) gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheKey := c.Request.URL.Path
		log.Println(cacheKey)
		value, err := px.Redis.GetValue(cacheKey)
		if px.Redis.IsExist(err) {
			log.Println("CACGE HIT")
			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			c.String(200, value)
			return
		} else if px.Redis.IsError(err) {
			panic(err)
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w
		h(c)
		log.Println("CACHE MISS")
		px.Redis.SetValueExpire(cacheKey, w.body.String(), expire)
	}
}

func FilterKey(id uint64, prefix string) []byte {
	return config.StringIn(prefix + config.Itoa(id))
}

func (px *Pixivel) WITDatabase(illust *pixiv.Illust) {
	px.Database.CreateIllust(illust)
	px.CuckooFilter.Insert(FilterKey(illust.ID, "i"))
}

func (px *Pixivel) WMITDatabase(illusts []pixiv.Illust) {
	for j := 0; j < len(illusts); j++ {
		px.WITDatabase(&illusts[j])
	}
}

func (px *Pixivel) SingleIllust(id uint64) interface{} {
	cfkey := FilterKey(id, "i")
	exist := px.CuckooFilter.Lookup(cfkey)
	if exist {
		log.Println("FILTER HIT")
		res, err := px.Database.QueryIllust(id)
		if err != nil {
			if err != database.RECORD_NOT_FOUND {
				panic(err)
			}
		} else {
			log.Println("DB HIT")
			return res
		}
	}
	log.Println("FILTER MISS")
	res, err := px.PixivAPI.IllustDetail(id)
	if err != nil {
		panic(err)
	}
	px.WITDatabase(res)
	return res
}

func (px *Pixivel) Close() {
	log.Print("Saving BitMap...")
	px.CuckooFilter.SaveFile()
}
