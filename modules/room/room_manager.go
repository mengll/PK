package room

type Roomer interface{
	JoinRoom() 		(bool,error)   	//加入房间
	OutRoom()  		(bool,error)   	//离开房间
	RoomLimit() 	(int ,error)	//房间人数限制
	IsFull()   		(bool,error)	//房间是否满员
	ExistsRoom() 	(bool,error)	//创建的房间是否还在
	NumRoom() 		(int ,error)    //当前房间的人数
	CloseRoom() 	error   		//关闭当前的房间

	CreateRoom()   (string,error)   // 返回当前的房间名称和错误信息
}

type RoomManager struct {}

/*
@return string the create room name
@return error  error message
@author mengll
*/
func (self *RoomManager) CreateRoom()(string,error){
	return "",nil
}

/*
@return error message
*/
func (self *RoomManager) CloseRoom() error{
	return nil
}

/*
@return int the num of the room
*/
func (self *RoomManager) NumRoom() (int,error){
	return 0,nil
}

/*
@bool 当前房间时候还存在
*/
func (self *RoomManager) ExistsRoom() (bool,error){
	return true,nil
}

/*
@return bool the room num is more than room limit
*/
func (self *RoomManager) IsFull() (bool,error){
	return true,nil
}

/*
@return the max user in the room
*/
func (self *RoomManager) RoomLimit() (int,error){
	return 0,nil
}

/*
@return bool is sucess out of the room
*/
func (self *RoomManager)OutRoom() (bool,error){
	return false,nil
}

/*
@bool 加入房间
*/
func (self *RoomManager)JoinRoom() (bool,error){
	return true,nil
}

//返回room manager
func NewRoomManager() Roomer{
	return new(RoomManager)
}
