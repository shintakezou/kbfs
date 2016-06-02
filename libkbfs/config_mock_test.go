// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package libkbfs

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/keybase/client/go/logger"
	"golang.org/x/net/context"
)

type FakeObserver struct {
	localChange  Node
	batchChanges []NodeChange
	ctx          context.Context
}

func (fn *FakeObserver) LocalChange(ctx context.Context,
	node Node, write WriteRange) {
	fn.localChange = node
	fn.ctx = ctx
}

func (fn *FakeObserver) BatchChanges(
	ctx context.Context, nodeChanges []NodeChange) {
	fn.batchChanges = nodeChanges
	fn.ctx = ctx
}

func (fn *FakeObserver) TlfHandleChange(ctx context.Context,
	newHandle *TlfHandle) {
	return
}

type ConfigMock struct {
	ConfigLocal

	// local references to the proper mock type
	mockKbfs       *MockKBFSOps
	mockKbpki      *MockKBPKI
	mockKbd        *MockKeybaseDaemon
	mockKeyman     *MockKeyManager
	mockRep        *MockReporter
	mockMdcache    *MockMDCache
	mockKcache     *MockKeyCache
	mockBcache     *MockBlockCache
	mockCrypto     *MockCrypto
	mockCodec      *MockCodec
	mockMdops      *MockMDOps
	mockKops       *MockKeyOps
	mockBops       *MockBlockOps
	mockMdserv     *MockMDServer
	mockKserv      *MockKeyServer
	mockBserv      *MockBlockServer
	mockBsplit     *MockBlockSplitter
	mockNotifier   *MockNotifier
	mockClock      *MockClock
	mockRekeyQueue *MockRekeyQueue
	observer       *FakeObserver
	ctr            *SafeTestReporter
}

func NewConfigMock(c *gomock.Controller, ctr *SafeTestReporter) *ConfigMock {
	config := &ConfigMock{}
	config.mockKbfs = NewMockKBFSOps(c)
	config.SetKBFSOps(config.mockKbfs)
	config.mockKbd = NewMockKeybaseDaemon(c)
	config.SetKeybaseDaemon(config.mockKbd)
	config.mockKbpki = NewMockKBPKI(c)
	config.SetKBPKI(config.mockKbpki)
	config.mockKeyman = NewMockKeyManager(c)
	config.SetKeyManager(config.mockKeyman)
	config.mockRep = NewMockReporter(c)
	config.SetReporter(config.mockRep)
	config.mockMdcache = NewMockMDCache(c)
	config.SetMDCache(config.mockMdcache)
	config.mockKcache = NewMockKeyCache(c)
	config.SetKeyCache(config.mockKcache)
	config.mockBcache = NewMockBlockCache(c)
	config.SetBlockCache(config.mockBcache)
	config.mockCrypto = NewMockCrypto(c)
	config.SetCrypto(config.mockCrypto)
	config.mockCodec = NewMockCodec(c)
	config.mockCodec.EXPECT().RegisterType(gomock.Any(), gomock.Any()).
		AnyTimes().Return()
	config.mockCodec.EXPECT().RegisterIfaceSliceType(gomock.Any(),
		gomock.Any(), gomock.Any()).AnyTimes().Return()
	config.SetCodec(config.mockCodec)
	config.mockMdops = NewMockMDOps(c)
	config.SetMDOps(config.mockMdops)
	config.mockKops = NewMockKeyOps(c)
	config.SetKeyOps(config.mockKops)
	config.mockBops = NewMockBlockOps(c)
	config.SetBlockOps(config.mockBops)
	config.mockMdserv = NewMockMDServer(c)
	config.SetMDServer(config.mockMdserv)
	config.mockKserv = NewMockKeyServer(c)
	config.SetKeyServer(config.mockKserv)
	config.mockBserv = NewMockBlockServer(c)
	config.SetBlockServer(config.mockBserv)
	config.mockBsplit = NewMockBlockSplitter(c)
	config.SetBlockSplitter(config.mockBsplit)
	config.mockNotifier = NewMockNotifier(c)
	config.SetNotifier(config.mockNotifier)
	config.mockClock = NewMockClock(c)
	config.SetClock(config.mockClock)
	config.mockRekeyQueue = NewMockRekeyQueue(c)
	config.SetRekeyQueue(config.mockRekeyQueue)
	config.observer = &FakeObserver{}
	config.ctr = ctr
	config.SetLoggerMaker(func(m string) logger.Logger {
		return logger.NewTestLogger(ctr.t)
	})
	// turn off background flushing by default during tests
	config.noBGFlush = true

	config.maxFileBytes = maxFileBytesDefault
	config.maxNameBytes = maxNameBytesDefault
	config.maxDirBytes = maxDirBytesDefault
	config.rwpWaitTime = rekeyWithPromptWaitTimeDefault

	config.qrPeriod = 0 * time.Second // no auto reclamation
	config.qrUnrefAge = qrUnrefAgeDefault

	return config
}

// CheckStateOnShutdown implements the Config interface for ConfigLocal.
func (c *ConfigMock) CheckStateOnShutdown() bool {
	return false
}