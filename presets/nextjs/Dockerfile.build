FROM kooldev/node:14 AS build

COPY . /app

RUN npm install && npm run build

FROM kooldev/node:14

COPY --from=build --chown=kool:kool /app /app

EXPOSE 3000

CMD [ "npm", "start" ]
