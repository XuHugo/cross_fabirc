package main

import(
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb"github.com/hyperledger/fabric-protos-go/peer"
)

type A struct {

}

const INIT_DATA = "init_data"
const GATEWAY_PRE = "gateway"
const ASSET_PRE = "asset"
const PROTOCOL_VERSION = "1.0.0"

const TxResultEnum = {'INIT':'0',   'ACK_SUCCESS':'1',  'ACK_FAIL':'2',     'ACK_TIMEOUT':'3'}
const TxRefundedEnum = {'NONE':'0',     'TODO':'1',     'REFUNDED': '2'}

type InitData struct{
	ChainCode     string `json:"chainCode"`
	CreateTime    string `json:"createTime"`
	CreateAccount string `json:"createAccount"`
	Managers      []string `json:"managers"`
	IsRelay       string `json:"isRelay"`
	BlockHeight   string `json:"blockHeight"`
}

type CrossTx struct {
	CrossTxNo    string `json:"crossTxNo"`
	SrcChainCode string `json:"srcChainCode"`
	DestChainCode string `json:"destChainCode"`
	SrcAddress  string `json:"srcAddress"`
	DestAddress string `json:"destAddress"`
	PayloadType string `json:"payloadType"`
	Payload   string `json:"payload"`
	Remark    string `json:"remark"`
	Result    string `json:"result"`
	Refunded  string `json:"refunded"`
	Extension  string `json:"extension"`
	SendProofs string `json:"sendProofs"`
	AckProofs string `json:"ackProofs"`
	Version   string `json:"version"`
	origin string `json:"origin"`
}

func (t *A) Init (stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters() // Input: "init, args[0], args[1] ..."
	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting a key and a value0")
	}
	chainCode := args[0]
	isRelay := args[1]
	managersList := args[2]

	err :=  stub.PutState(args[0],[]byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	managers := strings.Split(managersList, "&")

	CreateTime, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}

	creator, err := stub.GetCreator()
	certStart := bytes.IndexAny(creator, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("no certificate found")
	}

	var initdata = InitData{
		ChainCode:     chainCode,
		CreateTime:    string(CreateTime),
		CreateAccount: string(creator[certStart:]),
		Managers:      managers,
		IsRelay:       isRelay,
		BlockHeight:   "fabric",
	}
	initdataJon, err := json.Marshal(initdata)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(INIT_DATA,initdataJon)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *A) Invoke (stub shim.ChaincodeStubInterface) pb.Response {

	fn, args := stub.GetFunctionAndParameters ()

	if fn == "setGateway" {
		return t.setGateway(stub, args)
	} else if fn == "startTx" {
		return t.substartTx(stub,args)
	}else if fn == "sendTx" {
		return t.sendTx(stub,args)
	}else if fn == "sendAcked" {
		return t.sendAcked(stub,args)
	}else if fn == "takeOut" {
		return t.subTakeOut(stub,args)
	}else if fn == "transfer" {
		return t.transfer(stub,args)
	}else if fn == "getCrossTx"{
		return t.getCrossTx(stub,args)
	}else if fn == "balanceOf"{
		return t.balanceOf(stub,args)
	}else if fn == "version" {
		return t.version(stub,args)
	}

	return shim.Error("Invoke fn error")
}

func verifyManager(stub shim.ChaincodeStubInterface,){
	manager_list,err := stub.GetState("init_data")
	if err != nil {
		return shim.Error(err.Error)
	}
	usr,err := stub.getCreator()
	if err != nil {
		return shim.Error(err.Error)
	}
}

func (t *A) setGateway (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 3 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}
	chainChode := args[0]
	gatewayList := args[1]
	remark := args[2]

	//判断超级节点
	verifyManager(stub)
	//校验参数
	//checkSetGatewayParamObj();
	//保存公证人
	err := stub.PutState("init_data",[]byte(initdata))
	if err != nil {
		return shim.Error(err.Error())
	}

	//log
	var event = chainChode + gatewayList
	stub.setEvent("setGateway",event)

	return shim.Success(nil)
}

func _checkParamObj(paramObj){

}

func _subStartTxCall(paramObj){
	_checkParamObj(paramObj)
}
func _subStartTxData(paramObj){
	_checkParamObj(paramObj)
}
func _subStartTxTransferSgas(paramObj){
	_checkParamObj(paramObj)
}

func _createCrossTxObj(crossTxNo, srcChainCode, destChainCode, origin, paramObj)CrossTx{
	var crosstx = CrossTx{
		CrossTxNo:     crossTxNo,
		SrcChainCode:  srcChainCode,
		DestChainCode: destChainCode,
		SrcAddress:    paramObj.srcAddress,
		DestAddress:   paramObj.destAddress,
		PayloadType:   paramObj.payloadType,
		Payload:       paramObj.payload,
		Remark:        paramObj.remark,
		Result:        0,
		Refunded:      0,
		Extension:     paramObj.extension,
		SendProofs:    nil,
		AckProofs:     nil,
		Version:       PROTOCOL_VERSION,
		origin:        origin,
	}

	return crosstx
}

func (t *A) substartTx (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 2 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}

	switch crossTxObj.payloadType {
	case "2":
		_subStartTxCall(paramObj)
	case "3":
		_subStartTxData(paramObj)
	case "4":
		_subStartTxTransferSgas(paramObj)
	}

	var crossTxObj = _createCrossTxObj("", _getChainCode(), paramObj.destChainCode, TxOriginEnum.SRC, paramObj)
	err := stub.PutState(crossTxObj.CrossTxNo,[]byte(crossTxObj))
	if err != nil {
		return shim.Error(err.Error())
	}

	//log
	var event = args[1]+args[2];
	stub.setEvent("setGateway",event);

	return shim.Success(nil)
}

func _checkGateNode(stub shim.ChaincodeStubInterface, chainCode){
	nodegatekey, err := stub.CreateCompositeKey("gateway", []string{chainCode,} )
	if err != nil{
		return shim.Error(err.Error())
	}
	nodeGateways, err := stub.GetState(nodegatekey)
	if err != nil{
		return shim.Error(err.Error())
	}
}

func _subSendTxCall(paramObj, crossTxObj){

}

func _subSendTxData(paramObj, crossTxObj){

}

func _subSendTxTransferSgas(paramObj, crossTxObj){

}

func (t *A) sendTx (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 2 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}

	_checkGateNode(paramObj.destChainCode)

	var crossTxObj = _createCrossTxObj(paramObj.crossTxNo, paramObj.srcChainCode, paramObj.destChainCode, TxOriginEnum.DEST, paramObj)
	switch crossTxObj.payloadType {
	case "2":
		_subSendTxCall(paramObj, crossTxObj)
	case "3":
		_subSendTxData(paramObj, crossTxObj)
	case "4":
		_subSendTxTransferSgas(paramObj, crossTxObj)
	}
	crossTxObj.SendProofs.push(paramObj.proof)
	err := stub.PutState(crossTxObj.CrossTxNo,[]byte(crossTxObj))
	if err != nil {
		return shim.Error(err.Error())
	}

	//log
	var event = args[1]+args[2];
	stub.setEvent("setGateway",event);

	return shim.Success(nil)
}

func _getOrgin(origin) string{
	if origin == 0{
		return "src"
	} else if origin == 1 {
		return "dest"
	} else if origin == 2{
		return "relay"
	}
	return "err"
}

func _subSendAckSrcCall(paramObj, crossTxObj){

	_checkGateNode(crossTxObj.srcChainCode)
	if(paramObj.result !== TxResultEnum.ACK_SUCCESS){

		crossTxObj.refunded = TxRefundedEnum.TODO;
	}

	return crossTxObj;
}

func _subSendAckSrcData(paramObj, crossTxObj){
	_checkGateNode(crossTxObj.srcChainCode)
	return crossTxObj;
}

func _subSendAckSrcTransferSgas(paramObj, crossTxObj){

	_checkGateNode(crossTxObj.srcChainCode)
	if(paramObj.result !== TxResultEnum.ACK_SUCCESS){

		crossTxObj.refunded = TxRefundedEnum.TODO;
	}

	return crossTxObj;
}

func _subSendAckDestCommon(paramObj, crossTxObj){
	_checkGateNode(crossTxObj.srcChainCode)
	return crossTxObj;
}

func (t *A) sendAcked (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 5 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}

	crosstxno := args[0]
	result := args[1]
	version := args[2]
	ledgerSeq := args[3]
	txHash := args[4]

	crossTxObj, err := stub.GetState(crosstxno)
	if err != nil{
		return shim.Error(err.Error())
	}

	orgin := _getOrgin(crossTxObj.origin)

	switch orgin {
		case "src":
			switch crossTxObj.payloadType {
			case "2":
				_subSendAckSrcCall(paramObj, crossTxObj)
			case "3":
				_subSendAckSrcData(paramObj, crossTxObj)
			case "4":
				_subSendAckSrcTransferSgas(paramObj, crossTxObj)
			}
		case"dest":
			switch crossTxObj.payloadType {
			case "2":
				_subSendAckDestCommon(paramObj, crossTxObj)
			case "3":
				_subSendAckDestCommon(paramObj, crossTxObj)
			case "4":
				_subSendAckDestCommon(paramObj, crossTxObj)
			}
	}
	crossTxObj.result = paramObj.result
	paramObj.proof.verifierBid = Chain.msg.sender
	crossTxObj.ackProofs.push(paramObj.proof)

	err = stub.PutState(crossTxObj.crossTxNo,[]byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.setEvent("sendAcked",crosstxno)
	return shim.Success(nil)
}

func (t *A) subTakeOut (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 2 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}
	crosstxno := args[0]
	toaddr := args[1]

	crossTxObj, err := stub.GetState(crosstxno)
	if err != nil{
		return shim.Error(err.Error())
	}

	crossTxObj.refuned = 2

	err = stub.PutState(crosstxno, []byte(crossTxObj))
	if err != nil {
		return shim.Error(err.Error())
	}

	stub.setEvent("takeOut",crosstxno)

	return shim.Success(nil)
}

func (t *A) transfer (stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 3 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}
	asset := args[0]
	var to = args[1]

	value, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("expecting a integer value")
	}

	sender := stub.GetCreator()

	senderKey, err := stub.CreateCompositeKey("asset", []string{asset, sender} )
	if err != nil{
		return shim.Error(err.Error())
	}
	senderValue, err := stub.GetState(senderKey)
	if err != nil{
		return shim.Error(err.Error())
	}
	if senderValue < value{
		return shim.Error("balance insufficient!")
	}

	toKey, err := stub.CreateCompositeKey("asset", []string{asset, to} )
	if err != nil{
		return shim.Error(err.Error())
	}
	toValue, err := stub.GetState(toKey)
	if err != nil{
		return shim.Error(err.Error())
	}
	var newtovalue = toValue + value

	err = stub.PutState(toKey, []byte(newtovalue))
	if err != nil {
		return shim.Error(err.Error())
	}

	var newsendervalue = senderValue - value

	err = stub.PutState(senderKey, []byte(newsendervalue))
	if err != nil {
		return shim.Error(err.Error())
	}
	event := sender + to + value
	stub.setEvent("transfer",event)

	return shim.Success(nil)
}

func (t *A)  balanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}
	asset := args[0]
	owner := args[1]

	addrKey, err := stub.CreateCompositeKey("asset", []string{asset, owner} )
	if err != nil{
		return shim.Error(err.Error())
	}
	Value, err := stub.GetState(addrKey)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(Value)
}

func (t *A)  getCrossTx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return fmt.Errorf("Incorrect arguments. Expecting a key and value")
	}
	crosstx := args[0]

	data, err := stub.GetState(crosstx)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (t *A)  version (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(PROTOCOL_VERSION)
}


func main() {
	err := shim.Start(new(A))
	if err != nil {
		fmt.Printf("start error: %s",err)
	}
}
