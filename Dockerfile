FROM golang
RUN apt update
RUN apt install -y libasound2 alsa-utils libasound-dev portaudio19-dev libportaudio2 libportaudiocpp0
COPY ./bin/playback ./playback
#ENTRYPOINT ./playback server
ENTRYPOINT speaker-test
