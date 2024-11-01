package chaindispatcher

import (
	"context"
	"github.com/CavnHan/wallet-chain-account/chain"
	"github.com/CavnHan/wallet-chain-account/chain/ethereum"
	"github.com/CavnHan/wallet-chain-account/config"
	"github.com/CavnHan/wallet-chain-account/rpc/account"
	"github.com/CavnHan/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strings"
)

//RPC请求拦截器

type CommonRequest interface {
	//获取链名
	GetChain() string
}

type CommonReply = account.SupportChainsResponse
type ChainType = string

type ChainDispatcher struct {
	//注册链,key:链名,value:链适配器
	registry map[ChainType]chain.IChainAdaptor
}

/**
 * @description: 创建调度器
 * @param config 配置文件
 */
func New(conf *config.Config) (*ChainDispatcher, error) {
	//初始化调度器
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}
	//链适配器工厂 map
	chainAdaptorFactorMap := map[string]func(conf *config.Config) (chain.IChainAdaptor, error){
		//链名:工厂方法
		//add 支持的链以及对应的工厂方法
		ethereum.ChainName: ethereum.NewChainAdaptor,
	}

	supportedChains := []string{
		ethereum.ChainName,
	}
	//遍历配置文件中的链，根据工厂返回对应的链的适配器，即为链对应的chainAdaptor的实现
	for _, c := range conf.Chains {
		if factory, ok := chainAdaptorFactorMap[c]; ok {
			//调用工厂方法，返回链适配器
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			//注册链适配器
			dispatcher.registry[c] = adaptor

		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil

}

/**
* @description: 拦截器
* @param ctx 上下文
* @param req 请求
* @param info 服务信息
* @param handler 处理器
* @return resp 响应
* @return err 错误
 */
func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			//打印堆栈信息
			log.Debug(string(debug.Stack()))
			//返回内部错误
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	//获取请求方法名
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	//获取链名
	chainName := req.(CommonRequest).GetChain()
	log.Info(method, "chain", chainName, "req", req)

	//调用handler处理请求
	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp:", resp, "err:", err)
	return
}

/**
* @description: 预处理
* @param req 请求
* @return resp 响应
 */
func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply) {
	chainName := req.(CommonRequest).GetChain()
	if _, ok := d.registry[chainName]; !ok {
		return &CommonReply{
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}
	}
	return nil
}

func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	//前置处理
	resp := d.preHandler(request)
	if resp != nil {
		return &account.SupportChainsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	//调用链适配器的方法
	return d.registry[request.Chain].GetSupportChains(request)
}

func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "covert address fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].ConvertAddress(request)
}
func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *account.BlockNumberRequest) (*account.BlockResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by number fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *account.BlockHashRequest) (*account.BlockResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by number fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByNumber(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *account.AccountRequest) (*account.AccountResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account information fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetAccount(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *account.FeeRequest) (*account.FeeResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get fee fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetFee(request)
}

func (d *ChainDispatcher) SendTx(ctx context.Context, request *account.SendTxRequest) (*account.SendTxResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by address fail pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *account.TxHashRequest) (*account.TxHashResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetTxByHash(request)
}

func (d *ChainDispatcher) GetBlockByRange(ctx context.Context, request *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.BlockByRangeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get blcok by range fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetBlockByRange(request)
}

func (d *ChainDispatcher) CreateUnSignTransaction(ctx context.Context, request *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un sign tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].CreateUnSignTransaction(request)
}

func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "signed tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].BuildSignedTransaction(request)
}

func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.DecodeTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "decode tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].DecodeTransaction(request)
}

func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.VerifyTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "verify tx fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].VerifySignedTransaction(request)
}

func (d *ChainDispatcher) GetExtraData(ctx context.Context, request *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	resp := d.preHandler(request)
	if resp != nil {
		return &account.ExtraDataResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get extra data fail at pre handle",
		}, nil
	}
	return d.registry[request.Chain].GetExtraData(request)
}
