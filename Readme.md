# pg diagnose

# checks

* any querires > 1 min
* any locks > 1 min
* any idle in transaction > 1 min
* any bloat > 50MB and a factor > 10x
* any unused indexs
* cache hit rates < 0.98
* load higher than number of cores for plan
* connections near plan limit or high for other plans
* dataset + connections * 5b > plan memory

## api

start a job:
  POST /create , body: {'url': 'postgres://...'}
  response: url of result

view result:
  GET /results/:id



## license
MIT

