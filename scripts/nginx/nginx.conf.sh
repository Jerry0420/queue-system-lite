#!/bin/sh

cat <<EOF > nginx.conf
events {

}

http {
    
    access_log off;
    error_log /var/log/nginx/error.log crit;
    server_tokens off; 
    sendfile on;

    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_buffers 16 8k;
    gzip_min_length 128;
    gzip_comp_level 6;
    gzip_types
        # text/event-stream
        text/plain 
        text/css 
        text/js 
        text/xml 
        text/javascript 
        application/javascript 
        application/x-javascript 
        application/json 
        application/xml 
        application/rss+xml 
        image/svg+xml
        
    send_timeout 75s;
    keepalive_timeout  75s 75s;
    keepalive_requests 10000;
    etag on;
    
    add_header Cache-Control "no-cache";

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    limit_req_zone \$binary_remote_addr zone=reqlimit:5m rate=5r/s;
    limit_req_status 429;

    limit_conn_zone \$binary_remote_addr zone=connlimit:5m;
    limit_conn_status 429;

    upstream  upstream_backend {
        server   $BACKEND_HOST:8000;
        keepalive_timeout 75s;
        keepalive 10000;
    }

    upstream  upstream_frontend {
        server   $FRONTEND_HOST:3000;
        keepalive_timeout 75s;
        keepalive 10000;
    }

    server {
        listen 80 default_server;
        listen [::]:80 default_server;

        limit_req zone=reqlimit burst=10;
        limit_conn connlimit 5;
        
        server_name _;

        proxy_connect_timeout 75s;
        proxy_read_timeout 86400s;
        proxy_send_timeout 75s;

        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }

        location / {
            return 301 https://queue-system.vip\$request_uri;
        }
    }

    server {
        
        listen 443 default_server ssl http2;
        listen [::]:443 ssl http2;

        limit_req zone=reqlimit burst=10;
        limit_conn connlimit 5;

        server_name _;

        ssl_certificate /run/secrets/nginx.crt;
        ssl_certificate_key /run/secrets/nginx.key;

        proxy_connect_timeout 75s;
        proxy_read_timeout 86400s;
        proxy_send_timeout 75s;

        location /api {
            proxy_pass http://upstream_backend;

            proxy_set_header Connection '';
            proxy_http_version 1.1;
            chunked_transfer_encoding off;
            
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$http_x_forwarded_proto;
        }

        location / {
            proxy_pass http://upstream_frontend;

            if (\$request_uri ~* \.(css|gif|jpe?g|png|svg|ico)) {
                add_header Cache-Control "max-age=31536000";
            }

            if (\$request_uri ~* \.(js)) {
                add_header Cache-Control "private, max-age=31536000";
            }

            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$http_x_forwarded_proto;
        }
    }
}
EOF