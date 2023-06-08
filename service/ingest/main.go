package main

import (
	"emperror.dev/errors"
	"flag"
	"github.com/je4/filesystem/v2/pkg/sftpfs"
	"github.com/je4/filesystem/v2/pkg/zipasfolder"
	"github.com/je4/mediaserver/v2/pkg/config"
	lm "github.com/je4/utils/v2/pkg/logger"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
	"io"
	"io/fs"
	"os"
)

const VERSION = "v1.0-beta.1"

const LOGFORMAT = `%{time:2006-01-02T15:04:05.000} %{shortpkg}::%{longfunc} [%{shortfile}] > %{level:.5s} - %{message}`

var configFile = flag.String("config", "./database.toml", "configuration file")

func main() {
	flag.Parse()

	cfgData, err := os.ReadFile(*configFile)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read configuration from '%s'", *configFile))
	}
	conf, err := config.LoadIngestConfig(cfgData)
	if err != nil {
		panic(errors.Wrapf(err, "cannot unmarshal config toml data from '%s'", *configFile))
	}

	daLogger, lf := lm.CreateLogger("ocfl", conf.LogFile, nil, conf.LogLevel, LOGFORMAT)
	defer lf.Close()

	privKeyName := "C:/daten/keys/syncthing/putty_ed25519.priv.openssh"
	pem, err := os.ReadFile(privKeyName)
	if err != nil {
		daLogger.Panicf("cannot open private key file '%s'", privKeyName)
	}
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		daLogger.Panicf("cannot parse private key file '%s'", privKeyName)
	}

	sshConf := &ssh.ClientConfig{
		User:            "enge0000",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var addr = "localhost:2002"
	var baseDir = "/digispace"
	var numSessions = uint(5)
	sftpFS, err := sftpfs.NewSFTPFSRW(addr, sshConf, baseDir, numSessions)
	if err != nil {
		daLogger.Panicf("cannot create sftpFS(%s, %s)", addr, baseDir)
	}
	defer sftpFS.Close()

	zipFolderFS, err := zipasfolder.NewFS(sftpFS, 2)
	if err != nil {
		daLogger.Panicf("cannot create zipAsFolderFS(%v)", sftpFS)
	}

	startFolder := "ub-reprofiler/mets-container/bau1/2020"

	fs.WalkDir(zipFolderFS, startFolder, func(path string, d fs.DirEntry, err error) error {
		daLogger.Infof("path: %s", path)
		return nil
	})

	fp, err := zipFolderFS.Open("ub-reprofiler/mets-container/bau1/2020/BAU_1_007097043_20190726T001152_master_ver1.zip/007097043/image/2316616.tif")
	if err != nil {
		daLogger.Panicf("cannot open tif in zip file")
	}
	defer fp.Close()
	fpw, err := os.Create("c:/temp/test.tif")
	if err != nil {
		daLogger.Panicf("cannot create temp file")
	}
	defer fpw.Close()
	if _, err := io.Copy(fpw, fp); err != nil {
		daLogger.Panicf("error copying tif from zip")
	}

	/*
		dbService, err := database.NewService(db, conf.Postgres.Schema)
		if err != nil {
			daLogger.Panicf("cannot create database service: %v", err)
		}

		listener, err := net.Listen("tcp", conf.Addr)
		if err != nil {
			daLogger.Panicf("cannot listen to tcp %s", conf.Addr)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterDatabaseServer(grpcServer, dbService)

		fmt.Printf("starting grpc server at %s", conf.Addr)
		grpcServer.Serve(listener)

	*/
}
