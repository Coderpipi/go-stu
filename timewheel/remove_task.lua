local delete_key = KEYS[1]
local task_key = ARGV[1]



redis.call("sadd", delete_key, task_key)

local scnt = redis.call("scard", delete_key)
-- 如果是第一次创建集合, 设置 2m超时时间
if (tonumber(scnt) == 1)
then
    redis.call('expire', delete_key, 120)
end
return scnt
