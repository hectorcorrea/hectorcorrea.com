#! /bin/bash

rm -fr /tmp/proxy_cache_dir/
rm -fr /tmp/proxy_temp_dir/
nginx -s reload
