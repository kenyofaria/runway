FROM node:18

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

EXPOSE 3000

ENV HOST=0.0.0.0
ENV REACT_APP_API_URL=http://backend-service:8080

CMD ["npm", "start"]