package chaincode

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
	currentMemory StateMemory
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type ElementState int

const (
	DISABLE = iota
	ENABLE
	DONE
)

type Message struct {
	MessageID     string       `json:"messageID"`
	SendMspID     string       `json:"sendMspID"`
	ReceiveMspID  string       `json:"receiveMspID"`
	FireflyTranID string       `json:"fireflyTranID"`
	MsgState      ElementState `json:"msgState"`
}

type Gateway struct {
	GatewayID    string       `json:"gatewayID"`
	GatewayState ElementState `json:"gatewayState"`
}

type ActionEvent struct {
	EventID    string       `json:"eventID"`
	EventState ElementState `json:"eventState"`
}

type StateMemory struct {
	Confirm bool `json:"confirm"`
	Cancel  bool `json:"cancel"`
}

// Construct
func NewMessage(messageID, sendMspID, receiveMspID, fireflyTranID string, msgState ElementState) *Message {
	return &Message{
		MessageID:     messageID,
		SendMspID:     sendMspID,
		ReceiveMspID:  receiveMspID,
		FireflyTranID: fireflyTranID,
		MsgState:      msgState,
	}
}

func NewGateway(gatewayID string, gatewayState ElementState) *Gateway { //返回实际值是新建一个对象副本
	return &Gateway{
		GatewayID:    gatewayID,
		GatewayState: gatewayState,
	}
}

func NewStateMemory(confirm, cancel bool) *StateMemory {
	return &StateMemory{
		Confirm: confirm,
		Cancel:  cancel,
	}
}

// Create function
func (cc *SmartContract) CreateMessage(ctx contractapi.TransactionContextInterface, messageID string, sendMspID string, receiveMspID string, fireflyTranID string, msgState ElementState) (*Message, error) {
	stub := ctx.GetStub()

	// 检查是否存在具有相同ID的记录
	existingData, err := stub.GetState(messageID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData != nil {
		return nil, fmt.Errorf("消息 %s 已存在", messageID)
	}

	// 创建消息对象
	msg := &Message{
		MessageID:     messageID,
		SendMspID:     sendMspID,
		ReceiveMspID:  receiveMspID,
		FireflyTranID: fireflyTranID,
		MsgState:      msgState,
	}

	// 将消息对象序列化为JSON字符串并保存在状态数据库中
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("序列化消息数据时出错: %v", err)
	}
	err = stub.PutState(messageID, msgJSON)
	if err != nil {
		return nil, fmt.Errorf("保存消息数据时出错: %v", err)
	}

	return msg, nil
}

func (cc *SmartContract) CreateGateway(ctx contractapi.TransactionContextInterface, gatewayID string, gatewayState ElementState) (*Gateway, error) {
	stub := ctx.GetStub()

	// 检查是否存在具有相同ID的记录
	existingData, err := stub.GetState(gatewayID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData != nil {
		return nil, fmt.Errorf("网关 %s 已存在", gatewayID)
	}

	// 创建网关对象
	gtw := &Gateway{
		GatewayID:    gatewayID,
		GatewayState: gatewayState,
	}

	// 将网关对象序列化为JSON字符串并保存在状态数据库中
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		return nil, fmt.Errorf("序列化网关数据时出错: %v", err)
	}
	err = stub.PutState(gatewayID, gtwJSON)
	if err != nil {
		return nil, fmt.Errorf("保存网关数据时出错: %v", err)
	}

	return gtw, nil
}

func (cc *SmartContract) CreateActionEvent(ctx contractapi.TransactionContextInterface, eventID string, eventState ElementState) (*ActionEvent, error) {
	stub := ctx.GetStub()

	// 创建ActionEvent对象
	actionEvent := &ActionEvent{
		EventID:    eventID,
		EventState: eventState,
	}

	// 将ActionEvent对象序列化为JSON字符串并保存在状态数据库中
	actionEventJSON, err := json.Marshal(actionEvent)
	if err != nil {
		return nil, fmt.Errorf("序列化事件数据时出错: %v", err)
	}
	err = stub.PutState(eventID, actionEventJSON)
	if err != nil {
		return nil, fmt.Errorf("保存事件数据时出错: %v", err)
	}

	return actionEvent, nil
}

// Read function
func (c *SmartContract) ReadMsg(ctx contractapi.TransactionContextInterface, messageID string) (*Message, error) {
	msgJSON, err := ctx.GetStub().GetState(messageID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if msgJSON == nil {
		errorMessage := fmt.Sprintf("Message %s does not exist", messageID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var msg Message
	err = json.Unmarshal(msgJSON, &msg)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &msg, nil
}

func (c *SmartContract) ReadGtw(ctx contractapi.TransactionContextInterface, gatewayID string) (*Gateway, error) {
	gtwJSON, err := ctx.GetStub().GetState(gatewayID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if gtwJSON == nil {
		errorMessage := fmt.Sprintf("Gateway %s does not exist", gatewayID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var gtw Gateway
	err = json.Unmarshal(gtwJSON, &gtw)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &gtw, nil
}

func (c *SmartContract) ReadEvent(ctx contractapi.TransactionContextInterface, eventID string) (*ActionEvent, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if eventJSON == nil {
		errorMessage := fmt.Sprintf("Event state %s does not exist", eventID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var event ActionEvent
	err = json.Unmarshal(eventJSON, &event)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &event, nil
}

// Change State  function
func (c *SmartContract) ChangeMsgState(ctx contractapi.TransactionContextInterface, messageID string, msgState ElementState) error {
	stub := ctx.GetStub()

	msg, err := c.ReadMsg(ctx, messageID)
	if err != nil {
		return err
	}

	msg.MsgState = msgState

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(messageID, msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (c *SmartContract) ChangeGtwState(ctx contractapi.TransactionContextInterface, gatewayID string, gtwState ElementState) error {
	stub := ctx.GetStub()

	gtw, err := c.ReadGtw(ctx, gatewayID)
	if err != nil {
		return err
	}

	gtw.GatewayState = gtwState

	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(gatewayID, gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (c *SmartContract) ChangeEventState(ctx contractapi.TransactionContextInterface, eventID string, eventState ElementState) error {
	stub := ctx.GetStub()

	actionEvent, err := c.ReadEvent(ctx, eventID)
	if err != nil {
		return err
	}

	actionEvent.EventState = eventState

	actionEventJSON, err := json.Marshal(actionEvent)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(eventID, actionEventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

//get all message

func (cc *SmartContract) GetAllMessages(ctx contractapi.TransactionContextInterface) ([]*Message, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err) //直接err也行
	}
	defer resultsIterator.Close()

	var messages []*Message
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代状态数据时出错: %v", err)
		}

		var message Message
		err = json.Unmarshal(queryResponse.Value, &message)
		if err != nil {
			return nil, fmt.Errorf("反序列化消息数据时出错: %v", err)
		}

		// 可以添加更多的筛选条件来仅获取特定类型或状态的消息
		messages = append(messages, &message)
	}

	return messages, nil
}

// InitLedger adds a base set of assets to the ledger

var isInited bool = false

func (cc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// Determines whether the chain code is initialized
	if isInited {
		errorMessage := "Chaincode has already been initialized"
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.CreateActionEvent(ctx, "StartEvent_1jtgn3j", ENABLE)

	cc.CreateGateway(ctx, "ExclusiveGateway_0hs3ztq", DISABLE)
	cc.CreateGateway(ctx, "ExclusiveGateway_106je4z", DISABLE)
	cc.CreateGateway(ctx, "EventBasedGateway_1fxpmyn", DISABLE)
	cc.CreateGateway(ctx, "ExclusiveGateway_0nzwv7v", DISABLE)
	// cc.CreateGateway(ctx, "EndEvent_0366pfz", DISABLE)

	// mspid    hotel:Participant_0sktaei       client:Participant_1080bkg
	cc.CreateMessage(ctx, "Message_045i10y", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // Check_room(string date, uint bedrooms)"
	cc.CreateMessage(ctx, "Message_0r9lypd", "Participant_0sktaei", "Participant_1080bkg", "", DISABLE) // Give_availability(bool confirm)
	cc.CreateMessage(ctx, "Message_1em0ee4", "Participant_0sktaei", "Participant_1080bkg", "", DISABLE) // Price_quotation(uint quotation)
	cc.CreateMessage(ctx, "Message_1nlagx2", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // Book_room(bool confirmation)
	cc.CreateMessage(ctx, "Message_0o8eyir", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // payment0(address payable to)
	cc.CreateMessage(ctx, "Message_1ljlm4g", "Participant_0sktaei", "Participant_1080bkg", "", DISABLE) // Give_ID(string booking_id)
	cc.CreateMessage(ctx, "Message_0m9p3da", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // cancel_order(bool cancel)
	cc.CreateMessage(ctx, "Message_1joj7ca", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // ask_refund(string ID)
	cc.CreateMessage(ctx, "Message_1etcmvl", "Participant_0sktaei", "Participant_1080bkg", "", DISABLE) // payment1(address payable to)
	cc.CreateMessage(ctx, "Message_1xm9dxy", "Participant_1080bkg", "Participant_0sktaei", "", DISABLE) // Cancel_order(string motivation)

	cc.CreateActionEvent(ctx, "EndEvent_146eii4", DISABLE)
	cc.CreateActionEvent(ctx, "EndEvent_08edp7f", DISABLE)
	cc.CreateActionEvent(ctx, "EndEvent_0366pfz", DISABLE)

	isInited = true

	stub.SetEvent("initLedgerEvent", []byte("Contract has been initialized successfully"))
	return nil
}

// =================================================================================================
func (cc *SmartContract) StartEvent_1jtgn3j(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	actionEvent, err := cc.ReadEvent(ctx, "StartEvent_1jtgn3j")
	if err != nil {
		return err
	}

	if actionEvent.EventState != ENABLE {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", actionEvent.EventID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	actionEvent.EventState = DONE
	actionEventJSON, err := json.Marshal(actionEvent)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("StartEvent_1jtgn3j", actionEventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stub.SetEvent("StartEvent_1jtgn3j", []byte("Contract has been started successfully"))

	gtw, err := cc.ReadGtw(ctx, "ExclusiveGateway_0hs3ztq")
	if err != nil {
		return err
	}

	gtw.GatewayState = ENABLE
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("ExclusiveGateway_0hs3ztq", gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	cc.ExclusiveGateway_0hs3ztq(ctx)

	return nil
}

func (cc *SmartContract) ExclusiveGateway_0hs3ztq(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, "ExclusiveGateway_0hs3ztq")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLE {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	gtw.GatewayState = DONE
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("ExclusiveGateway_0hs3ztq", gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stub.SetEvent("ExclusiveGateway_0hs3ztq", []byte("ExclusiveGateway_0hs3ztq has been done"))

	msg2, err := cc.ReadMsg(ctx, "Message_045i10y")
	if err != nil {
		return err
	}

	msg2.MsgState = ENABLE
	msg2JSON, err := json.Marshal(msg2)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_045i10y", msg2JSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (cc *SmartContract) Message_045i10y(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, "Message_045i10y")
	if err != nil {
		return err
	}

	// TODO: 待确认如何确认有权限的msp ID
	clientMspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMspID != msg.SendMspID {
		errorMessage := fmt.Sprintf("Msp denied")
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_045i10y", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stub.SetEvent("Message_045i10y", []byte("Message_045i10y has been done"))

	msg2, err := cc.ReadMsg(ctx, "Message_0r9lypd")
	if err != nil {
		return err
	}
	msg2.MsgState = ENABLE
	msg2JSON, err := json.Marshal(msg2)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_0r9lypd", msg2JSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (cc *SmartContract) Message_0r9lypd(ctx contractapi.TransactionContextInterface, fireflyTranID string, confirm bool) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, "Message_0r9lypd")
	if err != nil {
		return err
	}

	// 获取客户端MSP ID
	clientMspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMspID != msg.SendMspID {
		errorMessage := fmt.Sprintf("Msp denied")
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_0r9lypd", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stub.SetEvent("Message_0r9lypd", []byte("Message_0r9lypd has been done"))

	// 设置当前内存的确认字段
	cc.currentMemory.Confirm = confirm

	gtw, err := cc.ReadGtw(ctx, "ExclusiveGateway_106je4z")
	if err != nil {
		return err
	}
	gtw.GatewayState = ENABLE
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("ExclusiveGateway_106je4z", gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 调用ExclusiveGateway_106je4z函数
	cc.ExclusiveGateway_106je4z(ctx)

	return nil
}

func (c *SmartContract) ExclusiveGateway_106je4z(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	gtw, err := c.ReadGtw(ctx, "ExclusiveGateway_106je4z")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLE {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf("%s", errorMessage)
	}

	gtw.GatewayState = DONE
	sortedJson, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState("ExclusiveGateway_106je4z", sortedJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	stub.SetEvent("ExclusiveGateway_106je4z", []byte("ExclusiveGateway_106je4z has been done"))

	if c.currentMemory.Confirm {
		msg2, err := c.ReadMsg(ctx, "Message_1em0ee4")
		if err != nil {
			return err
		}
		msg2.MsgState = ENABLE
		sortedJson2, err := json.Marshal(msg2)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = stub.PutState("Message_1em0ee4", sortedJson2)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		gtw2, err := c.ReadGtw(ctx, "ExclusiveGateway_0hs3ztq")
		if err != nil {
			return err
		}
		gtw2.GatewayState = ENABLE
		sortedJson2, err := json.Marshal(gtw2)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = stub.PutState("ExclusiveGateway_0hs3ztq", sortedJson2)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		err = c.ExclusiveGateway_0hs3ztq(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SmartContract) Message_1em0ee4(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息
	msg, err := s.ReadMsg(ctx, "Message_1em0ee4")
	if err != nil {
		return err
	}

	// 获取客户端MSP ID
	clientMspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMspID != msg.SendMspID {
		errorMessage := fmt.Sprintf("Msp denied")
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Msg state %s does not allowed", msg.MessageID))
	}

	// 更新消息状态
	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_1em0ee4", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_1em0ee4", []byte("Message_1em0ee4 has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 更新消息状态为ENABLE
	err = s.ChangeMsgState(ctx, "Message_1nlagx2", ENABLE)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) Message_1nlagx2(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息
	msg, err := s.ReadMsg(ctx, "Message_1nlagx2")
	if err != nil {
		return err
	}

	// 获取客户端身份
	clientIdentity := ctx.GetClientIdentity()
	clientMspID, _ := clientIdentity.GetMSPID()

	// 检查权限
	if clientMspID != msg.SendMspID {
		errorMessage := fmt.Sprintf("Msp denied")
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Msp denied"))
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Msg state %s does not allowed", msg.MessageID))
	}

	// 更新消息状态
	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_1nlagx2", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_1nlagx2", []byte("Message_1nlagx2 has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 更新网关状态为ENABLE
	err = s.ChangeGtwState(ctx, "EventBasedGateway_1fxpmyn", ENABLE)
	if err != nil {
		return err
	}

	// 调用EventBasedGateway_1fxpmyn方法
	err = s.EventBasedGateway_1fxpmyn(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) EventBasedGateway_1fxpmyn(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// 读取网关状态
	gtw, err := s.ReadGtw(ctx, "EventBasedGateway_1fxpmyn")
	if err != nil {
		return err
	}

	// 检查网关状态
	if gtw.GatewayState != ENABLE {
		errorMessage := fmt.Sprintf("Gateway state %s does not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Gateway state %s does not allowed", gtw.GatewayID))
	}

	// 更新网关状态为DONE
	gtw.GatewayState = DONE

	// 序列化并保存网关状态
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EventBasedGateway_1fxpmyn", gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("EventBasedGateway_1fxpmyn", []byte("EventBasedGateway_1fxpmyn has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 更新消息状态为ENABLE
	err = s.ChangeMsgState(ctx, "Message_0o8eyir", ENABLE)
	if err != nil {
		return err
	}

	err = s.ChangeMsgState(ctx, "Message_1xm9dxy", ENABLE)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) Message_0o8eyir(ctx contractapi.TransactionContextInterface, cancel bool, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息状态
	msg, err := s.ReadMsg(ctx, "Message_0o8eyir")
	if err != nil {
		return err
	}

	// 检查客户端MspId
	clientIdentity := ctx.GetClientIdentity()
	clientMspId, _ := clientIdentity.GetMSPID()
	if clientMspId != msg.SendMspID {
		errorMessage := fmt.Sprintf("Msp denied")
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Msp denied"))
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Msg state %s does not allowed", msg.MessageID))
	}

	// 更新消息状态为DONE
	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息状态
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_0o8eyir", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_0o8eyir", []byte("Message_0o8eyir has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 更新消息状态为DISABLE
	err = s.ChangeMsgState(ctx, "Message_1xm9dxy", DISABLE)
	if err != nil {
		return err
	}

	// 更新网关状态为ENABLE
	err = s.ChangeGtwState(ctx, "ExclusiveGateway_0nzwv7v", ENABLE)
	if err != nil {
		return err
	}

	// 设置当前内存状态
	s.currentMemory.Cancel = cancel

	// 跳转到ExclusiveGateway_0nzwv7v
	return s.ExclusiveGateway_0nzwv7v(ctx)
}

func (s *SmartContract) ExclusiveGateway_0nzwv7v(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// 读取网关状态
	gtw, err := s.ReadGtw(ctx, "ExclusiveGateway_0nzwv7v")
	if err != nil {
		return err
	}

	// 检查网关状态
	if gtw.GatewayState != ENABLE {
		errorMessage := fmt.Sprintf("Gateway state %s does not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return errors.New(fmt.Sprintf("Gateway state %s does not allowed", gtw.GatewayID))
	}

	// 更新网关状态为DONE
	gtw.GatewayState = DONE

	// 序列化并保存网关状态
	gtwJSON, err := json.Marshal(gtw)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("ExclusiveGateway_0nzwv7v", gtwJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("ExclusiveGateway_0nzwv7v", []byte("ExclusiveGateway_0nzwv7v has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if s.currentMemory.Cancel {
		// 如果取消标志为true，则启用消息
		msg2, err := s.ReadMsg(ctx, "Message_1joj7ca")
		if err != nil {
			return err
		}
		msg2.MsgState = ENABLE

		// 序列化并保存消息状态
		msg2JSON, err := json.Marshal(msg2)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = stub.PutState("Message_1joj7ca", msg2JSON)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		// 启用结束事件
		event, err := s.ReadEvent(ctx, "EndEvent_08edp7f")
		if err != nil {
			return err
		}
		event.EventState = ENABLE

		// 序列化并保存事件状态
		eventJSON, err := json.Marshal(event)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = stub.PutState("EndEvent_08edp7f", eventJSON)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// 跳转到EndEvent_08edp7f
		return s.EndEvent_08edp7f(ctx)
	}

	return nil
}

func (s *SmartContract) Message_1joj7ca(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息状态
	msg, err := s.ReadMsg(ctx, "Message_1joj7ca")
	if err != nil {
		return err
	}

	// 获取客户端身份信息
	clientIdentity := ctx.GetClientIdentity()
	clientMspID, err := clientIdentity.GetMSPID()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 检查MSPID是否匹配
	if clientMspID != msg.SendMspID {
		errorMessage := "Msp denied"
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新消息状态为DONE
	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息状态
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_1joj7ca", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_1joj7ca", []byte("Message_1joj7ca has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 启用下一条消息状态
	return s.ChangeMsgState(ctx, "Message_1etcmvl", ENABLE)
}

func (s *SmartContract) Message_1etcmvl(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息状态
	msg, err := s.ReadMsg(ctx, "Message_1etcmvl")
	if err != nil {
		return err
	}

	// 获取客户端身份信息
	clientIdentity := ctx.GetClientIdentity()
	clientMspID, err := clientIdentity.GetMSPID()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 检查MSPID是否匹配
	if clientMspID != msg.SendMspID {
		errorMessage := "Msp denied"
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新消息状态为DONE
	msg.MsgState = DONE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息状态
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_1etcmvl", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_1etcmvl", []byte("Message_1etcmvl has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 完成事件
	event, _ := s.ReadEvent(ctx, "EndEvent_146eii4")
	event.EventState = ENABLE

	// 序列化并保存事件状态
	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EndEvent_146eii4", eventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 执行EndEvent_146eii4方法
	return s.EndEvent_146eii4(ctx)
}

func (s *SmartContract) Message_1xm9dxy(ctx contractapi.TransactionContextInterface, fireflyTranID string) error {
	stub := ctx.GetStub()

	// 读取消息状态
	msg, err := s.ReadMsg(ctx, "Message_1xm9dxy")
	if err != nil {
		return err
	}

	// 获取客户端身份信息
	clientIdentity := ctx.GetClientIdentity()
	clientMspID, err := clientIdentity.GetMSPID()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 检查MSPID是否匹配
	if clientMspID != msg.SendMspID {
		errorMessage := "Msp denied"
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 检查消息状态
	if msg.MsgState != ENABLE {
		errorMessage := fmt.Sprintf("Msg state %s does not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新消息状态为ENABLE
	msg.MsgState = ENABLE
	msg.FireflyTranID = fireflyTranID

	// 序列化并保存消息状态
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("Message_1xm9dxy", msgJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("Message_1xm9dxy", []byte("Message_1xm9dxy has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 完成事件
	event, _ := s.ReadEvent(ctx, "EndEvent_0366pfz")
	event.EventState = ENABLE

	// 序列化并保存事件状态
	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EndEvent_0366pfz", eventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 执行EndEvent_0366pfz方法
	return s.EndEvent_0366pfz(ctx)
}

func (s *SmartContract) EndEvent_08edp7f(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// 读取事件状态
	event, err := s.ReadEvent(ctx, "EndEvent_08edp7f")
	if err != nil {
		return err
	}

	// 检查事件状态
	if event.EventState != ENABLE {
		errorMessage := fmt.Sprintf("Event state %s does not allowed", event.EventID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新事件状态为DONE
	event.EventState = DONE

	// 序列化并保存事件状态
	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EndEvent_08edp7f", eventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("EndEvent_08edp7f", []byte("EndEvent_08edp7f has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (s *SmartContract) EndEvent_146eii4(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// 读取事件状态
	event, err := s.ReadEvent(ctx, "EndEvent_146eii4")
	if err != nil {
		return err
	}

	// 检查事件状态
	if event.EventState != ENABLE {
		errorMessage := fmt.Sprintf("Event state %s does not allowed", event.EventID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新事件状态为DONE
	event.EventState = DONE

	// 序列化并保存事件状态
	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EndEvent_146eii4", eventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("EndEvent_146eii4", []byte("EndEvent_146eii4 has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (s *SmartContract) EndEvent_0366pfz(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// 读取事件状态
	event, err := s.ReadEvent(ctx, "EndEvent_0366pfz")
	if err != nil {
		return err
	}

	// 检查事件状态
	if event.EventState != ENABLE {
		errorMessage := fmt.Sprintf("Event state %s does not allowed", event.EventID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	// 更新事件状态为DONE
	event.EventState = DONE

	// 序列化并保存事件状态
	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stub.PutState("EndEvent_0366pfz", eventJSON)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 设置事件
	err = stub.SetEvent("EndEvent_0366pfz", []byte("EndEvent_0366pfz has been done"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
