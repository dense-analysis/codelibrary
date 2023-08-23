FROM node:18.17.1

WORKDIR /code

COPY ./package.json ./yarn.lock ./
RUN yarn

COPY ./tsconfig.json ./tsconfig.json
COPY ./webpack.config.ts ./webpack.config.ts
COPY ./src ./src
COPY ./public ./public

EXPOSE 9000

CMD ["yarn", "run", "dev-server"]
