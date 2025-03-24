local minute_slice_key = KEYS[1]
local delete_key = KEYS[2]
local score1 = ARGV[1]
local score2 = ARGV[2]

-- 获取已删除集合
local del_set = redis.call('smembers', delete_key)

-- 获取要执行的 tasks数据
local targets = redis.call('zrange', minute_slice_key, score1, score2, 'byscore')

redis.call('zremrangebyscore', minute_slice_key, score1, score2)

local reply = {}

reply[1] = del_set

for i, v in ipairs(targets) do
    reply[#reply + 1] = v
end

return reply
