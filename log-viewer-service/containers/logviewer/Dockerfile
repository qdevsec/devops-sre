FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package.json and install dependencies
COPY package*.json ./
RUN npm install --production

# Copy the rest of the app (including server.js)
COPY . .

# Expose port
EXPOSE 8080

# Start the app
CMD ["node", "server.js"]
