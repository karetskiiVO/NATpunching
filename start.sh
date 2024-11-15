apt-get update && apt install git && \
apt install make && apt-get install golang && \
git clone https://github.com/karetskiiVO/NATpunching && \
cd ./NATpunching && \
cd ./NATpunch/ && find ./ -maxdepth 1 -mindepth 1 -exec mv -t ../natpunch {} + && \
cd .. && rm -r ./NATpunch/ && make