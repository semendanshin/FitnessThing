FROM node:slim AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY tsconfig.json ./
COPY src ./src

RUN npm run build

FROM node:slim AS production

WORKDIR /app

COPY package*.json ./
RUN npm install --omit=dev

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/src/templates ./dist/templates

ENV NODE_ENV=production

CMD ["node", "dist/index.js"]
