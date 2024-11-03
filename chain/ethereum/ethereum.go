package ethereum

import (
	"context"
	"math/big"
	"strconv"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/CavnHan/wallet-chain-account/chain"
	"github.com/CavnHan/wallet-chain-account/config"
	"github.com/CavnHan/wallet-chain-account/rpc/account"
	"github.com/CavnHan/wallet-chain-account/rpc/common"
)

//对接所有的RPC

const ChainName = "Ethereum"

type ChainAdaptor struct {
	ethClient     EthClient
	ethDataClient *EthData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	ethClient, err := DialEthClient(context.Background(), conf.WalletNode.Eth.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := NewEthDataClient(conf.WalletNode.Eth.DataApiUrl, conf.WalletNode.Eth.DataApiKey, time.Duration(conf.WalletNode.Eth.TimeOut))
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil
}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	addressCommon := ethcommon.BytesToAddress(crypto.Keccak256(req.PublicKey[1:])[12:])
	return &account.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "convert address successs",
		Address: addressCommon.String(),
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	var blockHash ethcommon.Hash
	if req.Hash == "" {
		blockHash = ethcommon.Hash{}
	} else {
		blockHash = ethcommon.HexToHash(req.Hash)
	}
	blockInfo, err := c.ethClient.BlockHeaderByHash(blockHash)
	if err != nil {
		log.Error("get latest block header fail", "err", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get latest block header fail",
		}, nil
	}
	blockHeader := &account.BlockHeader{
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.String(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: blockInfo.ParentBeaconRoot.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            string(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  blockInfo.WithdrawalsHash.String(),
		BlobGasUsed:      *blockInfo.BlobGasUsed,
		ExcessBlobGas:    *blockInfo.ExcessBlobGas,
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil
	} else {
		blockNumber = big.NewInt(req.Height)
	}
	blockInfo, err := c.ethClient.BlockHeaderByNumber(blockNumber)
	if err != nil {
		log.Error("get latest block header fail !", "err:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get lates block number fail",
		}, nil
	}
	blockHeader := &account.BlockHeader{
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.String(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: blockInfo.ParentBeaconRoot.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            string(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  blockInfo.WithdrawalsHash.String(),
		BlobGasUsed:      *blockInfo.BlobGasUsed,
		ExcessBlobGas:    *blockInfo.ExcessBlobGas,
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByNumber(big.NewInt(req.Height))
	if err != nil {
		log.Error("block by number error", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitItem := &account.BlockInfoTransactionList{
			From:   "0x000",
			To:     v.To,
			Hash:   v.Hash,
			Time:   "0",
			Amount: "10",
			Fee:    "0",
			Status: "1",
		}
		txListRet = append(txListRet, bitItem)
	}
	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Msg:          "get block by number success",
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	return nil, nil
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	nonceResult, err := c.ethClient.TxCountByAddress(ethcommon.HexToAddress(req.Address))
	if err != nil {
		log.Error("get nonce by account fail", "err:", err)
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get nonce by address fail",
		}, nil
	}
	//获取balance
	//合约地址或者普通账户
	balanceResult, err := c.ethDataClient.getBalanceByaddress(req.ContractAddress, req.Address)
	if err != nil {
		log.Error("get balance by address fail", "err:", err)
		return &account.AccountResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance by address fail",
			Balance: "0",
		}, err
	}
	return &account.AccountResponse{
		Code:          common.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      nonceResult.String(),
		Balance:       balanceResult.BalanceStr,
	}, nil

}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
