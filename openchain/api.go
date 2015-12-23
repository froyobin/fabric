/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package openchain

import (
	"fmt"

	"golang.org/x/net/context"

	google_protobuf1 "google/protobuf"

	"github.com/openblockchain/obc-peer/openchain/ledger"
	pb "github.com/openblockchain/obc-peer/protos"
)

// ServerOpenchain defines the Openchain server object, which holds the
// Ledger data structure.
type ServerOpenchain struct {
	ledger *ledger.Ledger
}

// NewOpenchainServer creates a new instance of the ServerOpenchain.
func NewOpenchainServer() (*ServerOpenchain, error) {
	// Get a handle to the Ledger singleton.
	ledger, err := ledger.GetLedger()
	if err != nil {
		return nil, err
	}
	s := &ServerOpenchain{ledger: ledger}

	/*
		num := ledger.GetBlockchainSize()
		for i := uint64(0); i < num; i++ {
			block, err := ledger.GetBlockByNumber(i)
			if err != nil {
				log.Info("\nError retrieving block from blockchain: %s\n", err)
				return nil, err
			}
			log.Info("\n\nBlock %d:\n\n%s\n\n", i, block)
		}
	*/
	return s, nil
}

// GetBlockchainInfo returns information about the blockchain ledger such as
// height, current block hash, and previous block hash.
func (s *ServerOpenchain) GetBlockchainInfo(ctx context.Context, e *google_protobuf1.Empty) (*pb.BlockchainInfo, error) {
	blockchainInfo, err := s.ledger.GetBlockchainInfo()
	if blockchainInfo.Height == 0 {
		return nil, fmt.Errorf("No blocks in blockchain.")
	}
	return blockchainInfo, err
}

// GetBlockByNumber returns the data contained within a specific block in the
// blockchain. The genesis block is block zero.
func (s *ServerOpenchain) GetBlockByNumber(ctx context.Context, num *pb.BlockNumber) (*pb.Block, error) {
	block, err := s.ledger.GetBlockByNumber(num.Number)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving block from blockchain: %s", err)
	}
	if block == nil {
		return nil, fmt.Errorf("Requested block not in blockchain.")
	}
	return block, nil
}

// GetBlockCount returns the current number of blocks in the blockchain data
// structure.
func (s *ServerOpenchain) GetBlockCount(ctx context.Context, e *google_protobuf1.Empty) (*pb.BlockCount, error) {
	// Total number of blocks in the blockchain.
	size := s.ledger.GetBlockchainSize()

	// Check the number of blocks in the blockchain. If the blockchain is empty,
	// return error. There will always be at least one block in the blockchain,
	// the genesis block.
	if size > 0 {
		count := &pb.BlockCount{Count: size}
		return count, nil
	}

	return nil, fmt.Errorf("No blocks in blockchain.")
}

// GetState returns the value for a particular chaincode ID and key
func (s *ServerOpenchain) GetState(ctx context.Context, chaincodeID, key string) ([]byte, error) {
	return s.ledger.GetState(chaincodeID, key, true)
}
