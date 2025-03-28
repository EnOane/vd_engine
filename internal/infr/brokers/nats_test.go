package brokers

import (
	"os"
	"testing"
	"vd_engine/internal/util"
)

func TestExe(t *testing.T) {
	MustConnect()
	br := NewNatsBroker()

	in := ps(br.Conn, "https://youtube.com/shorts/X-xPsJfIWK0")

	chr := util.NewChannelReader(in)

	file, _ := os.Create("asd.mp4")
	defer file.Close()
	chr.ToFile(file, 1024)
}
