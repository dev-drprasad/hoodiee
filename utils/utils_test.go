package utils

import "testing"

import "log"

func TestGetHashFromMagnetURI(t *testing.T) {
	magnetURI := "magnet:?xt=urn:btih:bac2c9d9c552ab2465485fd37c11877f9af051db&dn=Rick and Morty S04E03 One Crew Over The Crewcoos Morty 1080p AMZN WEBRip DDP5 1 x264 CtrlHD [rartv]&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.pirateparty.gr:6969&tr=udp://9.rarbg.to:2710&tr=udp://9.rarbg.me:2710"
	hash := GetHashFromMagnet(magnetURI)
	log.Printf("%s", hash)
}
