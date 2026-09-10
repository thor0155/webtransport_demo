FROM node:22 AS build

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM nginx:1.29-alpine

COPY --from=build /app/dist \
    /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY 40-generate-config.sh /docker-entrypoint.d/40-generate-config.sh

RUN chmod +x /docker-entrypoint.d/40-generate-config.sh \
    && chown -R nginx:nginx /usr/share/nginx/html \
    && chmod -R u+rwX /usr/share/nginx/html

EXPOSE 80