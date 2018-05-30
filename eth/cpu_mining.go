// Copyright 2015 The go-ethereum Authors
// This file is part of Webchain.
//
// Webchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Webchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Webchain. If not, see <http://www.gnu.org/licenses/>.

// +build !opencl

package eth

import (
	"errors"
	"fmt"

	"github.com/webchain-network/webchaind/logger"
	"github.com/webchain-network/webchaind/logger/glog"
)

const disabledInfo = "Set GO_OPENCL and re-build to enable."

func (s *Ethereum) StartMining(threads int, gpus string) error {
	eb, err := s.Etherbase()
	if err != nil {
		err = fmt.Errorf("Cannot start mining without etherbase address: %v", err)
		glog.V(logger.Error).Infoln(err)
		return err
	}

	if gpus != "" {
		return errors.New("GPU mining disabled. " + disabledInfo)
	}

	// CPU mining
	go s.miner.Start(eb, threads)
	return nil
}

func GPUBench(gpuid uint64) {
	fmt.Println("GPU mining disabled. " + disabledInfo)
}

func PrintOpenCLDevices() {
	fmt.Println("OpenCL disabled. " + disabledInfo)
}
