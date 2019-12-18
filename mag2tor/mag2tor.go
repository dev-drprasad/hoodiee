// Converts magnet URIs and info hashes into torrent metainfo files.
package mag2tor

import (
	"log"
	"os"

	_ "github.com/anacrolix/envpprof"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

var cl, err = torrent.NewClient(nil)

func Mag2Tor(magnetURI string) string {

	t, err := cl.AddMagnet(magnetURI)
	if err != nil {
		log.Printf("error adding magnet to client: %s", err)
		return ""
	}

	<-t.GotInfo()
	mi := t.Metainfo()
	t.Drop()
	filename := t.Info().Name + ".torrent"
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("error creating torrent metainfo file: %s", err)
		return ""
	}
	defer f.Close()
	err = bencode.NewEncoder(f).Encode(mi)
	if err != nil {
		log.Printf("error writing torrent metainfo file: %s", err)
		return ""
	}
	return filename
}
