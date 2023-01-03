FROM node:18.12.1-slim

WORKDIR /usr/src/app

RUN npm install -g npm@9.2.0

COPY run.sh ..
RUN (sed -i 's/\r$//' ../run.sh) && \
    (chmod +x ../run.sh)

# run the script
CMD ../run.sh
