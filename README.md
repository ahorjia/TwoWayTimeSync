# TwoWayTimeSync

Setup:

sudo apt-get install portaudio19-dev

export GOPATH='/home/agah/gopath' 

mkdir /home/agah/gopath/src/portaudio

git clone https://github.com/gordonklaus/portaudio.git

copy pa.c and portaudio.go in /home/agah/gopath/src/portaudio

go build portaudio

change ./examples/stereoSine.go to portaudio

go run stereoSine
================================================

To generate A:

go run stereoSineProj.go 440

    Status API Training Shop Blog About 

