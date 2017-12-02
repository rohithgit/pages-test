package helper
import (
  "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/specnl/spectre-base-microservice/redis"
  "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
  "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
  "golang.org/x/net/context"
  "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
)

type cacheKey string

var (
	ctxCacheKey cacheKey = "resultsCache"
)

func FetchRCacheFromContext(ctx context.Context) (connCache *redis.RedisCache, err error) {
	if (ctx == nil) {
		ctx = CreateContext()
	}
	connCache, ok := ctx.Value(ctxCacheKey).(*redis.RedisCache)
	if ( !ok) {
		utils.SpectreLog.Debugln("global.Options.RedisServers ", global.Options.RedisServers)
		if global.Options.EnableSvcDiscovery  {
			utils.SpectreLog.Debugln("Service discovery enabled")
			redisServer := utils.ParseServerAddress(ctx, global.Options.RedisServers, false)
			utils.SpectreLog.Debugln("Redis server")
			connCache, err = redis.CreateConn(redisServer, constants.SERVICENAME)
		} else {
			utils.SpectreLog.Debugln("Service discovery not enabled")
			utils.SpectreLog.Debugln("Redis server: ", global.Options.RedisServers)
			connCache, err = redis.CreateConn(global.Options.RedisServers, constants.SERVICENAME)
		}
		context.WithValue(ctx, ctxCacheKey, connCache)
	}
	return connCache, err
}

func CreateContext() (ctx context.Context){
  if( ctx == nil) {
    ctx = context.WithValue(context.TODO(), "common", "common")
  }
  return ctx
}