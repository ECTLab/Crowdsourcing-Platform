FROM node:22-alpine
WORKDIR /app

COPY ./ /app
RUN yarn

ENTRYPOINT [ "node", "index.js" ]
CMD [ "3000" ]
