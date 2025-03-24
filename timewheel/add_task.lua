local zset_key = KEYS[1]
local delete_key = KEYS[2]
local score = ARGV[1]
local task_key = ARGV[2]
local task_body = ARGV[3]

-- 每次添加任务时都从删除集合中移除该 key
redis.call('srem', delete_key, task_key)
redis.call('zadd', zset_key, score, task_body)
