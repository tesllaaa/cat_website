FROM nginx:stable-alpine

# Копируем статические файлы
COPY . /usr/share/nginx/html

# Копируем конфигурацию nginx
COPY nginx.conf /etc/nginx/nginx.conf

# Создаём папку для дополнительных конфигураций, если потребуется
RUN mkdir -p /etc/nginx/extra-conf.d

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
