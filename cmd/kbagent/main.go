/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	kzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/apecloud/kubeblocks/pkg/kbagent/server"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	// Set GOMAXPROCS
	_, _ = maxprocs.Set()

	// Initialize flags
	opts := kzap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(errors.Wrap(err, "fatal error viper bindPFlags"))
	}

	// Initialize logger
	kopts := []kzap.Opts{kzap.UseFlagOptions(&opts)}
	if strings.EqualFold("debug", viper.GetString("zap-log-level")) {
		kopts = append(kopts, kzap.RawZapOpts(zap.AddCaller()))
	}
	ctrl.SetLogger(kzap.New(kopts...))

	// initialize kbagent
	service, err := kbagent.Initialize(envInMapping())
	if err != nil {
		panic(errors.Wrap(err, "init action handlers failed"))
	}

	// start HTTP Server
	server := server.NewHttpServer(service)
	err = server.StartNonBlocking()
	if err != nil {
		panic(errors.Wrap(err, "failed to start HTTP server"))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop
}

func envInMapping() map[string]string {
	vars := map[string]string{}
	for _, v := range os.Environ() {
		pairs := strings.SplitN(v, "=", 2)
		vars[pairs[0]] = pairs[1]
	}
	return vars
}
