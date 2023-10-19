FROM node:latest
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init

CMD npx tailwindcss -c /src/tailwindcss.config.js - i /src/styles.css - dist/styles.css --watch --poll