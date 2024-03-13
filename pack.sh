rm goweb.zip
md5sum goweb
git archive --format=zip --output=goweb.zip HEAD -- . ':!/.git'

md5sum goweb.zip
