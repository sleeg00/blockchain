// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: blockchain.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *VersionRequest) Reset() {
	*x = VersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRequest) ProtoMessage() {}

func (x *VersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRequest.ProtoReflect.Descriptor instead.
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{0}
}

func (x *VersionRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *VersionRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Height int64  `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{1}
}

func (x *VersionResponse) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *VersionResponse) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

type SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeFrom string `protobuf:"bytes,1,opt,name=node_from,json=nodeFrom,proto3" json:"node_from,omitempty"`
	From     string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To       string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Amount   int32  `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
	NodeTo   string `protobuf:"bytes,5,opt,name=node_to,json=nodeTo,proto3" json:"node_to,omitempty"`
	MineNow  bool   `protobuf:"varint,6,opt,name=mine_now,json=mineNow,proto3" json:"mine_now,omitempty"`
	Block    *Block `protobuf:"bytes,7,opt,name=Block,proto3" json:"Block,omitempty"`
}

func (x *SendRequest) Reset() {
	*x = SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRequest) ProtoMessage() {}

func (x *SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRequest.ProtoReflect.Descriptor instead.
func (*SendRequest) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{2}
}

func (x *SendRequest) GetNodeFrom() string {
	if x != nil {
		return x.NodeFrom
	}
	return ""
}

func (x *SendRequest) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SendRequest) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *SendRequest) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *SendRequest) GetNodeTo() string {
	if x != nil {
		return x.NodeTo
	}
	return ""
}

func (x *SendRequest) GetMineNow() bool {
	if x != nil {
		return x.MineNow
	}
	return false
}

func (x *SendRequest) GetBlock() *Block {
	if x != nil {
		return x.Block
	}
	return nil
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp     int64          `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Transactions  []*Transaction `protobuf:"bytes,2,rep,name=transactions,proto3" json:"transactions,omitempty"`
	PrevBlockHash []byte         `protobuf:"bytes,3,opt,name=prev_block_hash,json=prevBlockHash,proto3" json:"prev_block_hash,omitempty"`
	Hash          []byte         `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Nonce         int32          `protobuf:"varint,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Height        int32          `protobuf:"varint,6,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{3}
}

func (x *Block) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Block) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Block) GetPrevBlockHash() []byte {
	if x != nil {
		return x.PrevBlockHash
	}
	return nil
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetNonce() int32 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Block) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   []byte      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Vin  []*TXInput  `protobuf:"bytes,2,rep,name=vin,proto3" json:"vin,omitempty"`
	Vout []*TXOutput `protobuf:"bytes,3,rep,name=vout,proto3" json:"vout,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{4}
}

func (x *Transaction) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Transaction) GetVin() []*TXInput {
	if x != nil {
		return x.Vin
	}
	return nil
}

func (x *Transaction) GetVout() []*TXOutput {
	if x != nil {
		return x.Vout
	}
	return nil
}

type TXInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid      []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Vout      int64  `protobuf:"varint,2,opt,name=vout,proto3" json:"vout,omitempty"`
	Signature []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	PubKey    []byte `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
}

func (x *TXInput) Reset() {
	*x = TXInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TXInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TXInput) ProtoMessage() {}

func (x *TXInput) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TXInput.ProtoReflect.Descriptor instead.
func (*TXInput) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{5}
}

func (x *TXInput) GetTxid() []byte {
	if x != nil {
		return x.Txid
	}
	return nil
}

func (x *TXInput) GetVout() int64 {
	if x != nil {
		return x.Vout
	}
	return 0
}

func (x *TXInput) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *TXInput) GetPubKey() []byte {
	if x != nil {
		return x.PubKey
	}
	return nil
}

type TXOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value      int64  `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	PubKeyHash []byte `protobuf:"bytes,2,opt,name=pub_key_hash,json=pubKeyHash,proto3" json:"pub_key_hash,omitempty"`
}

func (x *TXOutput) Reset() {
	*x = TXOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TXOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TXOutput) ProtoMessage() {}

func (x *TXOutput) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TXOutput.ProtoReflect.Descriptor instead.
func (*TXOutput) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{6}
}

func (x *TXOutput) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *TXOutput) GetPubKeyHash() []byte {
	if x != nil {
		return x.PubKeyHash
	}
	return nil
}

type SendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *SendResponse) Reset() {
	*x = SendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendResponse) ProtoMessage() {}

func (x *SendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendResponse.ProtoReflect.Descriptor instead.
func (*SendResponse) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{7}
}

func (x *SendResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

type CreateWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *CreateWalletRequest) Reset() {
	*x = CreateWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletRequest) ProtoMessage() {}

func (x *CreateWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletRequest.ProtoReflect.Descriptor instead.
func (*CreateWalletRequest) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{8}
}

func (x *CreateWalletRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

type CreateWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *CreateWalletResponse) Reset() {
	*x = CreateWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletResponse) ProtoMessage() {}

func (x *CreateWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletResponse.ProtoReflect.Descriptor instead.
func (*CreateWalletResponse) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{9}
}

func (x *CreateWalletResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type CreateBlockchainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	NodeId  string `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *CreateBlockchainRequest) Reset() {
	*x = CreateBlockchainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateBlockchainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBlockchainRequest) ProtoMessage() {}

func (x *CreateBlockchainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBlockchainRequest.ProtoReflect.Descriptor instead.
func (*CreateBlockchainRequest) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{10}
}

func (x *CreateBlockchainRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *CreateBlockchainRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// ---------------------------
type CreateBlockchainResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *CreateBlockchainResponse) Reset() {
	*x = CreateBlockchainResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateBlockchainResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBlockchainResponse) ProtoMessage() {}

func (x *CreateBlockchainResponse) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBlockchainResponse.ProtoReflect.Descriptor instead.
func (*CreateBlockchainResponse) Descriptor() ([]byte, []int) {
	return file_blockchain_proto_rawDescGZIP(), []int{11}
}

func (x *CreateBlockchainResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_blockchain_proto protoreflect.FileDescriptor

var file_blockchain_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x22, 0x44,
	0x0a, 0x0e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x22, 0x42, 0x0a, 0x0f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0xc3, 0x01, 0x0a, 0x0b, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64,
	0x65, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x6f, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x54, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x69,
	0x6e, 0x65, 0x5f, 0x6e, 0x6f, 0x77, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6d, 0x69,
	0x6e, 0x65, 0x4e, 0x6f, 0x77, 0x12, 0x27, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0xcc,
	0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3b, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x70, 0x72,
	0x65, 0x76, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x6e, 0x0a,
	0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x03,
	0x76, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x58, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x03,
	0x76, 0x69, 0x6e, 0x12, 0x28, 0x0a, 0x04, 0x76, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x54,
	0x58, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x04, 0x76, 0x6f, 0x75, 0x74, 0x22, 0x68, 0x0a,
	0x07, 0x54, 0x58, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x76, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x76, 0x6f, 0x75, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x70, 0x75, 0x62, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x22, 0x42, 0x0a, 0x08, 0x54, 0x58, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x70, 0x75, 0x62,
	0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0a, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0x2a, 0x0a, 0x0c, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2e, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x30, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x4c, 0x0a, 0x17, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x17,
	0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0x84, 0x02, 0x0a, 0x11, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x17, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x51, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x12, 0x1f, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x23, 0x2e, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24,
	0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockchain_proto_rawDescOnce sync.Once
	file_blockchain_proto_rawDescData = file_blockchain_proto_rawDesc
)

func file_blockchain_proto_rawDescGZIP() []byte {
	file_blockchain_proto_rawDescOnce.Do(func() {
		file_blockchain_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockchain_proto_rawDescData)
	})
	return file_blockchain_proto_rawDescData
}

var file_blockchain_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_blockchain_proto_goTypes = []interface{}{
	(*VersionRequest)(nil),           // 0: blockchain.VersionRequest
	(*VersionResponse)(nil),          // 1: blockchain.VersionResponse
	(*SendRequest)(nil),              // 2: blockchain.SendRequest
	(*Block)(nil),                    // 3: blockchain.Block
	(*Transaction)(nil),              // 4: blockchain.Transaction
	(*TXInput)(nil),                  // 5: blockchain.TXInput
	(*TXOutput)(nil),                 // 6: blockchain.TXOutput
	(*SendResponse)(nil),             // 7: blockchain.SendResponse
	(*CreateWalletRequest)(nil),      // 8: blockchain.CreateWalletRequest
	(*CreateWalletResponse)(nil),     // 9: blockchain.CreateWalletResponse
	(*CreateBlockchainRequest)(nil),  // 10: blockchain.CreateBlockchainRequest
	(*CreateBlockchainResponse)(nil), // 11: blockchain.CreateBlockchainResponse
}
var file_blockchain_proto_depIdxs = []int32{
	3,  // 0: blockchain.SendRequest.Block:type_name -> blockchain.Block
	4,  // 1: blockchain.Block.transactions:type_name -> blockchain.Transaction
	5,  // 2: blockchain.Transaction.vin:type_name -> blockchain.TXInput
	6,  // 3: blockchain.Transaction.vout:type_name -> blockchain.TXOutput
	2,  // 4: blockchain.BlockchainService.Send:input_type -> blockchain.SendRequest
	8,  // 5: blockchain.BlockchainService.CreateWallet:input_type -> blockchain.CreateWalletRequest
	10, // 6: blockchain.BlockchainService.CreateBlockchain:input_type -> blockchain.CreateBlockchainRequest
	7,  // 7: blockchain.BlockchainService.Send:output_type -> blockchain.SendResponse
	9,  // 8: blockchain.BlockchainService.CreateWallet:output_type -> blockchain.CreateWalletResponse
	11, // 9: blockchain.BlockchainService.CreateBlockchain:output_type -> blockchain.CreateBlockchainResponse
	7,  // [7:10] is the sub-list for method output_type
	4,  // [4:7] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_blockchain_proto_init() }
func file_blockchain_proto_init() {
	if File_blockchain_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blockchain_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TXInput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TXOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateBlockchainRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_blockchain_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateBlockchainResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_blockchain_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_blockchain_proto_goTypes,
		DependencyIndexes: file_blockchain_proto_depIdxs,
		MessageInfos:      file_blockchain_proto_msgTypes,
	}.Build()
	File_blockchain_proto = out.File
	file_blockchain_proto_rawDesc = nil
	file_blockchain_proto_goTypes = nil
	file_blockchain_proto_depIdxs = nil
}
