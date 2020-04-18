package main

import (
	"net"
	"os"

	"github.com/p2sub/p2sub/configuration"
	"github.com/p2sub/p2sub/logger"
	"github.com/p2sub/p2sub/p2p"
)

//Handle connect from peer-to-peer network
func startNotary(conf *configuration.Config) {
	sugar.Info("Start notary node pid: ", os.Getpid())
	p2pNode := p2p.CreatePeer("tcp", conf.BindHost+":"+conf.BindPort)
	go p2pNode.Listen()
	//HandleLoop peer main loop
	for {
		select {
		//Handle new connect
		case connect := <-p2pNode.NewConnections:
			p2pNode.ActivePeers[connect] = true
			sugar.Info("New connection:", connect.RemoteAddr(), "->", connect.LocalAddr())
			go func() {
				buf := make([]byte, p2p.BufferSize)
				for {
					nbyte, err := connect.Read(buf)
					if err != nil {
						p2pNode.DeadConnections <- connect
						break
					} else {
						chunk := make([]byte, nbyte)
						copy(chunk, buf[:nbyte])
						sugar.Info("Received data:", connect.RemoteAddr(), "->", connect.LocalAddr())
						p2pNode.DataChannel <- chunk
					}
				}
			}()
		//Handle dead connect
		case deadConnect := <-p2pNode.DeadConnections:
			sugar.Info("Close connect", deadConnect.RemoteAddr(), "->", deadConnect.LocalAddr())
			err := deadConnect.Close()
			if err != nil {
				sugar.Error("Could not close connect:", err)
			}
			delete(p2pNode.ActivePeers, deadConnect)
		//Handle receiving data
		case receivedData := <-p2pNode.DataChannel:
			sugar.Debugf("Received  %d bytes of data", len(receivedData))
			logger.HexDump("Dumped data:", receivedData)
			sugar.Info("Active peers:", len(p2pNode.ActivePeers))
			for connect, i := range p2pNode.ActivePeers {
				if i {
					go func(connect net.Conn) {
						totalWritten := 0
						for totalWritten < len(receivedData) {
							writtenBytes, err := connect.Write(receivedData[totalWritten:])
							if err != nil {
								p2pNode.DeadConnections <- connect
								break
							}
							totalWritten += writtenBytes
						}
						sugar.Info("Sent data:", connect.LocalAddr(), "->", connect.RemoteAddr())
					}(connect)
				} else {
					p2pNode.DeadConnections <- connect
				}
			}
		}
	}

}
